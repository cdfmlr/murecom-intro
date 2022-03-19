package spotify

import (
	"testing"
)

func TestSearchTrack(t *testing.T) {
	InitConfig("config.json")
	InitAuth() // TODO: InitAuth => InitSpotify

	client := NewSpotifyClient()

	tests := []struct {
		name    string
		query   TrackQuery
		wantErr bool
	}{
		{"simple: 山丘", TrackQuery{
			TrackName: "山丘",
			Artists:   []string{"李宗盛"},
		}, false},
		{"artists missing: 挪威的森林", TrackQuery{
			TrackName: "挪威的森林",
			Artists:   nil,
		}, false},
		{"简繁: 朋友", TrackQuery{
			TrackName: "朋友",
			Artists:   []string{"周华健"},
		}, false},
		{"multi-artists: 神话·情话", TrackQuery{
			TrackName: "神话·情话",
			Artists:   []string{"周华健", "齐豫"},
		}, false},
		{"english", TrackQuery{
			TrackName: "A Tender Moon Tempo",
			Artists:   []string{"Vivy (Vo.Kairi Yagi)"},
		}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SearchTrack(client, &tt.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchTrack() error = %v, wantEq %v", err, tt.wantErr)
				return
			}
			t.Logf("🔍  Query: %v-%v\n\tgot: (%v) %v-%v",
				tt.query.TrackName, tt.query.Artists,
				got.ID, got.Name, got.Artist)
		})
	}
}
