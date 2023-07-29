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

type BankAccountCreateTransactionInput struct {
	FromAccountNumber string `json:"from_account_number"`
	ToAccountNumber   string `json:"to_account_number"`
	Amount            int    `json:"amount"`
	AccessToken       string
}

func (d *Dependency) BankAccountCreateTransaction(ctx context.Context, input BankAccountCreateTransactionInput) error {
	payload := map[string]any{
		"senderAccountNo":   input.FromAccountNumber,
		"receiverAccountNo": input.ToAccountNumber,
		"amount":            input.Amount,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return errors.Wrap(err, "failed to marshal payload")
	}

	log.Debug().Interface("payload", payload).Msg("bank_account_create_transaction: create bank account payload")

	type apiResponse struct {
		TraceId string `json:"traceId"`
		Data    struct {
			UID               int    `json:"uid"`
			Amount            int    `json:"amount"`
			SenderAccountNo   string `json:"senderAccountNo"`
			TraxType          string `json:"traxType"`
			ReceiverAccountNo string `json:"receiverAccountNo"`
			TransactionDate   int    `json:"transactionDate"`
		} `json:"data"`
		ErrCode string `json:"errCode"`
		Success bool   `json:"success"`
		ErrMsg  string `json:"errMsg"`
	}

	header := http.Header{}
	header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	header.Add(fiber.HeaderAuthorization, "Bearer "+input.AccessToken)

	response, err := d.httpclient.Post(config.BankApiURL()+"/bankAccount/transaction/create", bytes.NewReader(jsonPayload), header)
	if err != nil {
		return errors.Wrap(err, "bank_account_create_transaction: failed to send request")
	}

	// if config.Environment() == "dev" {
	// 	body, _ := ioutil.ReadAll(response.Body)
	// 	log.Debug().Msg(string(body))
	// }
	var apiResponseData apiResponse
	err = json.NewDecoder(response.Body).Decode(&apiResponseData)
	if err != nil {
		log.Debug().Err(err).Msg("bank_account_create_transaction: failed to decode response")
		return errors.Wrap(err, "bank_account_create_transaction: failed to decode response")
	}

	log.Debug().Interface("apiResponseData", apiResponseData).Msg("bank_account_create_transaction: api response data")

	if !apiResponseData.Success {
		if apiResponseData.ErrCode == "1020" {
			return errors.Wrap(ErrBankRequestParameter, apiResponseData.ErrMsg)
		}
		if apiResponseData.ErrCode == "4874" {
			return ErrUserNotFound
		}
		if apiResponseData.ErrCode == "4002" {
			return ErrNotFoundBankAccount
		}
		if apiResponseData.ErrCode == "3005" {
			return ErrBankAmountNotEnough
		}

		return errors.Wrap(ErrBankUnknownError, apiResponseData.ErrMsg)
	}

	return nil
}
