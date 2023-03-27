package number

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/nyaruka/phonenumbers"
)

var (
	InvalidPhoneNumberError = errors.New("invalid phone number")
	InvalidCombination      = errors.New("invalid combination of phone number and country code")
	MissingValueError       = errors.New("required value is missing")
)

var (
	onlyNumberRegex = regexp.MustCompile(`^\d+$`)
)

type Info struct {
	PhoneNumber      string `json:"phoneNumber,omitempty"`
	RegionCode       string `json:"countryCode,omitempty"` // ISO 3166-1 alpha-2, e.g.: "US", "UK", "CA"
	CountryCode      string `json:"-"`                     // ITU-T E.164, e.g.: 1, 52, 34
	AreaCode         string `json:"areaCode,omitempty"`
	LocalPhoneNumber string `json:"localPhoneNumber,omitempty"`
}

// isPhoneNumberMadeOfDigits check if the given number is composed of only digits.
// this function is necessary since nyaruka/phonenumbers accepts letters as valid phone numbers.
// this function accepts phone numbers with starting '+' and white spaces
func isPhoneNumberMadeOfDigits(number string) bool {
	return onlyNumberRegex.MatchString(strings.TrimPrefix(strings.ReplaceAll(number, " ", ""), "+"))
}

// checkWhiteSpaces checks if the number has whitespaces in the right places, between country, area code and local phone number
// it also checks if the number has only one white space between country, area code and local phone number
func hasCorrectWhiteSpaces(info Info, number string) bool {
	numberWithoutPlus := strings.TrimPrefix(number, "+")
	parts := strings.Split(numberWithoutPlus, " ")
	nParts := len(parts)

	switch nParts {
	case 1:
		//doesn't have whitespaces
		return true
	case 2:
		if info.AreaCode == "" {
			if numberWithoutPlus == fmt.Sprintf("%s %s", info.CountryCode, info.LocalPhoneNumber) {
				return true
			}
		} else {
			if numberWithoutPlus == fmt.Sprintf("%s %s%s", info.CountryCode, info.AreaCode, info.LocalPhoneNumber) ||
				numberWithoutPlus == fmt.Sprintf("%s%s %s", info.CountryCode, info.AreaCode, info.LocalPhoneNumber) {
				return true
			}
		}

	case 3:
		expectedNumber := fmt.Sprintf("%s %s %s", info.CountryCode, info.AreaCode, info.LocalPhoneNumber)
		if info.AreaCode != "" && numberWithoutPlus == expectedNumber {
			return true
		}

	default:
		// the number should have at maximum 3 parts (2 whitespaces)
	}

	return false
}

// LookupNumber returns the information of a given phone number.
// the phone number must be a sequence of digits in the format [+][country code][area code][local phone number]
// the '+' is optional
// phones can have one white space between country code, area code and local phone number,
// any other white space is invalid
// if the phone number is missing the country code then the caller must provide user must provide countryCode parameter in ISO 3166-1 alpha-2 format.
func LookupNumber(rawNumber, countryCode string) (info Info, err error) {
	if !isPhoneNumberMadeOfDigits(rawNumber) {
		return info, InvalidPhoneNumberError
	}

	rawNumberPlus := rawNumber
	// if we don't have a country code the phone number must start with a '+' to be parsable
	if countryCode == "" && strings.HasPrefix(rawNumber, "+") {
		rawNumberPlus = "+" + rawNumber
	}

	number, err := phonenumbers.Parse(rawNumberPlus, countryCode)
	if err != nil {
		return info, err
	}

	if !phonenumbers.IsValidNumber(number) {
		return info, InvalidPhoneNumberError
	}

	nationalNumber := fmt.Sprintf("%d", number.GetNationalNumber())
	ACLength := phonenumbers.GetLengthOfGeographicalAreaCode(number)

	info.PhoneNumber = rawNumber
	info.RegionCode = phonenumbers.GetRegionCodeForNumber(number)
	info.CountryCode = strconv.Itoa(int(number.GetCountryCode()))
	info.AreaCode = nationalNumber[:ACLength]
	info.LocalPhoneNumber = nationalNumber[ACLength:]

	if !hasCorrectWhiteSpaces(info, rawNumber) {
		return Info{}, InvalidPhoneNumberError
	}

	if countryCode != "" && info.RegionCode != countryCode {
		return Info{PhoneNumber: rawNumber, RegionCode: info.RegionCode}, InvalidCombination
	}

	return info, nil
}
