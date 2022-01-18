package spotify

import (
	"github.com/zmb3/spotify"
	"testing"
)

//func TestNewClientAndAuth(t *testing.T) {
//	type Args struct {
//		ClientID     string `json:"client_id"`
//		ClientSecret string `json:"client_secret"`
//	}
//
//	var goodArg Args
//	f, _ := os.Open("config.json")
//	j, _ := ioutil.ReadAll(f)
//	err := json.Unmarshal(j, &goodArg)
//	if err != nil {
//		t.Error(err)
//	}
//	fmt.Println(goodArg)
//
//	badArg := Args{
//		ClientID:     "sdfsdfsdfsfsd",
//		ClientSecret: "dsfasdfsdfsdf",
//	}
//
//	tests := []struct {
//		name    string
//		args    Args
//		wantErr bool
//	}{
//		{"bad", badArg, true},
//		{"good", goodArg, false},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := NewClientAndAuth(tt.args.ClientID, tt.args.ClientSecret)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("NewClientAndAuth() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			t.Logf("got: %#v", got)
//			t.Logf("err: %#v", err)
//		})
//	}
//}

func TestFindPlaylist(t *testing.T) {

	client := NewSpotifyClient()

	type args struct {
		client *spotify.Client
		w      string
	}
	tests := []struct {
		name string
		args args
	}{
		{"findMabanua_LessThanOnePage", args{
			client: client,
			w:      "mabanua",
		}},
		{"findMegalobox_TwoPages", args{
			client: client,
			w:      "megalobox",
		}},
		{"findA_TooMany", args{
			client: client,
			w:      "a",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindPlaylist(tt.args.client, tt.args.w)
			t.Logf("got ch: %#v", got)

			counter := 0
			seen := make(map[string]int)

			for p := range got {
				t.Logf("yield(%v) %v %v", counter, p.ID, p.Name)

				if idx, alreadyExist := seen[p.ID.String()]; alreadyExist {
					// 偶尔会重复，spotify 的锅。
					t.Logf("❌ Duplicated with (%v)", idx)

				} else {
					seen[p.ID.String()] = counter
				}

				counter++
			}
			t.Logf("yield items: %v", counter)
		})
	}
}

func TestFetchTracks(t *testing.T) {
	client := NewSpotifyClient()

	// https://open.spotify.com/playlist/{ID}
	PlaylistMegaloBox := "0qweVVPXhB0VTVulzgxy4W"
	PlaylistMabanua := "37i9dQZF1DZ06evO0JPEL6"
	PlaylistAllTheFeels := "37i9dQZF1DX7gIoKXt0gmx"

	type args struct {
		client   *spotify.Client
		playlist *spotify.SimplePlaylist
	}
	tests := []struct {
		name      string
		args      args
		wantCount int
	}{
		{"MegaloBox", args{
			client:   client,
			playlist: &spotify.SimplePlaylist{ID: spotify.ID(PlaylistMegaloBox)},
		}, 61},
		{"Mabanua", args{
			client:   client,
			playlist: &spotify.SimplePlaylist{ID: spotify.ID(PlaylistMabanua)},
		}, 31},
		{"AllTheFeels", args{
			client:   client,
			playlist: &spotify.SimplePlaylist{ID: spotify.ID(PlaylistAllTheFeels)},
		}, 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FetchTracks(tt.args.client, tt.args.playlist)
			t.Logf("got ch: %#v", got)

			counter := 0
			seen := make(map[string]int)

			for p := range got {
				t.Logf("yield(%v) %v %v", counter, p.ID, p.Name)

				if idx, alreadyExist := seen[p.ID.String()]; alreadyExist {
					t.Logf("⚠️ Duplicated with (%v)", idx)

				} else {
					seen[p.ID.String()] = counter
				}

				counter++
			}

			if counter != tt.wantCount {
				t.Errorf("❌ yield items: %v (want %v)", counter, tt.wantCount)
			} else {
				t.Logf("✅ yield items: %v (want %v)", counter, tt.wantCount)
			}
		})
	}
}
