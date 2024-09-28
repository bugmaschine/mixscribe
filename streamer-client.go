package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	recordingDir  = "./recordings/"
	filePrefix    = "Jump-"
	fileExtension = ".mp3"
)

// StartRecording starts the recording process based on ControlRecording signal
func StartRecording() {
	for <-ControlRecording {
		currentTime := time.Now().Format("2006-01-02_15:04:05")

		filename := filePrefix + currentTime + fileExtension
		CurrentRecordingFilename = filename
		filename = recordingDir + filename
		radioMP3URL, err := getRadioURLFromM3U(StreamURL)
		if err != nil {
			log.Fatalf("Error getting radio URL: %v", err)
		}

		log.Printf("Parsed radio url %v", radioMP3URL)

		log.Printf("Started recording at %v", TimeStartedRecording)

		err = streamDownloadAndCompare(radioMP3URL, filename)
		if err != nil {
			log.Fatalf("Error during stream download: %v", err)
		}
	}
}

// getRadioURLFromM3U retrieves the actual radio URL from the M3U file
func getRadioURLFromM3U(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 && (line[0] == 'h' || line[0] == 'H') {
			return line, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", fmt.Errorf("no URL found in the M3U file")
}

func streamDownloadAndCompare(url, filePath string) error {
	// Observeration: during mix restarts there is about a 3 second time without any sound

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			log.Printf("Followed Redirect to %v", req.URL)
			return nil // Follow all redirects
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	initialSegment := make([]byte, 0, 20*1024*1024) // 20MB buffer for the first 10 seconds
	currentChunk := make([]byte, 1024)
	rollingSeconds := make([]byte, 0, 20*1024*1024) // 20MB buffer for rolling seconds
	finishedSnippetRecording := false
	rollingSecondsFull := false
	var initialFingerprint [][]float64

	for IsRecording {
		n, err := resp.Body.Read(currentChunk)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			IsRecording = false
			break
		}

		log.Printf("Working on chunk ending at byte %v", n)
		fmt.Println("Time:", uint(int(ComparisonSnippetTime.Seconds())))
		if !finishedSnippetRecording {
			initialSegment = append(initialSegment, currentChunk[:n]...)

			if time.Since(TimeStartedRecording).Seconds() >= ComparisonSnippetTime.Seconds() { // Since we're adding currentChunk[:n]s every 100ms
				finishedSnippetRecording = true
				log.Println("Finished recording comparison segment!")
				// calculate fingerprint only once to save performance
				initialFingerprint, err = calculateFingerprint(initialSegment, targetSampleRate, int(ComparisonSnippetTime.Seconds()))
				if err != nil {
					panic(err)
				}
			}
		} else {
			rollingSeconds = append(rollingSeconds, currentChunk[:n]...)

			if len(rollingSeconds) > len(initialSegment) {
				rollingSeconds = rollingSeconds[len(rollingSeconds)-len(initialSegment):]
				log.Printf("\nRolling seconds full.\n")
				rollingSecondsFull = true
			} else {
				rollingSecondsFull = false
			}
			if rollingSecondsFull {
				rollingSecondsFP, err := calculateFingerprint(rollingSeconds, targetSampleRate, int(ComparisonSnippetTime.Seconds()))
				if err != nil {
					panic(err)
				}

				if compareSongs(initialFingerprint, rollingSecondsFP, SimilarityThreshold) { // Seems to be slow // maybe somehow make it multi threaded
					log.Printf("\nSimilar Segment found, stopping recording!")
					IsRecording = false
					TimeStoppedRecording = time.Now()
					break
				}
			}
			log.Printf("\nToMuchLength: %v", len(rollingSeconds)-len(initialSegment)) //Debug logging

		}

		log.Printf("\n")
		if _, err := file.Write(currentChunk[:n]); err != nil {
			return err
		}
	}

	return nil
}

func startRecording() {
	TimeStartedRecording = time.Now()
	TimeStoppedRecording = time.Time{} // reset it to be empty
	IsRecording = true
	ControlRecording <- true
}

func stopRecording() {
	TimeStoppedRecording = time.Now()
	IsRecording = false
	ControlRecording <- false
}
