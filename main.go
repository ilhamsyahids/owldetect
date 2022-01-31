package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	re "regexp"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

var indexMatches map[string]bool = map[string]bool{}

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
	tokens := strings.Split(text, " ")
	res := []string{}
	for _, word := range tokens {
		if word == "" || word == " " {
			continue
		}

		res = append(res, word)
	}

	return res
}

func tokenizeToSentence(text string) []sentenceToken {
	// seperate sentence by dot, question mark, exclamation mark
	reg := re.MustCompile(`[.!?]`)
	sentences := reg.Split(text, -1)

	tokens := []sentenceToken{}
	start := 0
	end := 0

	for _, sentence := range sentences {
		if len(sentence) == 0 {
			continue
		}

		end = start + len(sentence)

		tokens = append(tokens, sentenceToken{
			Text:  sentence,
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
			keys := inputToken + " " + refTokens[i]

			val, exists := indexMatches[keys]

			if exists {
				if val {
					numMatch++
					start = i
					break
				}
			} else {
				if fuzzy.MatchNormalizedFold(inputToken, refTokens[i]) {
					indexMatches[keys] = true
					numMatch++
					start = i
					break
				}
				indexMatches[keys] = false
			}

			i++
		}
	}

	return (float32(2*numMatch) / float32(len(inputTokens)+len(refTokens))) >= float32(0.5)
}

func mergeIntervals(intervals []match) []match {
	// merge overlapping intervals index
	res := []match{}

	for idx, interval := range intervals {
		if idx == 0 {
			res = append(res, interval)
			continue
		}

		prev := res[len(res)-1]

		if prev.Input.EndIdx < interval.Input.StartIdx {
			res = append(res, interval)
		} else {
			if interval.Input.EndIdx > prev.Input.EndIdx {
				prev.Input.EndIdx = interval.Input.EndIdx
				prev.Reference.EndIdx = interval.Reference.EndIdx
			}
		}
	}

	return res
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
			if isPlagiatSentence(strings.ToLower(inputToken.Text), strings.ToLower(refToken.Text)) {
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

				// add same character at the end
				for (inputToken.End < len(input)) && (refToken.End < len(ref)) && (string(input[inputToken.End]) == string(ref[refToken.End])) {
					val := string(input[inputToken.End])
					inputToken.Text += val
					inputToken.End++
					refToken.Text += val
					refToken.End++
				}

				inputIdx := inputToken.Start
				refIdx := refToken.Start

				// add match to group
				groupMatch = append(groupMatch, match{
					Input: matchDetails{
						Text:     inputToken.Text,
						StartIdx: inputIdx,
						EndIdx:   inputIdx + len(inputToken.Text) - 1,
					},
					Reference: matchDetails{
						Text:     refToken.Text,
						StartIdx: refIdx,
						EndIdx:   refIdx + len(refToken.Text) - 1,
					},
				})
			}
		}
	}

	return mergeIntervals(groupMatch)
}
