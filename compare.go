package main

import (
	"log"
)

// Function to compare two songs from []byte buffers
func compareSongs(fp1 []int32, fp2 []int32, threshold float64) bool {
	s := compareFingerprint(fp1, fp2)
	log.Println("Correlation:", s)

	return s > threshold
}

func compareFingerprint(fp1, fp2 []int32) float64 {
	return 1 // 100% match
	//s, _ := fingerprint.Compare(fp1, fp2)
	//return s
}

func calculateFingerprint(data []byte, targetSampleRate int, ComparisonSnippetTime int) []int32 {
	return nil
	//reader := bytes.NewReader(data)
	//fpcalc := chromaprint.New(chromaprint.AlgorithmDefault)
	//defer fpcalc.Close()
	//fprint1, _ := fpcalc.RawFingerprint( // TODO: calculate fingerprint 1 only once at the start of recording, it dosen't make any sense to fingerprint it again, but i'm to lazy for now | Lowikian
	//	fingerprint.RawInfo{
	//		Src:        reader,
	//		Channels:   2,
	//		Rate:       uint(targetSampleRate),
	//		MaxSeconds: uint(ComparisonSnippetTime),
	//	})
	//
	//return fprint1
}
