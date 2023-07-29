package backend

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/BlahajXD/backend/config"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type BankAccountInfoOutput struct {
	AccountNumber string  `json:"account_number"`
	Balance       float64 `json:"balance"`
}

type BankAccountInfoInput struct {
	AccessToken   string `json:"access_token"`
	AccountNumber string `json:"account_number"`
}

func (d *Dependency) BankAccountInfo(ctx context.Context, input BankAccountInfoInput) (BankAccountInfoOutput, error) {
	payload := map[string]any{
		"accountNo": input.AccountNumber,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return BankAccountInfoOutput{}, errors.Wrap(err, "failed to marshal payload")
	}

	log.Debug().Interface("payload", payload).Msg("bank_account_info: bank account info payload")

	type apiResponse struct {
		TraceId string `json:"traceId"`
		Data    struct {
			ID          int     `json:"id"`
			UID         int     `json:"uid"`
			Balance     float64 `json:"balance"`
			AccountName string  `json:"accountName"`
			CreateTime  int     `json:"createTime"`
			AccountNo   string  `json:"accountNo"`
			UpdateTime  int     `json:"updateTime"`
			Status      string  `json:"status"`
		} `json:"data"`
		ErrCode string `json:"errCode"`
		Success bool   `json:"success"`
		ErrMsg  string `json:"errMsg"`
	}

	header := http.Header{}
	header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	header.Add(fiber.HeaderAuthorization, "Bearer "+input.AccessToken)

	response, err := d.httpclient.Post(config.BankApiURL()+"/bankAccount/info", bytes.NewReader(jsonPayload), header)
	if err != nil {
		return BankAccountInfoOutput{}, errors.Wrap(err, "bank_account_info: failed to send request")
	}

	// if config.Environment() == "dev" {
	// 	body, _ := ioutil.ReadAll(response.Body)
	// 	log.Debug().Msg(string(body))
	// }
	var apiResponseData apiResponse
	err = json.NewDecoder(response.Body).Decode(&apiResponseData)
	if err != nil {
		log.Debug().Err(err).Msg("bank_account_info: failed to decode response")
		return BankAccountInfoOutput{}, errors.Wrap(err, "bank_account_info: failed to decode response")
	}

	log.Debug().Interface("apiResponseData", apiResponseData).Msg("bank_account_create_transaction: api response data")

	if !apiResponseData.Success {
		if apiResponseData.ErrCode == "4002" {
			return BankAccountInfoOutput{}, ErrNotFoundBankAccount
		}

		return BankAccountInfoOutput{}, errors.Wrap(ErrBankUnknownError, apiResponseData.ErrMsg)
	}

	var output BankAccountInfoOutput

	output.AccountNumber = apiResponseData.Data.AccountNo
	output.Balance = apiResponseData.Data.Balance

	return output, nil
}
