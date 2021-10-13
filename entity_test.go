package drupal_go_client

import (
	"github.com/google/jsonapi"
	"github.com/wangxb07/drupal-go-client/fixture"
	"net/http"
	"reflect"
	"testing"
)

func TestLoad(t *testing.T) {
	c := fixture.NodePoHttpMockWithSingleData()
	em := &EntityManager{
		client: c,
	}
	got, err := em.Request("node", "po").Load("da58cbf5-83a4-4850-8a6f-8d7618483ff6")
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

	_, err = em2.Request("node", "po").Load("da58cbf5-83a4-4850-8a6f-8d7618483ff7")
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
		WithQuery(q).
		LoadMultiple()
	if err != nil {
		t.Fatal(err)
	}

	if len(entities) != 1 {
		t.Errorf("expect entities length 1, got %d", len(entities))
	}

	titleField, _ := entities[0].GetField("title")
	s, _ := titleField.String()
	if *s != "test" {
		t.Errorf("expect title is test, got %v", s)
	}

	c2 := fixture.NodeBannerHttpMockNotFound()
	em2 := &EntityManager{
		client: c2,
	}

	q2 := JQ().
		Include([]string{"field_banner_image1"}).
		Page(0, 10).
		Sort([]string{"created"})

	_, err = em2.Request("node", "banner").WithQuery(q2).LoadMultiple()
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
	got, err := em.Request("node", "po").Load("da58cbf5-83a4-4850-8a6f-8d7618483ff6")
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

func TestEntityJsonapiRequest_Create(t *testing.T) {
	c := fixture.CreateBannerJSONAPIHttpMock()

	stubs1, err := NewStubConfigsFromJSON(fixture.NodeBannerTestSubConfigsJSON())
	if err != nil {
		t.Fatal(err)
	}
	em := NewEM(c, stubs1)

	origReq, _ := http.NewRequest(http.MethodPost, "http://www.demo.com", nil)
	origReq.Header.Set("Authorization", "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImp0aSI6ImI1OWJmZGRlY2UyZTRmNDEyNGJjY2NmMmZmNDEyMzRlMTI1Y2EzNjgzMTAzZGUwZTY4YjYzZjViMzYyMzZhZWNmYmMxZWQ1ODcwOGJiYTg5In0.eyJhdWQiOiI0NmFmMzdiNy0xY2E0LTQ1MjYtYTc0My05NDA5NzFmNmMwNDMiLCJqdGkiOiJiNTliZmRkZWNlMmU0ZjQxMjRiY2NjZjJmZjQxMjM0ZTEyNWNhMzY4MzEwM2RlMGU2OGI2M2Y1YjM2MjM2YWVjZmJjMWVkNTg3MDhiYmE4OSIsImlhdCI6MTYzMzg5MTkxMSwibmJmIjoxNjMzODkxOTExLCJleHAiOjE2MzUxMDE1MTEsInN1YiI6IjQiLCJzY29wZXMiOlsiYXV0aGVudGljYXRlZCIsImJmZiJdLCJ1aWQiOiJjOGJmZjU0Ny0zNjI1LTRlZjgtOTI2OC0yMjc0YzBkMjIzOGIiLCJtYWlsIjoiYmZmQHFxLmNvbSIsInVzZXJuYW1lIjoiYmZmIiwib3BlbmlkIjoiIn0.QFsQTovgFvxmn_i00LzHjtEkGOwbg5LhPeFO-dLkngRbHgD9niKGl-4nq50Ecbx7fSZrRz5yzgMmuFc-bqGvjPG0vqlr62Qz0FZzGx_lC8UBmXsYebUx7EgtmtyiE-AIJMV69XQWfD7-BzU2D6ZjaXP8XXz2jw8U9VvPqcqlmuOgWPamO2qHIjMQxyvJ2sj-WPVFCONuqPvd59NxhGDxPewGrsbK09hJJsphWl78RX5NWWPdWelw0f8j7Mf3QhZauMA0m53oUKcOMeFHPI9Y9P89hj06UhvboPTaS6ZnOUntlDCkOYiW5do6QuZPXav6mJGSt-n-3ylIhglUleVnydu3tVx6AQvXblhDcYO1BLYm_Mfjv7kDpZRGonL2TfwGJe6V9XnHV7KINNGvialzIIdVFP7I5r3BQdKivtTbe88xbxhMEb-0KmWLe_zEECu8wKAaVrUK8siGw_cUiDoA3BTfGHLuxTTuqJHH4oEwd1DYDlO9J1ikA6v4Of8_VXOpCsMa0m1MaRciTshT4wclkeVxprRpXH1EPzaGT1GC8ku49EBz01TExUI7_YFpClYIK-kjtNZBQ5cgtJ5XbxVuy6D2DU8bSA6SybZkpWB2HYQHylKaqNlgSrA4VPrnuvCcX_t1dYXSJXe1GzIwH0R48CfUGv1CoXOXHn0blfXDmME")

	type fields struct {
		em         *EntityManager
		entityType string
		bundle     string
		Req        *http.Request
	}
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "create success",
			fields: fields{
				em:         em,
				entityType: "node",
				bundle:     "banner",
				Req:        origReq,
			},
			args: args{
				b: []byte(`{"type":"node--banner", "title": "banner1"}`),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EntityJsonapiRequest{
				em:         tt.fields.em,
				entityType: tt.fields.entityType,
				bundle:     tt.fields.bundle,
				Req:        tt.fields.Req,
			}
			if _, err := e.Create(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEntityJsonapiRequest_Delete(t *testing.T) {
	c := fixture.DeleteBannerJSONAPIHttpMock()

	stubs1, err := NewStubConfigsFromJSON(fixture.NodeBannerTestSubConfigsJSON())
	if err != nil {
		t.Fatal(err)
	}
	em := NewEM(c, stubs1)

	origReq, _ := http.NewRequest(http.MethodPost, "http://www.demo.com", nil)
	origReq.Header.Set("Authorization", "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImp0aSI6IjEzYzA4ZWYwZGIzNWI0NjNmYTgxMjdkYWJkN2IxZDQ1MzRkOWIxODg5NmIxMjk3Y2FmNDJjYjc2NDgyNTliZWZlZWI4ZGIzNTBlNzBiNjRlIn0.eyJhdWQiOiI0NmFmMzdiNy0xY2E0LTQ1MjYtYTc0My05NDA5NzFmNmMwNDMiLCJqdGkiOiIxM2MwOGVmMGRiMzViNDYzZmE4MTI3ZGFiZDdiMWQ0NTM0ZDliMTg4OTZiMTI5N2NhZjQyY2I3NjQ4MjU5YmVmZWViOGRiMzUwZTcwYjY0ZSIsImlhdCI6MTYzMzkyOTQwNiwibmJmIjoxNjMzOTI5NDA2LCJleHAiOjE2MzUxMzkwMDYsInN1YiI6IjQiLCJzY29wZXMiOlsiYXV0aGVudGljYXRlZCIsImJmZiJdLCJ1aWQiOiJjOGJmZjU0Ny0zNjI1LTRlZjgtOTI2OC0yMjc0YzBkMjIzOGIiLCJtYWlsIjoiYmZmQHFxLmNvbSIsInVzZXJuYW1lIjoiYmZmIiwib3BlbmlkIjoiIn0.vHLk_xzF_Wpkat4Pm5XU_gDJbcRqGNDHs6347KakO7Px9kMv48-dFFFgHpx82ioztp1fHQ226B07pvQCjGr1otWuvDJV8yl4vGUrtLTCttfzZ-ZNtH8RqNf5AHGTvs-veyKssWQ4BEd0-gf5sctj3HHfncBx5V5yM7YHO1-WLKwNgGHZRbq5JVAOu6oemb_40dOhcsKMJ4ogmbhQ-rGtBEc_C6zVLv0M7Hlau_pb1KA4Hc113YKEAffLqU2XG-H2D9OSkrWrBYJeZajuDrQB3oTdIGkj6uSlg5pJYCCefTPO6OQbeu1Bgz49J_A34Adt5U4YxCPIhdvy0pbRHIplZ2jkxFYDV815hmZsh6Zf0-aaAIhzPGSTmH4C3kKjEaDp9ze7iX9x02l5rGo1mI6PtGJu57SHpOmD5jCc1XgUqQ8vOxGm_72N5cnvoBzNtgWHgGi5xpWlwrDmaiN6QzKKon467SzWSpL-7jKQE1DZ6zyTpz8x2NaQzAT6Uu-P8Gts6QmURSPrjJOYbIQsSv_Or1kYD0KTObKO-kPDVZZrtlA4vmQ5To5oT2tpsZV6AepFIF-l4WJg5trfQrngzxx1N9gPqg2UAmJxUDswyNxVMVci00eJ-yTQgr_R-mNVwDLJn9fHds97KYkwu-b4qBFjLpXrilxm5rBtUv1fXH_ZMvA")

	type fields struct {
		em         *EntityManager
		entityType string
		bundle     string
		Req        *http.Request
		Query      JsonapiQuery
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "delete",
			fields: fields{
				em:         em,
				entityType: "node",
				bundle:     "banner",
			},
			args: args{
				id: "d44cb65f-00a1-4eb2-a038-5960833654f1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EntityJsonapiRequest{
				em:         tt.fields.em,
				entityType: tt.fields.entityType,
				bundle:     tt.fields.bundle,
				Req:        tt.fields.Req,
				Query:      tt.fields.Query,
			}
			if err = e.WithRequest(origReq).Delete(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEntityJsonapiRequest_Update(t *testing.T) {
	c := fixture.UpdateBannerJSONAPIHttpMock()

	stubs1, err := NewStubConfigsFromJSON(fixture.NodeBannerTestSubConfigsJSON())
	if err != nil {
		t.Fatal(err)
	}
	em := NewEM(c, stubs1)

	origReq, _ := http.NewRequest(http.MethodPost, "http://www.demo.com", nil)
	origReq.Header.Set("Authorization", "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImp0aSI6IjBhZDA5OTlkZmQ0ZjRjZDVjNjdjY2NmODA1MzgzYzBmOTg2MWI1OTM4NTMzODQ5Nzg4YTJjOTNhMjRmNzkxMjIzODZlZjEzOTcwYTM3ZWRhIn0.eyJhdWQiOiI0NmFmMzdiNy0xY2E0LTQ1MjYtYTc0My05NDA5NzFmNmMwNDMiLCJqdGkiOiIwYWQwOTk5ZGZkNGY0Y2Q1YzY3Y2NjZjgwNTM4M2MwZjk4NjFiNTkzODUzMzg0OTc4OGEyYzkzYTI0Zjc5MTIyMzg2ZWYxMzk3MGEzN2VkYSIsImlhdCI6MTYzMzkzMDYxMywibmJmIjoxNjMzOTMwNjEzLCJleHAiOjE2MzUxNDAyMTMsInN1YiI6IjQiLCJzY29wZXMiOlsiYXV0aGVudGljYXRlZCIsImJmZiJdLCJ1aWQiOiJjOGJmZjU0Ny0zNjI1LTRlZjgtOTI2OC0yMjc0YzBkMjIzOGIiLCJtYWlsIjoiYmZmQHFxLmNvbSIsInVzZXJuYW1lIjoiYmZmIiwib3BlbmlkIjoiIn0.dI7o5TvInq49F0rTh0WTTlx3Dnn5wZE4bMiar7OO2is3QGlH4Jux3KSTv5af98XrmsVZNeMFz1SKTNOuIEuTNdb7mxCTB1AciQRYzCI7Si5zPiAqvOM0ekSAV2c9XjdfM2iT6fA6-dPhWfG7KPZ5St1ZAzyzua9uUjrQxQoWmP64YEJ6l41tLW7gcABeLxtEP29xeZvMoouWV1K-0j17jzOlAot23_vteu7j6dcpKjTa3XzeFqXzAQlLBNbLTfs7S6HGHZzwGg9H9YIF8dYdsJ8xD9Ps_eXKKOyWmMIH78CWqsRgCjykePQwVI88Hu6bwOrtgdzwpq-nzzFae21oqOyhWbSUSXLf_cXwu-W8ZcJjg9089fIdFYLv3BexXqd4ghG7UjHIxY5vRXBeGcXRcjXr6E70bllhrrUcUbq0YgPxooSdvjrE8EDkCCZXOejSTUb_lX1lvO7w5HlPW69WREHOZVUTbhS8i9Q05YfLe2E7pViQ9NXmW1inecjvJSZdcEohbYw_Apq2Z1dI7tlI7qUN3EjnEYLPeqcugwi_5W4cLLY4m8357VQHxKyyLeJqePd8PqrGnFfkUQrmJesgUTNQaOc7DTLszHdnyLnX20OyLLlOTV1wcrb9LDa8JjXnRMnTSN7xsQ8--f2P5rTbvryjuauND8uWiuEu0RzYKnE")

	type fields struct {
		em         *EntityManager
		entityType string
		bundle     string
		Req        *http.Request
		Query      JsonapiQuery
	}
	type args struct {
		id string
		b  []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "update",
			fields: fields{
				em:         em,
				entityType: "node",
				bundle:     "banner",
			},
			args: args{
				id: "34fe2569-18f0-40d9-a727-a274e300d7d6",
				b:  []byte(`{"id": "34fe2569-18f0-40d9-a727-a274e300d7d6", "type": "node--banner", "title": "banner2"}`),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EntityJsonapiRequest{
				em:         tt.fields.em,
				entityType: tt.fields.entityType,
				bundle:     tt.fields.bundle,
				Req:        tt.fields.Req,
				Query:      tt.fields.Query,
			}
			if _, err := e.WithRequest(origReq).Update(tt.args.id, tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
