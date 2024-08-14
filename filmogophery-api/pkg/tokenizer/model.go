package tokenizer

import "github.com/ikawaha/kagome/v2/tokenizer"

// 辞書について
// Ref: https://zenn.dev/ikawaha/books/kagome-v2-japanese-tokenizer/viewer/dictionary
type NEologd struct {
	Surface          string `comment:"表層形"`
	POS0             string `comment:"品詞大分類"`
	POS1             string `comment:"品詞中分類"`
	POS2             string `comment:"品詞小分類"`
	POS3             string `comment:"品詞細分類"`
	InflectionalType string `comment:"活用型"`
	InflectionalForm string `comment:"活用形"`
	BaseForm         string `comment:"原形"`
	Reading          string `comment:"読み"`
	Pronunciation    string `comment:"発音"`
}

func NewNEologd(token *tokenizer.Token) *NEologd {
	var err = true
	var inflectionalType string
	var inflectionalForm string
	var baseForm string
	var reading string
	var pronunciation string
	pos := token.POS()
	pos0 := pos[0]
	pos1 := pos[1]
	pos2 := pos[2]
	pos3 := pos[3]

	surface := token.Surface
	if inflectionalType, err = token.InflectionalType(); !err {
		inflectionalType = ""
	}
	if inflectionalForm, err = token.InflectionalForm(); !err {
		inflectionalForm = ""
	}
	if baseForm, err = token.BaseForm(); !err {
		baseForm = ""
	}
	if reading, err = token.Reading(); !err {
		reading = ""
	}
	if pronunciation, err = token.Pronunciation(); !err {
		pronunciation = ""
	}

	return &NEologd{
		Surface:          surface,
		POS0:             pos0,
		POS1:             pos1,
		POS2:             pos2,
		POS3:             pos3,
		InflectionalType: inflectionalType,
		InflectionalForm: inflectionalForm,
		BaseForm:         baseForm,
		Reading:          reading,
		Pronunciation:    pronunciation,
	}
}
