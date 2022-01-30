package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	re "regexp"
	"strings"
)

func main() {
	// define handlers
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/analysis", func(w http.ResponseWriter, r *http.Request) {
		// check http method
		if r.Method != http.MethodPost {
			WriteAPIResp(w, NewErrorResp(NewErrMethodNotAllowed()))
			return
		}

		// parse request body
		var reqBody analyzeReqBody
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			WriteAPIResp(w, NewErrorResp(NewErrBadRequest(err.Error())))
			return
		}

		// validate request body
		err = reqBody.Validate()
		if err != nil {
			WriteAPIResp(w, NewErrorResp(err))
			return
		}

		// do analysis
		matches := doAnalysis(reqBody.InputText, reqBody.RefText)

		// output success response
		WriteAPIResp(w, NewSuccessResp(map[string]interface{}{
			"matches": matches,
		}))
	})

	// define port, we need to set it as env for Heroku deployment
	port := os.Getenv("PORT")
	if port == "" {
		port = "9056"
	}

	// run server
	log.Printf("server is listening on :%v", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("unable to run server due: %v", err)
	}
}

func tokenizeToWord(text string) []string {
	// seperate word by space
	return strings.Split(text, " ")
}

func tokenizeToSentence(text string) []sentenceToken {
	// seperate sentence by dot, question mark, exclamation mark
	reg := re.MustCompile(`[.!?]`)
	sentence := reg.Split(text, -1)

	tokens := []sentenceToken{}
	start := 0
	end := 0

	for _, word := range sentence {
		if len(word) == 0 {
			continue
		}
		end = start + len(word)
		tokens = append(tokens, sentenceToken{
			Text:  word,
			Start: start,
			End:   end,
		})
		start = end + 1
	}
	return tokens
}

func isPlagiatSentence(input, ref string) bool {
	// count number of words matches
	inputTokens := tokenizeToWord(input)
	refTokens := tokenizeToWord(ref)

	nRef := len(refTokens)
	numMatch := 0
	start := 0

	for _, inputToken := range inputTokens {
		i := start

		for i < nRef {
			if refTokens[i] == inputToken {
				numMatch++
				start = i
				break
			}
			i++
		}
	}

	return (float32(2*numMatch) / float32(len(inputTokens)+len(refTokens))) >= float32(0.5)
}

func doAnalysis(input, ref string) []match {

	// tokenize input and ref
	inputTokens := tokenizeToSentence(input)
	refTokens := tokenizeToSentence(ref)

	groupMatch := []match{}

	// flags for sentence that already matched
	flags := ItemSet{}

	for _, inputToken := range inputTokens {
		for idx, refToken := range refTokens {
			// already matched before
			if flags.Has(idx) {
				continue
			}

			// check if sentence is plagiat
			if isPlagiatSentence(inputToken.Text, refToken.Text) {
				flags.Add(idx)

				// remove front whitespace
				for inputToken.Text[0] == ' ' {
					inputToken.Text = inputToken.Text[1:]
					inputToken.Start++
				}

				for refToken.Text[0] == ' ' {
					refToken.Text = refToken.Text[1:]
					refToken.Start++
				}

				// add match to group
				groupMatch = append(groupMatch, match{
					Input: matchDetails{
						Text:     inputToken.Text,
						StartIdx: inputToken.Start,
						EndIdx:   inputToken.End,
					},
					Reference: matchDetails{
						Text:     refToken.Text,
						StartIdx: refToken.Start,
						EndIdx:   refToken.End,
					},
				})
			}
		}
	}

	return groupMatch
}
