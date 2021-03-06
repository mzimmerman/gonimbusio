package nimbusapi

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

func StartConjoined(requester Requester, collectionName string, key string) (
	string, error) {
	method := "POST"
	hostName := requester.CollectionHostName(collectionName)
	values := url.Values{}
	values.Set("action", "start")
	path := fmt.Sprintf("/conjoined/%s?%s", url.QueryEscape(key),
		values.Encode())

	request, err := requester.CreateRequest(method, hostName, path, nil)
	if err != nil {
		return "", err
	}

	response, err := requester.Do(request)
	if err != nil {
		return "", err
	}

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("POST %s %s failed (%d) %s", hostName, path,
			response.StatusCode, response.Body)
		return "", err
	}

	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	var conjoinedMap map[string]string
	err = json.Unmarshal(responseBody, &conjoinedMap)
	if err != nil {
		return "", err
	}

	return conjoinedMap["conjoined_identifier"], nil
}

func AbortConjoined(requester Requester, collectionName string, key string,
	conjoinedIdentifier string) error {
	method := "POST"
	hostName := requester.CollectionHostName(collectionName)
	values := url.Values{}
	values.Set("action", "abort")
	values.Set("conjoined_identifier", conjoinedIdentifier)
	path := fmt.Sprintf("/conjoined/%s?%s", url.QueryEscape(key),
		values.Encode())

	request, err := requester.CreateRequest(method, hostName, path, nil)
	if err != nil {
		return err
	}

	response, err := requester.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("POST %s %s failed (%d) %s", hostName, path,
			response.StatusCode, response.Body)
		return err
	}

	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	var conjoinedMap map[string]bool
	err = json.Unmarshal(responseBody, &conjoinedMap)
	if err != nil {
		return err
	}

	if !conjoinedMap["success"] {
		return fmt.Errorf("conjoined action=abort returned false")
	}

	return nil
}
func FinishConjoined(requester Requester, collectionName string, key string,
	conjoinedIdentifier string) error {
	method := "POST"
	hostName := requester.CollectionHostName(collectionName)
	values := url.Values{}
	values.Set("action", "finish")
	values.Set("conjoined_identifier", conjoinedIdentifier)
	path := fmt.Sprintf("/conjoined/%s?%s", url.QueryEscape(key),
		values.Encode())

	request, err := requester.CreateRequest(method, hostName, path, nil)
	if err != nil {
		return err
	}

	response, err := requester.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("POST %s %s failed (%d) %s", hostName, path,
			response.StatusCode, response.Body)
		return err
	}

	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	var conjoinedMap map[string]bool
	err = json.Unmarshal(responseBody, &conjoinedMap)
	if err != nil {
		return err
	}

	if !conjoinedMap["success"] {
		return fmt.Errorf("conjoined action=finish returned false")
	}

	return nil
}

type ConjoinedParams struct {
	ConjoinedIdentifier string
	ConjoinedPart       int
}

func Archive(requester Requester, collectionName string, key string,
	conjoinedParams *ConjoinedParams,
	requestBody io.Reader) (string, error) {
	method := "POST"
	hostName := requester.CollectionHostName(collectionName)

	var path string
	if conjoinedParams != nil {
		values := url.Values{}
		values.Set("conjoined_identifier", conjoinedParams.ConjoinedIdentifier)
		values.Set("conjoined_part",
			strconv.Itoa(conjoinedParams.ConjoinedPart))
		path = fmt.Sprintf("/data/%s?%s", url.QueryEscape(key), values.Encode())
	} else {
		path = fmt.Sprintf("/data/%s", url.QueryEscape(key))
	}

	request, err := requester.CreateRequest(method, hostName, path, requestBody)
	if err != nil {
		return "", err
	}
	// request.contentLength is set in http.NewRequest

	response, err := requester.Do(request)
	if err != nil {
		return "", err
	}

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("POST %s %s failed (%d) %s", hostName, path,
			response.StatusCode, response.Body)
		return "", err
	}

	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	var versionMap map[string]string
	err = json.Unmarshal(responseBody, &versionMap)
	if err != nil {
		return "", err
	}

	return versionMap["version_identifier"], nil
}

type RetrieveParams struct {
	VersionID      string
	SliceOffset    int
	SliceSize      int
	ModifiedSince  interface{}
	UnodifiedSince interface{}
}

func Retrieve(requester Requester, collectionName string, key string,
	retrieveParams RetrieveParams) (io.ReadCloser, error) {

	if retrieveParams.VersionID != "" {
		return nil, fmt.Errorf("not implemented: retrieveParams.VersionID")
	}
	if retrieveParams.ModifiedSince != nil {
		return nil, fmt.Errorf("not implemented: retrieveParams.ModifiedSince")
	}
	if retrieveParams.UnodifiedSince != nil {
		return nil, fmt.Errorf("not implemented: retrieveParams.UnodifiedSince")
	}

	method := "GET"
	hostName := requester.CollectionHostName(collectionName)
	path := fmt.Sprintf("/data/%s", url.QueryEscape(key))

	request, err := requester.CreateRequest(method, hostName, path, nil)
	if err != nil {
		return nil, err
	}

	successfulStatusCode := http.StatusOK

	if retrieveParams.SliceOffset > 0 || retrieveParams.SliceSize > 0 {
		var rangeArg string

		if retrieveParams.SliceSize > 0 {
			rangeArg = fmt.Sprintf("bytes=%d-%d", retrieveParams.SliceOffset,
				retrieveParams.SliceOffset+retrieveParams.SliceSize)
		} else {
			rangeArg = fmt.Sprintf("bytes=%d-", retrieveParams.SliceOffset)
		}

		request.Header.Add("Range", rangeArg)
		successfulStatusCode = http.StatusPartialContent
	}

	response, err := requester.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != successfulStatusCode {
		err = fmt.Errorf("GET %s %s failed (%d) %s", hostName, path,
			response.StatusCode, response.Body)
		return nil, err
	}

	return response.Body, nil
}

func DeleteKey(requester Requester, collectionName string, key string) error {
	method := "DELETE"
	hostName := requester.CollectionHostName(collectionName)
	path := fmt.Sprintf("/data/%s", url.QueryEscape(key))

	request, err := requester.CreateRequest(method, hostName, path, nil)
	if err != nil {
		return err
	}

	response, err := requester.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("POST %s %s failed (%d) %s", hostName, path,
			response.StatusCode, response.Body)
		return err
	}

	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var resultMap map[string]bool
	err = json.Unmarshal(responseBody, &resultMap)
	if err != nil {
		return err
	}
	if !resultMap["success"] {
		err = fmt.Errorf("unexpected 'false' for 'success'")
		return err
	}

	return nil
}
