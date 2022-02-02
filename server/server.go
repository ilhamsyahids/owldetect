package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/ilhamsyahids/owldetect/helpers/api"
	"github.com/ilhamsyahids/owldetect/helpers/errors"
	"github.com/ilhamsyahids/owldetect/service/plagiarism"
)

func ServeServer(port string) error {
	log.Printf("server is listening on :%v", port)
	return http.ListenAndServe(":"+port, nil)
}

func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9056"
	}

	return port
}

func HomeRoot() http.Handler {
	return http.FileServer(http.Dir("./static"))
}

func AnalysisRoot(w http.ResponseWriter, r *http.Request) {
	// check http method
	if r.Method != http.MethodPost {
		api.WriteAPIResp(w, api.NewErrorResp(errors.NewErrMethodNotAllowed()))
		return
	}

	// parse request body
	var reqBody api.AnalyzeReqBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		api.WriteAPIResp(w, api.NewErrorResp(errors.NewErrBadRequest(err.Error())))
		return
	}

	// validate request body
	err = reqBody.Validate()
	if err != nil {
		api.WriteAPIResp(w, api.NewErrorResp(err))
		return
	}

	// do analysis
	matches := plagiarism.Analyze(reqBody.InputText, reqBody.RefText)

	// output success response
	api.WriteAPIResp(w, api.NewSuccessResp(map[string]interface{}{
		"matches": matches,
	}))
}

func InitHandler() {
	http.Handle("/", HomeRoot())
	http.HandleFunc("/analysis", AnalysisRoot)
}
