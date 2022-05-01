package main

import (
	"flag"
	"fmt"
	spotify "github.com/cdfmlr/murecom-intro/spotify"
	"os"
	"path"
	"runtime/pprof"
	"strings"
	"time"
)

var configFile = flag.String("config", "config.json", "/path/to/config/file")
var runComplementTracksImage = flag.Bool("complementTracksImage", false, "run ComplementTracksImage instead of CrawlPlaylistsAndTracks")

// Profile for debug
const (
	Profile = false
)

func main() {
	flag.Parse()

	fmt.Printf("spotify init all, using config file: %v\n", *configFile)
	spotify.InitAll(*configFile)

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

	switch {
	case *runComplementTracksImage:
		spotify.ComplementTracksImage()
	default:
		spotify.CrawlPlaylistsAndTracks(
			spotify.Config.Seed,
			spotify.Config.Playlists)
	}

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
