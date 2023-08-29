package main

import "fmt"

const (
	SUCCESSFUL_RESULT_CODE = 0
)

type BizResponse struct {
	ResultCode int                    `json:"resultCode"`
	ErrorCode  int                    `json:"errorCode,omitempty"`
	ErrorMsg   string                 `json:"errorMsg,omitempty"`
	BizResp    map[string]interface{} `json:"bizResp,omitempty"`
}

func SetupErrBizResp(res *BizResponse, err BizError) {
	res.ResultCode = -1
	res.ErrorCode = err.ErrorCode
	res.ErrorMsg = err.ErrorMsg
}

func SetupSuccess(res *BizResponse, bizResp map[string]interface{}) {
	res.ResultCode = SUCCESSFUL_RESULT_CODE
	res.BizResp = bizResp
}

type BizError struct {
	ErrorCode int
	ErrorMsg  string
}

func (e BizError) Error(err error) BizError {
	e.ErrorMsg = fmt.Sprintf("%s : %s", e.ErrorMsg, err.Error())
	return e
}
