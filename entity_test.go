package drupal_go_client

import (
	"github.com/google/jsonapi"
	"github.com/wangxb07/drupal-go-client/fixture"
	"reflect"
	"testing"
)

func TestLoad(t *testing.T) {
	c := fixture.NodePoHttpMockWithSingleData()
	em := &EntityManager{
		client: c,
	}
	got, err := em.Request("node", "po").Load("da58cbf5-83a4-4850-8a6f-8d7618483ff6", JQ())
	if err != nil {
		t.Fatal(err)
	}

	if got.Type() != "node--po" {
		t.Errorf("Entity Type() want bundle type node--po, got %s", got.Type())
	}

	f, _ := got.GetField("title")
	if f.Raw() != "月饼DIY制作活动" {
		t.Errorf("Entity Type() want title is \"月饼DIY制作活动\", got %s", f.raw)
	}

	c2 := fixture.NodePoHttpMockNotFound()
	em2 := &EntityManager{
		client: c2,
	}

	_, err = em2.Request("node", "po").Load("da58cbf5-83a4-4850-8a6f-8d7618483ff7", JQ())
	if jsonapiErr, ok := err.(*jsonapi.ErrorObject); !ok {
		t.Errorf("not found error object expected, but got: %v", err)
	} else {
		if jsonapiErr.Title != "Not Found" {
			t.Errorf("err title want Not found, got: %v", jsonapiErr.Title)
		}
	}
}

func TestLoadMultiple(t *testing.T) {
	c := fixture.NodeBannerHttpMockWithMultipleData()

	em := &EntityManager{
		client: c,
	}

	q := JQ().
		Include([]string{"field_banner_image"}).
		Page(0, 10).
		Sort([]string{"created"})

	entities, err := em.
		Request("node", "banner").
		LoadMultiple(q)
	if err != nil {
		t.Fatal(err)
	}

	if len(entities) != 1 {
		t.Errorf("expect entities length 1, got %d", len(entities))
	}

	titleField, _ := entities[0].GetField("title")
	s, _ := titleField.String()
	if s != "test" {
		t.Errorf("expect title is test, got %s", s)
	}

	c2 := fixture.NodeBannerHttpMockNotFound()
	em2 := &EntityManager{
		client: c2,
	}

	q2 := JQ().
		Include([]string{"field_banner_image1"}).
		Page(0, 10).
		Sort([]string{"created"})

	_, err = em2.Request("node", "banner").LoadMultiple(q2)
	if jsonapiErr, ok := err.(*jsonapi.ErrorObject); !ok {
		t.Errorf("not found error object expected, but got: %v", err)
	} else {
		if jsonapiErr.Title != "Bad Request" {
			t.Errorf("err title want Not found, got: %v", jsonapiErr.Title)
		}
	}
}

func TestEntity_GetSchema(t *testing.T) {
	c := fixture.SimpleJSONAPIHttpMockWithSingleData()

	em := &EntityManager{
		client: c,
	}
	got, err := em.Request("node", "po").Load("da58cbf5-83a4-4850-8a6f-8d7618483ff6", JQ())
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		payload *jsonapi.OnePayload
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Schema
		wantErr bool
	}{
		{
			name: "node banner schema test",
			fields: fields{
				payload: got.(*Entity).Payload(),
			},
			want: &Schema{
				fields: []*FieldType{
					{
						t:    reflect.TypeOf("s"),
						name: "title",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Entity{
				payload: tt.fields.payload,
			}
			got, err := e.GetSchema()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSchema() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got.fields[0].t.String() != "string" {
				t.Errorf("first field type want string, but %v", got.fields[0].t.String())
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSchema() got = %v, want %v", got, tt.want)
			}
		})
	}
}
