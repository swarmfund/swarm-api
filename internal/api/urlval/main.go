package urlval

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"

	"github.com/pkg/errors"
)

// Populate parses fields from url.Values according to `url` struct tag on dest
func Decode(values url.Values, dest interface{}) error {
	rval := reflect.Indirect(reflect.ValueOf(dest))
	rtyp := rval.Type()
	for i := 0; i < rval.NumField(); i++ {
		tag := rtyp.Field(i).Tag.Get("url")
		if tag == "" {
			continue
		}
		if values.Get(tag) == "" {
			continue
		}
		switch rval.Field(i).Interface().(type) {
		case uint64:
			uint, err := strconv.ParseUint(values.Get(tag), 0, 64)
			if err != nil {
				return errors.Wrapf(err, "failed to parse %s to uint64", tag)
			}
			rval.Field(i).Set(reflect.ValueOf(uint))
		case *uint64:
			uint, err := strconv.ParseUint(values.Get(tag), 0, 64)
			if err != nil {
				return errors.Wrapf(err, "failed to parse %s to uint64", tag)
			}
			rval.Field(i).Set(reflect.ValueOf(&uint))
		default:
			fmt.Println("unknown type")
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

func encodeUint64Pointer(queries []url.Values, tag string, value reflect.Value) {
	if value.IsNil() {
		return
	}
	uint := reflect.Indirect(value).Uint()
	// prev
	queries[0].Add(tag, fmt.Sprintf("%d", uint))
	// self
	queries[1].Add(tag, fmt.Sprintf("%d", uint))
	// next
	queries[2].Add(tag, fmt.Sprintf("%d", uint))
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
