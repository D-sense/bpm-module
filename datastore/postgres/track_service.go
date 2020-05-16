package postgres

import (
	"context"
	"github.com/go-apps/bpm-module/modules"
	"time"
)

type TrackService struct {
	client *Client
}

func NewTrackService(ctx context.Context, client *Client) *TrackService {
	ts := &TrackService{client}
	return ts
}


func (ts TrackService) Tracks(ctx context.Context, excludeIds []string) ([]*modules.Track, error) {
	var track []*modules.Track
	err := ts.client.db.Select("id, counter, original_resource, metadata").Where("bpm is null AND deleted_at is null").Order("counter desc").Find(&track).Error
	return track, err
}

func (ts TrackService) UpdateTrackBPM(ctx context.Context, track *modules.Track) error {
	track.UpdatedAt = time.Now()
	return ts.client.db.Save(&track).Error
}

func (ts TrackService) FindTrackByIdAndCounter(ctx context.Context, track *modules.Track) (*modules.Track, error) {
	output := &modules.Track{}
	err := ts.client.db.First(output, "id = ? AND counter = ? ", track.ID, track.Counter).Error
	return output, err
}
