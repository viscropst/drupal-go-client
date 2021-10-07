package drupal_go_client

import (
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
			args: args{bytes: []byte(jsonCfg)},
			wantErr: false,
		},
		{
			name: "not support field type",
			want: nil,
			args: args{bytes: []byte(errJsonCfg)},
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
	entity, err := em.Request("node", "po").Load("da58cbf5-83a4-4850-8a6f-8d7618483ff6", JQ())
	if err != nil {
		t.Fatal(err)
	}

	stubs, err := NewStubConfigsFromJSON(fixture.NodeBannerTestSubConfigsJSON())
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
				stubs: *stubs,
			},
			want: []byte(`{"banner_image":{"fid":"db3b76f9-5020-47fb-beb0-5c5966c9740c","langcode":"en","filename":"WechatIMG8660.jpeg","uri":{"value":"public://2021-09/WechatIMG8660.jpeg","url":"/sites/default/files/2021-09/WechatIMG8660.jpeg"},"filemime":"image/jpeg","filesize":296160,"status":true,"created":"2021-09-29T17:49:45+00:00","changed":"2021-09-29T17:50:19+00:00"},"link":{"options":[],"title":"","uri":"internal:/pages/topic/topic"}}`),
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