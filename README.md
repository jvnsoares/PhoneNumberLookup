# PhoneNumberLookup

[![Build Status](https://github.com/jvnsoares/PhoneNumberLookup/actions/workflows/CI.yml/badge.svg)](https://github.com/jvnsoares/PhoneNumberLookup/actions/workflows/CI.yml)
[![codecov](https://codecov.io/gh/jvnsoares/PhoneNumberLookup/branch/main/graph/badge.svg)](https://codecov.io/gh/jvnsoares/PhoneNumberLookup)
[![Deployment Status](https://github.com/jvnsoares/PhoneNumberLookup/actions/workflows/CD.yml/badge.svg)](https://github.com/jvnsoares/PhoneNumberLookup/actions/workflows/CD.yml)

`PhoneNumberLookup` is a microservices that gives the `Country Code` in ISO 3166-1 alpha-2 format, `Area Code` and the
`Local Phone Number` of a given `Phone Number`.

The service has only endpoint `/v1/phone-numbers` that accepts GET requests. It accepts the following parameters:

- `phoneNumber`: sequence of digits in the format `[+][country code][area code][local phone number]`, the `+` is
  optional.

  Phone numbers can have one white space between country code and area code, and one white space between area code and
  the local phone number. Any other white space or characters invalid the Phone number.

  The user must provide the `countryCode` if the phone number is missing the country code.


- `countryCode`: Country code in ISO 3166-1 alpha-2 format.

#### Return codes:

- `200`: It received valid `phoneNumber` and `countryCode` and was able to find expected information.
- `400`: It received an empty or invalid `phoneNumber`, or an invalid combination of phone number + country code
- `405`: It received a method other than `GET` request.

### Building

You can build a docker image of `PhoneNumberLookup` using the following command:

```shell
docker build -f build/docker/Dockerfile -t phone-number-lookup:local .
```

It will create a docker imaged called `phone-number-lookup` tagged `local` in your local machine.

Alternatively you can fetch the latest build of `PhoneNumberLookup` from dockerhub using the command

```shell
  <PLACE HOLDER>
```

### Running

You can start the `PhoneNumberLookup` service from using the command:

```shell
 docker run -p 8008:8008 <PLACE HOLDER>
```

### Deployment

You can deploy this service to production using `Docker` or a container orchestration service of you choice.

You must ensure that the clients are in the same subnet of the service.


### Assumptions

- Malformed of invalid combinations of `phoneNumber` + `countryCode` are seen as client side error and return
the error `400 Bad Request`.
- Endpoint processes only `GET` requests and returns `405 Method Not Allowed` for any other type of request.

### Possible Improvements

- Return more information about the given number like:
  - Carrier
  - Number in international format
  - Number in national format
  - Phone number type
  - Accept more than just white spaces
  - Accept letter mapping

### FAQ
 - Why golang?

This is the language I am most familiar.

 - Why didn't you use a web framework?

Go `net/http` is very complete and easy to use on small services. I didn't find a need for a framework.

- Why did you use `nyaruka/phonenumbers`?

As a port of Google's `libphonenumber`, `nyaruka/phonenumbers` offers all features need to develop this service.
It also uses `go.mod`, so it is easily importable.

### Examples

- Input: `phoneNumber: "+12125690123"`

  Status Code: `200`

  Output:
  ```json
  {
    "phoneNumber": "+12125690123",
    "countryCode": "US",
    "areaCode": "212",
    "localPhoneNumber": "5690123"
  }
  ```

- Input: `phoneNumber: "+1 212 5690123"`

  Status Code: `200`

  Output:
  ```json
  {
    "phoneNumber": "+1 212 5690123",
    "countryCode": "US",
    "areaCode": "212",
    "localPhoneNumber": "5690123"
  }
  ```

- Input: `phoneNumber: "+1 212 569 0123"`

  Status Code: `400`

  Output:
  ```json
  {
    "phoneNumber": "+1 212 569 0123",
    "error": {
      "phoneNumber": "invalid phone number"
    }
  }
  ```

- Input: `phoneNumber: "2125690123"`

  Status Code: `400`

  Output:
  ```json
  {
    "phoneNumber": "2125690123",
    "error": {
      "countryCode": "required value is missing"
    }
  }
  ```

- Input: `phoneNumber: "+12125690123", countryCode: "CA"`

  Status Code: `400`

  Output:
  ```json
  {
    "phoneNumber": "+12125690123",
    "countryCode": "CA",
    "error": {
      "countryCode": "invalid combination of phone number and country code"
    }
  }
  ```

- Input: `phoneNumber: ""`

  Status Code: `400`

  Output:
  ```json
  {
    "error": {
      "phoneNumber": "required value is missing"
    }
  }
  ```