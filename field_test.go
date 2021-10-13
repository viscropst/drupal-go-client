package drupal_go_client

import (
	"encoding/json"
	"github.com/google/jsonapi"
	"github.com/wangxb07/drupal-go-client/fixture"
	"reflect"
	"testing"
)

func TestFieldToTypeValue(t *testing.T) {
	c := fixture.NodeBannerHttpMockWithIncluded()
	em := &EntityManager{
		client: c,
	}
	entity, err := em.Request("node", "banner").
		WithQuery(JQ().Include([]string{"field_banner_image"})).
		Load("6085d170-5ec1-4a22-b69e-ecdd41242eab")
	if err != nil {
		t.Fatal(err)
	}

	entityRelationshipIsSlice, err := em.Request("node", "banner").
		WithQuery(JQ().Include([]string{"field_banner_image"})).
		Load("6085d170-5ec1-4a22-b69e-ecdd41242eac")
	if err != nil {
		t.Fatal(err)
	}

	titleField, _ := entity.GetField("title")
	s, _ := titleField.String()
	if *s != "test" {
		t.Errorf("titleField.String not equal test, got %v", s)
	}

	_, err = titleField.Int32()
	if err == nil {
		t.Errorf("expect an error, return nil")
	}

	nidField, _ := entity.GetField("drupal_internal__nid")
	filesize, err := nidField.Float64()
	if *filesize != 9999 {
		t.Errorf("expect filesize 9999, got %v", filesize)
	}

	stickyField, _ := entity.GetField("sticky")
	sticky, _ := stickyField.Bool()
	if *sticky != false {
		t.Errorf("expect sticky false, got %v", sticky)
	}

	fileField, _ := entity.GetField("field_banner_image")
	file, err := fileField.File()
	if err != nil {
		t.Fatal(err)
	}

	fileWant := &File{
		FID:      "db3b76f9-5020-47fb-beb0-5c5966c9740c",
		LangCode: "en",
		Filename: "WechatIMG8660.jpeg",
		URI: struct {
			Value string `json:"value"`
			URL   string `json:"url"`
		}{
			Value: "public://2021-09/WechatIMG8660.jpeg",
			URL:   "/sites/default/files/2021-09/WechatIMG8660.jpeg",
		},
		FileMime: "image/jpeg",
		Filesize: 296160,
		Status:   true,
		Created:  "2021-09-29T17:49:45+00:00",
		Changed:  "2021-09-29T17:50:19+00:00",
	}
	if !reflect.DeepEqual(file, fileWant) {
		t.Errorf("File() got = %v, want %v", file, fileWant)
	}

	linkField, _ := entity.GetField("field_banner_link")
	l := new(Link)
	err = linkField.Unmarshal(l)
	if err != nil {
		t.Fatal(err)
	}

	if l.URI != "internal:/pages/topic/topic" {
		t.Errorf("Link field got = %v", l)
	}

	fileField1, _ := entityRelationshipIsSlice.GetField("field_banner_image")
	file1, err := fileField1.File()
	if err != nil {
		t.Fatal(err)
	}

	fileWant1 := &File{
		FID:      "db3b76f9-5020-47fb-beb0-5c5966c9740c",
		LangCode: "en",
		Filename: "WechatIMG8660.jpeg",
		URI: struct {
			Value string `json:"value"`
			URL   string `json:"url"`
		}{
			Value: "public://2021-09/WechatIMG8660.jpeg",
			URL:   "/sites/default/files/2021-09/WechatIMG8660.jpeg",
		},
		FileMime: "image/jpeg",
		Filesize: 296160,
		Status:   true,
		Created:  "2021-09-29T17:49:45+00:00",
		Changed:  "2021-09-29T17:50:19+00:00",
	}
	if !reflect.DeepEqual(file1, fileWant1) {
		t.Errorf("File() got = %v, want %v", file, fileWant1)
	}

	// raw is null handle
	fLog, _ := entity.GetField("revision_log")
	fLogStr, _ := fLog.String()
	if fLogStr != nil {
		t.Errorf("null field to string should return empty, but %v", fLogStr)
	}

	fLogInt64, _ := fLog.Int64()
	if fLogInt64 != nil {
		t.Errorf("null field to int64 should return 0, but %v", fLogInt64)
	}

	fLogInt32, _ := fLog.Int32()
	if fLogInt32 != nil {
		t.Errorf("null field to Int32 should return 0, but %v", fLogInt32)
	}

	fLogFloat64, _ := fLog.Float64()
	if fLogFloat64 != nil {
		t.Errorf("null field to Float64 should return 0, but %v", fLogFloat64)
	}

	fLogFloat32, _ := fLog.Float32()
	if fLogFloat32 != nil {
		t.Errorf("null field to Float32 should return 0, but %v", fLogFloat32)
	}

	fLogBool, _ := fLog.Bool()
	if fLogBool != nil {
		t.Errorf("null field to Bool should return 0, but %v", fLogBool)
	}
}

func TestField_Relation(t *testing.T) {
	oPayloadJSON := `{
	  "links": {
		"self": "http://example.com/articles/1/relationships/author",
		"related": "http://example.com/articles/1/author"
	  },
	  "data": {
		"type": "people",
		"id": "9"
	  }
	}`

	oPayloadRaw := new(interface{})
	json.Unmarshal([]byte(oPayloadJSON), oPayloadRaw)

	oPayload := new(jsonapi.OnePayload)
	json.Unmarshal([]byte(oPayloadJSON), oPayload)

	mPayloadJSON := `{
	  "links": {
		"self": "http://example.com/articles/1/relationships/comments",
		"related": "http://example.com/articles/1/comments"
	  },
	  "data": [
		{
		  "type": "comments",
		  "id": "5"
		},
		{
		  "type": "comments",
		  "id": "12"
		}
	  ]
	}`
	mPayloadRaw := new(interface{})
	json.Unmarshal([]byte(mPayloadJSON), mPayloadRaw)

	mPayload := new(jsonapi.ManyPayload)
	json.Unmarshal([]byte(mPayloadJSON), mPayload)

	simpleData := fixture.SimpleOnePayload()
	e := &Entity{payload: simpleData}

	stubs, err := NewStubConfigsFromJSON(fixture.SimpleTestSubConfigsJSON())
	if err != nil {
		t.Fatal(err)
	}

	fComments, _ := e.GetField("comments")
	fAuthor, _ := e.GetField("author")

	type fields struct {
		raw            interface{}
		name           string
		refPayload     *jsonapi.OnePayload
		IsRelationship bool
	}
	type args struct {
		include bool
		stubs   *StubConfigs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "one payload relationship",
			fields: fields{
				raw:            *oPayloadRaw,
				name:           "author",
				refPayload:     nil,
				IsRelationship: true,
			},
			args:    args{include: false, stubs: stubs},
			want:    oPayload,
			wantErr: false,
		}, {
			name: "many payload relationship",
			fields: fields{
				raw:            *mPayloadRaw,
				name:           "author",
				refPayload:     nil,
				IsRelationship: true,
			},
			args:    args{include: false, stubs: stubs},
			want:    mPayload,
			wantErr: false,
		}, {
			name: "author relationship",
			fields: fields{
				raw:            fAuthor.raw,
				name:           fAuthor.name,
				refPayload:     fAuthor.refPayload,
				IsRelationship: true,
			},
			args: args{include: true, stubs: stubs},
			want: map[string]interface{}{
				"id":        "9",
				"type":      "people",
				"firstName": "Dan",
				"lastName":  "Gebhardt",
				"twitter":   "dgeb",
			},
			wantErr: false,
		}, {
			name: "comments relationship",
			fields: fields{
				raw:            fComments.raw,
				name:           fComments.name,
				refPayload:     fComments.refPayload,
				IsRelationship: true,
			},
			args: args{include: true, stubs: stubs},
			want: []map[string]interface{}{
				{
					"author": nil,
					"id":     "5",
					"type":   "comments",
					"body":   "First!",
				},
				{
					"author": map[string]interface{}{
						"id":        "9",
						"type":      "people",
						"firstName": "Dan",
						"lastName":  "Gebhardt",
						"twitter":   "dgeb",
					},
					"id":   "12",
					"type": "comments",
					"body": "I like XML better",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Field{
				raw:            tt.fields.raw,
				name:           tt.fields.name,
				refPayload:     tt.fields.refPayload,
				IsRelationship: tt.fields.IsRelationship,
			}
			got, err := f.Relation(tt.args.include, tt.args.stubs)
			if (err != nil) != tt.wantErr {
				t.Errorf("Relation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.name == "comments relationship" {
				gotJson, err := json.Marshal(got)
				if err != nil {
					t.Fatal(err)
				}
				wantJson, _ := json.Marshal(tt.want)
				if string(gotJson) != string(wantJson) {
					t.Errorf("Relation() got = %s, want %s", gotJson, wantJson)
				}
			} else {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Relation() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
