package drupal_go_client

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/google/jsonapi"
	"net/http"
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
	Marshal(stubs *StubConfigs) ([]byte, error)
}

type Entity struct {
	payload *jsonapi.OnePayload
}

func (e *Entity) Marshal(stubs *StubConfigs) ([]byte, error) {
	return entityStubMarshal(e, stubs)
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

func NewEM(client *resty.Client, stubs *StubConfigs) *EntityManager {
	return &EntityManager{client: client, stubs: stubs}
}

func (e *EntityManager) GetClient() *resty.Client {
	return e.client
}

func (e *EntityManager) GetStubs() *StubConfigs {
	return e.stubs
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
	Update(id string, b []byte) error
	Delete(id string) error
	Load(id string) (EntityCompatible, error)
	LoadMultiple() ([]EntityCompatible, error)

	WithRequest(req *http.Request) *EntityJsonapiRequest
	WithQuery(query JsonapiQuery) *EntityJsonapiRequest
}

type EntityJsonapiRequest struct {
	em         *EntityManager
	entityType string
	bundle     string
	Req        *http.Request
	Query      JsonapiQuery
}

func (e *EntityJsonapiRequest) WithRequest(req *http.Request) *EntityJsonapiRequest {
	e.Req = req
	return e
}

func (e *EntityJsonapiRequest) WithQuery(query JsonapiQuery) *EntityJsonapiRequest {
	e.Query = query
	return e
}

func (e *EntityJsonapiRequest) Update(id string, b []byte) error {
	payload, err := entityStubUnmarshal(b, e.em.stubs)
	if err != nil {
		return fmt.Errorf("entity unmarshal with stub: %v", err)
	}

	r := e.em.client.R().
		SetHeader("Content-Type", "application/vnd.api+json").
		SetHeader("Accept", "application/vnd.api+json").
		SetError(&jsonapi.ErrorsPayload{}).
		SetBody(payload)

	if e.Req != nil && e.Req.Header.Get("Authorization") != "" {
		r = r.SetHeader("Authorization", e.Req.Header.Get("Authorization"))
	}

	resp, err := r.Patch(fmt.Sprintf("/%s/%s/%s", e.entityType, e.bundle, id))
	if err != nil {
		return fmt.Errorf("update %s", err)
	}

	jsonapiErr, ok := resp.Error().(*jsonapi.ErrorsPayload)
	if ok && len(jsonapiErr.Errors) > 0 {
		return jsonapiErr.Errors[0]
	}

	return nil
}

func (e *EntityJsonapiRequest) Delete(id string) error {
	r := e.em.client.R().
		SetHeader("Accept", "application/vnd.api+json").
		SetError(&jsonapi.ErrorsPayload{})
	if e.Req != nil && e.Req.Header.Get("Authorization") != "" {
		r = r.SetHeader("Authorization", e.Req.Header.Get("Authorization"))
	}
	resp, err := r.Delete(fmt.Sprintf("/%s/%s/%s", e.entityType, e.bundle, id))
	if err != nil {
		return fmt.Errorf("delete %s", err)
	}

	jsonapiErr, ok := resp.Error().(*jsonapi.ErrorsPayload)
	if ok && len(jsonapiErr.Errors) > 0 {
		return jsonapiErr.Errors[0]
	}

	return nil
}

func (e *EntityJsonapiRequest) Create(b []byte) error {
	payload, err := entityStubUnmarshal(b, e.em.stubs)
	if err != nil {
		return fmt.Errorf("entity unmarshal with stub: %v", err)
	}

	r := e.em.client.R().
		SetHeader("Content-Type", "application/vnd.api+json").
		SetHeader("Accept", "application/vnd.api+json").
		SetError(&jsonapi.ErrorsPayload{}).
		SetBody(payload)

	if e.Req != nil && e.Req.Header.Get("Authorization") != "" {
		r = r.SetHeader("Authorization", e.Req.Header.Get("Authorization"))
	}

	resp, err := r.Post(fmt.Sprintf("/%s/%s", e.entityType, e.bundle))
	if err != nil {
		return fmt.Errorf("create %s", err)
	}

	jsonapiErr, ok := resp.Error().(*jsonapi.ErrorsPayload)
	if ok && len(jsonapiErr.Errors) > 0 {
		return jsonapiErr.Errors[0]
	}

	return nil
}

func (e *EntityJsonapiRequest) Load(id string) (EntityCompatible, error) {
	r := e.em.client.R().
		SetHeader("Accept", "application/vnd.api+json").
		SetError(&jsonapi.ErrorsPayload{}).
		SetResult(&jsonapi.OnePayload{})
	if e.Query != nil {
		r.SetQueryParams(e.Query.QueryParams())
	}
	if e.Req != nil && e.Req.Header.Get("Authorization") != "" {
		r = r.SetHeader("Authorization", e.Req.Header.Get("Authorization"))
	}
	resp, err := r.Get(fmt.Sprintf("/%s/%s/%s", e.entityType, e.bundle, id))
	if err != nil {
		return nil, fmt.Errorf("load %s", err)
	}

	jsonapiErr, ok := resp.Error().(*jsonapi.ErrorsPayload)
	if ok && len(jsonapiErr.Errors) > 0 {
		return nil, jsonapiErr.Errors[0]
	}

	return &Entity{payload: resp.Result().(*jsonapi.OnePayload)}, nil
}

func (e *EntityJsonapiRequest) LoadMultiple() ([]EntityCompatible, error) {
	r := e.em.client.R().
		SetQueryParams(e.Query.QueryParams()).
		SetHeader("Accept", "application/vnd.api+json").
		SetError(&jsonapi.ErrorsPayload{}).
		SetResult(&jsonapi.ManyPayload{})
	if e.Query != nil {
		r.SetQueryParams(e.Query.QueryParams())
	}
	if e.Req != nil && e.Req.Header.Get("Authorization") != "" {
		r = r.SetHeader("Authorization", e.Req.Header.Get("Authorization"))
	}
	resp, err := r.Get(fmt.Sprintf("/%s/%s", e.entityType, e.bundle))
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
