package spotify

import (
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/zmb3/spotify"
	"gorm.io/gorm"
	"sync"
)

const (
	UpdateTracksBatchSize = 32 // < 50
	//UpdateTracksWorkers   = 128 // use config
)

// UpdateTracks 请求 spotify，更新数据库中所有 Track 的数据。
// where 是数据库取 Track 的 where 查询条件。
//
// 也可用于在数据库新增字段后，补充之前的记录。
func UpdateTracks(where string) {
	client := NewSpotifyClient()

	fmt.Printf("Counting... ")

	var count int64 //  = 4244758
	DB.Model(&Track{}).Where(where).Count(&count)
	batchCount := count / UpdateTracksBatchSize
	progress := pb.Full.Start(int(batchCount))

	fmt.Printf("UpdateTracks in DB (%v): total_tracks=%v in %v batchs\n", Config.DB, count, batchCount)

	workers := make(chan struct{}, Config.UpdateTracksWorkers) // limit workers count
	wg := &sync.WaitGroup{}

	var tracks []Track
	DB.Model(&Track{}).Where(where).FindInBatches(&tracks, UpdateTracksBatchSize, func(tx *gorm.DB, batch int) error {
		workers <- struct{}{} // P
		wg.Add(1)

		ids := tracks2ids(tracks)

		go func(trackIDs []spotify.ID) {
			//println("in closure tracks:", trackIDs)
			_ = updateTracksWorker(client, trackIDs)
			<-workers // V
			wg.Done()
			progress.Increment()
		}(ids)
		return nil
	})
	wg.Wait()
	progress.Finish()
}

func tracks2ids(tracks []Track) []spotify.ID {
	var ids []spotify.ID
	for _, t := range tracks {
		if t.ID == "" {
			continue
		}
		ids = append(ids, spotify.ID(t.ID))
	}
	return ids
}

func updateTracksWorker(client *spotify.Client, trackIDs []spotify.ID) error {
	//fmt.Printf("debug updateTracksWorker: ids: %v\n", trackIDs)

	// XXX: but why this happens
	if trackIDs == nil {
		//println("ids == nil")
		return nil
	}

	newTracks, err := client.GetTracks(trackIDs...)
	//println("debug updateTracksWorker GetTracks done")
	if err != nil {
		fmt.Printf("UpdateTracks failed: error=%v, tracks=%#v\n", err, trackIDs)
		return err
	}

	for _, st := range newTracks {
		t := trackFromSpotify(*st)
		//fmt.Println("debug save", t.ID, t.ImageUrl, st.Album.Images)
		DB.Save(&t)
	}
	//println("debug updateTracksWorker")
	return err
}

// ComplementTracksImage is an instance of UpdateTracks to complement Tracks' ImageUrl fields.
func ComplementTracksImage() {
	fmt.Println("Complement Tracks ImageUrl field")

	UpdateTracks("image_url is null or image_url = ''")
}
