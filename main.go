package main

import (
	"io"
	"log"
	"os"
	"time"
)

var (
	StreamURL             = "https://avw.mdr.de/streams/284321-2_mp3_high.m3u"
	SongNameURL           = "https://www.mdr.de/XML/titellisten/mdr_piraten_2.json" // lol
	IsRunning             = true
	ControlRecording      = make(chan bool)
	IsRecording           = false // has to be here because working with channels is shit
	SongStart             = ""
	SongEnd               = ""
	Duration              = ""
	ImageURL              = ""
	CurrentSong           *Song
	TimeStartedRecording  time.Time
	TimeStoppedRecording  time.Time
	SimilarityThreshold   = 0.95
	ComparisonSnippetTime = 10 * time.Second

	EstimatedRuntime         = 8 * time.Hour
	CurrentRecordingFilename = ""

	targetSampleRate = 44100 // yolo
)

func main() {

	logFile, err := os.Create("server.log")
	if err != nil {
		log.Fatal("Error creating log file:", err)
	}
	defer logFile.Close()
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	log.Println("mixscribe | © lowikian@einfachzocken & contributors")

	log.Println("Starting Webserver")

	go startWeb()

	log.Println("Starting Recorder")

	go StartRecording()

	log.Println("Starting UpdateInfoLoop")
	go UpdateSongInfoLoop()

	select {}
}
