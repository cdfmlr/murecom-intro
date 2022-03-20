package spotify

import "github.com/zmb3/spotify"

type AudioFeatures = spotify.AudioFeatures

func GetAudioFeatures(client *spotify.Client, trackIds ...string) ([]*AudioFeatures, error) {
	if len(trackIds) == 0 {
		return nil, nil
	}

	ids := make([]spotify.ID, len(trackIds))
	for i, v := range trackIds {
		ids[i] = spotify.ID(v)
	}

	return client.GetAudioFeatures(ids...)
}
