package spotify

import (
	"github.com/zmb3/spotify"
	"testing"
)

func TestSave(t *testing.T) {
	client := NewSpotifyClient()

	chPl := FindPlaylist(client, "fool")
	playlist := <-chPl

	chTk := FetchTracks(client, &playlist)
	var tracks []spotify.SimpleTrack
	for t := range chTk {
		tracks = append(tracks, t)
	}

	p := PlaylistFromSpotify(playlist, tracks)
	t.Logf("playlist(%v) %v, tracks %v", p.ID, p.Name, len(p.Tracks))

	type args struct {
		playlist *Playlist
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"save", args{
			playlist: p,
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Save(tt.args.playlist)
		})
	}
}
