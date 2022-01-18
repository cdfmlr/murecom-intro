package spotify

import (
	"fmt"
	"github.com/zmb3/spotify"
)

// 给我慢一点！！太快会 API rate limit exceeded
// 用 zmb3/spotify 的 AutoRetry 就不必这个了
func slowDown() {
	//numGoroutine := runtime.NumGoroutine()
	//time.Sleep(time.Duration(rand.Intn(numGoroutine)) * 10 * time.Millisecond)
}

//// Playlist 是要获取的播放列表信息
//// 包含 SimplePlaylist 提供的列表信息，以及 SimpleTracks 是歌曲信息列表
//type Playlist struct {
//	SimplePlaylist spotify.SimplePlaylist //https://pkg.go.dev/github.com/zmb3/spotify#SimplePlaylist
//	SimpleTracks   []spotify.SimpleTrack  // https://pkg.go.dev/github.com/zmb3/spotify#SimpleTrack
//}

// Paging of playlist searching
// https://developer.spotify.com/documentation/web-api/reference/#/operations/search
const (
	PlaylistPageLimit     = 50
	PlaylistPageMaxOffset = 1000
)

// FindPlaylist 获取匹配搜索项 w 的所有播放列表
//
// 返回一个出 spotify.SimplePlaylist 的 chan （带缓冲: size = PlaylistPageLimit），结果这个 chan 传
func FindPlaylist(client *spotify.Client, w string) <-chan spotify.SimplePlaylist {
	ch := make(chan spotify.SimplePlaylist, PlaylistPageLimit)

	go func() {
		counter := 0 // 记录出了多少了

		limit := PlaylistPageLimit // spotify.Options 要取 *int
		// Search API
		res, err := client.SearchOpt(w, spotify.SearchTypePlaylist, &spotify.Options{
			Limit: &limit,
		})
		if err != nil {
			fmt.Printf("FindPlaylist faild: %v\n", err)
			close(ch)
			return
		}
		//println("total: ", res.Playlists.Total)

		// Yield first page
		for _, p := range res.Playlists.Playlists {
			ch <- p
			counter++
		}

		for counter < res.Playlists.Total {
			// 下一页
			if err = client.NextPlaylistResults(res); err != nil {
				// 404 算正常的结束
				// 接口类型选择，防止其他错误 panic 崩掉
				switch e := err.(type) {
				case spotify.Error:
					if e.Status != 404 {
						fmt.Printf("FetchTracks faild: %v\n", err)
					}
				default:
					fmt.Printf("FetchTracks faild: %v\n", err)
				}
				close(ch)
				return
			}

			// yield
			for _, p := range res.Playlists.Playlists {
				ch <- p
				counter++
			}

			slowDown()
		}

		// Done
		close(ch)
	}()

	return ch
}

const (
	TracksPageLimit = 50
)

// FetchTracks 获取给定的 playlist 中的曲目
//
// 返回一个出 spotify.SimpleTrack 的 chan （带缓冲: size = TracksPageLimit），结果这个 chan 传
func FetchTracks(client *spotify.Client, playlist *spotify.SimplePlaylist) <-chan spotify.SimpleTrack {
	ch := make(chan spotify.SimpleTrack, TracksPageLimit)

	go func() {
		counter := 0 // 记录出了多少了

		limit := TracksPageLimit // spotify.Options 要取 *int
		// Search API
		res, err := client.GetPlaylistTracksOpt(playlist.ID, &spotify.Options{
			Limit: &limit,
		}, "")
		if err != nil {
			fmt.Printf("FetchTracks faild: %v\n", err)
			close(ch)
			return
		}
		//println("total: ", res.Total)

		// Yield first page
		for _, t := range res.Tracks {
			ch <- t.Track.SimpleTrack
			counter++
		}

		for counter < res.Total {
			// 下一页
			if err = client.NextPage(res); err != nil {
				// 404 算正常的结束
				// 接口类型选择，防止其他错误 panic 崩掉
				switch e := err.(type) {
				case spotify.Error:
					if e.Status != 404 {
						fmt.Printf("FetchTracks faild: %v\n", err)
					}
				default:
					fmt.Printf("FetchTracks faild: %v\n", err)
				}

				close(ch)
				return
			}

			// yield
			for _, t := range res.Tracks {
				ch <- t.Track.SimpleTrack
				counter++
			}

			slowDown()
		}

		// Done
		close(ch)
	}()

	return ch
}
