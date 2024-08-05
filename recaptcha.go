package router

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ReCaptcher interface {
	GetReCaptcha() string
}

const recaptchaServerName = "https://www.google.com/recaptcha/api/siteverify"

var secretKey string

// RecaptchaResponse is struct that contains recaptcha validating response fields
type RecaptchaResponse struct {
	Success     bool      `json:"success"`
	Score       float64   `json:"score"`
	Action      string    `json:"action"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
}

// checkReCaptcha check if recaptcha is valid ot raise error
func checkReCaptcha[T ReCaptcher](data T) error {
	if err := recaptchaResponse(data.GetReCaptcha()); err != nil {
		return err
	}
	return nil
}

// SetRecaptchaKey is function that set recaptcha secret key
func SetRecaptchaKey(key string) {
	secretKey = key
}

// VerifyReCaptcha is ...
func VerifyReCaptcha(recaptcha string) (RecaptchaResponse, error) {
	resp, err := http.PostForm(recaptchaServerName,
		url.Values{"secret": {secretKey}, "response": {recaptcha}})
	if err != nil {
		log.Error().Err(err).Msg("VerifyReCaptcha PostForm error")
		return RecaptchaResponse{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error().Err(err).Msg("verifyReCaptcha Body.Close() error")
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("verifyReCaptcha io.ReadAll(resp.Body) error")
		return RecaptchaResponse{}, err
	}

	var responseData RecaptchaResponse
	if err := json.Unmarshal(body, &responseData); err != nil {
		log.Error().Err(err).Msg("verifyReCaptcha json.Unmarshal(body, &responseData) error")
		return RecaptchaResponse{}, err
	}

	return responseData, nil
}

func recaptchaResponse(recaptcha string) error {
	response, err := VerifyReCaptcha(recaptcha)
	if err != nil {
		log.Error().Err(err).Msg("signUpUser verifyReCaptcha error")
		return err
	}

	if !response.Success {
		err = errors.New(fmt.Sprintf("Recaptcha error: success: %v, score: %v, error codes: %s",
			response.Success, response.Score, strings.Join(response.ErrorCodes, ", ")))
		log.Error().Err(err).Msg(err.Error())
		return err
	}
	return nil
}
