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

type BankTransactionInfoInput struct {
	AccountNumber  string `json:"accountNumber"`
	Page           int    `json:"page"`
	RecordsPerPage int    `json:"recordsPerPage"`
	AccessToken    string `json:"accessToken"`
}

type BankTransactionInfoItem struct {
	UID               int     `json:"uid"`
	Amount            float64 `json:"amount"`
	AccountNo         string  `json:"accountNo"`
	ReceiverAccountNo string  `json:"receiverAccountNo"`
	TransactionType   string  `json:"transactionType"`
}

type BankTransactionInfoOutput struct {
	TotalItems int                       `json:"total_items"`
	Items      []BankTransactionInfoItem `json:"items"`
}

func (d *Dependency) BankTransactionInfo(input BankTransactionInfoInput) (BankTransactionInfoOutput, error) {
	payload := map[string]any{
		"accountNo":      input.AccountNumber,
		"traxType":       []string{"TRANSFER_IN", "TRANSFER_OUT"},
		"pageNumber":     input.Page,
		"recordsPerPage": input.RecordsPerPage,
	}

	log.Debug().Interface("payload", payload).Msg("bank_transaction_info: generate token payload")

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return BankTransactionInfoOutput{}, errors.Wrap(err, "failed to marshal payload")
	}

	type apiResponse struct {
		TraceId string `json:"traceId"`
		Data    struct {
			Transactions []struct {
				UID               int     `json:"uid"`
				Amount            float64 `json:"amount"`
				TransactionDate   int     `json:"transactionDate"`
				TraxID            int     `json:"traxId"`
				SenderAccountNo   string  `json:"senderAccountNo"`
				TraxType          string  `json:"traxType"`
				ReceiverAccountNo string  `json:"receiverAccountNo"`
			} `json:"transactions"`
		} `json:"data"`
		ErrCode string `json:"errCode"`
		Success bool   `json:"success"`
		ErrMsg  string `json:"errMsg"`
	}

	header := http.Header{}
	header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	header.Add(fiber.HeaderAuthorization, "Bearer "+input.AccessToken)

	response, err := d.httpclient.Post(config.BankApiURL()+"/bankAccount/transaction/info", bytes.NewReader(jsonPayload), header)
	if err != nil {
		return BankTransactionInfoOutput{}, errors.Wrap(err, "bank_transaction_info: failed to send request")
	}

	// if config.Environment() == "dev" {
	// 	body, _ := ioutil.ReadAll(response.Body)
	// 	log.Debug().Msg(string(body))
	// }
	var apiResponseData apiResponse
	err = json.NewDecoder(response.Body).Decode(&apiResponseData)
	if err != nil {
		return BankTransactionInfoOutput{}, errors.Wrap(err, "bank_transaction_info: failed to decode response")
	}

	log.Debug().Interface("apiResponseData", apiResponseData).Msg("bank_transaction_info: api response data")

	if !apiResponseData.Success {
		if apiResponseData.ErrCode == "4874" {
			return BankTransactionInfoOutput{}, ErrUserNotFound
		}
		if apiResponseData.ErrCode == "4503" {
			return BankTransactionInfoOutput{}, ErrBankTransactionNotExists
		}

		return BankTransactionInfoOutput{}, errors.Wrap(ErrBankUnknownError, apiResponseData.ErrMsg)
	}

	var output BankTransactionInfoOutput
	output.TotalItems = len(apiResponseData.Data.Transactions)
	output.Items = make([]BankTransactionInfoItem, output.TotalItems)

	for i, transaction := range apiResponseData.Data.Transactions {
		output.Items[i].UID = transaction.UID
		output.Items[i].Amount = transaction.Amount
		output.Items[i].AccountNo = transaction.SenderAccountNo
		output.Items[i].ReceiverAccountNo = transaction.ReceiverAccountNo
		output.Items[i].TransactionType = transaction.TraxType
	}

	return output, nil
}
