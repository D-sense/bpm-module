package track

import (
	"context"
	"fmt"
	"github.com/go-apps/bpm-module/bpm"
	"github.com/go-apps/bpm-module/logger"
	"github.com/go-apps/bpm-module/modules"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var standardLogger = logger.NewLogger()

type Handler struct {
	trackService    modules.TrackService
	bpmEngine       bpm.InterfaceBpm
}

func NewHandler(trackService modules.TrackService, bpm bpm.InterfaceBpm) *Handler {
	return &Handler {
		trackService:    trackService,
		bpmEngine: bpm,
	}
}

func (h *Handler) StartBpmService(ctx context.Context, logger log.Entry, excludeCounters []string) {
	success := 0
	failure := 0

	// fetch tracks
	tracks, err := h.Tracks(ctx, excludeCounters, logger)
	if err != nil {
		standardLogger.GetError(fmt.Sprintf("error getting tracks: %v ", err))
		log.Fatal()
	}

	standardLogger.GetNotice(fmt.Sprintf("********* Starting New Process **********"))

	for _, track := range tracks {
		//obtain URLs
		trackUrl, err := GetTrackUrl(*track.OriginalResource)
		if err != nil {
			standardLogger.GetError(fmt.Sprintf("cannot get track url: ID: %v    |  Counter: %v  ", track.ID, track.Counter))
			log.Fatal()
		}

		// get track's bpm
		bpmResult, status, err := h.bpmEngine.ExtractBpm(string(trackUrl))
		if err != nil && status == false {
			excludeCounters = append(excludeCounters, track.ID)
			failure++

			standardLogger.GetError(fmt.Sprintf(" BPM extraction failed: ID: %v    |  Counter: %v   |  track_url: %v | Error: %v", track.ID, track.Counter, trackUrl, err))
		}else if status {
			tr, err := h.trackService.FindTrackByIdAndCounter(ctx, track)

			//update track
			tr.BPM = bpmResult
			err = h.UpdateTrackBPM(ctx, tr, logger)
			if err != nil{
				standardLogger.GetError(fmt.Sprintf("cannot update track row: ID: %v    |  Counter: %v   |    bpm: %v", track.ID, track.Counter, bpmResult))
			}
			success++

			standardLogger.GetInfo(fmt.Sprintf(" BPM extraction passed: ID: %v    |  Counter: %v   |    bpm: %v", track.ID, track.Counter, bpmResult))
		}
	}
	standardLogger.GetNotice(fmt.Sprintf("********* Ended The Process **********"))
}

func (h *Handler) Tracks(ctx context.Context, excluded []string, logger log.Entry) ([]*modules.Track, error) {
	//logger.Info("get tracks")

	track, err := h.trackService.Tracks(ctx, excluded)
	if err != nil {
		return nil, errors.Wrap(err, "error getting tracks:")
	}

	return track, err
}

func (h *Handler) UpdateTrackBPM(ctx context.Context, input *modules.Track, logger log.Entry) error {
	//logger.Info("update track")
	err := h.trackService.UpdateTrackBPM(ctx, input)
	if err != nil {
		return errors.Wrap(err, "error updating a track:")
	}

	return err
}

func GetTrackUrl(url modules.UrlResource) (string, error) {
	u, err := url.Value()
	if err != nil {
		return "", errors.WithMessage(err, "error while converting original_resource struct to driver.value:")
	}

	err = url.Scan(u)
	if err != nil {
		return "", errors.WithMessage(err, "error while processing original_resource struct:")
	}

	return url.URL.String(), nil
}

