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

type BankAddBalanceInput struct {
	ParentID    int `json:"parent_id"`
	Balance     int `json:"balance"`
	AccessToken string
}

func (d *Dependency) BankAddBalance(ctx context.Context, input BankAddBalanceInput) error {
	parent, err := d.repo.FindParent(ctx, "id", input.ParentID)
	if err != nil {
		return errors.Wrap(err, "bank_add_balance: failed to find parent")
	}

	payload := map[string]any{
		"amount":            input.Balance,
		"receiverAccountNo": parent.AccountNumber,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return errors.Wrap(err, "bank_add_balance: failed to marshal payload")
	}

	log.Debug().Interface("payload", input).Msg("bank_add_balance: create bank account payload")

	type apiResponse struct {
		TraceId string `json:"traceId"`
		Data    struct {
			UID               int    `json:"uid"`
			Amount            int    `json:"amount"`
			CreateTime        int    `json:"createTime"`
			TraxType          string `json:"traxType"`
			ReceiverAccountNo string `json:"receiverAccountNo"`
		} `json:"data"`
		ErrCode string `json:"errCode"`
		Success bool   `json:"success"`
		ErrMsg  string `json:"errMsg"`
	}

	header := http.Header{}
	header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	header.Add(fiber.HeaderAuthorization, "Bearer "+input.AccessToken)

	response, err := d.httpclient.Post(config.BankApiURL()+"/bankAccount/addBalance", bytes.NewReader(jsonPayload), header)
	if err != nil {
		return errors.Wrap(err, "bank_add_balance: failed to send request")
	}

	var apiResponseData apiResponse
	err = json.NewDecoder(response.Body).Decode(&apiResponseData)
	if err != nil {
		log.Debug().Err(err).Msg("bank_add_balance: failed to decode response")
		return errors.Wrap(err, "bank_add_balance: failed to decode response")
	}

	log.Debug().Interface("apiResponseData", apiResponseData).Msg("bank_create_bank_account: api response data")

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
	}

	return nil
}
