package droxolite

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/louisevanderlith/droxolite/bodies"
)

//DoGET does a GET request and will update the container with the reponse's values.
//token: this is the access_token/avosession
//container: the object that will be populated with the results
//instanceID: instance of the application making the request
//serviceName: the name of the service being requested
//controller: the Controller to call
//params: additional path variables
//returns int : httpStatusCode
//return error: error
func DoGET(token string, container interface{}, instanceID, serviceName, controller string, params ...string) (int, error) {
	url, err := GetServiceURL(instanceID, serviceName, false)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	fullURL := fmt.Sprintf("%s%s/%s", url, controller, strings.Join(params, "/"))

	req, err := http.NewRequest("GET", fullURL, nil)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	if len(token) > 0 {
		req.Header.Add("Authorization", "Bearer "+token)
	}

	client := new(http.Client)
	resp, err := client.Do(req)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return resp.StatusCode, nil
	}

	rest, err := bodies.MarshalToResult(contents, container)

	if err != nil {
		msg := fmt.Errorf("Invalid JSON; Body:\n%s\nError:\n%s", string(contents), err)
		return http.StatusInternalServerError, msg
	}

	if len(rest.Reason) > 0 {
		return resp.StatusCode, rest
	}

	return resp.StatusCode, nil
}

//DoSEND is able to do a POST or PUT request and will update the container with the reponse's values.
//token: this is the access_token/avosession
//container: the object that will be populated with the results
//instanceID: instance of the application making the request
//serviceName: the name of the service being requested
//controller: the Controller to call
//data: the data to be sent with the request
//params: additional path variables
//returns int : httpStatusCode
//return error: error
func DoSEND(method, token string, container interface{}, instanceID, serviceName, controller string, data interface{}, params ...string) (int, error) {
	url, err := GetServiceURL(instanceID, serviceName, false)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	fullURL := fmt.Sprintf("%s%s/%s", url, controller, strings.Join(params, "/"))

	bits, err := json.Marshal(data)

	if err != nil {
		return http.StatusBadRequest, err
	}

	req, err := http.NewRequest(method, fullURL, bytes.NewBuffer(bits))

	if err != nil {
		return http.StatusInternalServerError, err
	}

	if len(token) > 0 {
		req.Header.Add("Authorization", "Bearer "+token)
	}

	client := new(http.Client)
	resp, err := client.Do(req)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return resp.StatusCode, nil
	}

	rest, err := bodies.MarshalToResult(contents, container)

	if err != nil {
		msg := fmt.Errorf("Invalid JSON; Body:\n%s\nError:\n%s", string(contents), err)
		return http.StatusInternalServerError, msg
	}

	if len(rest.Reason) > 0 {
		return resp.StatusCode, rest
	}

	return resp.StatusCode, nil
}
