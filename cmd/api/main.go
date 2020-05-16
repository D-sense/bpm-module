package main

import (
	"context"
	"github.com/go-apps/bpm-module/bpm/audio"
	"github.com/go-apps/bpm-module/datastore/postgres"
	"github.com/go-apps/bpm-module/features/track"
	"github.com/go-apps/bpm-module/vault"
	"github.com/jasonlvhit/gocron"
	log "github.com/sirupsen/logrus"
)

func main(){

	log.Println("Starting Gbedu's BPM-Module Service!")

	initContext := context.Background()
	tracksTempFolder := "home/ubuntu/tracks_folder/"

	postgresCred, err := vault.GetCredentials()
	if err != nil {
		log.Panic("error fetching postgres credentials err: %v", err)
	}

	database := postgres.NewPGX(initContext, postgresCred)
	defer database.Close()

	trackService := postgres.NewTrackService(initContext, database)
	audioBpm := audio.BpmAudio{
		Filepath: tracksTempFolder,
	}

	trackHandler := track.NewHandler(trackService, audioBpm)

	var logger log.Entry
	excludeCounters := make([]string, 0)

	gocron.Every(1).Second().DoSafely(trackHandler.StartBpmService, initContext, logger, excludeCounters)

	<- gocron.Start()
}

