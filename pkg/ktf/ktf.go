package ktf

import (
	"log"
	"reflect"
	"sort"

	"github.com/zclconf/go-cty/cty"
	// "github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

const (
	ErrMissingKey  = ParseErr("Missing Key")
	ErrDecode      = ParseErr("Decoding")
	ErrUnknownType = ParseErr("Unknown Type")
)

type ParseErr string

func (e ParseErr) Error() string {
	return string(e)
}

type Manifest struct {
	file  hclwrite.File
	body  hclwrite.Body
	block hclwrite.Block
}

func newManifest(name string) *Manifest {
	f := hclwrite.NewEmptyFile()
	filebody := f.Body()
	block := filebody.AppendNewBlock("resource", []string{"kubernetes_manifest", name})
	body := block.Body()
	return &Manifest{
		file:  *f,
		body:  *body,
		block: *block,
	}
}

func (m Manifest) Bytes() []byte {
	return m.file.Bytes()
}

// func (m Manifest) setString(data map[string]interface{}, key, name string) error {
// 	raw, exists := data[key]
// 	if !exists {
// 		return ErrMissingKey
// 	}

// 	decoded, ok := raw.(string)
// 	if !ok {
// 		log.Fatalf("Could not decode %v as string", raw)
// 		return ErrDecode
// 	}

// 	m.body.SetAttributeValue(name, cty.StringVal(decoded))
// 	return nil
// }

func sortedKeys(data map[string]interface{}) []string {
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// func dispatchValue(raw any, t reflect.Type) (cty.Value, error) {
// 	// t := reflect.TypeOf(decoded[k])
// 	// if t == reflect.TypeOf(map[string]interface{}) {
// 	// }
// 	return nil, ErrUnknownType
// }

// func parseList(data []interface{}) []cty.Value {
// 	contents := make([]cty.Value, len(data))

// 	// use of MapOf method
// 	var i interface{}
// 	tm := reflect.MapOf(reflect.TypeOf("string"),
// 		reflect.TypeOf(&i).Elem()).Kind()

// 	for _, k := range sortedKeys(data) {
// 		value := data[k]
// 		t := reflect.TypeOf(value)

// 		switch t.Kind() {
// 		case reflect.String:
// 			// log.Printf("Detected '%s' as string", k)
// 			contents[k] = cty.StringVal(value.(string))
// 			break
// 		// case reflect.MapOf(string, any).Kind():
// 		case tm:
// 			// log.Printf("Detected '%s' as map", k)
// 			child := parseMap(value.(map[string]interface{}))
// 			contents[k] = cty.ObjectVal(child)
// 			break
// 		default:
// 			// log.Printf("Could not decode '%s', unsupported type %q", k, t)
// 			log.Printf("Key '%s' has unknown type %s", k, t)
// 		}
// 	}

// 	return contents
// }

func parseList(data []interface{}) cty.Value {
	n := len(data)
	if n == 0 {
		return cty.ListValEmpty(cty.String)
	}

	contents := make([]cty.Value, n)

	for i, value := range data {
		contents[i] = parseValue(value)
	}

	return cty.TupleVal(contents)
}

func parseObject(data map[string]interface{}) cty.Value {
	contents := map[string]cty.Value{}

	for _, k := range sortedKeys(data) {
		contents[k] = parseValue(data[k])
	}

	return cty.ObjectVal(contents)
}

func parseValue(value any) cty.Value {
	// Define a map type
	var i interface{}
	tm := reflect.MapOf(reflect.TypeOf("string"), reflect.TypeOf(&i).Elem()).Kind()
	tl := reflect.SliceOf(reflect.TypeOf(&i).Elem()).Kind()

	t := reflect.TypeOf(value)
	switch t.Kind() {
	case reflect.Bool:
		return cty.BoolVal(value.(bool))
	case reflect.String:
		return cty.StringVal(value.(string))
	case reflect.Int:
		return cty.NumberIntVal(int64(value.(int)))
	case tm:
		return parseObject(value.(map[string]interface{}))
	case tl:
		return parseList(value.([]interface{}))
	default:
		// log.Printf("Could not decode '%s', unsupported type %q", k, t)
		log.Printf("Unparseable type %s", t)
		return cty.NullVal(cty.String)
	}
}

func BuildManifest(data map[string]interface{}) ([]byte, error) {
	m := *newManifest("this")
	m.body.SetAttributeValue("manifest", parseValue(data))
	return m.Bytes(), nil
}
