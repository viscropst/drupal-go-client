package drupal_go_client

import (
	"fmt"
	"strings"
)

type JsonapiQuery interface {
	QueryParams() map[string]string
	Include([]string) JsonapiQuery
	Sort([]string) JsonapiQuery
	Page(offset, limit int) JsonapiQuery
}

type Query struct {
	params map[string]string
}

func JQ() *Query {
	return &Query{params: make(map[string]string)}
}

func (q *Query) QueryParams() map[string]string {
	return q.params
}

func (q *Query) SetQueryParams(p map[string]string) {
	q.params = p
}

func (q *Query) Include(s []string) JsonapiQuery {
	q.params["include"] = strings.Join(s, ",")
	return q
}

func (q *Query) Sort(s []string) JsonapiQuery {
	q.params["sort"] = strings.Join(s, ",")
	return q
}

func (q *Query) Page(offset, limit int) JsonapiQuery {
	q.params["page[offset]"] = fmt.Sprint(offset)
	q.params["page[limit]"] = fmt.Sprint(limit)
	return q
}
