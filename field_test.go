package drupal_go_client

import (
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"reflect"
	"testing"
)

func TestFieldToTypeValue(t *testing.T) {
	c := resty.New()
	httpmock.ActivateNonDefault(c.GetClient())

	fixture := `{
  "jsonapi": {
    "version": "1.0"
  },
  "data": {
    "type": "node--banner",
    "id": "6085d170-5ec1-4a22-b69e-ecdd41242eab",
    "links": {
    },
    "attributes": {
      "drupal_internal__nid": 9999,
      "drupal_internal__vid": 10263,
      "langcode": "en",
      "revision_timestamp": "2021-09-29T17:50:53+00:00",
      "revision_log": null,
      "status": true,
      "title": "test",
      "created": "2021-09-29T17:49:27+00:00",
      "changed": "2021-09-29T17:50:53+00:00",
      "promote": true,
      "sticky": false,
      "default_langcode": true,
      "revision_translation_affected": true,
      "body": {
        "value": "<p>test</p>\r\n",
        "format": "fuwenben",
        "processed": "<p>test</p>\n",
        "summary": ""
      },
      "field_banner_link": {
        "uri": "internal:/pages/topic/topic",
        "title": "",
        "options": []
      }
    },
    "relationships": {
      "uid": {
        "data": {
          "type": "user--user",
          "id": "c862c0f4-9a5b-42ff-be6f-e5d323e90ed9"
        },
        "links": {
        }
      },
      "field_banner_image": {
        "data": {
          "type": "file--file",
          "id": "db3b76f9-5020-47fb-beb0-5c5966c9740c",
          "meta": {
            "alt": "test banner",
            "title": "",
            "width": 1920,
            "height": 960
          }
        },
        "links": {
        }
      }
    }
  },
  "included": [
    {
      "type": "file--file",
      "id": "db3b76f9-5020-47fb-beb0-5c5966c9740c",
      "links": {
      },
      "attributes": {
        "drupal_internal__fid": 16950,
        "langcode": "en",
        "filename": "WechatIMG8660.jpeg",
        "uri": {
          "value": "public://2021-09/WechatIMG8660.jpeg",
          "url": "/sites/default/files/2021-09/WechatIMG8660.jpeg"
        },
        "filemime": "image/jpeg",
        "filesize": 296160,
        "status": true,
        "created": "2021-09-29T17:49:45+00:00",
        "changed": "2021-09-29T17:50:19+00:00"
      },
      "relationships": {
        "uid": {
          "data": {
            "type": "user--user",
            "id": "c862c0f4-9a5b-42ff-be6f-e5d323e90ed9"
          },
          "links": {
          }
        }
      }
    }
  ],
  "links": {
  }
}`
	responder := httpmock.NewStringResponder(200, fixture)
	fakeUrl := "https://milliface-base.beehomeplus.cn/jsonapi/node/banner/6085d170-5ec1-4a22-b69e-ecdd41242eab?include=field_banner_image"
	httpmock.RegisterResponder("GET", fakeUrl, responder)

	c.SetHostURL("https://milliface-base.beehomeplus.cn/jsonapi")
	em := &EntityManager{
		client: c,
	}
	entity, err := em.GetRequest("node", "banner", map[string]string{
		"include": "field_banner_image",
	}).Load("6085d170-5ec1-4a22-b69e-ecdd41242eab")
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
}
