package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	err "github.com/ilhamsyahids/owldetect/helpers/errors"
)

type ApiResp struct {
	StatusCode int         `json:"-"`
	OK         bool        `json:"ok"`
	Data       interface{} `json:"data,omitempty"`
	ErrCode    string      `json:"err,omitempty"`
	Message    string      `json:"msg,omitempty"`
	Timestamp  int64       `json:"ts"`
}

func NewSuccessResp(data interface{}) ApiResp {
	return ApiResp{
		StatusCode: http.StatusOK,
		OK:         true,
		Data:       data,
		Timestamp:  time.Now().Unix(),
	}
}

func NewErrorResp(er error) ApiResp {
	var e *err.Error
	if !errors.As(er, &e) {
		e = err.NewErrInternalError(er)
	}
	return ApiResp{
		StatusCode: e.StatusCode,
		OK:         false,
		ErrCode:    e.ErrCode,
		Message:    e.Message,
		Timestamp:  time.Now().Unix(),
	}
}

func WriteAPIResp(w http.ResponseWriter, resp ApiResp) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)

	b, _ := json.Marshal(resp)
	w.Write(b)
}

type AnalyzeReqBody struct {
	InputText string `json:"input_text"`
	RefText   string `json:"ref_text"`
}

func (rb AnalyzeReqBody) Validate() error {
	if len(rb.InputText) == 0 {
		return err.NewErrBadRequest("missing `input_text`")
	}
	if len(rb.RefText) == 0 {
		return err.NewErrBadRequest("missing `ref_text`")
	}
	if len(rb.InputText) > len(rb.RefText) {
		return err.NewErrBadRequest("`ref_text` must be longer than `input_text`")
	}
	return nil
}
