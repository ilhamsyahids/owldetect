package plagiarism

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MatchDetails_To_JSON(t *testing.T) {
	var md matchDetails = matchDetails{
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
	inp := matchDetails{
		Text:     `hello`,
		StartIdx: 0,
		EndIdx:   4,
	}
	ref := matchDetails{
		Text:     `helloo`,
		StartIdx: 5,
		EndIdx:   10,
	}

	expected_inp := `{"text":"hello","start_idx":0,"end_idx":4}`
	expected_ref := `{"text":"helloo","start_idx":5,"end_idx":10}`

	m := match{
		Input:     inp,
		Reference: ref,
	}

	expected := `{"input":` + expected_inp + `,"ref":` + expected_ref + `}`

	b, err := json.Marshal(m)

	assert.Nil(t, err)
	assert.Equal(t, expected, string(b))
}
