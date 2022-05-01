package spotify

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

var Config struct {
	ClientID            string   `json:"client_id"`
	ClientSecret        string   `json:"client_secret"`
	DB                  string   `json:"db"`
	Profile             string   `json:"profile"`
	Seed                string   `json:"seed"`
	Playlists           int      `json:"playlists"`
	SeenWords           []string `json:"seen_words"`
	UpdateTracksWorkers int      `json:"update_tracks_workers"`
}

// InitConfig read the config file
func InitConfig(configfile string) {
	f, err := os.Open(configfile)
	if err != nil {
		panic("failed to load config: " + err.Error())
	}
	j, err := ioutil.ReadAll(f)
	if err != nil {
		panic("failed to load config: " + err.Error())
	}
	err = json.Unmarshal(j, &Config)
	if err != nil {
		panic("failed to load config: " + err.Error())
	}
}
