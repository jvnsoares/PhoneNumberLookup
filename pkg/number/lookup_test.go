package number

import (
	"errors"
	"reflect"
	"testing"

	"github.com/nyaruka/phonenumbers"
)

func Test_isNumberMadeOfDigits(t *testing.T) {
	tests := []struct {
		number        string
		hasOnlyDigits bool
	}{
		{number: "+1    212     5690123", hasOnlyDigits: true},
		{number: "+1 212 5690123", hasOnlyDigits: true},
		{number: "1 212 5690123", hasOnlyDigits: true},
		{number: "123ccd", hasOnlyDigits: false},
		{number: "hello world", hasOnlyDigits: false},
		{number: "+1 (212) 5690123", hasOnlyDigits: false},
		{number: "1-212-5690123", hasOnlyDigits: false},
	}
	for _, tt := range tests {
		t.Run(tt.number, func(t *testing.T) {
			if got := isPhoneNumberMadeOfDigits(tt.number); got != tt.hasOnlyDigits {
				t.Errorf("isNumberMadeOfDigits() = %v, want %v", got, tt.hasOnlyDigits)
			}
		})
	}
}

func Test_hasCorrectWhiteSpaces(t *testing.T) {
	usPhoneInfo := Info{CountryCode: "1", AreaCode: "212", LocalPhoneNumber: "5690123"}
	shPhoneInfo := Info{CountryCode: "290", AreaCode: "", LocalPhoneNumber: "21234"}
	tests := []struct {
		info      Info
		number    string
		isCorrect bool
	}{
		{info: usPhoneInfo, number: "+12125690123", isCorrect: true},
		{info: usPhoneInfo, number: "+1 2125690123", isCorrect: true},
		{info: usPhoneInfo, number: "+121 25690123", isCorrect: false},
		{info: usPhoneInfo, number: "+1 212 5690123", isCorrect: true},
		{info: usPhoneInfo, number: "+1 212   5690123", isCorrect: false},
		{info: usPhoneInfo, number: "+1   212 5690123", isCorrect: false},
		{info: usPhoneInfo, number: "+1212   5690123", isCorrect: false},
		{info: usPhoneInfo, number: "+1   2125690123", isCorrect: false},
		{info: usPhoneInfo, number: "+12125690 123", isCorrect: false},
		{info: usPhoneInfo, number: "+1 212 5690 123", isCorrect: false},
		{info: usPhoneInfo, number: "+1212 5690 123", isCorrect: false},
		{info: usPhoneInfo, number: "+1212 5690 123", isCorrect: false},

		{info: shPhoneInfo, number: "+29021234", isCorrect: true},
		{info: shPhoneInfo, number: "+290 21234", isCorrect: true},
		{info: shPhoneInfo, number: "+290  21234", isCorrect: false},
		{info: shPhoneInfo, number: "+290 2 1234", isCorrect: false},
		{info: shPhoneInfo, number: "+2 9021234", isCorrect: false},
		{info: shPhoneInfo, number: "+2 902 1234", isCorrect: false},
	}
	for _, tt := range tests {
		t.Run(tt.number, func(t *testing.T) {
			if got := hasCorrectWhiteSpaces(tt.info, tt.number); got != tt.isCorrect {
				t.Errorf("hasCorrectWhiteSpaces() = %v, want %v", got, tt.isCorrect)
			}
		})
	}
}

func TestLookupNumber(t *testing.T) {
	tests := []struct {
		rawNumber     string
		countryCode   string
		expectedInfo  Info
		expectedError error
	}{
		{
			rawNumber: "+12125690123", countryCode: "", expectedError: nil,
			expectedInfo: Info{PhoneNumber: "+12125690123", RegionCode: "US", CountryCode: "1", AreaCode: "212", LocalPhoneNumber: "5690123"},
		},

		{
			rawNumber: "+1  2125690123", countryCode: "", expectedError: InvalidPhoneNumberError,
			expectedInfo: Info{},
		},
		{
			rawNumber: "+1212569012", countryCode: "", expectedError: InvalidPhoneNumberError,
			expectedInfo: Info{},
		},
		{
			rawNumber: "+1212569ABCD", countryCode: "", expectedError: InvalidPhoneNumberError,
			expectedInfo: Info{},
		},
		{
			rawNumber: "2125690123", countryCode: "", expectedError: phonenumbers.ErrInvalidCountryCode,
			expectedInfo: Info{},
		},
		{
			rawNumber: "+12125690123", countryCode: "CA", expectedError: InvalidCombination,
			expectedInfo: Info{PhoneNumber: "+12125690123", RegionCode: "US"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.rawNumber+tt.countryCode, func(t *testing.T) {
			gotInfo, err := LookupNumber(tt.rawNumber, tt.countryCode)
			if err != nil {
				if !errors.Is(err, tt.expectedError) {
					t.Errorf("LookupNumber() error = %v, expected error %v", err, tt.expectedError)
					return
				}
				return
			}
			if !reflect.DeepEqual(gotInfo, tt.expectedInfo) {
				t.Errorf("LookupNumber() gotInfo = %v, want %v", gotInfo, tt.expectedInfo)
			}
		})
	}
}
