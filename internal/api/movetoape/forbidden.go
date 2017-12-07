package movetoape

import (
	"fmt"
	"net/http"

	"github.com/google/jsonapi"
)

func Forbidden(code string) *jsonapi.ErrorObject {
	return &jsonapi.ErrorObject{
		Title:  http.StatusText(http.StatusForbidden),
		Status: fmt.Sprintf("%d", http.StatusForbidden),
		Code:   code,
	}
}
