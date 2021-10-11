package drupal_go_client

import (
	"github.com/google/jsonapi"
	"github.com/wangxb07/drupal-go-client/fixture"
	"reflect"
	"testing"
)

func TestNewStubConfigsFromJSON(t *testing.T) {
	jsonCfg := `
{
	"node--article": {
		"entity_type": "node",
		"bundle": "article",
		"mapping": {
			"field_image": {
				"type": "file",
				"name": "image"
			},
			"field_category": {
				"type": "string",
				"name": "category"
			}
		}
	}
}
`

	errJsonCfg := `
{
	"node--article": {
		"entity_type": "node",
		"bundle": "article",
		"mapping": {
			"field_image": {
				"type": "image",
				"name": "image"
			},
			"field_category": {
				"type": "string",
				"name": "category"
			}
		}
	}
}
`
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *StubConfigs
		wantErr bool
	}{
		{
			name: "base unmarshal",
			want: &StubConfigs{
				"node--article": Stub{
					EntityType: "node",
					Bundle:     "article",
					Mapping: map[string]StubFieldMapper{
						"field_image": {
							Type: "file",
							Name: "image",
						},
						"field_category": {
							Type: "string",
							Name: "category",
						},
					},
				},
			},
			args:    args{bytes: []byte(jsonCfg)},
			wantErr: false,
		},
		{
			name:    "not support field type",
			want:    nil,
			args:    args{bytes: []byte(errJsonCfg)},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewStubConfigsFromJSON(tt.args.bytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewStubConfigsFromJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStubConfigsFromJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_entityStubMarshal(t *testing.T) {
	em := &EntityManager{
		client: fixture.NodeBannerHttpMockWithSingleData(),
	}
	entity, err := em.Request("node", "po").
		Load("da58cbf5-83a4-4850-8a6f-8d7618483ff6")
	if err != nil {
		t.Fatal(err)
	}

	stubs1, err := NewStubConfigsFromJSON(fixture.NodeBannerTestSubConfigsJSON())
	if err != nil {
		t.Fatal(err)
	}

	stubs2, err := NewStubConfigsFromJSON(fixture.NodeBannerTestNoMappingIgnoreSubConfigsJSON())
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		entity EntityCompatible
		stubs  StubConfigs
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "node article marshal",
			args: args{
				entity: entity,
				stubs:  *stubs1,
			},
			want:    []byte(`{"banner_image":{"fid":"db3b76f9-5020-47fb-beb0-5c5966c9740c","langcode":"en","filename":"WechatIMG8660.jpeg","uri":{"value":"public://2021-09/WechatIMG8660.jpeg","url":"/sites/default/files/2021-09/WechatIMG8660.jpeg"},"filemime":"image/jpeg","filesize":296160,"status":true,"created":"2021-09-29T17:49:45+00:00","changed":"2021-09-29T17:50:19+00:00"},"body":{"format":"fuwenben","processed":"\u003cp\u003etest\u003c/p\u003e\n","summary":"","value":"\u003cp\u003etest\u003c/p\u003e\r\n"},"changed":"2021-09-29T17:50:53+00:00","created":"2021-09-29T17:49:27+00:00","default_langcode":true,"drupal_internal__nid":9999,"drupal_internal__vid":10263,"id":"6085d170-5ec1-4a22-b69e-ecdd41242eab","langcode":"en","link":{"options":[],"title":"","uri":"internal:/pages/topic/topic"},"promote":true,"revision_timestamp":"2021-09-29T17:50:53+00:00","revision_translation_affected":true,"status":true,"sticky":false,"title":"test","type":"node--banner"}`),
			wantErr: false,
		}, {
			name: "no mapping mode ignore",
			args: args{
				entity: entity,
				stubs:  *stubs2,
			},
			want:    []byte(`{"banner_image":{"fid":"db3b76f9-5020-47fb-beb0-5c5966c9740c","langcode":"en","filename":"WechatIMG8660.jpeg","uri":{"value":"public://2021-09/WechatIMG8660.jpeg","url":"/sites/default/files/2021-09/WechatIMG8660.jpeg"},"filemime":"image/jpeg","filesize":296160,"status":true,"created":"2021-09-29T17:49:45+00:00","changed":"2021-09-29T17:50:19+00:00"},"id":"6085d170-5ec1-4a22-b69e-ecdd41242eab","link":{"options":[],"title":"","uri":"internal:/pages/topic/topic"},"type":"node--banner"}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := entityStubMarshal(tt.args.entity, tt.args.stubs)
			if (err != nil) != tt.wantErr {
				t.Errorf("entityStubMarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("entityStubMarshal() got = %s, want %s", got, tt.want)
			}
		})
	}
}

func Test_shallowNodeFromMap(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *jsonapi.Node
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				v: map[string]interface{}{
					"id":   "1",
					"type": "node--article",
				},
			},
			want: &jsonapi.Node{
				Type: "node--article",
				ID:   "1",
			},
			wantErr: false,
		}, {
			name: "key not right",
			args: args{
				v: map[string]interface{}{
					"uuid": "1",
					"type": "node--article",
				},
			},
			want:    nil,
			wantErr: true,
		}, {
			name: "extra key",
			args: args{
				v: map[string]interface{}{
					"id":    "1",
					"type":  "node--article",
					"extra": false,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := shallowNodeFromMap(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("shallowNodeFromMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("shallowNodeFromMap() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_entityStubUnmarshal(t *testing.T) {
	stubJsonConfigs := `
{
	"node--article": {
		"entity_type": "node",
		"bundle": "article",
		"mapping": {
			"foo": { "name": "bar", "type": "string" }
		}
	}
}
`
	stubs1, err := NewStubConfigsFromJSON([]byte(stubJsonConfigs))
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		b     []byte
		stubs StubConfigs
	}
	tests := []struct {
		name    string
		args    args
		want    *jsonapi.OnePayload
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				b:     []byte(`{"id": "1", "type": "node--article", "title": "JSON:API paints my bikeshed!","bar": "hello world", "author": {"type": "people", "id": "9"}}`),
				stubs: *stubs1,
			},
			want: &jsonapi.OnePayload{
				Data: &jsonapi.Node{
					Type: "node--article",
					ID:   "1",
					Attributes: map[string]interface{}{
						"title": "JSON:API paints my bikeshed!",
						"foo":   "hello world",
					},
					Relationships: map[string]interface{}{
						"author": jsonapi.OnePayload{
							Data: &jsonapi.Node{
								Type: "people",
								ID:   "9",
							},
						},
					},
				},
			},
			wantErr: false,
		}, {
			name: "many payload relationship",
			args: args{
				b:     []byte(`{"id": "1", "type": "node--article", "title": "JSON:API paints my bikeshed!","bar": "hello world", "author": {"type": "people", "id": "9"}, "tags": [{"type": "tag", "id": "1"},{"type": "tag", "id": "2"}]}`),
				stubs: *stubs1,
			},
			want: &jsonapi.OnePayload{
				Data: &jsonapi.Node{
					Type: "node--article",
					ID:   "1",
					Attributes: map[string]interface{}{
						"title": "JSON:API paints my bikeshed!",
						"foo":   "hello world",
					},
					Relationships: map[string]interface{}{
						"author": jsonapi.OnePayload{
							Data: &jsonapi.Node{
								Type: "people",
								ID:   "9",
							},
						},
						"tags": jsonapi.ManyPayload{
							Data: []*jsonapi.Node{
								{
									Type: "tag",
									ID:   "1",
								}, {
									Type: "tag",
									ID:   "2",
								},
							},
						},
					},
				},
			},
			wantErr: false,
		}, {
			name: "lose id",
			args: args{
				b:     []byte(`{"type": "node--article", "title": "JSON:API paints my bikeshed!","bar": "hello world", "author": {"type": "people", "id": "9"}}`),
				stubs: *stubs1,
			},
			want: &jsonapi.OnePayload{
				Data: &jsonapi.Node{
					Type: "node--article",
					Attributes: map[string]interface{}{
						"title": "JSON:API paints my bikeshed!",
						"foo":   "hello world",
					},
					Relationships: map[string]interface{}{
						"author": jsonapi.OnePayload{
							Data: &jsonapi.Node{
								Type: "people",
								ID:   "9",
							},
						},
					},
				},
			},
			wantErr: false,
		},{
			name: "lose type",
			args: args{
				b:     []byte(`{"title": "JSON:API paints my bikeshed!","bar": "hello world", "author": {"type": "people", "id": "9"}}`),
				stubs: *stubs1,
			},
			want: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := entityStubUnmarshal(tt.args.b, tt.args.stubs)
			if (err != nil) != tt.wantErr {
				t.Errorf("entityStubUnmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("entityStubUnmarshal() got = %v, want %v", got, tt.want)
			}
		})
	}
}
