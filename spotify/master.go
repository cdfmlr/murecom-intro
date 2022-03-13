package spotify

import (
	"fmt"
	"github.com/cdfmlr/murecom-intro/spotify/nlp"
	"github.com/zmb3/spotify"
	"sync"
	"time"
)

// CrawlPlaylistsAndTracks 从 seed 开始，获取 maxNum 个包含这个单词的播放列表。
//
// 计算这些列表中所有单词出现的频次。
// 然后对列表中出现次数最多的单词执行相同的操作。
// 重复，知道获取到足够的列表。
func CrawlPlaylistsAndTracks(seed string, maxNum int) {
	countPlaylists := 0
	startTime := time.Now()

	wordsCounter := nlp.NewEnglishWordCounterPQ()
	wordsCounter.AddSentence(seed)

	for countPlaylists < maxNum {
		var w string
		for {
			w = wordsCounter.PopMostCommon()
			//if _, ok := seenWords[w]; !ok {
			//	seenWords[w] = struct{}{}
			//	break
			//}
			if !SeenWord(w) {
				break
			}
		}

		fmt.Printf("%v/%v (%v%%) word: %v\n", countPlaylists, maxNum, countPlaylists*100/maxNum, w)

		client := NewSpotifyClient()
		wg := sync.WaitGroup{}

		for playlist := range FindPlaylist(client, w) {
			if SeenPlaylist(playlist.ID.String()) {
				continue
			}

			wg.Add(1)

			// fetch tracks and save playlist
			playlist := playlist
			ct := countPlaylists
			go func() {
				defer wg.Done()
				defer logProgress(ct, maxNum, startTime)

				var tracks []spotify.SimpleTrack

				cli := NewSpotifyClient()

				for track := range FetchTracks(cli, &playlist) {
					tracks = append(tracks, track)
				}

				Save(PlaylistFromSpotify(playlist, tracks))
			}()

			// counts and words
			countPlaylists++
			if countPlaylists >= maxNum {
				break
			}
			wordsCounter.AddSentence(playlist.Name)
		}

		wg.Wait()
		//fmt.Printf("%v/%v (%v%%) word: %v\n", countPlaylists, maxNum, countPlaylists*100/maxNum, w)
	}
}

func logProgress(countPlaylists, maxNum int, startTime time.Time) {
	// mod = max(maxNum/1000, PlaylistPageLimit)
	thousandth := maxNum / 1000
	mod := thousandth
	if mod < PlaylistPageLimit {
		mod = PlaylistPageLimit
	}

	now := time.Now()
	elapse := now.Sub(startTime)

	var left time.Duration
	f := float32(countPlaylists) / float32(maxNum)
	if f != 0 {
		left = time.Duration(
			float32(elapse)/f - float32(elapse),
		)
	} else {
		left = time.Hour * 999
	}

	elapse = elapse.Round(time.Second)
	left = left.Round(time.Second)

	if countPlaylists%mod == 0 {
		fmt.Printf("[%v ... %v] %v/%v (%v%%)\n",
			elapse, left,
			countPlaylists, maxNum, countPlaylists*100/maxNum)
	}
}
