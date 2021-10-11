package drupal_go_client

import (
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
	if s != "test" {
		t.Errorf("titleField.String not equal test, got %s", s)
	}

	_, err = titleField.Int32()
	if err == nil {
		t.Errorf("expect an error, return nil")
	}

	nidField, _ := entity.GetField("drupal_internal__nid")
	filesize, err := nidField.Float64()
	if filesize != 9999 {
		t.Errorf("expect filesize 9999, got %f", filesize)
	}

	stickyField, _ := entity.GetField("sticky")
	sticky, _ := stickyField.Bool()
	if sticky != false {
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

}
