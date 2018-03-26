package urlval

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"

	"gitlab.com/distributed_lab/figure"
)

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
		fieldValue := rval.Field(fi)
		fieldType := rtyp.Field(fi)

		isSet, err := figure.Out(&dest).With(figure.BaseHooks).From(parseFieldTag(values)).SetField(fieldValue, fieldType, `url`)
		if err != nil {
			return err
		}

		if isSet {
			tag := fieldType.Tag.Get("url")
			values.Del(tag)
		}
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
	for i := 0; i < rval.NumField(); i++ {
		tag := rtyp.Field(i).Tag.Get("url")
		if tag == "" {
			continue
		}
		if tag == "page" {
			encodePage(queries, tag, rval.Field(i))
		}
		switch rval.Field(i).Interface().(type) {
		case *uint64:
			encodeUint64Pointer(queries, tag, rval.Field(i))
		case *string:
			encodeStringPointer(queries, tag, rval.Field(i))
		}
	}
	links := FilterLinks{
		Self: fmt.Sprintf("%s?%s", r.URL.Path, queries[1].Encode()),
		Next: fmt.Sprintf("%s?%s", r.URL.Path, queries[2].Encode()),
	}
	if queries[0].Get("page") != "	" && queries[1].Get("page") != "" {
		links.Prev = fmt.Sprintf("%s?%s", r.URL.Path, queries[0].Encode())
	}
	return links
}

func encodeStringPointer(queries []url.Values, tag string, value reflect.Value) {
	if value.IsNil() {
		return
	}
	str := reflect.Indirect(value).String()
	for _, query := range queries {
		query.Add(tag, str)
	}
}

func encodeUint64Pointer(queries []url.Values, tag string, value reflect.Value) {
	if value.IsNil() {
		return
	}
	uint := reflect.Indirect(value).Uint()
	//uint := reflect.Indirect(value).String()

	for _, query := range queries {
		query.Add(tag, fmt.Sprintf("%d", uint))
	}
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
