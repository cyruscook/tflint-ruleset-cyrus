package main

import (
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/plugin/host2plugin"
)

// Query the built executable through the same protocol used by TFLint. This
// catches linker flags that silently target the wrong Go symbol.
func TestPluginVersion(t *testing.T) {
	for _, tc := range []struct {
		name, ldflags, want string
	}{
		{name: "source", want: version},
		{name: "release", ldflags: "-s -w -X main.version=1.2.3", want: "1.2.3"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			binary := filepath.Join(t.TempDir(), "tflint-ruleset-cyrus.exe")
			cmd := exec.Command("go", "build", "-ldflags="+tc.ldflags, "-o", binary, ".")
			if output, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("build plugin: %v\n%s", err, output)
			}
			client := host2plugin.NewClient(&host2plugin.ClientOpts{Cmd: exec.Command(binary)})
			t.Cleanup(client.Kill)
			rpc, err := client.Client()
			if err != nil {
				t.Fatal(err)
			}
			raw, err := rpc.Dispense("ruleset")
			if err != nil {
				t.Fatal(err)
			}
			got, err := raw.(*host2plugin.Client).RuleSetVersion()
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.want {
				t.Fatalf("plugin version = %q, want %q", got, tc.want)
			}
		})
	}
}
