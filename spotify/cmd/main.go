package main

import (
	"fmt"
	"os"
	"path"
	"runtime/pprof"
	spotify "spotifyplaylist"
	"strings"
	"time"
)

// Profile for debug
const (
	Profile = false
)

func main() {
	st := strings.ReplaceAll(time.Now().Format(time.RFC3339), ":", "-")
	if Profile {
		f, err := os.Create(path.Join(
			spotify.Config.Profile,
			fmt.Sprintf("%v.prof", st)))
		if err != nil {
			panic("Profile Failed: " + err.Error())
		}
		err = pprof.StartCPUProfile(f)
		if err != nil {
			panic("Profile Failed: " + err.Error())
		}
		defer pprof.StopCPUProfile()
	}

	spotify.CrawlPlaylistsAndTracks(
		spotify.Config.Seed,
		spotify.Config.Playlists)

	if Profile {
		f, err := os.Create(path.Join(
			spotify.Config.Profile,
			fmt.Sprintf("%v.mprof", st)))
		if err != nil {
			panic("Profile Failed: " + err.Error())
		}
		_ = pprof.WriteHeapProfile(f)
		_ = f.Close()
	}
}
