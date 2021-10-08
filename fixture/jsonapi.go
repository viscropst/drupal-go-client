package fixture

import (
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
)

func NodePoHttpMockWithSingleData() *resty.Client {
	c := resty.New()
	httpmock.ActivateNonDefault(c.GetClient())
	fixture := `{"jsonapi":{"version":"1.0","meta":{"links":{"self":{"href":"http:\/\/jsonapi.org\/format\/1.0\/"}}}},"data":{"type":"node--po","id":"da58cbf5-83a4-4850-8a6f-8d7618483ff6","links":{"self":{"href":"http:\/\/milliface-base.beehomeplus.cn\/jsonapi\/node\/po\/da58cbf5-83a4-4850-8a6f-8d7618483ff6?resourceVersion=id%3A9934"}},"attributes":{"drupal_internal__nid":9547,"drupal_internal__vid":9934,"langcode":"en","revision_timestamp":"2021-09-14T15:58:08+00:00","revision_log":null,"status":true,"title":"\u6708\u997cDIY\u5236\u4f5c\u6d3b\u52a8","created":"2021-09-10T09:44:06+00:00","changed":"2021-09-14T15:58:08+00:00","promote":true,"sticky":false,"default_langcode":true,"revision_translation_affected":true,"field_address_name":"\u4e1c\u5149\u6751\u6587\u5316\u793c\u5802","field_green_check_state":"pass","field_map":{"value":"POINT (121.54257745 29.764702899)","geo_type":"Point","lat":29.764702899,"lon":121.54257745,"left":121.54257745,"top":29.764702899,"right":121.54257745,"bottom":29.764702899,"geohash":"wtq3mf2xgqf2y34t","latlon":"29.764702899,121.54257745"},"field_text":{"value":"\u003Cp\u003E2021\u5e749\u670821\u65e5\u662f\u6211\u4eec\u7684\u4f20\u7edf\u8282\u65e5\u4e2d\u79cb\u8282\uff0c\u4e3a\u5f18\u626c\u6c11\n"}},"relationships":{"node_type":{"data":{"type":"node_type--node_type","id":"8a85371e-61ad-498b-9172-55ff156028d5"},"links":{"related":{"href":"http:\/\/milliface-base.beehomeplus.cn\/jsonapi\/node\/po\/da58cbf5-83a4-4850-8a6f-8d7618483ff6\/node_type?resourceVersion=id%3A9934"},"self":{"href":"http:\/\/milliface-base.beehomeplus.cn\/jsonapi\/node\/po\/da58cbf5-83a4-4850-8a6f-8d7618483ff6\/relationships\/node_type?resourceVersion=id%3A9934"}}},"revision_uid":{"data":{"type":"user--user","id":"c862c0f4-9a5b-42ff-be6f-e5d323e90ed9"},"links":{"related":{"href":"http:\/\/milliface-base.beehomeplus.cn\/jsonapi\/node\/po\/da58cbf5-83a4-4850-8a6f-8d7618483ff6\/revision_uid?resourceVersion=id%3A9934"},"self":{"href":"http:\/\/milliface-base.beehomeplus.cn\/jsonapi\/node\/po\/da58cbf5-83a4-4850-8a6f-8d7618483ff6\/relationships\/revision_uid?resourceVersion=id%3A9934"}}},"uid":{"data":{"type":"user--user","id":"2e291629-e88e-4db8-8b88-4c786c6ec948"},"links":{"related":{"href":"http:\/\/milliface-base.beehomeplus.cn\/jsonapi\/node\/po\/da58cbf5-83a4-4850-8a6f-8d7618483ff6\/uid?resourceVersion=id%3A9934"},"self":{"href":"http:\/\/milliface-base.beehomeplus.cn\/jsonapi\/node\/po\/da58cbf5-83a4-4850-8a6f-8d7618483ff6\/relationships\/uid?resourceVersion=id%3A9934"}}},"field_category":{"data":[],"links":{"related":{"href":"http:\/\/milliface-base.beehomeplus.cn\/jsonapi\/node\/po\/da58cbf5-83a4-4850-8a6f-8d7618483ff6\/field_category?resourceVersion=id%3A9934"},"self":{"href":"http:\/\/milliface-base.beehomeplus.cn\/jsonapi\/node\/po\/da58cbf5-83a4-4850-8a6f-8d7618483ff6\/relationships\/field_category?resourceVersion=id%3A9934"}}},"field_column":{"data":[],"links":{"related":{"href":"http:\/\/milliface-base.beehomeplus.cn\/jsonapi\/node\/po\/da58cbf5-83a4-4850-8a6f-8d7618483ff6\/field_column?resourceVersion=id%3A9934"},"self":{"href":"http:\/\/milliface-base.beehomeplus.cn\/jsonapi\/node\/po\/da58cbf5-83a4-4850-8a6f-8d7618483ff6\/relationships\/field_column?resourceVersion=id%3A9934"}}},"field_cover":{"data":{"type":"file--file","id":"fc72f787-106b-43a8-a0ed-d0e695f3c030","meta":{"alt":null,"title":null,"width":3543,"height":6496}},"links":{"related":{"href":"http:\/\/milliface-base.beehomeplus.cn\/jsonapi\/node\/po\/da58cbf5-83a4-4850-8a6f-8d7618483ff6\/field_cover?resourceVersion=id%3A9934"},"self":{"href":"http:\/\/milliface-base.beehomeplus.cn\/jsonapi\/node\/po\/da58cbf5-83a4-4850-8a6f-8d7618483ff6\/relationships\/field_cover?resourceVersion=id%3A9934"}}},"field_pictures":{"data":[{"type":"file--file","id":"fc72f787-106b-43a8-a0ed-d0e695f3c030","meta":{"alt":null,"title":null,"width":3543,"height":6496}}],"links":{"related":{"href":"http:\/\/milliface-base.beehomeplus.cn\/jsonapi\/node\/po\/da58cbf5-83a4-4850-8a6f-8d7618483ff6\/field_pictures?resourceVersion=id%3A9934"},"self":{"href":"http:\/\/milliface-base.beehomeplus.cn\/jsonapi\/node\/po\/da58cbf5-83a4-4850-8a6f-8d7618483ff6\/relationships\/field_pictures?resourceVersion=id%3A9934"}}},"field_po_tags":{"data":[],"links":{"related":{"href":"http:\/\/milliface-base.beehomeplus.cn\/jsonapi\/node\/po\/da58cbf5-83a4-4850-8a6f-8d7618483ff6\/field_po_tags?resourceVersion=id%3A9934"},"self":{"href":"http:\/\/milliface-base.beehomeplus.cn\/jsonapi\/node\/po\/da58cbf5-83a4-4850-8a6f-8d7618483ff6\/relationships\/field_po_tags?resourceVersion=id%3A9934"}}},"field_video":{"data":[],"links":{"related":{"href":"http:\/\/milliface-base.beehomeplus.cn\/jsonapi\/node\/po\/da58cbf5-83a4-4850-8a6f-8d7618483ff6\/field_video?resourceVersion=id%3A9934"},"self":{"href":"http:\/\/milliface-base.beehomeplus.cn\/jsonapi\/node\/po\/da58cbf5-83a4-4850-8a6f-8d7618483ff6\/relationships\/field_video?resourceVersion=id%3A9934"}}},"field_video_cover":{"data":null,"links":{"related":{"href":"http:\/\/milliface-base.beehomeplus.cn\/jsonapi\/node\/po\/da58cbf5-83a4-4850-8a6f-8d7618483ff6\/field_video_cover?resourceVersion=id%3A9934"},"self":{"href":"http:\/\/milliface-base.beehomeplus.cn\/jsonapi\/node\/po\/da58cbf5-83a4-4850-8a6f-8d7618483ff6\/relationships\/field_video_cover?resourceVersion=id%3A9934"}}}}},"links":{"self":{"href":"http:\/\/milliface-base.beehomeplus.cn\/jsonapi\/node\/po\/da58cbf5-83a4-4850-8a6f-8d7618483ff6"}}}`
	responder := httpmock.NewStringResponder(200, fixture)
	fakeUrl := "https://milliface-base.beehomeplus.cn/jsonapi/node/po/da58cbf5-83a4-4850-8a6f-8d7618483ff6"
	httpmock.RegisterResponder("GET", fakeUrl, responder)
	c.SetHostURL("https://milliface-base.beehomeplus.cn/jsonapi")
	return c
}

func NodeBannerHttpMockWithSingleData() *resty.Client {
	c := resty.New()
	httpmock.ActivateNonDefault(c.GetClient())
	fixture := `{"jsonapi":{"version":"1.0","meta":{"links":{"self":{"href":"http:\/\/jsonapi.org\/format\/1.0\/"}}}},"data":{"type":"node--banner","id":"6085d170-5ec1-4a22-b69e-ecdd41242eab","links":{"self":{"href":"http:\/\/milliface-base.beehomeplus.cn\/jsonapi\/node\/banner\/6085d170-5ec1-4a22-b69e-ecdd41242eab?resourceVersion=id%3A10263"}},"attributes":{"drupal_internal__nid":9999,"drupal_internal__vid":10263,"langcode":"en","revision_timestamp":"2021-09-29T17:50:53+00:00","revision_log":null,"status":true,"title":"test","created":"2021-09-29T17:49:27+00:00","changed":"2021-09-29T17:50:53+00:00","promote":true,"sticky":false,"default_langcode":true,"revision_translation_affected":true,"body":{"value":"\u003Cp\u003Etest\u003C\/p\u003E\r\n","format":"fuwenben","processed":"\u003Cp\u003Etest\u003C\/p\u003E\n","summary":""},"field_banner_link":{"uri":"internal:\/pages\/topic\/topic","title":"","options":[]}},"relationships":{"node_type":{"data":{"type":"node_type--node_type","id":"39f44cf3-59d5-4f3f-a953-aa9738ea6f61"},"links":{"related":{"href":"http:\/\/milliface-base.beehomeplus.cn\/jsonapi\/node\/banner\/6085d170-5ec1-4a22-b69e-ecdd41242eab\/node_type?resourceVersion=id%3A10263"},"self":{"href":"http:\/\/milliface-base.beehomeplus.cn\/jsonapi\/node\/banner\/6085d170-5ec1-4a22-b69e-ecdd41242eab\/relationships\/node_type?resourceVersion=id%3A10263"}}},"revision_uid":{"data":{"type":"user--user","id":"c862c0f4-9a5b-42ff-be6f-e5d323e90ed9"},"links":{"related":{"href":"http:\/\/milliface-base.beehomeplus.cn\/jsonapi\/node\/banner\/6085d170-5ec1-4a22-b69e-ecdd41242eab\/revision_uid?resourceVersion=id%3A10263"},"self":{"href":"http:\/\/milliface-base.beehomeplus.cn\/jsonapi\/node\/banner\/6085d170-5ec1-4a22-b69e-ecdd41242eab\/relationships\/revision_uid?resourceVersion=id%3A10263"}}},"uid":{"data":{"type":"user--user","id":"c862c0f4-9a5b-42ff-be6f-e5d323e90ed9"},"links":{"related":{"href":"http:\/\/milliface-base.beehomeplus.cn\/jsonapi\/node\/banner\/6085d170-5ec1-4a22-b69e-ecdd41242eab\/uid?resourceVersion=id%3A10263"},"self":{"href":"http:\/\/milliface-base.beehomeplus.cn\/jsonapi\/node\/banner\/6085d170-5ec1-4a22-b69e-ecdd41242eab\/relationships\/uid?resourceVersion=id%3A10263"}}},"field_banner_image":{"data":{"type":"file--file","id":"db3b76f9-5020-47fb-beb0-5c5966c9740c","meta":{"alt":"test banner","title":"","width":1920,"height":960}},"links":{"related":{"href":"http:\/\/milliface-base.beehomeplus.cn\/jsonapi\/node\/banner\/6085d170-5ec1-4a22-b69e-ecdd41242eab\/field_banner_image?resourceVersion=id%3A10263"},"self":{"href":"http:\/\/milliface-base.beehomeplus.cn\/jsonapi\/node\/banner\/6085d170-5ec1-4a22-b69e-ecdd41242eab\/relationships\/field_banner_image?resourceVersion=id%3A10263"}}}}},"included":[{"type":"file--file","id":"db3b76f9-5020-47fb-beb0-5c5966c9740c","links":{"self":{"href":"http:\/\/milliface-base.beehomeplus.cn\/jsonapi\/file\/file\/db3b76f9-5020-47fb-beb0-5c5966c9740c"}},"attributes":{"drupal_internal__fid":16950,"langcode":"en","filename":"WechatIMG8660.jpeg","uri":{"value":"public:\/\/2021-09\/WechatIMG8660.jpeg","url":"\/sites\/default\/files\/2021-09\/WechatIMG8660.jpeg"},"filemime":"image\/jpeg","filesize":296160,"status":true,"created":"2021-09-29T17:49:45+00:00","changed":"2021-09-29T17:50:19+00:00"},"relationships":{"uid":{"data":{"type":"user--user","id":"c862c0f4-9a5b-42ff-be6f-e5d323e90ed9"},"links":{"related":{"href":"http:\/\/milliface-base.beehomeplus.cn\/jsonapi\/file\/file\/db3b76f9-5020-47fb-beb0-5c5966c9740c\/uid"},"self":{"href":"http:\/\/milliface-base.beehomeplus.cn\/jsonapi\/file\/file\/db3b76f9-5020-47fb-beb0-5c5966c9740c\/relationships\/uid"}}}}}],"links":{"self":{"href":"http:\/\/milliface-base.beehomeplus.cn\/jsonapi\/node\/banner\/6085d170-5ec1-4a22-b69e-ecdd41242eab?include=field_banner_image"}}}`
	responder := httpmock.NewStringResponder(200, fixture)
	fakeUrl := "https://milliface-base.beehomeplus.cn/jsonapi/node/po/da58cbf5-83a4-4850-8a6f-8d7618483ff6"
	httpmock.RegisterResponder("GET", fakeUrl, responder)
	c.SetHostURL("https://milliface-base.beehomeplus.cn/jsonapi")
	return c
}


func SimpleJSONAPIHttpMockWithSingleData() *resty.Client {
	c := resty.New()
	httpmock.ActivateNonDefault(c.GetClient())
	fixture := `{
  "links": {
    "self": "http://example.com/articles",
    "next": "http://example.com/articles?page[offset]=2",
    "last": "http://example.com/articles?page[offset]=10"
  },
  "data": {
    "type": "articles",
    "id": "1",
    "attributes": {
      "title": "JSON:API paints my bikeshed!"
    },
    "relationships": {
      "author": {
        "links": {
          "self": "http://example.com/articles/1/relationships/author",
          "related": "http://example.com/articles/1/author"
        },
        "data": { "type": "people", "id": "9" }
      },
      "comments": {
        "links": {
          "self": "http://example.com/articles/1/relationships/comments",
          "related": "http://example.com/articles/1/comments"
        },
        "data": [
          { "type": "comments", "id": "5" },
          { "type": "comments", "id": "12" }
        ]
      }
    },
    "links": {
      "self": "http://example.com/articles/1"
    }
  },
  "included": [{
    "type": "people",
    "id": "9",
    "attributes": {
      "firstName": "Dan",
      "lastName": "Gebhardt",
      "twitter": "dgeb"
    },
    "links": {
      "self": "http://example.com/people/9"
    }
  }, {
    "type": "comments",
    "id": "5",
    "attributes": {
      "body": "First!"
    },
    "relationships": {
      "author": {
        "data": { "type": "people", "id": "2" }
      }
    },
    "links": {
      "self": "http://example.com/comments/5"
    }
  }, {
    "type": "comments",
    "id": "12",
    "attributes": {
      "body": "I like XML better"
    },
    "relationships": {
      "author": {
        "data": { "type": "people", "id": "9" }
      }
    },
    "links": {
      "self": "http://example.com/comments/12"
    }
  }]
}`
	responder := httpmock.NewStringResponder(200, fixture)
	fakeUrl := "https://milliface-base.beehomeplus.cn/jsonapi/node/po/da58cbf5-83a4-4850-8a6f-8d7618483ff6"
	httpmock.RegisterResponder("GET", fakeUrl, responder)
	c.SetHostURL("https://milliface-base.beehomeplus.cn/jsonapi")
	return c
}

func NodeBannerHttpMockWithMultipleData() *resty.Client {
	c := resty.New()
	httpmock.ActivateNonDefault(c.GetClient())
	fixture := `{
 "jsonapi": {
   "version": "1.0"
 },
 "data": [
   {
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
   }
 ],
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
 "meta": {
   "count": "1"
 },
 "links": {
 }
}`
	responder := httpmock.NewStringResponder(200, fixture)
	fakeUrl := "https://milliface-base.beehomeplus.cn/jsonapi/node/banner?include=field_banner_image&page%5Blimit%5D=10&page%5Boffset%5D=0&sort=created"
	httpmock.RegisterResponder("GET", fakeUrl, responder)
	c.SetHostURL("https://milliface-base.beehomeplus.cn/jsonapi")
	return c
}