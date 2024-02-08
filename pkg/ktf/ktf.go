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

func (m Manifest) setString(data map[string]interface{}, key, name string) error {
	raw, exists := data[key]
	if !exists {
		return ErrMissingKey
	}

	decoded, ok := raw.(string)
	if !ok {
		log.Fatalf("Could not decode %v as string", raw)
		return ErrDecode
	}

	m.body.SetAttributeValue(name, cty.StringVal(decoded))
	return nil
}

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

func (m Manifest) setObject(data map[string]interface{}, key, name string) error {
	raw, exists := data[key]
	if !exists {
		return ErrMissingKey
	}

	decoded, ok := raw.(map[string]interface{})
	if !ok {
		log.Fatalf("Could not decode %v as string", raw)
		return ErrDecode
	}

	contents := map[string]cty.Value{}

	for _, k := range sortedKeys(decoded) {
		value := decoded[k]
		t := reflect.TypeOf(value)

		switch t.Kind() {
		case reflect.String:
			// log.Printf("Detected '%s' as string", k)
			contents[k] = cty.StringVal(value.(string))
			break
		case reflect.MapOf(string, interface):
			log.Printf("Detected '%s' as map", k)

		default:
			// log.Printf("Could not decode '%s', unsupported type %q", k, t)
			log.Printf("Key '%s' has unknown type %s", k, t)
		}
	}

	m.body.SetAttributeValue(name, cty.ObjectVal(contents))
	return nil
}

func BuildManifest(data map[string]interface{}) ([]byte, error) {
	m := *newManifest("this")

	// body := m.block.Body()
	var err error

	err = m.setString(data, "apiVersion", "apiVersion")
	if err != nil {
		return nil, err
	}

	err = m.setString(data, "kind", "kind")
	if err != nil {
		return nil, err
	}

	err = m.setObject(data, "metadata", "metadata")
	if err != nil {
		return nil, err
	}

	// annotations := map[string]cty.Value{"controller-gen.kubebuilder.io/version": cty.StringVal("v0.13.0")}
	// metadata := map[string]cty.Value{
	// 	"name":        cty.StringVal("workerpools.workers.spacelift.io"),
	// 	"annotations": cty.MapVal(annotations),
	// }
	// body.SetAttributeValue("metadata", cty.ObjectVal(metadata))

	// body.SetAttributeValue("kind", cty.StringVal(data["kind"]))
	// body.SetAttributeValue("metadata", cty.StringVal("boop"))
	// body.SetAttributeValue("spec", cty.StringVal("boop"))

	// bazBody.SetAttributeValue("foo", cty.NumberIntVal(10))
	// bazBody.SetAttributeValue("beep", cty.StringVal("boop"))
	// bazBody.SetAttributeValue("baz", cty.ListValEmpty(cty.String))

	// keys := make([]string, 0, len(data))
	// for k := range data {
	// 	keys = append(keys, k)
	// }
	// sort.Strings(keys)

	// for k := range keys {
	// 	body.append
	// }
	// log.Printf("t1: %s\n", reflect.TypeOf(data["apiVersion"]))

	return m.Bytes(), nil
}
