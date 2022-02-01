package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_ApiResp_Create(t *testing.T) {
	now := time.Now().Unix()
	api := ApiResp{
		StatusCode: http.StatusOK,
		OK:         true,
		Data:       "data",
		Timestamp:  now,
	}

	assert.True(t, api.OK)
	assert.Equal(t, http.StatusOK, api.StatusCode)
	assert.Equal(t, "data", api.Data)
	assert.Equal(t, now, api.Timestamp)
}

func Test_NewSuccessResp(t *testing.T) {
	data := map[string]interface{}{
		"matches": []string{"match1", "match2"},
	}
	api := NewSuccessResp(data)

	assert.Equal(t, api.Data, data)
}

func Test_NewErrorResp_Missing_Input(t *testing.T) {
	rb := AnalyzeReqBody{
		InputText: "",
		RefText:   "",
	}

	err := rb.Validate()
	expectedErrorMsg := "missing `input_text`"

	api := NewErrorResp(err)

	assert.Equal(t, api.StatusCode, http.StatusBadRequest)
	assert.Equal(t, api.OK, false)
	assert.Equal(t, api.ErrCode, "ERR_BAD_REQUEST")
	assert.Equal(t, api.Message, expectedErrorMsg)
}

func Test_NewErrorResp_Missing_Ref(t *testing.T) {
	rb := AnalyzeReqBody{
		InputText: "Halo",
		RefText:   "",
	}

	err := rb.Validate()
	expectedErrorMsg := "missing `ref_text`"

	api := NewErrorResp(err)

	assert.Equal(t, api.StatusCode, http.StatusBadRequest)
	assert.Equal(t, api.OK, false)
	assert.Equal(t, api.ErrCode, "ERR_BAD_REQUEST")
	assert.Equal(t, api.Message, expectedErrorMsg)
}

func Test_AnalyzeReqBody_Input_Greater_Than_Ref(t *testing.T) {
	rb := AnalyzeReqBody{
		InputText: "Halo",
		RefText:   "H",
	}

	err := rb.Validate()
	expectedErrorMsg := "`ref_text` must be longer than `input_text`"

	api := NewErrorResp(err)

	assert.Equal(t, api.StatusCode, http.StatusBadRequest)
	assert.Equal(t, api.OK, false)
	assert.Equal(t, api.ErrCode, "ERR_BAD_REQUEST")
	assert.Equal(t, api.Message, expectedErrorMsg)
}

func Test_AnalyzeReqBody_Succss(t *testing.T) {
	rb := AnalyzeReqBody{
		InputText: "Halo",
		RefText:   "Haloo",
	}

	err := rb.Validate()

	assert.Equal(t, err, nil)
}

func Test_WriteAPIResp_Success(t *testing.T) {
	api := NewSuccessResp(map[string]interface{}{
		"matches": []string{"match1", "match2"},
	})

	now := time.Now().Unix()

	respw := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		WriteAPIResp(w, api)
	})

	handler.ServeHTTP(respw, nil)

	assert.Equal(t, respw.Header().Get("Content-Type"), "application/json")
	assert.Equal(t, respw.Code, http.StatusOK)
	assert.Equal(t, respw.Body.String(), `{"ok":true,"data":{"matches":["match1","match2"]},"ts":`+fmt.Sprintf("%d", now)+`}`)
}

func Test_WriteAPIResp_Error(t *testing.T) {
	err := fmt.Errorf("some error")
	api := NewErrorResp(err)

	now := time.Now().Unix()

	respw := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		WriteAPIResp(w, api)
	})

	handler.ServeHTTP(respw, nil)

	assert.Equal(t, respw.Header().Get("Content-Type"), "application/json")
	assert.Equal(t, respw.Code, http.StatusInternalServerError)
	assert.Equal(t, respw.Body.String(), `{"ok":false,"err":"ERR_INTERNAL_ERROR","msg":"some error","ts":`+fmt.Sprintf("%d", now)+`}`)
}
