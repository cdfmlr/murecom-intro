package spotify

import (
	"fmt"
	"testing"
)

func TestSpotifyRecommendTracks(t *testing.T) {
	var samples = map[string]string{ // ID: track_name
		"4sZDrNBr9EsEvoOj7Z9Gr7": "山丘",
		"5KdZLYyoBP0qsjDTDP28do": "挪威的森林",
		"3Gl0Y21R2b5fPPxge18pTu": "朋友",
		"7m6CUc95QqwoaNtJ82YbTO": "神话情话",
		"4M4kt9jB1nOMmyA83rDy98": "A Tender Moon Tempo",
	}

	InitConfig("config.json")
	InitAuth()

	client := NewSpotifyClient()

	type args struct {
		seedTrackIds []string
		limit        int
	}
	tests := []struct {
		name    string
		args    args
		wantLen int
		wantErr bool
	}{
		{"nil seed", args{nil, 10}, 0, true},
		{"single seed", args{
			seedTrackIds: []string{"4sZDrNBr9EsEvoOj7Z9Gr7"}, //山丘
			limit:        10,
		}, 10, false},
		{"limit 20", args{
			seedTrackIds: []string{"4sZDrNBr9EsEvoOj7Z9Gr7"}, //山丘
			limit:        20,
		}, 20, false},
		{"multi seeds", args{
			seedTrackIds: []string{"5KdZLYyoBP0qsjDTDP28do", "4M4kt9jB1nOMmyA83rDy98"}, // 挪威的森林, A Tender Moon Tempo
			limit:        10,
		}, 10, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SpotifyRecommendTracks(client, tt.args.seedTrackIds, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("SpotifyRecommendTracks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.wantLen {
				t.Errorf("❌ SpotifyRecommendTracks() len(got) = %v, wantLen = %v: got = %v", len(got), tt.wantLen, got)
			}
			var seedNames []string
			for _, id := range tt.args.seedTrackIds {
				seedNames = append(seedNames, samples[id])
			}
			t.Logf("✅ seed %v: got len %v, err = %v", seedNames, len(got), err)
			for i, s := range got {
				fmt.Printf("\t%v %v\n", i, s.Name)
			}
		})
	}
}
