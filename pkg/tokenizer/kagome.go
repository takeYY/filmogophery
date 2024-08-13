package tokenizer

import (
	"bufio"
	"io"
	"sync"

	neologd "github.com/ikawaha/kagome-dict-ipa-neologd"
	"github.com/ikawaha/kagome/v2/filter"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

var Tokenizer *tokenizer.Tokenizer

func NewTokenizer() error {
	var err error
	if Tokenizer, err = tokenizer.New(neologd.Dict(), tokenizer.OmitBosEos()); err != nil {
		return err
	}

	return nil
}

func SyncTokenize(ch chan<- *NEologd, r io.Reader) {
	posFilter := filter.NewPOSFilter([]filter.POS{
		{"名詞"},
		{"形容詞"},
		{"動詞"},
	}...)

	scanner := bufio.NewScanner(r)
	scanner.Split(filter.ScanSentences)

	var wg sync.WaitGroup
	for scanner.Scan() {
		wg.Add(1)
		go func(s string) {
			defer wg.Done()

			tokens := Tokenizer.Analyze(s, tokenizer.Search)
			posFilter.Keep(&tokens)
			for _, t := range tokens {
				kagome := NewNEologd(&t)
				ch <- kagome
			}
		}(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		close(ch)
	}
	wg.Wait()
	close(ch)
}
