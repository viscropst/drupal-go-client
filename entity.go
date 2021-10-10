package drupal_go_client

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/google/jsonapi"
	"reflect"
)

type Schema struct {
	fields []*FieldType
}

type EntityWritable interface {
	Save() error
}

type EntityCompatible interface {
	Type() string
	ID() string
	GetField(f string) (*Field, error)
	GetSchema() (*Schema, error)
}

type Entity struct {
	payload *jsonapi.OnePayload
}

func (e *Entity) Marshal(v *Stub) ([]byte, error) {
	panic("implement me")
}

func (e *Entity) Type() string {
	return e.payload.Data.Type
}

func (e *Entity) ID() string {
	return e.payload.Data.ID
}

func (e *Entity) GetField(f string) (*Field, error) {
	a, ok := e.payload.Data.Attributes[f]

	// attributes fields
	if ok {
		return &Field{
			raw:            a,
			name:           f,
			IsRelationship: false,
			refPayload:     e.payload,
		}, nil
	}

	// relationship fields
	r, ok := e.payload.Data.Relationships[f]
	if ok {
		return &Field{
			raw:            r,
			name:           f,
			IsRelationship: true,
			refPayload:     e.payload,
		}, nil
	}

	return nil, fmt.Errorf("field %s not existed", f)
}

func (e *Entity) GetSchema() (*Schema, error) {
	schema := new(Schema)

	schema.fields = make([]*FieldType, 0)
	for name, attr := range e.payload.Data.Attributes {
		if attr == nil {
			continue
		}

		f := &FieldType{
			t:    reflect.TypeOf(attr),
			name: name,
		}
		schema.fields = append(schema.fields, f)
	}

	return schema, nil
}

func (e *Entity) Payload() *jsonapi.OnePayload {
	return e.payload
}

type EntityManager struct {
	// Resty Client
	client *resty.Client
	stubs  *StubConfigs
}

func NewEM(client *resty.Client) *EntityManager {
	return &EntityManager{client: client}
}

func (e *EntityManager) Request(t, b string) EntityRequest {
	return &EntityJsonapiRequest{
		em:         e,
		entityType: t,
		bundle:     b,
	}
}

type EntityRequest interface {
	Create(b []byte) error
	Update(entity EntityCompatible) error
	Delete(entity EntityCompatible) error
	Load(id string, query JsonapiQuery) (EntityCompatible, error)
	LoadMultiple(query JsonapiQuery) ([]EntityCompatible, error)
}

type EntityJsonapiRequest struct {
	em         *EntityManager
	entityType string
	bundle     string
}

func (e *EntityJsonapiRequest) Update(entity EntityCompatible) error {
	panic("implement me")
}

func (e *EntityJsonapiRequest) Delete(entity EntityCompatible) error {
	panic("implement me")
}

func (e *EntityJsonapiRequest) Create(b []byte) error {
	payload, err := entityStubUnmarshal(b, *e.em.stubs)
	if err != nil {
		return fmt.Errorf("entity unmarshal with stub: %v", err)
	}

	resp, err := e.em.client.R().
		SetHeader("Content-Type", "application/vnd.api+json").
		SetHeader("Accept", "application/vnd.api+json").
		SetError(&jsonapi.ErrorsPayload{}).
		SetBody(payload).
		Post(fmt.Sprintf("/%s/%s", e.entityType, e.bundle))
	if err != nil {
		return fmt.Errorf("load %s", err)
	}

	jsonapiErr, ok := resp.Error().(*jsonapi.ErrorsPayload)
	if ok && len(jsonapiErr.Errors) > 0 {
		return jsonapiErr.Errors[0]
	}

	return nil
}

func (e *EntityJsonapiRequest) Load(id string, q JsonapiQuery) (EntityCompatible, error) {
	resp, err := e.em.client.R().
		SetQueryParams(q.QueryParams()).
		SetHeader("Accept", "application/vnd.api+json").
		SetError(&jsonapi.ErrorsPayload{}).
		SetResult(&jsonapi.OnePayload{}).
		Get(fmt.Sprintf("/%s/%s/%s", e.entityType, e.bundle, id))
	if err != nil {
		return nil, fmt.Errorf("load %s", err)
	}

	jsonapiErr, ok := resp.Error().(*jsonapi.ErrorsPayload)
	if ok && len(jsonapiErr.Errors) > 0 {
		return nil, jsonapiErr.Errors[0]
	}

	return &Entity{payload: resp.Result().(*jsonapi.OnePayload)}, nil
}

func (e *EntityJsonapiRequest) LoadMultiple(q JsonapiQuery) ([]EntityCompatible, error) {
	resp, err := e.em.client.R().
		SetQueryParams(q.QueryParams()).
		SetHeader("Accept", "application/vnd.api+json").
		SetError(&jsonapi.ErrorsPayload{}).
		SetResult(&jsonapi.ManyPayload{}).
		Get(fmt.Sprintf("/%s/%s", e.entityType, e.bundle))
	if err != nil {
		return nil, fmt.Errorf("loadMultiple %s", err)
	}

	jsonapiErr, ok := resp.Error().(*jsonapi.ErrorsPayload)
	if ok && len(jsonapiErr.Errors) > 0 {
		return nil, jsonapiErr.Errors[0]
	}

	// ManyPayload to OnePayload slice
	res := make([]EntityCompatible, 0)
	p := resp.Result().(*jsonapi.ManyPayload)
	for _, n := range p.Data {
		entity := &Entity{
			payload: &jsonapi.OnePayload{
				Data:     n,
				Included: p.Included,
				Links:    nil,
				Meta:     nil,
			},
		}
		res = append(res, entity)
	}

	return res, nil
}
