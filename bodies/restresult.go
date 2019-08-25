package bodies

import (
	"encoding/json"
	"net/http"
)

//RESTResult is the base object of every response.
type RESTResult struct {
	Reason string      `json:"Error"`
	Data   interface{} `json:"Data"`
}

var client = &http.Client{}

//NewRESTResult is used to wrap responses in a consistent manner. data will be tested for the 'error' interface{}
func NewRESTResult(data interface{}) *RESTResult {
	result := &RESTResult{}

	if data == nil {
		return result
	}

	if err, isErr := data.(error); isErr {
		result.Reason = err.Error()
		return result
	}

	result.Data = data

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
