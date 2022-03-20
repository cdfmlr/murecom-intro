package spotify

import (
	"fmt"
	"testing"
)

func TestGetAudioFeatures(t *testing.T) {
	InitConfig("config.json")
	InitAuth()

	client := NewSpotifyClient()

	tests := []struct {
		name    string
		ids     []string
		wantErr bool
		wantLen int
	}{
		{"nil", nil, false, 0},
		{"empty", []string{}, false, 0},
		{"single", []string{"4sZDrNBr9EsEvoOj7Z9Gr7"}, false, 1},
		{"multi", []string{ // from TestSearchTrack: 山丘，挪威的森林，朋友，神话情话，A Tender Moon Tempo
			"4sZDrNBr9EsEvoOj7Z9Gr7", "5KdZLYyoBP0qsjDTDP28do",
			"3Gl0Y21R2b5fPPxge18pTu", "7m6CUc95QqwoaNtJ82YbTO",
			"4M4kt9jB1nOMmyA83rDy98"}, false, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAudioFeatures(client, tt.ids...)
			if (err != nil) != tt.wantErr {
				t.Errorf("❌ GetAudioFeatures() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.wantLen {
				t.Errorf("❌ got len = %v, wantLen = %v", len(got), tt.wantLen)
			}
			t.Log("✅ got:", got)
			for i, f := range got {
				fmt.Printf("\t%v\t%#v\n", i, f)
			}
		})
	}
}
