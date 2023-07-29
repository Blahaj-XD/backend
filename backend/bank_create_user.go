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

type BankCreateUserInput struct {
	KTPID         string `json:"ktpId"`
	Username      string `json:"username"`
	LoginPassword string `json:"loginPassword"`
	PhoneNumber   string `json:"phoneNumber"`
	BirthDate     string `json:"birthDate"`
	Gender        int    `json:"gender"`
	Email         string `json:"email"`
}

func (d *Dependency) BankCreateUser(input BankCreateUserInput) (int, error) {
	payload := map[string]any{
		"ktpId":         input.KTPID,
		"username":      input.Username,
		"loginPassword": input.LoginPassword,
		"phoneNumber":   input.PhoneNumber,
		"birthDate":     input.BirthDate,
		"gender":        input.Gender,
		"email":         input.Email,
	}

	log.Debug().Interface("payload", payload).Msg("bank_create_user: create user payload")

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal payload")
		return 0, errors.Wrap(err, "failed to marshal payload")
	}

	log.Debug().Str("payload", string(jsonPayload)).Msg("bank_create_user: create user payload")

	type apiResponse struct {
		TraceId string `json:"traceId"`
		Data    struct {
			UID         int    `json:"uid"`
			PhoneNumber string `json:"phoneNumber"`
			Gender      string `json:"gender"`
			CreateTime  int    `json:"createTime"`
			KTPID       string `json:"ktpId"`
			UpdateTime  int    `json:"updateTime"`
			BirthDate   string `json:"birthDate"`
			Email       string `json:"email"`
			Username    string `json:"username"`
		} `json:"data"`
		ErrCode string `json:"errCode"`
		Success bool   `json:"success"`
		ErrMsg  string `json:"errMsg"`
	}

	header := http.Header{}
	header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	response, err := d.httpclient.Post(config.BankApiURL()+"/user/auth/create", bytes.NewReader(jsonPayload), header)
	if err != nil {
		log.Error().Err(err).Msg("bank_create_user: failed to send request")
		return 0, errors.Wrap(err, "bank_create_user: failed to send request")
	}

	log.Debug().Msg("bank_create_user: sent request")

	var apiResponseData apiResponse
	err = json.NewDecoder(response.Body).Decode(&apiResponseData)
	if err != nil {
		return 0, errors.Wrap(err, "bank_create_user: failed to decode response")
	}

	log.Debug().Interface("apiResponseData", apiResponseData).Msg("bank_create_user: api response data")

	if !apiResponseData.Success {
		if apiResponseData.ErrCode == "1020" {
			return 0, errors.Wrap(ErrBankRequestParameter, apiResponseData.ErrMsg)
		}

		if apiResponseData.ErrCode == "4873" {
			return 0, ErrUserAlreadyExists
		}

		return 0, errors.Wrap(ErrBankUnknownError, apiResponseData.ErrMsg)
	}

	return apiResponseData.Data.UID, nil
}
