package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type RadioData struct {
	Resulttype string          `json:"Resulttype"`
	Info       string          `json:"Info"`
	Songs      map[string]Song `json:"Songs"`
}

type Song struct {
	Status           string        `json:"status"`
	IDTitel          string        `json:"id_titel"`
	Title            string        `json:"title"`
	Subtitle         string        `json:"subtitle"`
	Starttime        string        `json:"starttime"`
	Author           string        `json:"author"`
	AvID             string        `json:"av_id"`
	AvNextID         string        `json:"av_next_id"`
	Duration         string        `json:"duration"`
	Interpret        string        `json:"interpret"`
	Kurzinfo         string        `json:"kurzinfo"`
	Metadatentext    string        `json:"metadatentext"`
	InterpretURL     string        `json:"interpret_url"`
	ArtistImageID    ArtistImageID `json:"artist_image_id"`
	Transmissiontype string        `json:"transmissiontype"`
	Audioasset       AudioAsset    `json:"audioasset"`
	Komponist        string        `json:"komponist"`
	Label            string        `json:"label"`
	Tontraeger       string        `json:"tontraeger"`
}

type ArtistImageID struct {
	ImageVariant []ImageVariant `json:"imageVariant"`
}

type ImageVariant struct {
	Attributes ImageAttributes `json:"@attributes"`
}

type ImageAttributes struct {
	Name     string `json:"name"`
	Width    string `json:"width"`
	Height   string `json:"height"`
	MimeType string `json:"mimeType"`
	URL      string `json:"url"`
}

type AudioAsset struct {
	Asset []Asset `json:"asset"`
}

type Asset struct {
}

var SleepTime = 0

func UpdateSongInfoLoop() {
	for {
		UpdateInfo()
		time.Sleep(time.Duration(SleepTime) * time.Second)
	}
}

func UpdateInfo() {
	resp, err := http.Get(SongNameURL)
	if err != nil {
		log.Fatalf("Failed to make the request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read the response body: %v", err)
	}

	var radioData RadioData
	err = json.Unmarshal(body, &radioData)
	if err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	playingSong, err := findCurrentPlayingSong(radioData)
	if err != nil {
		log.Println("Failed to find current song, retrying later")
		return
	}

	CurrentSong = playingSong

	// generate time to sleep to save network bandwith  | Lowikian
	SleepTime = getSecondsUntilNextSongs(*playingSong) + 5
	//log.Printf("Current playing song: %+v\n", CurrentSong)
}

func generateCurrentDateTime() string {
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	return formattedTime
}

func findCurrentPlayingSong(radioData RadioData) (*Song, error) {
	for _, song := range radioData.Songs {
		if song.Status == "now" {
			return &song, nil
		}
	}
	return nil, fmt.Errorf("no current playing song found")
}

func getSecondsUntilNextSongs(song Song) int {
	// Parse the start time
	startTime, err := time.Parse("2006-01-02 15:04:05", song.Starttime)
	if err != nil {
		log.Println("Error parsing start time:", err)
		return 10
	}

	// Parse the duration
	durationParts := strings.Split(song.Duration, ":")
	if len(durationParts) != 3 {
		log.Println("Invalid duration format")
		return 10
	}

	hours, err := strconv.Atoi(durationParts[0])
	if err != nil {
		log.Println("Error parsing hours:", err)
		return 10
	}

	minutes, err := strconv.Atoi(durationParts[1])
	if err != nil {
		log.Println("Error parsing minutes:", err)
		return 10
	}

	seconds, err := strconv.Atoi(durationParts[2])
	if err != nil {
		log.Println("Error parsing seconds:", err)
		return 10
	}

	// all this witchcraft is needed so that the time it calculates is right | Lowikian
	loc, _ := time.LoadLocation("Europe/Berlin")

	// the time the song ends
	endTime := startTime.Add(time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second)

	endTime = time.Date(endTime.Year(), endTime.Month(), endTime.Day(), endTime.Hour(), endTime.Minute(), endTime.Second(), endTime.Nanosecond(), loc)

	now := time.Now().In(loc)

	diff := endTime.Sub(now)

	if diff < 0 {
		return 0
	}

	return int(diff.Seconds())
}

func getEndTime(song Song) time.Time {
	// Parse the start time
	startTime, err := time.Parse("2006-01-02 15:04:05", song.Starttime)
	if err != nil {
		log.Println("Error parsing start time:", err)
		return time.Now()
	}

	// Parse the duration
	durationParts := strings.Split(song.Duration, ":")
	if len(durationParts) != 3 {
		log.Println("Invalid duration format")
		return time.Now()
	}

	hours, err := strconv.Atoi(durationParts[0])
	if err != nil {
		log.Println("Error parsing hours:", err)
		return time.Now()
	}

	minutes, err := strconv.Atoi(durationParts[1])
	if err != nil {
		log.Println("Error parsing minutes:", err)
		return time.Now()
	}

	seconds, err := strconv.Atoi(durationParts[2])
	if err != nil {
		log.Println("Error parsing seconds:", err)
		return time.Now()
	}

	// all this witchcraft is needed so that the time it calculates is right
	loc, _ := time.LoadLocation("Europe/Berlin")

	// the time the song ends
	endTime := startTime.Add(time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second)

	endTime = time.Date(endTime.Year(), endTime.Month(), endTime.Day(), endTime.Hour(), endTime.Minute(), endTime.Second(), endTime.Nanosecond(), loc)

	return endTime
}
