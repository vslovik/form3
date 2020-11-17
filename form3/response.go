package form3

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Response is a Form3 API response. This wraps the standard http.Response.
type Response struct {
	*http.Response
}

/*
An ErrorResponse reports an error caused by an API request.
*/
type ErrorResponse struct {
	ErrorMessage string `json:"error_message"` // error message
}

// newResponse creates a new Response for the provided http.Response.
// r must not be nil.
func newResponse(r *http.Response) *Response {
	return &Response{Response: r}
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v", r.ErrorMessage)
}

// CheckResponse checks the API response for errors, and returns them if
// present. A response is considered an error if it has a status code outside
// the 200 range
// API error responses are expected to have response
// body, and a JSON response body that maps to ErrorResponse.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{}
	data, err := ioutil.ReadAll(r.Body)

	if err == nil && data != nil {
		err = json.Unmarshal(data, errorResponse)
		if err != nil && r.StatusCode == 404 {
			return &ErrorResponse{ErrorMessage: "Account not found"}
		}
	}
	return errorResponse
}
