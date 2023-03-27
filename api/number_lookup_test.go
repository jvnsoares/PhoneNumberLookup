package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"

	"phone-number-lookup/pkg/number"
)

func getHTTPTestRequestResponse(method, url string, parameters map[string]string,
) (*http.Request, *httptest.ResponseRecorder) {

	request := httptest.NewRequest(method, url, nil)

	query := request.URL.Query()
	for key, value := range parameters {
		query.Add(key, value)
	}
	request.URL.RawQuery = query.Encode()

	response := httptest.NewRecorder()

	return request, response
}

func Test_numberLookupHandler(t *testing.T) {
	url := "localhost:8080" + numberLookupEndpoint
	tests := []struct {
		method             string
		expectedStatusCode int
		parameters         map[string]string
		expectedOutput     numberLookupResponse
	}{
		{method: http.MethodPost, expectedStatusCode: http.StatusMethodNotAllowed},
		{
			method: http.MethodGet, expectedStatusCode: http.StatusOK,
			parameters: map[string]string{phoneNumberField: "+12125690123"},
			expectedOutput: numberLookupResponse{
				Info: number.Info{PhoneNumber: "+12125690123", RegionCode: "US", AreaCode: "212", LocalPhoneNumber: "5690123"},
			},
		},
		{
			method: http.MethodGet, expectedStatusCode: http.StatusOK,
			parameters: map[string]string{phoneNumberField: "2125690123", countryCodeField: "US"},
			expectedOutput: numberLookupResponse{
				Info: number.Info{PhoneNumber: "2125690123", RegionCode: "US", AreaCode: "212", LocalPhoneNumber: "5690123"},
			},
		},
		{
			method: http.MethodGet, expectedStatusCode: http.StatusOK,
			parameters: map[string]string{phoneNumberField: "2125690123", countryCodeField: "US"},
			expectedOutput: numberLookupResponse{
				Info: number.Info{PhoneNumber: "2125690123", RegionCode: "US", AreaCode: "212", LocalPhoneNumber: "5690123"},
			},
		},
		{
			method: http.MethodGet, expectedStatusCode: http.StatusBadRequest,
			parameters: map[string]string{phoneNumberField: "", countryCodeField: "US"},
			expectedOutput: numberLookupResponse{
				Err: &number.Info{PhoneNumber: number.MissingValueError.Error()},
			},
		},
		{
			method: http.MethodGet, expectedStatusCode: http.StatusBadRequest,
			parameters: map[string]string{phoneNumberField: "2125690123", countryCodeField: ""},
			expectedOutput: numberLookupResponse{
				Info: number.Info{PhoneNumber: "2125690123"},
				Err:  &number.Info{RegionCode: number.MissingValueError.Error()},
			},
		},
		{
			method: http.MethodGet, expectedStatusCode: http.StatusBadRequest,
			parameters: map[string]string{phoneNumberField: "+1212569ABCD", countryCodeField: ""},
			expectedOutput: numberLookupResponse{
				Info: number.Info{PhoneNumber: "+1212569ABCD"},
				Err:  &number.Info{PhoneNumber: number.InvalidPhoneNumberError.Error()},
			},
		},
		{
			method: http.MethodGet, expectedStatusCode: http.StatusBadRequest,
			parameters: map[string]string{phoneNumberField: "631 311 8150", countryCodeField: ""},
			expectedOutput: numberLookupResponse{
				Info: number.Info{PhoneNumber: "631 311 8150"},
				Err:  &number.Info{RegionCode: number.MissingValueError.Error()},
			},
		},
		{
			method: http.MethodGet, expectedStatusCode: http.StatusBadRequest,
			parameters: map[string]string{phoneNumberField: "+12125690123", countryCodeField: "CA"},
			expectedOutput: numberLookupResponse{
				Info: number.Info{PhoneNumber: "+12125690123", RegionCode: "CA"},
				Err:  &number.Info{RegionCode: number.InvalidCombination.Error()},
			},
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			httpRequest, httpResponse := getHTTPTestRequestResponse(tt.method, url, tt.parameters)
			numberLookupHandler(httpResponse, httpRequest)
			if httpResponse.Code != tt.expectedStatusCode {
				t.Errorf("numberLookupHandler() StatusCode= %v, want %v", httpResponse.Code, tt.expectedStatusCode)
				return
			}

			var bodyJson numberLookupResponse
			err := json.Unmarshal(httpResponse.Body.Bytes(), &bodyJson)
			if err != nil {
				t.Errorf("numberLookupHandler() failes to parse json body: %v", err)
				return
			}

			if !reflect.DeepEqual(bodyJson, tt.expectedOutput) {
				t.Errorf("numberLookupHandler()\nbodyJson = %+v\nwant =     %+v", bodyJson, tt.expectedOutput)
			}

		})
	}
}
