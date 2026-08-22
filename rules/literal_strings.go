package rules

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

func walkLiteralStrings(runner tflint.Runner, visit func(value string, rng hcl.Range) error) error {
	var visitErr error

	diags := runner.WalkExpressions(tflint.ExprWalkFunc(func(expr hcl.Expression) hcl.Diagnostics {
		if visitErr != nil {
			return nil
		}

		if node, ok := expr.(hclsyntax.Node); ok {
			literal, ok := node.(*hclsyntax.LiteralValueExpr)
			if !ok || !literal.Val.IsKnown() || literal.Val.IsNull() || literal.Val.Type() != cty.String {
				return nil
			}
			visitErr = visit(literal.Val.AsString(), literal.SrcRange)
			return nil
		}

		var value cty.Value
		if err := runner.EvaluateExpr(expr, &value, nil); err != nil {
			visitErr = err
			return nil
		}
		if !value.IsKnown() || value.IsNull() {
			return nil
		}

		visitJSONStrings(value, expr.Range(), func(value string, rng hcl.Range) error {
			return visit(value, rng)
		}, &visitErr)
		return nil
	}))

	if diags.HasErrors() {
		return diags
	}
	return visitErr
}

func visitJSONStrings(value cty.Value, rng hcl.Range, visit func(value string, rng hcl.Range) error, visitErr *error) {
	if *visitErr != nil || !value.IsKnown() || value.IsNull() {
		return
	}
	if value.Type() == cty.String {
		*visitErr = visit(value.AsString(), rng)
		return
	}
	if value.Type().IsTupleType() || value.Type().IsListType() || value.Type().IsSetType() {
		it := value.ElementIterator()
		for it.Next() {
			_, element := it.Element()
			visitJSONStrings(element, rng, visit, visitErr)
			if *visitErr != nil {
				return
			}
		}
		return
	}
	if value.Type().IsObjectType() || value.Type().IsMapType() {
		it := value.ElementIterator()
		for it.Next() {
			_, element := it.Element()
			visitJSONStrings(element, rng, visit, visitErr)
			if *visitErr != nil {
				return
			}
		}
	}
}

func evaluatedBool(runner tflint.Runner, expr hcl.Expression) (bool, bool, error) {
	var value cty.Value
	if err := runner.EvaluateExpr(expr, &value, nil); err != nil {
		return false, false, err
	}
	if !value.IsKnown() || value.IsNull() || value.Type() != cty.Bool {
		return false, false, nil
	}
	unmarked, _ := value.Unmark()
	return unmarked.True(), true, nil
}

func evaluatedString(runner tflint.Runner, expr hcl.Expression) (string, bool, error) {
	var value cty.Value
	if err := runner.EvaluateExpr(expr, &value, nil); err != nil {
		return "", false, err
	}
	if !value.IsKnown() || value.IsNull() || value.Type() != cty.String {
		return "", false, nil
	}
	unmarked, _ := value.Unmark()
	return unmarked.AsString(), true, nil
}
