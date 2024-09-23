package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type WebStatus struct {
	Recording            bool
	RecordedLength       int
	CurrentSong          string
	Starttime            string
	Length               string
	SecondsUntilNextSong int
	EndTime              string
	EstimatedProgress    int
	NewestRecording      string
}

type Settings struct {
	DoRecord              bool
	StreamingURL          string
	SongNameURL           string
	ComparisonSnippetTime int
	SimilarityThreshold   float64
	TargetSampleRate      int
}

func startWeb() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.Static("/recordings", "./recordings")

	//r.StaticFS("/recordings", http.Dir("recodings"))
	// Route to display the form

	r.GET("/", func(c *gin.Context) {

		Status := WebStatus{
			Recording:            IsRecording,
			RecordedLength:       getRecordedTime(),
			CurrentSong:          CurrentSong.Title + " - " + CurrentSong.Author,
			Starttime:            CurrentSong.Starttime,
			Length:               CurrentSong.Duration,
			SecondsUntilNextSong: getSecondsUntilNextSongs(*CurrentSong) + 7,
			EndTime:              getEndTime(*CurrentSong).Format("2006-01-02 15:04:05"),
			EstimatedProgress:    getProgress(),
			NewestRecording:      CurrentRecordingFilename,
		}

		c.HTML(http.StatusOK, "status.html", Status)
	})

	r.GET("/settings", func(c *gin.Context) {

		Settings := Settings{
			DoRecord:              IsRecording,
			SongNameURL:           SongNameURL,
			StreamingURL:          StreamURL,
			SimilarityThreshold:   SimilarityThreshold,
			ComparisonSnippetTime: int(ComparisonSnippetTime / time.Second),
			TargetSampleRate:      targetSampleRate,
		}
		c.HTML(http.StatusOK, "settings.html", Settings)
	})
	// // Route to handle form submission
	r.POST("/settings", func(c *gin.Context) {
		songNameURL := c.PostForm("songnameURL")
		doRecord := c.PostForm("dorecord")
		streamingURL := c.PostForm("StreamURL")
		similarityThreshold := c.PostForm("SimilarityThreshold")
		comparisonSnippetTime := c.PostForm("ComparisonSnippetTime")
		TargetSampleRate := c.PostForm("targetSampleRate")

		SongNameURL = songNameURL
		StreamURL = streamingURL

		//convert that thing because of thing
		i, err := strconv.ParseFloat(similarityThreshold, 2)
		if err != nil {
			// ... handle error
			//c.HTML(http.StatusBadRequest, "", nil)
			//return
		} else {
			SimilarityThreshold = i
		}

		d, err := strconv.Atoi(comparisonSnippetTime)
		if err != nil {
			// ... handle error
			c.HTML(http.StatusBadRequest, "", nil)
			return
		} else {
			ComparisonSnippetTime = time.Duration(d) * time.Second
		}

		a, err := strconv.Atoi(TargetSampleRate)
		if err != nil {
			// ... handle error
			c.HTML(http.StatusBadRequest, "", nil)
			return
		} else {
			targetSampleRate = a
		}

		// we do this so we don't get blocked, otherwise it could happen that it hangs
		if doRecord == "on" {
			go startRecording()
		} else {
			go stopRecording()
		}
		// set new settings to variables
		log.Println("Updated settings from webui!")

		c.Redirect(http.StatusFound, "/")
	})
	// Run the server
	r.Run(":6969")
}

func getRecordedTime() int {
	var duration time.Duration
	var endTime time.Time
	var startTime time.Time
	if TimeStoppedRecording.IsZero() {
		endTime = time.Now()
	} else {
		endTime = TimeStoppedRecording
	}

	if TimeStartedRecording.IsZero() { // This is here so we dont get any weird numbers in the ui, instead it will display as 00:00:00
		startTime = time.Now()
	} else {
		startTime = TimeStartedRecording
	}

	duration = endTime.Sub(startTime)
	//hours := int(duration.Hours())
	//minutes := int(duration.Minutes()) % 60
	//seconds := int(duration.Seconds()) % 60
	return int(duration.Seconds())
}

func getProgress() int {

	duration := getRecordedTime()

	progress := (float64(duration) / EstimatedRuntime.Seconds()) * 100

	if progress > 100 {
		progress = 100
	}

	return int(progress)
}
