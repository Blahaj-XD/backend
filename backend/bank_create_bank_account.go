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

type BankCreateBankAccountInput struct {
	Balance     float64 `json:"balance"`
	AccessToken string
}

func (d *Dependency) BankCreateBankAccount(input BankCreateBankAccountInput) (string, error) {
	payload := map[string]any{
		"balance":      input.Balance,
		"access_token": input.AccessToken,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal payload")
	}

	log.Debug().Interface("payload", payload).Msg("bank_create_bank_account: create bank account payload")

	type apiResponse struct {
		TraceId string `json:"traceId"`
		Data    struct {
			UID         int    `json:"uid"`
			Balance     int    `json:"balance"`
			AccountName string `json:"accountName"`
			CreateTime  int    `json:"createTime"`
			AccountNo   string `json:"accountNo"`
			UpdateTime  int    `json:"updateTime"`
			Status      string `json:"status"`
		} `json:"data"`
		ErrCode string `json:"errCode"`
		Success bool   `json:"success"`
		ErrMsg  string `json:"errMsg"`
	}

	header := http.Header{}
	header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	header.Add(fiber.HeaderAuthorization, "Bearer "+input.AccessToken)

	response, err := d.httpclient.Post(config.BankApiURL()+"/bankAccount/create", bytes.NewReader(jsonPayload), header)
	if err != nil {
		return "", errors.Wrap(err, "bank_create_bank_account: failed to send request")
	}

	log.Debug().Msg("bank_create_bank_account: create bank account sent request")

	// if config.Environment() == "dev" {
	// 	body, _ := ioutil.ReadAll(response.Body)
	// 	log.Debug().Msg(string(body))
	// }

	var apiResponseData apiResponse
	err = json.NewDecoder(response.Body).Decode(&apiResponseData)
	if err != nil {
		log.Debug().Msg("bank_create_bank_account: failed to decode response")
		return "", errors.Wrap(err, "bank_create_bank_account: failed to decode response")
	}

	log.Debug().Interface("apiResponseData", apiResponseData).Msg("bank_create_bank_account: api response data")

	if !apiResponseData.Success {
		if apiResponseData.ErrCode == "1025" {
			return "", errors.Wrap(ErrHasNoDataPermission, apiResponseData.ErrMsg)
		}
		if apiResponseData.ErrCode == "4874" {
			return "", ErrUserNotFound
		}
	}

	return apiResponseData.Data.AccountNo, nil
}
