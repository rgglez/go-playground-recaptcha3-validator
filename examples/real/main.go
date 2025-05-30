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
	"os"

	"github.com/go-playground/validator/v10"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	recaptcha3 "github.com/rgglez/go-playground-recaptcha3-validator"
)

// Struct to bind and validate
type ContactForm struct {
	Name           string `json:"name" validate:"required"`
	Email          string `json:"email" validate:"required,email"`
	RecaptchaToken string `json:"recaptcha_token" validate:"required,recaptcha"`
}

func main() {
	// Load reCAPTCHA secret
	secret := os.Getenv("RECAPTCHA3_SECRET_KEY")
	if secret == "" {
		fmt.Println("Missing RECAPTCHA3_SECRET_KEY env variable")
		os.Exit(1)
	}

	// Create real Google verifier
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

	// Start Fiber app
	app := fiber.New()

	app.Use(cors.New(cors.Config{AllowOrigins: "*"}))

	app.Post("/contact", func(c *fiber.Ctx) error {
		var form ContactForm
		if err := c.BodyParser(&form); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid JSON",
			})
		}

		if err := validate.Struct(form); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"message": "Form sent correctly",
		})
	})

	app.Listen(":3000")
}
