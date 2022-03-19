package spotify

import "testing"

func TestChineseStringEqual(t *testing.T) {
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
		{
			wantEq: true,
			T:      "周華健",
			S:      "周华健",
		},
		{
			wantEq: true,
			T:      "齊豫",
			S:      "齐豫",
		},
	}

	type testcase struct {
		name   string
		A      string
		B      string
		wantEq bool
	}

	// testcases = examples * {Traditional, Simplified}^2 in cartesian product
	var testcases []testcase
	for _, e := range examples {
		k := "P"
		if e.wantEq == false {
			k = "N"
		}
		k = e.S + k
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

	for _, tt := range testcases {
		name := tt.name
		t.Run(name, func(t *testing.T) {
			got := ChineseStringEqual(tt.A, tt.B)
			if got != tt.wantEq {
				t.Errorf("name=%v\n\tA=%v\n\tB=%v\n\tgotEq=%v, wantEq=%v", name, tt.A, tt.B, got, tt.wantEq)
			}
		})
	}
}

func TestRemovePunct(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"神話.情話", args{"神話.情話"}, "神話情話"},
		{"神话·情话", args{"神话·情话"}, "神话情话"},
		{"自然语言", args{"自然語言處理!是人工智能！領域中？的一個重要方向，。？；：「」【】｜"}, "自然語言處理是人工智能領域中的一個重要方向｜"},
		{"English", args{"English1234567890!@#$%^&*()-=_+[]{};':,.<>/?|"}, "English1234567890$^=+<>|"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemovePunct(tt.args.s); got != tt.want {
				t.Errorf("RemovePunct() = %v, want %v", got, tt.want)
			}
		})
	}
}
