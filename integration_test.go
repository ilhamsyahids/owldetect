package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ilhamsyahids/owldetect/helpers/api"
	"github.com/ilhamsyahids/owldetect/server"
	"github.com/ilhamsyahids/owldetect/service/plagiarism"
	"github.com/stretchr/testify/suite"
)

type IntegrationTest struct {
	suite.Suite
}

func (suite *IntegrationTest) createTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(server.AnalysisRoot))
}

func (suite *IntegrationTest) TestAnalyzeApi() {
	server := suite.createTestServer()
	defer server.Close()

	var payload api.AnalyzeReqBody = api.AnalyzeReqBody{
		InputText: "Hello World! I'm a cat!",
		RefText:   "Welcome to my world! Hello World! I'm a ca! But I eat pizza!",
	}

	p, _ := json.Marshal(payload)
	res, err := http.Post(server.URL, "application/json", bytes.NewBuffer(p))
	if err != nil {
		suite.FailNow("Failed do request: " + err.Error())
		return
	}

	var successResponse api.ApiResp
	json.NewDecoder(res.Body).Decode(&successResponse)

	suite.Equal(200, res.StatusCode)
	suite.True(successResponse.OK)

	var matches []plagiarism.Match
	matchesRaw := successResponse.Data.(map[string]interface{})["matches"]
	sr, _ := json.Marshal(matchesRaw)
	json.Unmarshal(sr, &matches)

	expectedInput1 := plagiarism.MatchDetails{
		Text:     "Hello World!",
		StartIdx: 0,
		EndIdx:   11,
	}
	expectedRef1 := plagiarism.MatchDetails{
		Text:     "Hello World!",
		StartIdx: 21,
		EndIdx:   32,
	}
	expectedInput2 := plagiarism.MatchDetails{
		Text:     "I'm a cat!",
		StartIdx: 13,
		EndIdx:   22,
	}
	expectedRef2 := plagiarism.MatchDetails{
		Text:     "I'm a ca!",
		StartIdx: 34,
		EndIdx:   42,
	}

	// Test Match
	suite.Equal(len(matches), 2)

	suite.Equal(matches[0].Input, expectedInput1)
	suite.Equal(matches[0].Reference, expectedRef1)

	suite.Equal(matches[1].Input, expectedInput2)
	suite.Equal(matches[1].Reference, expectedRef2)
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTest))
}
