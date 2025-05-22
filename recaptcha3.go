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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

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

	// Optional http.Client object.
	HTTPClient *http.Client
}

//-----------------------------------------------------------------------------

type GoogleVerifier struct {
	secret         string
	expectedAction string
	minScore       float64
	client         *http.Client
}

//-----------------------------------------------------------------------------

func NewGoogleVerifier(cfg Config) (*GoogleVerifier, error) {
	if cfg.Secret == "" {
		return nil, errors.New("recaptcha secret is required")
	}
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = &http.Client{Timeout: 5 * time.Second}
	}
	return &GoogleVerifier{
		secret:         cfg.Secret,
		expectedAction: cfg.ExpectedAction,
		minScore:       cfg.MinScore,
		client:         cfg.HTTPClient,
	}, nil
}

//-----------------------------------------------------------------------------

func (g *GoogleVerifier) Verify(token string) (bool, error) {
	req, _ := http.NewRequest("POST", "https://www.google.com/recaptcha/api/siteverify", nil)
	q := req.URL.Query()
	q.Add("secret", g.secret)
	q.Add("response", token)
	req.URL.RawQuery = q.Encode()

	resp, err := g.client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var body struct {
		Success    bool     `json:"success"`
		Score      float64  `json:"score"`
		Action     string   `json:"action"`
		ErrorCodes []string `json:"error-codes"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return false, err
	}

	if !body.Success {
		return false, fmt.Errorf("recaptcha failed: %v", body.ErrorCodes)
	}
	if body.Action != g.expectedAction {
		return false, fmt.Errorf("unexpected action: %s", body.Action)
	}
	if body.Score < g.minScore {
		return false, fmt.Errorf("low score: %f", body.Score)
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
