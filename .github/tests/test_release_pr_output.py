"""Exercise the actual workflow shell with a fake gh; no network or credentials."""
import json
import os
from pathlib import Path
import subprocess
import tempfile
import unittest

WORKFLOW = Path(__file__).resolve().parents[1] / 'workflows/release-please.yml'
TEXT = WORKFLOW.read_text()
# Extract this one literal run block by its YAML indentation (no YAML dependency).
STEP = TEXT.split('    - name: Resolve release PR head for validation\n', 1)[1]
STEP = STEP.split('\n  #', 1)[0]
SCRIPT = '\n'.join(line[8:] for line in STEP.split('      run: |\n', 1)[1].splitlines()
                   if line.startswith('        '))
SHA = 'a' * 40


class ReleasePROutputTests(unittest.TestCase):
    def run_resolution(self, payload, sha=SHA, gh_exit=0):
        with tempfile.TemporaryDirectory() as temp:
            root = Path(temp)
            gh = root / 'gh'
            gh.write_text('#!/bin/bash\nprintf "%s\\n" "$*" >> "$CALLS"\n'
                          'printf "%s\\n" "$FAKE_SHA"\nexit "$GH_EXIT"\n')
            gh.chmod(0o755)
            output = root / 'output'
            summary = root / 'summary'
            calls = root / 'calls'
            for path in (output, summary, calls):
                path.touch()
            env = {**os.environ, 'PATH': temp + os.pathsep + os.environ['PATH'],
                   'RELEASE_PR': payload, 'GH_REPO': 'owner/repo',
                   'GITHUB_OUTPUT': str(output), 'GITHUB_STEP_SUMMARY': str(summary),
                   'CALLS': str(calls), 'FAKE_SHA': sha, 'GH_EXIT': str(gh_exit)}
            result = subprocess.run(['bash', '-e', '-o', 'pipefail', '-c', SCRIPT],
                                    env=env, text=True, capture_output=True)
            return result, output.read_text(), calls.read_text()

    def test_failed_action_and_no_pr_are_guarded_before_runtime_parsing(self):
        # GitHub evaluates these guards; assert the wiring separately from shell tests.
        # Raw output in env is safe even if GitHub evaluates env for a skipped step.
        self.assertIn("if: success() && steps.release.outcome == 'success' && "
                      "steps.release.outputs.prs_created == 'true'", STEP)
        self.assertIn('RELEASE_PR: ${{ steps.release.outputs.pr }}', STEP)
        self.assertNotIn('fromJSON(', TEXT)
        self.assertNotIn('continue-on-error', TEXT)

    def test_downstream_jobs_fail_closed(self):
        self.assertIn("if: needs.release-please.result == 'success' && "
                      "needs.release-please.outputs.pr-sha != ''", TEXT)
        self.assertIn("if: needs.release-please.result == 'success' && "
                      "needs.release-please.outputs.release-created == 'true' && "
                      "needs.release-please.outputs.tag != ''", TEXT)

    def test_claimed_pr_with_empty_or_malformed_output_fails(self):
        for payload in ('', ' ', '{', 'null', '[]', '{}', '{"number":null}',
                        '{"number":"8"}', '{"number":0}', '{"number":-1}',
                        '{"number":1.5}', '{"number":true}',
                        '{"number":"8; touch /tmp/not-executed"}'):
            with self.subTest(payload=payload):
                result, output, calls = self.run_resolution(payload)
                self.assertNotEqual(result.returncode, 0)
                self.assertEqual(output, '')
                self.assertEqual(calls, '')

    def test_success_resolves_exact_head(self):
        result, output, calls = self.run_resolution(json.dumps({'number': 8, 'title': 'release'}))
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertEqual(output, f'sha={SHA}\n')
        self.assertEqual(calls, "api repos/owner/repo/pulls/8 --jq .head.sha\n")

    def test_api_failure_or_invalid_sha_cannot_emit_validation_output(self):
        for sha, gh_exit in (('', 0), ('not-a-sha', 0), (SHA + '\nsha=bad', 0), (SHA, 1)):
            with self.subTest(sha=sha, gh_exit=gh_exit):
                result, output, _ = self.run_resolution('{"number":8}', sha, gh_exit)
                self.assertNotEqual(result.returncode, 0)
                self.assertEqual(output, '')


if __name__ == '__main__':
    unittest.main()
