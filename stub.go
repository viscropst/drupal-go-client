package drupal_go_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/jsonapi"
	"reflect"
)

const (
	NoMappingModeIgnore   = "ignore"
	NoMappingModeOriginal = "original"
)

//{
//   "entity_type": "node",
//   "bundle": "article",
//   "no_mapping_mode": "ignore|original"
//   "mapping": {
//     "field_image": {
//        "type": "file"
//        "name": "image"
//     },
//     "field_category": {
//        "type": "string"
//        "name": "category"
//     }
//   }
// }

type StubFieldMapper struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type Stub struct {
	EntityType    string `json:"entity_type"`
	Bundle        string `json:"bundle"`
	NoMappingMode string `json:"no_mapping_mode,omitempty"`

	Mapping map[string]StubFieldMapper `json:"mapping"`
}

func (s *Stub) Validate() error {
	for k, v := range s.Mapping {
		support := false
		for _, st := range supportDataTypes {
			if st == v.Type {
				support = true
				break
			}
		}

		if !support {
			return fmt.Errorf("%s type %s not support", k, v.Type)
		}
	}
	return nil
}

type StubConfigs map[string]Stub

func NewStubConfigsFromJSON(bytes []byte) (*StubConfigs, error) {
	sc := new(StubConfigs)
	err := json.Unmarshal(bytes, sc)
	if err != nil {
		return nil, err
	}

	for _, v := range *sc {
		if err = v.Validate(); err != nil {
			return nil, fmt.Errorf("stub validate: %v", err)
		}
	}

	return sc, nil
}

type EntityStubMarshaler interface {
	Marshal(v *Stub) ([]byte, error)
}

func entityStubUnmarshal(b []byte, stubs *StubConfigs) (*jsonapi.OnePayload, error) {
	srcMap := make(map[string]interface{})
	if err := json.Unmarshal(b, &srcMap); err != nil {
		return nil, fmt.Errorf("unarmshal src map: %v", err)
	}

	entityType, ok := srcMap["type"]
	if !ok {
		return nil, errors.New("src []byte not include entity type")
	}

	t, ok := entityType.(string)
	if !ok {
		return nil, errors.New("entity type in src must be string")
	}

	stub, ok := (*stubs)[t]
	if !ok {
		return nil, errors.New("stub config not existed")
	}

	payload := &jsonapi.OnePayload{
		Data: &jsonapi.Node{
			Type: t,
		},
	}

	entityID, ok := srcMap["id"]
	if ok {
		payload.Data.ID, ok = entityID.(string)
		if !ok {
			return nil, errors.New("entity id in src must be string")
		}
	}

	fields := make(map[string]interface{})
	for k, v := range srcMap {
		noMapping := true
		for sk, sv := range stub.Mapping {
			if k == sv.Name {
				noMapping = false
				fields[sk] = v
				break
			}
		}

		if noMapping {
			fields[k] = v
		}
	}

	delete(fields, "id")
	delete(fields, "type")

	attrs := make(map[string]interface{})
	relationships := make(map[string]interface{})
	for k, v := range fields {
		vt := reflect.TypeOf(v)

		switch vt.Kind() {
		case reflect.Map:
			if n, err := shallowNodeFromMap(v); err == nil {
				relationships[k] = jsonapi.OnePayload{
					Data: n,
				}
				break
			} else {
				// map but can't unmarshal to shallow node, return raw
				attrs[k] = v
			}
		case reflect.Slice:
			s := reflect.ValueOf(v)

			data := make([]*jsonapi.Node, 0)
			for i := 0; i < s.Len(); i++ {
				if n, err := shallowNodeFromMap(s.Index(i).Interface()); err == nil {
					data = append(data, n)
				}
			}

			if len(data) == s.Len() {
				relationships[k] = jsonapi.ManyPayload{
					Data: data,
				}
				break
			}
		default:
			attrs[k] = v
		}
	}

	payload.Data.Attributes = attrs
	payload.Data.Relationships = relationships

	return payload, nil
}

func shallowNodeFromMap(v interface{}) (*jsonapi.Node, error) {
	vt := reflect.TypeOf(v)
	vv := reflect.ValueOf(v)

	if vt.Kind() != reflect.Map {
		return nil, errors.New("value is not a Map")
	}

	for _, k := range vv.MapKeys() {
		if k.String() != "id" && k.String() != "type" {
			return nil, errors.New("jsonapi node must include id and type")
		}
	}

	shallowNode := new(jsonapi.Node)
	buf := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buf).Encode(v); err != nil {
		return nil, err
	}

	if err := json.NewDecoder(buf).Decode(shallowNode); err != nil {
		return nil, err
	}
	return shallowNode, nil
}

func entityStubTransform(entity EntityCompatible, stubs *StubConfigs) (map[string]interface{}, error) {
	stub, ok := (*stubs)[entity.Type()]
	if !ok {
		return nil, errors.New("stub config not existed")
	}

	if stub.NoMappingMode == "" {
		stub.NoMappingMode = NoMappingModeOriginal
	}

	resMap := make(map[string]interface{})
	resMap["id"] = entity.ID()
	resMap["type"] = entity.Type()

	for s, d := range stub.Mapping {
		if field, err := entity.GetField(s); err != nil {
			continue
		} else {
			if r, err := getEntityFieldValue(d.Type, field, stubs); err != nil {
				return nil, fmt.Errorf("get entity field value: %v", err)
			} else {
				resMap[d.Name] = r
			}
		}
	}

	switch stub.NoMappingMode {
	case NoMappingModeOriginal:
		schema, err := entity.GetSchema()
		if err != nil {
			return nil, fmt.Errorf("get schema: %v", err)
		}

		for _, f := range schema.fields {
			noMapping := true
			for s, _ := range stub.Mapping {
				if f.name == s {
					noMapping = false
					break
				}
			}

			if noMapping {
				field, err := entity.GetField(f.name)
				if err != nil {
					return nil, fmt.Errorf("no mapping get field: %v", err)
				}
				r, err := getEntityFieldValue(f.t.String(), field, stubs)
				if err != nil {
					return nil, fmt.Errorf("no mapping get entity field value: %v", err)
				}
				resMap[f.name] = r
			}
		}
		break
	case NoMappingModeIgnore:
		break
	}

	return resMap, nil
}

func entityStubMarshal(entity EntityCompatible, stubs *StubConfigs) ([]byte, error) {
	resMap, err := entityStubTransform(entity, stubs)
	if err != nil {
		return nil, fmt.Errorf("transform: %v", err)
	}
	if res, err := json.Marshal(resMap); err != nil {
		return nil, fmt.Errorf("result map marshal: %v", err)
	} else {
		return res, nil
	}
}

func getEntityFieldValue(t string, field *Field, stub *StubConfigs) (interface{}, error) {
	var err error
	var r interface{}
	switch t {
	case "string":
		if r, err = field.String(); err != nil {
			return nil, fmt.Errorf("field to string: %v", err)
		}
		break
	case "int32":
		if r, err = field.Int32(); err != nil {
			return nil, fmt.Errorf("field to int32: %v", err)
		}
		break
	case "int64":
		if r, err = field.Int64(); err != nil {
			return nil, fmt.Errorf("field to int64: %v", err)
		}
		break
	case "float32":
		if r, err = field.Float32(); err != nil {
			return nil, fmt.Errorf("field to float32: %v", err)
		}
		break
	case "float64":
		if r, err = field.Float64(); err != nil {
			return nil, fmt.Errorf("field to float64: %v", err)
		}
		break
	case "file":
		if r, err = field.File(); err != nil {
			return nil, fmt.Errorf("field to file: %v", err)
		}
		break
	case "bool":
		if r, err = field.Bool(); err != nil {
			return nil, fmt.Errorf("field to bool: %v", err)
		}
		break
	case "relation":
		if r, err = field.Relation(true, stub); err != nil {
			return nil, fmt.Errorf("field to relation: %v", err)
		}
		return r, nil
	default:
		r = field.Raw()
		return r, nil
	}

	if reflect.ValueOf(r).IsNil() {
		return nil, nil
	}

	return reflect.ValueOf(r).Elem().Interface(), nil
}
