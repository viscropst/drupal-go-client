package drupal_go_client

import (
	"encoding/json"
	"errors"
	"fmt"
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

func entityStubMarshal(entity EntityCompatible, stubs StubConfigs) ([]byte, error) {
	stub, ok := stubs[entity.Type()]
	if !ok {
		return nil, errors.New("stub config not existed")
	}

	if stub.NoMappingMode == "" {
		stub.NoMappingMode = NoMappingModeOriginal
	}

	resMap := make(map[string]interface{})
	for s, d := range stub.Mapping {
		if field, err := entity.GetField(s); err != nil {
			return nil, fmt.Errorf("entity get field: %v", err)
		} else {
			if r, err := getEntityFieldValue(d.Type, field); err != nil {
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
				r, err := getEntityFieldValue(f.t.String(), field)
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

	if res, err := json.Marshal(resMap); err != nil {
		return nil, fmt.Errorf("result map marshal: %v", err)
	} else {
		return res, nil
	}
}

func getEntityFieldValue(t string, field *Field) (interface{}, error) {
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
			return nil, fmt.Errorf("field to int32: %v", err)
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
	default:
		r = field.Raw()
	}

	return r, nil
}
