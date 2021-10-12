package drupal_go_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/jsonapi"
	"reflect"
)

var supportDataTypes = []string{"string", "int32", "int64", "float32", "float64", "file", "bool", "raw", "relation"}

type Field struct {
	raw        interface{}
	name       string
	refPayload *jsonapi.OnePayload

	IsRelationship bool
}

type FieldType struct {
	t    reflect.Type
	name string
}

type File struct {
	FID      string `json:"fid"`
	LangCode string `json:"langcode"`
	Filename string `json:"filename"`
	URI      struct {
		Value string `json:"value"`
		URL   string `json:"url"`
	} `json:"uri"`
	FileMime string  `json:"filemime"`
	Filesize float64 `json:"filesize"`
	Status   bool    `json:"status"`
	Created  string  `json:"created"`
	Changed  string  `json:"changed"`
}

type Link struct {
	URI   string `json:"uri"`
	Title string `json:"title"`
}

type Body struct {
	Value     string `json:"value"`
	Format    string `json:"format"`
	Processed string `json:"processed"`
	Summary   string `json:"summary"`
}

func (f *Field) Raw() interface{} {
	return f.raw
}

func (f *Field) String() (string, error) {
	if f.IsRelationship {
		return "", fmt.Errorf("field is relatiionship")
	}
	s, ok := f.raw.(string)
	if ok {
		return s, nil
	}

	return "", fmt.Errorf("field is not a string")
}

func (f *Field) Int32() (int32, error) {
	if f.IsRelationship {
		return 0, fmt.Errorf("field is relatiionship")
	}
	s, ok := f.raw.(int32)
	if ok {
		return s, nil
	}

	return 0, fmt.Errorf("field is not int32")
}

func (f *Field) Int64() (int64, error) {
	if f.IsRelationship {
		return 0, fmt.Errorf("field is relatiionship")
	}
	s, ok := f.raw.(int64)
	if ok {
		return s, nil
	}

	return 0, fmt.Errorf("field is not int64")
}

func (f *Field) Bool() (bool, error) {
	if f.IsRelationship {
		return false, fmt.Errorf("field is relatiionship")
	}
	s, ok := f.raw.(bool)
	if ok {
		return s, nil
	}

	return false, fmt.Errorf("field is not bool")
}

func (f *Field) Float32() (float32, error) {
	if f.IsRelationship {
		return 0, fmt.Errorf("field is relatiionship")
	}
	s, ok := f.raw.(float32)
	if ok {
		return s, nil
	}

	return 0, fmt.Errorf("field is not float32")
}

func (f *Field) Float64() (float64, error) {
	if f.IsRelationship {
		return 0, fmt.Errorf("field is relatiionship")
	}
	s, ok := f.raw.(float64)
	if ok {
		return s, nil
	}

	return 0, fmt.Errorf("field is not float64")
}

func (f *Field) Relation(include bool, stubs *StubConfigs) (interface{}, error) {
	if f.IsRelationship {
		fv := reflect.ValueOf(f.raw)

		dk := reflect.ValueOf("data")
		dv := fv.MapIndex(dk)

		isSlice := reflect.ValueOf(dv.Interface()).Kind() == reflect.Slice
		var payload interface{}
		{
			if isSlice {
				payload = new(jsonapi.ManyPayload)
			} else {
				payload = new(jsonapi.OnePayload)
			}

			buf := bytes.NewBuffer(nil)
			err := json.NewEncoder(buf).Encode(f.raw)
			if err != nil {
				return nil, fmt.Errorf("raw encode: %v", err)
			}
			err = json.NewDecoder(buf).Decode(payload)
			if err != nil {
				return nil, fmt.Errorf("raw decode: %v", err)
			}
		}

		if include {
			if isSlice {
				res := make([]interface{}, 0)
				p := payload.(*jsonapi.ManyPayload)
				if p.Data == nil || len(p.Data) == 0 {
					return nil, nil
				}

				for _, d := range p.Data {
					for _, n := range f.refPayload.Included {
						if n.ID == d.ID && n.Type == d.Type {
							nestedEntity := &Entity{
								payload: &jsonapi.OnePayload{Data: n, Included: f.refPayload.Included},
							}
							r, err := entityStubTransform(nestedEntity, stubs)
							if err != nil {
								return nil, fmt.Errorf("nested entity stub marhsal: %v", err)
							}

							res = append(res, r)
						}
					}
				}

				return res, nil
			} else {
				p := payload.(*jsonapi.OnePayload)
				if p.Data == nil {
					return nil, nil
				}

				for _, n := range f.refPayload.Included {
					if n.ID == p.Data.ID && n.Type == p.Data.Type {
						nestedEntity := &Entity{
							payload: &jsonapi.OnePayload{
								Data: n, Included: f.refPayload.Included,
							},
						}
						r, err := entityStubTransform(nestedEntity, stubs)
						if err != nil {
							return nil, fmt.Errorf("nested entity stub marhsal: %v", err)
						}

						return r, nil
					}
				}

				return nil, nil
			}
		}

		return payload, nil
	}

	return nil, errors.New("field does not belong to relationship")
}

// deprecated
func (f *Field) File() (*File, error) {
	var node *jsonapi.Node
	{
		payload := new(jsonapi.OnePayload)

		buf := bytes.NewBuffer(nil)
		err := json.NewEncoder(buf).Encode(f.raw)
		if err != nil {
			return nil, fmt.Errorf("raw encode: %v", err)
		}
		err = json.NewDecoder(buf).Decode(payload)
		if err != nil {
			// try to decode by ManyPayload
			mPayload := new(jsonapi.ManyPayload)

			buf = bytes.NewBuffer(nil)
			err = json.NewEncoder(buf).Encode(f.raw)
			if err != nil {
				return nil, fmt.Errorf("raw many encode: %v", err)
			}
			err = json.NewDecoder(buf).Decode(mPayload)
			if err != nil {
				return nil, fmt.Errorf("raw many decode: %v", err)
			}
			if mPayload.Data == nil || len(mPayload.Data) == 0 {
				return nil, nil
			}
			for _, n := range f.refPayload.Included {
				if n.ID == mPayload.Data[0].ID {
					node = n
					break
				}
			}
		} else {
			if payload.Data == nil {
				return nil, nil
			}
			for _, n := range f.refPayload.Included {
				if n.ID == payload.Data.ID {
					node = n
					break
				}
			}
		}
	}

	if node == nil || node.Attributes == nil {
		return nil, fmt.Errorf("not found in included")
	}

	file := new(File)

	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(node.Attributes)
	if err != nil {
		return nil, fmt.Errorf("attr encode: %v", err)
	}
	err = json.NewDecoder(buf).Decode(file)
	file.FID = node.ID
	if err != nil {
		return nil, fmt.Errorf("file decode: %v", err)
	}
	return file, nil
}

func (f *Field) Unmarshal(model interface{}) error {
	if f.IsRelationship {
		return fmt.Errorf("field is relatiionship")
	}

	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(f.raw)
	if err != nil {
		return err
	}

	err = json.NewDecoder(buf).Decode(model)
	if err != nil {
		return err
	}
	return nil
}
