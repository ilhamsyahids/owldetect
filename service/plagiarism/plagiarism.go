package plagiarism

import (
	"strings"

	"regexp"

	"github.com/ilhamsyahids/owldetect/helpers/set"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type match struct {
	Input     matchDetails `json:"input"`
	Reference matchDetails `json:"ref"`
}

type matchDetails struct {
	Text     string `json:"text"`
	StartIdx int    `json:"start_idx"`
	EndIdx   int    `json:"end_idx"`
}

type sentenceToken struct {
	Text  string
	Start int
	End   int
}

var indexMatches map[string]bool = map[string]bool{}

func tokenizeToSentence(text string) []sentenceToken {
	// seperate sentence by dot, question mark, exclamation mark
	reg := regexp.MustCompile(`[.!?]`)
	sentences := reg.Split(text, -1)

	tokens := []sentenceToken{}
	start := 0
	end := 0

	for _, sentence := range sentences {
		if len(sentence) == 0 {
			continue
		}

		// remove front whitespace
		i := 0
		for {
			if sentence[0] == ' ' {
				sentence = sentence[1:]
			} else {
				break
			}
		}

		// add back sentence
		end = start + len(sentence)
		for end < len(text) {
			if text[end] == '!' || text[end] == '.' || text[end] == '?' {
				sentence = text[start : end+1]
				end++
			} else {
				break
			}

			i++
		}

		tokens = append(tokens, sentenceToken{
			Text:  sentence,
			Start: start,
			End:   end - 1,
		})
		start = end + 1

		for start < len(text) {
			if text[start] == ' ' {
				start++
			} else {
				break
			}

			i++
		}
	}
	return tokens
}

func tokenizeToWord(text string) []string {
	// seperate word by space
	reg := regexp.MustCompile(`[^\w+]`)
	tokens := reg.Split(text, -1)
	res := []string{}
	for _, word := range tokens {
		if word == "" || word == " " {
			continue
		}

		res = append(res, word)
	}

	return res
}

func checkPlagiarismSentence(input, ref string) bool {
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

func Analyze(input, ref string) []match {
	// remove front and back whitespace
	input = strings.Trim(input, " ")
	ref = strings.Trim(ref, " ")

	// tokenize input and ref
	inputTokens := tokenizeToSentence(input)
	refTokens := tokenizeToSentence(ref)

	groupMatch := []match{}

	// flags for sentence that already matched
	flags := set.ItemSet{}

	for _, inputToken := range inputTokens {
		for idx, refToken := range refTokens {
			// already matched before
			if flags.Has(idx) {
				continue
			}

			// check if sentence is plagiat
			if checkPlagiarismSentence(strings.ToLower(inputToken.Text), strings.ToLower(refToken.Text)) {
				flags.Add(idx)

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
