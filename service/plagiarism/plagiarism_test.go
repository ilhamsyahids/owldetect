package plagiarism

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MatchDetails_To_JSON(t *testing.T) {
	var md MatchDetails = MatchDetails{
		Text:     `text`,
		StartIdx: 0,
		EndIdx:   3,
	}

	expected := `{"text":"text","start_idx":0,"end_idx":3}`

	b, err := json.Marshal(md)

	assert.Nil(t, err)
	assert.Equal(t, expected, string(b))
}

func Test_Match_To_JSON(t *testing.T) {
	inp := MatchDetails{
		Text:     `hello`,
		StartIdx: 0,
		EndIdx:   4,
	}
	ref := MatchDetails{
		Text:     `helloo`,
		StartIdx: 5,
		EndIdx:   10,
	}

	expected_inp := `{"text":"hello","start_idx":0,"end_idx":4}`
	expected_ref := `{"text":"helloo","start_idx":5,"end_idx":10}`

	m := Match{
		Input:     inp,
		Reference: ref,
	}

	expected := `{"input":` + expected_inp + `,"ref":` + expected_ref + `}`

	b, err := json.Marshal(m)

	assert.Nil(t, err)
	assert.Equal(t, expected, string(b))
}

func Test_Tokenize_Sentence_1(t *testing.T) {
	text := `Why upon your first voyage as a passenger, did you yourself feel such a mystical vibration, when first told that you and your ship were now out of sight of land? Why did the old Persians hold the sea holy? But look!`

	expected := []sentenceToken{
		{
			Text:  `Why upon your first voyage as a passenger, did you yourself feel such a mystical vibration, when first told that you and your ship were now out of sight of land?`,
			Start: 0,
			End:   160,
		},
		{
			Text:  `Why did the old Persians hold the sea holy?`,
			Start: 162,
			End:   204,
		},
		{
			Text:  `But look!`,
			Start: 206,
			End:   214,
		},
	}

	tokens := tokenizeToSentence(text)

	assert.Equal(t, expected, tokens)
}

func Test_Tokenize_Sentence_2(t *testing.T) {
	text := `hello. world. Am I a sentence?     Yes!`

	expected := []sentenceToken{
		{
			Text:  `hello.`,
			Start: 0,
			End:   5,
		},
		{
			Text:  `world.`,
			Start: 7,
			End:   12,
		},
		{
			Text:  `Am I a sentence?`,
			Start: 14,
			End:   29,
		},
		{
			Text:  `Yes!`,
			Start: 35,
			End:   38,
		},
	}

	tokens := tokenizeToSentence(text)

	assert.Equal(t, expected, tokens)
}

func Test_Tokenize_Word(t *testing.T) {
	text := `hello world, I am a sentence     !`

	expected := []string{
		`hello`,
		`world`,
		`I`,
		`am`,
		`a`,
		`sentence`,
	}

	tokens := tokenizeToWord(text)

	assert.Equal(t, expected, tokens)
}

func Test_Check_Plagiarism_Sentence_All_Same(t *testing.T) {
	text := `hello world, I am a sentence     !`
	ref := `hello world, I am a sentence     !`

	result := checkPlagiarismSentence(text, ref)

	assert.True(t, result)
}

func Test_Check_Plagiarism_Sentence_Half_Same(t *testing.T) {
	text := `hello world, I am a sentence     !`
	ref := `hello word, I m sentence!`

	result := checkPlagiarismSentence(text, ref)

	assert.True(t, result)
}

func Test_Check_Plagiarism_Sentence_Below_Threshold_Same(t *testing.T) {
	text := `hello world, I am a sentence     !`
	ref := `Hlo word, I m sentence!`

	result := checkPlagiarismSentence(text, ref)

	assert.False(t, result)
}

func Test_Check_Plagiarism_Sentence_Not_Same(t *testing.T) {
	text := `hello world, I am a sentence     !`
	ref := `Good BYE!`

	result := checkPlagiarismSentence(text, ref)

	assert.False(t, result)
}
