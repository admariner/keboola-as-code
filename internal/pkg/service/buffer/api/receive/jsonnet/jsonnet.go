// Package jsonnet provides Jsonnet functions used by the Buffer API import endpoint.
//
// # Jsonnet Functions
//
//	Ip() string - formatted IP address of the client
//	Now() string - current datetime in UTC, in fixed length DefaultTimeFormat, for example "2006-01-01T08:04:05.123Z"
//	Now("%Y-%m-%d") string - current datetime in UTC timezone, in a custom "strftime" compatible format
//	HeaderStr() string - request headers as a string, each lines contains one "Header: value", the lines are sorted alphabetically
//	Header() object - all headers as a JSON object
//	Header("Header-Name") string - value of the header, if it is not found, then an error occurs and the record is not saved
//	Header("Header-Name", "default value") string - value of the header or the default value
//	BodyStr() string - raw request body as a string
//	Body() object - parsed JSON/form-data body as a JSON object
//	Body("some.key1[2].key2") mixed - value of the path in the parsed body, if it is not found, then an error occurs and the record is not saved
//	Body("some.key1[2].key2", "default value") mixed - value of the path in the parsed body or the default value
package jsonnet

import (
	"net/http"

	jsonnetLib "github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
	"github.com/lestrrat-go/strftime"

	"github.com/keboola/keboola-as-code/internal/pkg/encoding/jsonnet"
	"github.com/keboola/keboola-as-code/internal/pkg/service/buffer/api/receive/receivectx"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/errors"
)

const (
	// DefaultTimeFormat with fixed length, so it can be used for lexicographic sorting in the target database.
	// Value is RFC3339 and ISO8601 compatible.
	DefaultTimeFormat   = "%Y-%m-%dT%H:%M:%S.%fZ"
	ThrowErrOnUndefined = "<<~~errOnUndefined~~>>"
	nowFnName           = "_now"
	headersMapFnName    = "_headersMap"
	headerFnName        = "_header"
	bodyMapFnName       = "_bodyMap"
	bodyPathFnName      = "_bodyPath"
)

func Evaluate(reqCtx *receivectx.Context, template string) (string, error) {
	c := jsonnet.NewContext().WithCtx(reqCtx.Ctx).WithPretty(false)
	RegisterFunctions(c, reqCtx)
	out, err := jsonnet.Evaluate(template, c)
	if err != nil {
		var jsonnetErr jsonnetLib.RuntimeError
		if errors.As(err, &jsonnetErr) {
			// Trim "jsonnet error: RUNTIME ERROR: "
			err = errors.Wrap(err, jsonnetErr.Msg)
		}
	}
	return out, err
}

type Validator struct {
	ctx *jsonnet.Context
}

func NewValidator() *Validator {
	ctx := jsonnet.NewContext()
	// we don't actually call these functions, we only register them for enumeration later,
	// so the context can be empty, because it will never be used.
	RegisterFunctions(ctx, &receivectx.Context{})
	return &Validator{ctx}
}

func (v *Validator) Validate(template string) error {
	return v.ctx.Validate(template)
}

func RegisterFunctions(c *jsonnet.Context, reqCtx *receivectx.Context) {
	// Global functions
	c.NativeFunctionWithAlias(ipFn("Ip", reqCtx))
	c.NativeFunctionWithAlias(headerStrFn("HeaderStr", reqCtx))
	c.NativeFunctionWithAlias(bodyStrFn("BodyStr", reqCtx))
	c.GlobalBinding("Now", nowFn())
	c.GlobalBinding("Header", headerFn())
	c.GlobalBinding("Body", bodyFn())

	// Internal functions
	// Optional function parameters cannot be specified directly by the Go SDK,
	// so these partial functions are used by global functions above.
	c.NativeFunction(nowInternalFn(reqCtx))
	c.NativeFunction(headersMapInternalFn(reqCtx))
	c.NativeFunction(headerValueInternalFn(reqCtx))
	c.NativeFunction(bodyMapInternalFn(reqCtx))
	c.NativeFunction(bodyPathInternalFn(reqCtx))
}

func ipFn(fnName string, reqCtx *receivectx.Context) *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name: fnName,
		Func: func(params []interface{}) (any, error) {
			if len(params) != 0 {
				return nil, errors.Errorf("no parameter expected, found %d", len(params))
			}
			return jsonnet.ValueToJSONType(reqCtx.IP.String()), nil
		},
	}
}

func headerStrFn(fnName string, reqCtx *receivectx.Context) *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name: fnName,
		Func: func(params []interface{}) (any, error) {
			if len(params) != 0 {
				return nil, errors.Errorf("no parameter expected, found %d", len(params))
			}

			return jsonnet.ValueToJSONType(reqCtx.HeadersStr()), nil
		},
	}
}

func bodyStrFn(fnName string, reqCtx *receivectx.Context) *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name: fnName,
		Func: func(params []interface{}) (any, error) {
			if len(params) != 0 {
				return nil, errors.Errorf("no parameter expected, found %d", len(params))
			}
			return jsonnet.ValueToJSONType(reqCtx.Body), nil
		},
	}
}

func nowFn() ast.Node {
	formatParam := ast.Identifier("format")
	formatVar := &ast.Var{Id: formatParam}
	defaultFormat := jsonnet.ValueToLiteral(DefaultTimeFormat)

	var node ast.Node = &ast.Function{
		Parameters: []ast.Parameter{{Name: formatParam, DefaultArg: defaultFormat}},
		Body:       applyNativeFn(nowFnName, formatVar),
		NodeBase:   ast.NodeBase{FreeVars: []ast.Identifier{"std"}},
	}

	return node
}

// headerFn - if header == "" then std.native("_headersMap") else std.native("_header", header, defaultValue).
func headerFn() ast.Node {
	headerParam := ast.Identifier("header")
	headerVar := &ast.Var{Id: headerParam}
	defaultValParam := ast.Identifier("default")
	defaultValVar := &ast.Var{Id: defaultValParam}
	emptyStr := jsonnet.ValueToLiteral("")
	throwErrOnUndefined := jsonnet.ValueToLiteral(ThrowErrOnUndefined)
	var node ast.Node = &ast.Function{
		Parameters: []ast.Parameter{{Name: headerParam, DefaultArg: emptyStr}, {Name: defaultValParam, DefaultArg: throwErrOnUndefined}},
		Body: &ast.Conditional{
			Cond:        &ast.Binary{Right: headerVar, Left: emptyStr, Op: ast.BopManifestEqual},
			BranchTrue:  applyNativeFn(headersMapFnName),
			BranchFalse: applyNativeFn(headerFnName, headerVar, defaultValVar),
		},
		NodeBase: ast.NodeBase{FreeVars: []ast.Identifier{"std"}},
	}

	return node
}

// bodyFn - if path == "" then std.native("_bodyMap") else std.native("_bodyPath", path, defaultValue).
func bodyFn() ast.Node {
	pathParam := ast.Identifier("path")
	pathVar := &ast.Var{Id: pathParam}
	defaultValParam := ast.Identifier("default")
	defaultValVar := &ast.Var{Id: defaultValParam}
	emptyStr := jsonnet.ValueToLiteral("")
	throwErrOnUndefined := jsonnet.ValueToLiteral(ThrowErrOnUndefined)
	var node ast.Node = &ast.Function{
		Parameters: []ast.Parameter{{Name: pathParam, DefaultArg: emptyStr}, {Name: defaultValParam, DefaultArg: throwErrOnUndefined}},
		Body: &ast.Conditional{
			Cond:        &ast.Binary{Right: &ast.Var{Id: pathParam}, Left: emptyStr, Op: ast.BopManifestEqual},
			BranchTrue:  applyNativeFn(bodyMapFnName),
			BranchFalse: applyNativeFn(bodyPathFnName, pathVar, defaultValVar),
		},
		NodeBase: ast.NodeBase{FreeVars: []ast.Identifier{"std"}},
	}

	return node
}

func nowInternalFn(reqCtx *receivectx.Context) *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name:   nowFnName,
		Params: ast.Identifiers{"format"},
		Func: func(params []interface{}) (any, error) {
			if len(params) != 1 {
				return nil, errors.Errorf("one parameter expected, found %d", len(params))
			}

			format, ok := params[0].(string)
			if !ok {
				return nil, errors.New("parameter must be a string")
			}

			formatter, err := strftime.New(format, strftime.WithMilliseconds('f'))
			if err != nil {
				return nil, errors.Errorf(`datetime format "%s" is invalid: %w`, format, err)
			}

			return jsonnet.ValueToJSONType(formatter.FormatString(reqCtx.Now.UTC())), nil
		},
	}
}

func headersMapInternalFn(reqCtx *receivectx.Context) *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name: headersMapFnName,
		Func: func(params []interface{}) (any, error) {
			if len(params) != 0 {
				return nil, errors.Errorf("no parameter expected, found %d", len(params))
			}
			return jsonnet.ValueToJSONType(reqCtx.HeadersMap()), nil
		},
	}
}

func headerValueInternalFn(reqCtx *receivectx.Context) *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name:   headerFnName,
		Params: ast.Identifiers{"path", "default"},
		Func: func(params []interface{}) (any, error) {
			if len(params) != 2 {
				return nil, errors.Errorf("two parameters expected, found %d", len(params))
			}

			name, ok := params[0].(string)
			defaultVal := params[1]
			if !ok {
				return nil, errors.New("parameter must be a string")
			}

			value := reqCtx.Headers.Get(name)
			if value == "" {
				if defaultVal == ThrowErrOnUndefined {
					return nil, errors.Errorf(`header "%s" not found`, http.CanonicalHeaderKey(name))
				} else {
					return defaultVal, nil
				}
			}
			return jsonnet.ValueToJSONType(value), nil
		},
	}
}

func bodyMapInternalFn(reqCtx *receivectx.Context) *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name: bodyMapFnName,
		Func: func(params []interface{}) (any, error) {
			if len(params) != 0 {
				return nil, errors.Errorf("no parameter expected, found %d", len(params))
			}
			bodyMap, err := reqCtx.BodyMap()
			if err != nil {
				return nil, err
			}
			return jsonnet.ValueToJSONType(bodyMap), nil
		},
	}
}

func bodyPathInternalFn(reqCtx *receivectx.Context) *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name:   bodyPathFnName,
		Params: ast.Identifiers{"path", "default"},
		Func: func(params []interface{}) (any, error) {
			if len(params) != 2 {
				return nil, errors.Errorf("two parameters expected, found %d", len(params))
			}

			path, ok := params[0].(string)
			defaultVal := params[1]
			if !ok {
				return nil, errors.New("first parameter must be a string")
			}

			bodyMap, err := reqCtx.BodyMap()
			if err != nil {
				return nil, err
			}

			val, _, err := bodyMap.GetNested(path)
			if err != nil {
				if defaultVal == ThrowErrOnUndefined {
					return nil, errors.Errorf(`path "%s" not found in the body`, path)
				} else {
					return defaultVal, nil
				}
			}
			return jsonnet.ValueToJSONType(val), nil
		},
	}
}

func applyNativeFn(fnName string, args ...ast.Node) ast.Node {
	var freeVars []ast.Identifier
	var fnArgs []ast.CommaSeparatedExpr
	for _, item := range args {
		fnArgs = append(fnArgs, ast.CommaSeparatedExpr{Expr: item})
		// Build list of the freeVars manually, so we can skip parser.PreprocessAst step, it is faster
		if v, ok := item.(*ast.Var); ok {
			freeVars = append(freeVars, v.Id)
		}
	}

	nativeFn := &ast.Apply{
		Target: &ast.Index{
			Target: &ast.Var{Id: "std"},
			Index:  &ast.LiteralString{Value: "native"},
		},
		Arguments: ast.Arguments{Positional: []ast.CommaSeparatedExpr{{Expr: &ast.LiteralString{Value: fnName}}}},
	}

	return &ast.Apply{
		NodeBase:  ast.NodeBase{FreeVars: freeVars},
		Target:    nativeFn,
		Arguments: ast.Arguments{Positional: fnArgs},
	}
}
