package spotify

import "testing"

func Test_sumID(t *testing.T) {
	goodID1 := "37i9dQZF1DXe30kLtifvte"
	goodID1t := "37i9dQZF1DXeifvte30kLt"
	goodID1h := "9d37iQZF1DXe30kLtifvte"
	goodID2t := "37i9dQZF1DXe30kLt9fvte"
	goodID2h := "3ai9dQZF1DXe30kLtifvte"
	longID := "helloworld" + goodID1

	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		wantLen int
	}{
		{"realID", args{id: goodID1}, 12},
		{"realID尾重排", args{id: goodID1t}, 12},
		{"realID头重排", args{id: goodID1h}, 12},
		{"realID改尾一位", args{id: goodID2t}, 12},
		{"realID改头一位", args{id: goodID2h}, 12},
		{"empty", args{id: ""}, 0},
		{"tooShort", args{id: "asdf"}, 4},
		{"tooLong", args{id: longID}, len(longID) - 12},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sumID(tt.args.id); len(got) > tt.wantLen {
				t.Errorf("❌ sumID(%#v) = %v (len=%v), want len<=%v",
					tt.args.id, got, len(got), tt.wantLen)
			} else {
				t.Logf("✅ sumID(%#v) = %v (len=%v), want len=%v",
					tt.args.id, got, len(got), tt.wantLen)
			}
		})
	}
}
