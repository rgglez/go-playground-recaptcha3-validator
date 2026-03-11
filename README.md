# go-playground-recaptcha3-validator

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
![GitHub all releases](https://img.shields.io/github/downloads/rgglez/go-playground-recaptcha3-validator/total)
![GitHub issues](https://img.shields.io/github/issues/rgglez/go-playground-recaptcha3-validator)
![GitHub commit activity](https://img.shields.io/github/commit-activity/y/rgglez/go-playground-recaptcha3-validator)
[![Go Report Card](https://goreportcard.com/badge/github.com/rgglez/go-playground-recaptcha3-validator)](https://goreportcard.com/report/github.com/rgglez/go-playground-recaptcha3-validator)
[![GitHub release](https://img.shields.io/github/release/rgglez/go-playground-recaptcha3-validator.svg)](https://github.com/rgglez/go-playground-recaptcha3-validator/releases/)
![GitHub stars](https://img.shields.io/github/stars/rgglez/go-playground-recaptcha3-validator?style=social)
![GitHub forks](https://img.shields.io/github/forks/rgglez/go-playground-recaptcha3-validator?style=social)

**go-playground-recaptcha3-validator** is a custom validator for [playground](https://github.com/go-playground/validator/) implementing [Google reCAPTCHA v3](https://developers.google.com/recaptcha/docs/v3?hl=es-419) validation.

## Installation

```bash
go get github.com/rgglez/go-playground-recaptcha3-validator
```

## Usage

First import the module:

```go
import (
   recaptcha3 "github.com/rgglez/go-playground-recaptcha3-validator"
)
```

This is a sample structure to validate:

```go
// Struct to bind and validate
type ContactForm struct {
	Name           string `json:"name" validate:"required"`
	Email          string `json:"email" validate:"required,email"`
	RecaptchaToken string `json:"recaptcha_token" validate:"required,recaptcha"`
}
```

Set it up:

```go
// Load reCAPTCHA secret
secret := os.Getenv("RECAPTCHA3_SECRET_KEY") // this is just a sample name
if secret == "" {
    fmt.Println("Missing RECAPTCHA3_SECRET_KEY env variable")
    os.Exit(1)
}

// Create Google verifier
verifier, err := recaptcha3.NewGoogleVerifier(recaptcha3.Config{
    Secret:         secret,
    ExpectedAction: "contact",
    MinScore:       0.5,
})
if err != nil {
    panic(err)
}

// Set up validator
validate := validator.New()
err = recaptcha3.RegisterRecaptchaValidator(validate, "recaptcha", verifier)
if err != nil {
    panic(err)
}
```

Call it in your handler:

```go
if err := validate.Struct(form); err != nil {
    return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
        "error": err.Error(),
    })
}
```

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

Copyright (c) 2026 Rodolfo González González

Licensed under the [Apache 2.0](LICENSE) license. Read the [LICENSE](LICENSE) file.

