/*

   Copyright 2025 Rodolfo González González

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/

package main

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	recaptcha "github.com/rgglez/go-playground-recaptcha3-validator"
)

// Your form struct with a recaptcha_token field
type ContactForm struct {
	Name           string `json:"name" validate:"required"`
	Email          string `json:"email" validate:"required,email"`
	RecaptchaToken string `json:"recaptcha_token" validate:"required,recaptcha"`
}

func main() {
	// Create a mock verifier: success = true
	mock := &recaptcha.MockVerifier{ShouldPass: true}

	// Setup validator and register custom "recaptcha" rule
	validate := validator.New()
	err := recaptcha.RegisterRecaptchaValidator(validate, "recaptcha", mock)
	if err != nil {
		panic(err)
	}

	// Simulate incoming data
	form := ContactForm{
		Name:           "Alice",
		Email:          "alice@example.com",
		RecaptchaToken: "mock-token",
	}

	// Validate
	if err := validate.Struct(form); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
		return
	}

	fmt.Println("✅ Mock reCAPTCHA and form validated successfully!")
}
