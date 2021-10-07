package drupal_go_client

import (
	"github.com/wangxb07/drupal-go-client/fixture"
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
}
