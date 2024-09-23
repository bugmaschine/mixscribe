package main

import (
	"bytes"
	"log"

	"github.com/go-fingerprint/fingerprint"
	chromaprint "github.com/go-fingerprint/gochroma"
)

// Function to compare two songs from []byte buffers
func compareSongs(fp1 []int32, fp2 []int32, threshold float64) bool {
	s := compareFingerprint(fp1, fp2)
	log.Println("Correlation:", s)

	return s > threshold
}

func compareFingerprint(fp1, fp2 []int32) float64 {
	s, _ := fingerprint.Compare(fp1, fp2)
	return s
}

func calculateFingerprint(data []byte, targetSampleRate int, ComparisonSnippetTime int) []int32 {
	reader := bytes.NewReader(data)
	fpcalc := chromaprint.New(chromaprint.AlgorithmDefault)
	defer fpcalc.Close()
	fprint1, _ := fpcalc.RawFingerprint( // TODO: calculate fingerprint 1 only once at the start of recording, it dosen't make any sense to fingerprint it again, but i'm to lazy for now | Lowikian
		fingerprint.RawInfo{
			Src:        reader,
			Channels:   2,
			Rate:       uint(targetSampleRate),
			MaxSeconds: uint(ComparisonSnippetTime),
		})

	return fprint1
}
