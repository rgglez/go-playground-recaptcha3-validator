# go-playground-recaptcha3-validator

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
![GitHub all releases](https://img.shields.io/github/downloads/rgglez/go-playground-recaptcha3-validator/total)
![GitHub issues](https://img.shields.io/github/issues/rgglez/go-playground-recaptcha3-validator)
![GitHub commit activity](https://img.shields.io/github/commit-activity/y/rgglez/go-playground-recaptcha3-validator)
[![Go Report Card](https://goreportcard.com/badge/github.com/rgglez/go-playground-recaptcha3-validator)](https://goreportcard.com/report/github.com/rgglez/go-playground-recaptcha3-validator)
[![GitHub release](https://img.shields.io/github/release/rgglez/go-playground-recaptcha3-validator.svg)](https://github.com/rgglez/go-playground-recaptcha3-validator/releases/)

**go-playground-recaptcha3-validator** is a custom validator for [playground](https://github.com/go-playground/validator/) implementing [Google reCAPTCHA v3](https://developers.google.com/recaptcha/docs/v3?hl=es-419) validation.

## Installation

```bash
go get github.com/rgglez/go-playground-recaptcha3-validator
```

## Usage

See the [real example](examples/real/).

## Configuration

* **Secret** the secret key for reCAPTCHA v3. Required.
* **ExpectedAction** a string specifying the expected action. For example: "login" or "contact".
* **MinScore** a float value specifying the lowest acceptable score for the reCAPTCHA validation.
* **HTTPClient** a pointer to an optional http.Client object

## Examples

Two examples are provided:

* *[mock](examples/mock/)* for testing (it does not contact Google's servers).
* *[real](examples/real/)* to try the validator, using the sample *[web page](examples/frontend/index.html)*. You need to copy that webpage (or create a similar) to your web server, and use your own site key.

## Dependencies

* go get github.com/go-playground/validator/v10

## Security

* Remember that reCAPTCHA v3 does not stops bots directly, but just gives a risk score. It is responsability of the application to accept the request or deny it.
* Another captcha (such as reCAPTCHA v2) could be shown with a classic challenge as a "plan B", if the score is lower than expected. 

## License

Copyright (c) 2025 Rodolfo González González

Licensed under the [Apache 2.0](LICENSE) license. Read the [LICENSE](LICENSE) file.

