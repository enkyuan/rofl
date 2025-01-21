package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func MakeAPICall(url string, method string, headers map[string]interface{}, body interface{}) (*http.Response, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		logrus.Error("Error creating marshaling: ", err)
		return nil, err
	}
	_ = bodyBytes

	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		logrus.Error("Error creating http.NewRequest: ", err)
		return nil, err
	}

	// NOW WE CAN SET HEADERS
	// request.Header.Add("Content-Type", "application/json")
	// SETHeaders(request, headers)

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		logrus.Error("Error occur on HTTClient.Do()", err)
		return nil, err
	}

	logrus.Infof("Response Status: %v", response.Status)

	return response, nil
}

func SETHeaders(request *http.Request, headers map[string]interface{}) *http.Request {
	for key, value := range headers {
		request.Header.Add(key, fmt.Sprintf("%s", value))
	}
	return request
}

// HTTPResponseHandler
func HTTPResponseHandler(response *http.Response) map[string]interface{} {

	var msErrorStatuses = []int{
		http.StatusNotFound,
		http.StatusServiceUnavailable,
		http.StatusBadGateway,
		http.StatusInternalServerError,
		http.StatusBadRequest,
	}

	var r = make(map[string]interface{})

	if contains(msErrorStatuses, response.StatusCode) {
		r["microservice_status"] = response.StatusCode
		r["microservice_error"] = response.Status
		r["url"] = response.Request.URL.String()
		r["message"] = "NB: This is a Microservice Services Related Error!"

		body, _ := io.ReadAll(response.Body)
		defer response.Body.Close()

		var erroResp interface{}
		_ = json.Unmarshal(body, &erroResp)

		logrus.Errorf("[core.HTTPClient.HTTPResponseHandler] Decoded response error: %v", string(body))
		r["microservice_response"] = erroResp
		return r
	}
	// READ RESPONSE BODY.
	body, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		logrus.Errorf("[HTTPResponseHandler] Body ReadAll Error, %v", err.Error())
	}
	if err = json.Unmarshal(body, &r); err != nil {
		logrus.Errorf("Body parse error, %v", err)
	}

	return r
}

func HTTPErrorHandler(ctx echo.Context, err error, statusCode int) error {
	return ctx.JSON(statusCode, map[string]interface{}{
		"microservice_status": statusCode,
		"microservice_error":  http.StatusText(statusCode),
		"message":             err.Error(),
	})
}

func contains(statuses []int, status int) bool {

	for _, s := range statuses {
		if status == s {
			return true
		}
	}
	return false
}
