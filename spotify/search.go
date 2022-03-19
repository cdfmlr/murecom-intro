package spotify

import (
	"errors"
	"fmt"
	"github.com/zmb3/spotify"
	"sort"
	"strings"
)

var SearchTrackLimit = 5

type TrackQuery struct {
	TrackName string
	Artists   []string
}

func (q TrackQuery) String() string {
	return q.TrackName + " " + strings.Join(q.Artists, " ")
}

func SearchTrack(client *spotify.Client, query *TrackQuery) (*Track, error) {
	client.AcceptLanguage = "zh"
	// query spotify
	result, err := client.SearchOpt(query.String(), spotify.SearchTypeTrack,
		&spotify.Options{Limit: &SearchTrackLimit})
	if err != nil {
		return nil, err
	}

	// is result really matching?
	for _, t := range result.Tracks.Tracks {
		// XXX: 这里现在只检查歌名了，不看艺人，艺人在 Query 里就行。
		//      有些艺人日文的英文、假名、罗马音、汉字不好处理。
		if trackMatches(&t, query, trackNameMatch /*, trackArtistMatch*/) {
			track := trackFromSpotify(t.SimpleTrack)
			return track, nil
		}
	}

	// TODO: not matched: try fuzzing matches

	return nil, ErrTrackNotFound
}

var ErrTrackNotFound = errors.New("no matching track")

// region TrackMatcher

// trackMatcher checks whether track matches trackName and artists.
type trackMatcher func(track *spotify.FullTrack, query *TrackQuery) bool

// trackMatches checks track by chain
func trackMatches(track *spotify.FullTrack, query *TrackQuery, checks ...trackMatcher) bool {
	var as []string
	for _, a := range track.Artists {
		as = append(as, a.Name)
	}
	fmt.Printf("track=%v, artist=%v,  qurey=%v\n", track.Name, as, query)
	for _, f := range checks {
		ok := f(track, query)
		if !ok {
			return false
		}
	}
	return true
}

// true only if track names are literally equal
func trackNameMatch(track *spotify.FullTrack, query *TrackQuery) bool {
	return ChineseStringEqual(track.Name, query.TrackName)
}

// checks the first artist
func trackArtistMatch(track *spotify.FullTrack, query *TrackQuery) bool {
	artists := query.Artists
	if len(artists) == 0 {
		return true
	}

	sort.Strings(artists)
	for _, a := range track.Artists {
		aname, _ := T2S(a.Name)
		idx := sort.SearchStrings(artists, aname)
		if idx < len(artists) && ChineseStringEqual(artists[idx], aname) {
			return true
		}
	}
	return false
}

// endregion TrackMatcher

// TODO: fuzzing matcher
