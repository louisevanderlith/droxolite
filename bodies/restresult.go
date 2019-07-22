package bodies

import (
	"encoding/json"
	"net/http"
)

//RESTResult is the base object of every response.
type RESTResult struct {
	Code   int         `json:"Code"`
	Reason string      `json:"Error"`
	Data   interface{} `json:"Data"`
}

var client = &http.Client{}

//NewRESTResult is used to wrap responses in a consistent manner
func NewRESTResult(code int, reason error, data interface{}) *RESTResult {
	result := &RESTResult{
		Code: code,
		Data: data,
	}

	if reason != nil {
		result.Reason = reason.Error()
	}

	return result
}

func MarshalToResult(content []byte, dataObj interface{}) (*RESTResult, error) {
	result := &RESTResult{Data: dataObj}
	err := json.Unmarshal(content, result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r RESTResult) Error() string {
	return r.Reason
}
