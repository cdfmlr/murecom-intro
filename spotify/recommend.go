package spotify

import "github.com/zmb3/spotify"

func SpotifyRecommendTracks(client *spotify.Client, seedTrackIds []string, limit int) ([]*Track, error) {
	client.AcceptLanguage = "zh"

	ids := make([]spotify.ID, len(seedTrackIds))
	for i, v := range seedTrackIds {
		ids[i] = spotify.ID(v)
	}

	result, err := client.GetRecommendations(
		spotify.Seeds{Tracks: ids},
		nil,
		&spotify.Options{Limit: &limit})

	if err != nil {
		return nil, err
	}

	var tracks []*Track
	for _, t := range result.Tracks {
		tracks = append(tracks, trackFromSpotify(t))
	}
	return tracks, nil
}
