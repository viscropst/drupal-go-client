package drupal_go_client

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/google/jsonapi"
)

type EntityWritable interface {
	Save() error
}

type EntityCompatible interface {
	Type() string
	ID() string
	GetField(f string) (*Field, error)
}

type Entity struct {
	payload *jsonapi.OnePayload
}

func (e *Entity) Type() string {
	return e.payload.Data.Type
}

func (e *Entity) ID() string {
	return e.payload.Data.ID
}

func (e *Entity) GetField(f string) (*Field, error) {
	a, ok := e.payload.Data.Attributes[f]

	if ok {
		return &Field{
			raw:            a,
			name:           f,
			IsRelationship: false,
			refPayload:     e.payload,
		}, nil
	}

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

func (e *Entity) Payload() *jsonapi.OnePayload {
	return e.payload
}

type EntityManager struct {
	// Resty Client
	client *resty.Client
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
	Create(entity EntityCompatible) error
	Load(id string, query JsonapiQuery) (EntityCompatible, error)
	LoadMultiple(query JsonapiQuery) ([]EntityCompatible, error)
	GetQuery() JsonapiQuery
}

type EntityJsonapiRequest struct {
	em         *EntityManager
	entityType string
	bundle     string
}

func (e *EntityJsonapiRequest) Create(en EntityCompatible) error {
	panic("not implemented")
}

func (e *EntityJsonapiRequest) Load(id string, q JsonapiQuery) (EntityCompatible, error) {
	resp, err := e.em.client.R().
		SetQueryParams(q.QueryParams()).
		SetHeader("Accept", "application/json").
		Get(fmt.Sprintf("/%s/%s/%s", e.entityType, e.bundle, id))
	if err != nil {
		return nil, fmt.Errorf("load %s", err)
	}

	p := jsonapi.OnePayload{}
	if err := json.Unmarshal(resp.Body(), &p); err != nil {
		return nil, fmt.Errorf("unmarshal to one payload: %v", err)
	}

	return &Entity{payload: &p}, nil
}

func (e *EntityJsonapiRequest) LoadMultiple(q JsonapiQuery) ([]EntityCompatible, error) {
	resp, err := e.em.client.R().
		SetQueryParams(q.QueryParams()).
		SetHeader("Accept", "application/json").
		Get(fmt.Sprintf("/%s/%s", e.entityType, e.bundle))
	if err != nil {
		return nil, fmt.Errorf("loadMultiple %s", err)
	}

	p := jsonapi.ManyPayload{}
	if err := json.Unmarshal(resp.Body(), &p); err != nil {
		return nil, fmt.Errorf("unmarshal to one payload: %v", err)
	}

	// ManyPayload to OnePayload slice
	res := make([]EntityCompatible, 0)
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

func (e EntityJsonapiRequest) GetQuery() JsonapiQuery {
	panic("not implemented")
}
