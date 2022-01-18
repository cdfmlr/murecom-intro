package spotify

import (
	"testing"
)

func TestCrawlPlaylistsAndTracks(t *testing.T) {
	//time.Sleep(5 * time.Second)

	type args struct {
		seed   string
		maxNum int
	}
	tests := []struct {
		name string
		args args
	}{
		{"a-1", args{seed: "a", maxNum: 1}},
		{"a-10", args{seed: "a", maxNum: 10}},
		{"a-100", args{seed: "a", maxNum: 100}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CrawlPlaylistsAndTracks(tt.args.seed, tt.args.maxNum)
		})
	}

	//time.Sleep(10 * time.Second)
}

func TestWordsCounts(t *testing.T) {

	type args struct {
		seed   string
		maxNum int
	}
	tests := []struct {
		name string
		args args
	}{
		{"mabanua-1000", args{seed: "mabanua", maxNum: 1000}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CrawlPlaylistsAndTracks(tt.args.seed, tt.args.maxNum)
		})
	}

	//time.Sleep(10 * time.Second)
}
