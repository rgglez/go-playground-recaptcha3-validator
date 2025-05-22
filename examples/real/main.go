package main

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	recaptcha "github.com/rgglez/go-playground-recaptcha3-validator"
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
	verifier, err := recaptcha.NewGoogleVerifier(recaptcha.Config{
		Secret:         secret,
		ExpectedAction: "contact",
		MinScore:       0.5,
	})
	if err != nil {
		panic(err)
	}

	// Set up validator
	validate := validator.New()
	err = recaptcha.RegisterRecaptchaValidator(validate, "recaptcha", verifier)
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
