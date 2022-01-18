package nlp

import (
	"io"
	"os"
	"strings"
	"testing"
)

var tests = []struct {
	name        string
	strs        []string
	sortedWords []string
}{
	{"EnglishWordCounter",
		[]string{"Hello, world!", "hello, go.", "hello, foo. Let's go!"},
		[]string{"hello", "go"},
	},
}

func TestEnglishWordCounterPQ(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			counter := NewEnglishWordCounterPQ()
			for _, s := range tt.strs {
				counter.AddSentence(s)
			}

			for k, v := range counter.(*EnglishWordsCounterPQ).Map {
				t.Logf("item %v: %#v", k, *v)
			}

			for _, r := range tt.sortedWords {
				c := counter.PopMostCommon()
				if c != r {
					t.Errorf("❌ pop %v (want %v)", c, r)
				}
			}
		})
	}
}

func TestEnglishWordCounterQS(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			counter := NewEnglishWordCounterQS()
			for _, s := range tt.strs {
				counter.AddSentence(s)
			}

			for k, v := range counter.(*EnglishWordsCounterQS).List {
				t.Logf("item %v: %#v", k, *v)
			}

			for _, r := range tt.sortedWords {
				c := counter.PopMostCommon()
				if c != r {
					t.Errorf("❌ pop %v (want %v)", c, r)
				}
			}
		})
	}
}

const Document = "/Users/c/Desktop/nietzsche.txt"

func readFile(path string) string {
	f, err := os.Open(path)
	if err != nil {
		panic("readFile failed: " + err.Error())
	}
	t, err := io.ReadAll(f)
	if err != nil {
		panic("readFile failed: " + err.Error())
	}
	return string(t)
}

func splitSentences(s string) []string {
	return strings.Split(s, ".")
}

func benchmarkEnglishWordCounter(b *testing.B, counter WordsCounter) {
	ss := splitSentences(readFile(Document))
	res := make([]string, 10)

	for n := 0; n < b.N; n++ {
		for _, s := range ss {
			counter.AddSentence(s)
		}
		for i := 0; i < 10; i++ {
			res[i] = counter.PopMostCommon()
		}
	}

	b.Logf("most common: %v", res)
}

func BenchmarkEnglishWordCounterPQ(b *testing.B) {
	counter := NewEnglishWordCounterPQ()
	benchmarkEnglishWordCounter(b, counter)
}

func BenchmarkEnglishWordCounterQS(b *testing.B) {
	counter := NewEnglishWordCounterQS()
	benchmarkEnglishWordCounter(b, counter)
}
