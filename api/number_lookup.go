package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/nyaruka/phonenumbers"

	"phone-number-lookup/pkg/number"
)

const (
	numberLookupEndpoint = "/v1/phone-numbers"
	phoneNumberField     = "phoneNumber"
	countryCodeField     = "countryCode"
)

type numberLookupResponse struct {
	number.Info
	Err *number.Info `json:"error,omitempty"`
}

// numberLookupHandler accepts only GET requests
// this handler accepts two arguments:
// - phoneNumberField: sequence of digits in the format [+][country code][area code][local phone number], the '+' is optional
// phones can have one white space between country code, area code and local phone number,
// any other white space is invalid
// if the phone number is missing the country code then the caller must provide user must provide the countryCodeField.
// - countryCodeField: Country code in ISO 3166-1 alpha-2 format.
//
// Return codes:
// - 200: It received valid phoneNumberField and countryCodeField and was able to find expected information.
// - 400: It received an missing or invalid phone number, or an invalid combination of phone number + country code
// - 405: It received a method other than GET request
//
// Examples:
// Input: phoneNumberField: "+12125690123"
// Status Code: 200
// Output:
//{
//    "phoneNumber": "+12125690123",
//    "countryCode": "US",
//    "areaCode": "212",
//    "localPhoneNumber": "5690123"
//}
//
// Input: phoneNumberField: "+1 212 5690123"
// Status Code: 200
// Output:
//{
//    "phoneNumber": "+1 212 5690123",
//    "countryCode": "US",
//    "areaCode": "212",
//    "localPhoneNumber": "5690123"
//}
//
// Input: phoneNumberField: "+1 212 569 0123"
// Status Code: 400
// Output:
//{
//    "phoneNumber": "+1 212 569 0123",
//    "error": {
//        "phoneNumber": "invalid phone number"
//    }
//}
//
// Input: phoneNumberField: "2125690123"
// Status Code: 400
// Output:
//{
//    "phoneNumber": "2125690123",
//    "error": {
//        "countryCode": "required value is missing"
//    }
//}
//
//
// Input: phoneNumberField: "+12125690123", countryCodeField: "CA"
// Status Code: 400
//{
//    "phoneNumber": "+12125690123",
//    "countryCode": "CA",
//    "error": {
//        "countryCode": "invalid combination of phone number and country code"
//    }
//}
//
// Input: phoneNumberField: ""
//// Status Code: 400
//{
//    "error": {
//        "phoneNumber": "required value is missing"
//    }
//}

func numberLookupHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var response numberLookupResponse
	defer func() { json.NewEncoder(w).Encode(response) }()

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	log.Println("received GET request:", r.URL.Query())
	phoneNumber, countryCode := parseNumberLookupRequest(r)

	// we received an empty phoneNumberField
	if phoneNumber == "" {
		w.WriteHeader(http.StatusBadRequest)
		response.Err = &number.Info{PhoneNumber: number.MissingValueError.Error()}
		return
	}

	info, err := number.LookupNumber(phoneNumber, countryCode)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response.Info.PhoneNumber = phoneNumber
		if errors.Is(err, phonenumbers.ErrInvalidCountryCode) && countryCode == "" {
			response.Err = &number.Info{RegionCode: number.MissingValueError.Error()}
		} else if errors.Is(err, number.InvalidCombination) {
			response.Info.RegionCode = countryCode
			response.Err = &number.Info{RegionCode: number.InvalidCombination.Error()}
		} else {
			response.Err = &number.Info{PhoneNumber: err.Error()}
		}
	} else {
		w.WriteHeader(http.StatusOK)
		response.Info = info
	}

	return
}

func parseNumberLookupRequest(r *http.Request) (phoneNumber, countryCode string) {
	query := r.URL.Query()
	phoneNumber = query.Get(phoneNumberField)
	if phoneNumber == "" {
		return "", ""
	}

	countryCode = query.Get(countryCodeField)
	return phoneNumber, countryCode
}
