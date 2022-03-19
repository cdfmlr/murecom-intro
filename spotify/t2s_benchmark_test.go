package spotify

import "testing"

// Best: use this in search.go
//func chineseStringEqualOnlyNecessaryCC(a, b string) bool {
//	if len(a) != len(b) {
//		return false
//	}
//	if a == b {
//		return true
//	}
//	// opencc: t2s
//	sa, _ := T2S(a)
//	if sa == b {
//		return true
//	}
//	sb, _ := T2S(b)
//	if sa == sb {
//		return true
//	}
//	return false
//}
var chineseStringEqualOnlyNecessaryCC = ChineseStringEqual

func chineseStringEqualBothCC(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	if a == b {
		return true
	}
	// opencc: t2s
	sa, _ := T2S(a)
	sb, _ := T2S(b)
	if sa == sb {
		return true
	}
	return false
}

func chineseStringEqualCoroutineCC(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	if a == b {
		return true
	}

	chA := make(chan string)
	chB := make(chan string)

	go func() {
		sa, _ := T2S(a)
		chA <- sa
	}()
	go func() {
		sb, _ := T2S(b)
		chB <- sb
	}()

	if <-chA == <-chB {
		return true
	}
	return false
}

func chineseStringEqualAlwaysCC(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	sa, _ := T2S(a)
	sb, _ := T2S(b)
	if sa == sb {
		return true
	}
	return false
}

func chineseStringEqualAlwaysCoroutineCC(a, b string) bool {
	if len(a) != len(b) {
		return false
	}

	chA := make(chan string)
	chB := make(chan string)

	go func() {
		sa, _ := T2S(a)
		chA <- sa
	}()
	go func() {
		sb, _ := T2S(b)
		chB <- sb
	}()

	if <-chA == <-chB {
		return true
	}
	return false
}

func TestChineseStringEqualAlgos(t *testing.T) {
	algos := map[string]func(string, string) bool{
		"OnlyNecessaryCC":   chineseStringEqualOnlyNecessaryCC,
		"BothCC":            chineseStringEqualBothCC,
		"CoroutineCC":       chineseStringEqualCoroutineCC,
		"AlwaysCC":          chineseStringEqualAlwaysCC,
		"AlwaysCoroutineCC": chineseStringEqualAlwaysCoroutineCC,
	}

	examples := []struct {
		wantEq bool
		T      string
		S      string
	}{
		{
			wantEq: true,
			T:      "自然語言處理是人工智能領域中的一個重要方向。",
			S:      "自然语言处理是人工智能领域中的一个重要方向。",
		},
		{
			wantEq: false,
			T:      "自然語言",
			S:      "自然狗言",
		},
	}

	type testcase struct {
		name   string
		A      string
		B      string
		wantEq bool
	}

	var testcases []testcase

	for _, e := range examples {
		k := "P"
		if e.wantEq == false {
			k = "N"
		}
		testcases = append(testcases,
			testcase{name: k + "TT", A: e.T, B: e.T, wantEq: true},
			testcase{name: k + "SS", A: e.S, B: e.S, wantEq: true},
			testcase{name: k + "TS", A: e.T, B: e.S, wantEq: e.wantEq},
			testcase{name: k + "ST", A: e.S, B: e.T, wantEq: e.wantEq},
		)
	}

	//fmt.Println("testcases {name, A, B, wantEq}")
	//for i, c := range testcases {
	//	fmt.Println(i, "\t\t", c)
	//}

	for fn, f := range algos {
		for _, tt := range testcases {
			name := fn + "_" + tt.name
			t.Run(name, func(t *testing.T) {
				got := f(tt.A, tt.B)
				if got != tt.wantEq {
					t.Errorf("name=%v\n\tA=%v\n\tB=%v\n\tgotEq=%v, wantEq=%v", name, tt.A, tt.B, got, tt.wantEq)
				}
			})
		}
	}
}

var (
	bcmkCeqS1 = "自然語言處理是人工智能領域中的一個重要方向。"
	bcmkCeqS2 = "自然语言处理是人工智能领域中的一个重要方向。"
)

func benchmarkChineseStringEqual(b *testing.B, algo func(string, string) bool) {
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		algo(bcmkCeqS1, bcmkCeqS2)
	}
}

func BenchmarkChineseStringEqual_OnlyNecessaryCC(b *testing.B) {
	benchmarkChineseStringEqual(b, chineseStringEqualOnlyNecessaryCC)
}

func BenchmarkChineseStringEqual_BothCC(b *testing.B) {
	benchmarkChineseStringEqual(b, chineseStringEqualBothCC)
}

func BenchmarkChineseStringEqual_CoroutineCC(b *testing.B) {
	benchmarkChineseStringEqual(b, chineseStringEqualCoroutineCC)
}

func BenchmarkChineseStringEqual_AlwaysCC(b *testing.B) {
	benchmarkChineseStringEqual(b, chineseStringEqualAlwaysCC)
}

func BenchmarkChineseStringEqual_AlwaysCoroutineCC(b *testing.B) {
	benchmarkChineseStringEqual(b, chineseStringEqualAlwaysCoroutineCC)
}
