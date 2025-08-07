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

package recaptcha3

import (
	"errors"
	"fmt"
	"time"

	resty "resty.dev/v3"

	"github.com/go-playground/validator/v10"
)

//-----------------------------------------------------------------------------

// Verifier defines how to verify a token.
type Verifier interface {
	Verify(token string) (bool, error)
}

//-----------------------------------------------------------------------------

// Config holds the parameters to build a real verifier
type Config struct {
	// The secret key for reCAPTCHA v3. Required.
	Secret string

	// A string indicating the expected action. For example: "login" or "contact".
	ExpectedAction string

	// The lowest acceptable score for the reCAPTCHA validation.
	MinScore float64

	// Optional Resty client
	Client *resty.Client
}

//-----------------------------------------------------------------------------

type RestyConfig struct {
	RetryCount int
	RetryWaitTime int
	RetryMaxWaitTime int
}

//-----------------------------------------------------------------------------

type GoogleVerifier struct {
	Secret         string
	ExpectedAction string
	MinScore       float64
	Client         *resty.Client
}

//-----------------------------------------------------------------------------

// reCAPTCHA response
type RecaptchaResponse struct {
	Success     bool     `json:"success"`
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	Score       float64  `json:"score"`
	Action      string   `json:"action"`
	ErrorCodes  []string `json:"error-codes"`
}

//-----------------------------------------------------------------------------

func NewGoogleVerifier(cfg Config) (*GoogleVerifier, error) {
	var client *resty.Client

	if cfg.Secret == "" {
		return nil, errors.New("recaptcha secret is required")
	}
	if cfg.Client == nil {
		client = resty.New().
			SetRetryCount(3).
			SetRetryWaitTime(2 * time.Second).
			SetRetryMaxWaitTime(8 * time.Second).
			AddRetryConditions(
				func(r *resty.Response, err error) bool {
					return err != nil || r.StatusCode() >= 500
				})
	} else {
		client = cfg.Client
	}
	return &GoogleVerifier{
		Secret:         cfg.Secret,
		ExpectedAction: cfg.ExpectedAction,
		MinScore:       cfg.MinScore,
		Client:         client,
	}, nil
}

//-----------------------------------------------------------------------------

func (g *GoogleVerifier) Verify(token string) (bool, error) {
	resp := RecaptchaResponse{}
	_, err := g.Client.R().
		SetFormData(map[string]string{
			"secret":   g.Secret,
			"response": token,
		}).
		SetResult(&resp).
		Post("https://www.google.com/recaptcha/api/siteverify")

	if err != nil {
		return false, err
	}

	if !resp.Success || resp.Score < g.MinScore {
		return false, errors.New("reCAPTCHA verification failed")
	}

	if g.ExpectedAction != "" && resp.Action != g.ExpectedAction {
		return false, errors.New("unexpected reCAPTCHA action")
	}

	return true, nil
}

//-----------------------------------------------------------------------------

func RegisterRecaptchaValidator(v *validator.Validate, fieldTag string, verifier Verifier) error {
	return v.RegisterValidation(fieldTag, func(fl validator.FieldLevel) bool {
		token := fl.Field().String()
		if token == "" {
			return false
		}
		ok, err := verifier.Verify(token)
		if err != nil {
			fmt.Println("reCAPTCHA error:", err)
		}
		return ok
	})
}
