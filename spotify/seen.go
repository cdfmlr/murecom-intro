package spotify

import (
	"sync"
)

const SeenDebug = false

// seenWords 见过的词（seed）
var seenWords sync.Map

// InitSeenWords 初始化已见词表.
// After InitConfig
func InitSeenWords() {
	seenWords = sync.Map{}

	for _, w := range Config.SeenWords {
		// Seed should not be seen when init,
		// otherwise master popping from an empty wordsCounter.
		if w != Config.Seed {
			seenWords.Store(w, struct{}{})
		}
	}
	if SeenDebug {
		println("debug(seen) InitSeenWords ", len(Config.SeenWords))
	}
}

// SeenWord 返回 false 若没见过 w，同时把 w 标记为见过
//
// Example:
//     w := "notSeen"
//     SeenWord(w) // false
//     SeenWord(w) // true
func SeenWord(w string) bool {
	_, ok := seenWords.LoadOrStore(w, struct{}{})
	if SeenDebug {
		println("debug(seen): SeenWord", w, ": ", ok)
	}
	return ok
}

// seenPlaylists 见过的播放列表
var seenPlaylists sync.Map

// sumID 给 ID 做摘要：22 chars -> 10 chars
func sumID(id string) string {
	if len(id) < 16 {
		return id
	}

	s := id[len(id)-8:]

	var c uint8 = 0
	for i := len(id) - 8; i >= 0; i-- {
		c = (c + id[i]) % 255
		if i%8 == 0 {
			if c == 0 {
				c = 101
			}
			s += string(c)
			c = 0
		}
	}

	return s
}

// InitSeenPlaylists 初始化已见播放列表.
// After InitConfig & InitDB
func InitSeenPlaylists() {
	seenPlaylists = sync.Map{}

	type playlistID struct {
		ID string
	}
	var ids []playlistID

	DB.Model(&Playlist{}).Find(&ids)

	for _, id := range ids {
		seenPlaylists.Store(sumID(id.ID), struct{}{})
	}

	if SeenDebug {
		println("debug(seen) InitSeenPlaylists ", len(ids))
	}
}

func SeenPlaylist(id string) bool {
	_, ok := seenPlaylists.LoadOrStore(sumID(id), struct{}{})
	if SeenDebug {
		println("debug(seen): SeenPlaylist", id, ": ", ok)
	}
	return ok
}
