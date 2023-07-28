package backend

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/BlahajXD/backend/config"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type HackathonGenerateTokenInput struct {
	Username      string `json:"username"`
	LoginPassword string `json:"loginPassword"`
}

func (d *Dependency) HackathonGenerateToken(input HackathonGenerateTokenInput) (string, error) {
	payload := map[string]interface{}{
		"username":      input.Username,
		"loginPassword": input.LoginPassword,
	}

	log.Debug().Interface("payload", payload).Msg("hackathon: generate token payload")

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal payload")
	}

	type apiResponse struct {
		TraceId string `json:"traceId"`
		Data    struct {
			AccessToken string `json:"accessToken"`
		} `json:"data"`
		ErrCode string `json:"errCode"`
		Success bool   `json:"success"`
		ErrMsg  string `json:"errMsg"`
	}

	header := http.Header{}
	header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	response, err := d.httpclient.Post(config.HackathonApiURL()+"/user/auth/token", bytes.NewReader(jsonPayload), header)
	if err != nil {
		return "", errors.Wrap(err, "hackathon: failed to send request")
	}

	var apiResponseData apiResponse
	err = json.NewDecoder(response.Body).Decode(&apiResponseData)
	if err != nil {
		return "", errors.Wrap(err, "hackathon: failed to decode response")
	}

	log.Debug().Interface("apiResponseData", apiResponseData).Msg("hackathon: api response data")

	if !apiResponseData.Success {
		if apiResponseData.ErrCode == "1025" {
			return "", errors.Wrap(ErrHasNoDataPermission, apiResponseData.ErrMsg)
		}

		return "", errors.Wrap(ErrHackathonUnknownError, apiResponseData.ErrMsg)
	}

	return apiResponseData.Data.AccessToken, nil
}
