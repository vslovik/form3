package form3

import (
	"net/url"
	"reflect"
)

// ListOptions specifies the optional parameters to various List methods that
// support offset pagination.
type ListOptions struct {
	// For paginated result sets, page of results to retrieve.
	Page int `url:"page[number],omitempty"`

	// For paginated result sets, the number of results to include per page.
	PerPage int `url:"page[size],omitempty"`
}

// addOptions adds the parameters in opts as URL query parameters to s. opts
// must be a struct whose fields may contain "url" tags.
func addOptions(s string, opts interface{}) (string, error) {
	v := reflect.ValueOf(opts)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := Values(opts)

	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}
