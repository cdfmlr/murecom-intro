package spotify

import (
	"github.com/zmb3/spotify"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Playlist 是 spotify.SimplePlaylist 删掉不要的字段
type Playlist struct {
	Collaborative bool    `json:"collaborative"`
	Endpoint      string  `json:"href"`
	ID            string  `json:"id" gorm:"primaryKey"`
	Name          string  `json:"name"`
	IsPublic      bool    `json:"public"`
	SnapshotID    string  `json:"snapshot_id"`
	URI           string  `json:"uri"`
	Tracks        []Track `json:"my_tracks" gorm:"many2many:playlist_tracks"`
}

func PlaylistFromSpotify(sp spotify.SimplePlaylist, sts []spotify.FullTrack) *Playlist {
	p := Playlist{
		Collaborative: sp.Collaborative,
		Endpoint:      sp.Endpoint,
		ID:            string(sp.ID),
		Name:          sp.Name,
		IsPublic:      sp.IsPublic,
		SnapshotID:    sp.SnapshotID,
		URI:           string(sp.URI),
		Tracks:        []Track{},
	}
	for _, st := range sts {
		p.Tracks = append(p.Tracks, *trackFromSpotify(st))
	}

	return &p
}

type Artist struct {
	Name string `json:"name"`
	ID   string `json:"id" gorm:"primaryKey"`
	URI  string `json:"uri"`
}

func artistFromSpotify(sa spotify.SimpleArtist) *Artist {
	return &Artist{
		Name: sa.Name,
		ID:   string(sa.ID),
		URI:  string(sa.URI),
	}
}

type Track struct {
	Artist      []Artist `json:"artists" gorm:"many2many:track_artists"`
	DiscNumber  int      `json:"disc_number"`
	Duration    int      `json:"duration_ms"`
	Explicit    bool     `json:"explicit"`
	Endpoint    string   `json:"href"`
	ID          string   `json:"id"  gorm:"primaryKey"`
	Name        string   `json:"name"`
	PreviewURL  string   `json:"preview_url"`
	TrackNumber int      `json:"track_number"`
	URI         string   `json:"uri"`
	Type        string   `json:"type"`
	ImageUrl    string   `json:"image_url"`
}

func trackFromSpotify(st spotify.FullTrack) *Track {
	t := Track{
		Artist:      []Artist{},
		DiscNumber:  st.DiscNumber,
		Duration:    st.Duration,
		Explicit:    st.Explicit,
		Endpoint:    st.Endpoint,
		ID:          string(st.ID),
		Name:        st.Name,
		PreviewURL:  st.PreviewURL,
		TrackNumber: st.TrackNumber,
		URI:         string(st.URI),
		Type:        st.Type,
	}
	for _, a := range st.Artists {
		t.Artist = append(t.Artist, *artistFromSpotify(a))
	}
	if len(st.Album.Images) > 0 {
		t.ImageUrl = st.Album.Images[0].URL
	}

	return &t
}

// Save 在一次事物中保存播放列表及其中所有曲目
func Save(playlist *Playlist) {
	DB.Save(playlist)
}

// DB 全局的数据库单例
var DB *gorm.DB

// InitDB 初始化 DB
func InitDB() {
	var err error

	DB, err = gorm.Open(sqlite.Open(Config.DB), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	err = DB.AutoMigrate(&Artist{}, &Track{}, &Playlist{})
	if err != nil {
		panic("failed to AutoMigrate database: " + err.Error())
	}
}
