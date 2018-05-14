package urlval

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"

	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const urlvalKey = "url"

func parseFieldTag(values url.Values) map[string]interface{} {
	figValues := make(map[string]interface{})
	for k, v := range values {
		figValues[k] = v[0]
	}
	return figValues
}

// Populate parses fields from url.Values according to `url` struct tag on dest
// When values is set to dest, we delete it from values
// After execution function you need to check residual from values
func Decode(values url.Values, dest interface{}) error {
	rval := reflect.Indirect(reflect.ValueOf(dest))
	rtyp := rval.Type()
	for fi := 0; fi < rval.NumField(); fi++ {
		value := rval.Field(fi)
		field := rtyp.Field(fi)

		err := figure.Out(&dest).With(figure.BaseHooks).From(parseFieldTag(values)).SetField(value, field, urlvalKey)
		if err != nil {
			return errors.Wrap(err, "failed to decode field", logan.F{"field": field.Name})
		}

		values.Del(field.Tag.Get(urlvalKey))
	}
	return nil
}

type FilterLinks struct {
	Self string `json:"self"`
	Next string `json:"next"`
	Prev string `json:"prev,omitempty"`
}

func Encode(r *http.Request, filters interface{}) FilterLinks {
	rval := reflect.ValueOf(filters)
	rtyp := rval.Type()
	queries := []url.Values{
		{}, {}, {},
	}
	for fi := 0; fi < rval.NumField(); fi++ {
		encodeField(queries, rval.Field(fi), rtyp.Field(fi))
	}
	return createLinks(r, queries)

}

func encodeField(queries []url.Values, fieldValue reflect.Value, fieldType reflect.StructField) {
	tag := fieldType.Tag.Get(urlvalKey)
	if tag == "" {
		return
	}

	if tag == "page" {
		encodePage(queries, tag, fieldValue)
		return
	}

	if fieldValue.IsNil() {
		return
	}

	uint := reflect.Indirect(fieldValue)

	stringer, ok := uint.Interface().(fmt.Stringer)
	if ok {
		for _, query := range queries {
			query.Add(tag, stringer.String())
		}
	} else {
		for _, query := range queries {
			query.Add(tag, fmt.Sprintf("%v", uint))
		}
	}
}

func createLinks(r *http.Request, queries []url.Values) FilterLinks {
	links := FilterLinks{
		Self: fmt.Sprintf("%s?%s", r.URL.Path, queries[1].Encode()),
		Next: fmt.Sprintf("%s?%s", r.URL.Path, queries[2].Encode()),
	}

	if queries[0].Get("page") != "	" && queries[1].Get("page") != "" {
		links.Prev = fmt.Sprintf("%s?%s", r.URL.Path, queries[0].Encode())
	}

	return links
}

func encodePage(queries []url.Values, tag string, value reflect.Value) {
	uint := value.Uint()
	// prev
	if value.Uint() > 1 {
		queries[0].Add(tag, fmt.Sprintf("%d", uint-1))
	}
	// self
	queries[1].Add(tag, fmt.Sprintf("%d", uint))
	// next
	queries[2].Add(tag, fmt.Sprintf("%d", uint+1))
}
