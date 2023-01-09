package sHTTP

import (
	"context"

	"gitlab.com/soteapps/packages/v2023/sCustom"
	"gitlab.com/soteapps/packages/v2023/sError"
	"gitlab.com/soteapps/packages/v2023/sLogger"
)

// Response represents the final response output
type Response struct {
	MessageId string      `json:"message-id"`
	Message   interface{} `json:"message,omitempty"`
	Error     interface{} `json:"error,omitempty"`
}

// ProcessLeafResponse unmarshal the leaf service response to defined struct and Sote Error
func ProcessLeafResponse(ctx context.Context, message []byte, respPtr interface{}, testMode bool) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		tempRespByte    []byte
		leafResponsePtr = &Response{}
	)

	if soteErr = sCustom.JSONUnmarshal(ctx, message, leafResponsePtr); soteErr.ErrCode != nil {
		soteErr = sError.ConvertError(soteErr, testMode)
		return
	}

	if leafResponsePtr.Error != nil {
		if tempRespByte, soteErr = sCustom.JSONMarshal(leafResponsePtr.Error); soteErr.ErrCode != nil {
			return
		}

		if tSoteErr := sCustom.JSONUnmarshal(ctx, tempRespByte, &soteErr); tSoteErr.ErrCode != nil {
			soteErr = sError.ConvertError(tSoteErr, testMode)
			return
		}

		return
	}

	if soteErr.ErrCode == nil && respPtr != nil {
		if tempRespByte, soteErr = sCustom.JSONMarshal(leafResponsePtr.Message); soteErr.ErrCode != nil {
			return
		}

		if soteErr = sCustom.JSONUnmarshal(ctx, tempRespByte, respPtr); soteErr.ErrCode != nil {
			soteErr = sError.ConvertError(soteErr, testMode)
			return
		}
	}

	return
}
