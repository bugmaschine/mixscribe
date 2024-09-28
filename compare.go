package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math"
	"math/cmplx"

	mp3 "github.com/tcolgate/mp3"
)

var (
	binSize = 512
)

func compareSongs(initialFingerprint, rollingSeconds [][]float64, SimilarityThreshold float64) bool {
	if compareFingerprint(initialFingerprint, rollingSeconds) >= SimilarityThreshold {
		return true
	} else {
		return false
	}
}

// why is audio stuff always so complex?? (hopefully ai knows what it's doing because i don't)
func compareFingerprint(fp1, fp2 [][]float64) float64 {

	if len(fp1) != len(fp2) || len(fp1[0]) != len(fp2[0]) {
		log.Println("fp1" + string(len(fp1)))
		log.Println("fp2" + string(len(fp2)))
		log.Println("fp1 0" + string(len(fp1[0])))
		log.Println("fp2 0" + string(len(fp2[0])))
		panic("Fingerprints have different dimensions")
	}

	dotProduct := 0.0
	for i := 0; i < len(fp1); i++ {
		for j := 0; j < len(fp1[0]); j++ {
			dotProduct += fp1[i][j] * fp2[i][j]
		}
	}

	// calculate the magnitudes of the two fingerprints
	magnitude1 := 0.0
	magnitude2 := 0.0
	for i := 0; i < len(fp1); i++ {
		for j := 0; j < len(fp1[0]); j++ {
			magnitude1 += fp1[i][j] * fp1[i][j]
			magnitude2 += fp2[i][j] * fp2[i][j]
		}
	}

	// Calculate the correlation coefficient
	correlation := dotProduct / (math.Sqrt(magnitude1) * math.Sqrt(magnitude2))

	// debug
	log.Printf("Correlation: %f\n", correlation)

	return correlation
}

func calculateFingerprint(data []byte, targetSampleRate int, ComparisonSnippetTime int) ([][]float64, error) {
	convertedAudio, err := audiobytesToFloat64(data)
	if err != nil {
		panic(err) //normaly this shouldn't happen, only if the internet cuts out or if the radio station has problems
	}
	// check if sample is smaller than bins
	if len(convertedAudio) < binSize {
		return make([][]float64, 0), fmt.Errorf("audio is smaller than binSize")
	}
	fft, err := drawFFT(convertedAudio, binSize)
	if err != nil {
		return make([][]float64, 0), err
	}

	// debug log
	//log.Printf("Fingerprint: %v\n", fft)

	return fft, nil
}

// code taken from https://github.com/corny/spectrogram/blob/master/

func drawFFT(samples []float64, bins int) ([][]float64, error) {
	fftResults := make([][]float64, bins)
	for i := range fftResults {
		fftResults[i] = make([]float64, bins)
	}

	sub := make([]float64, bins*2)

	for x := 0; x < bins; x++ {
		n0 := int64(mapRange(float64(x+0), 0, float64(bins), 0, float64(len(samples))))
		n1 := int64(mapRange(float64(x+1), 0, float64(bins), 0, float64(len(samples))))
		c := int(n0 + (n1-n0)/2)

		for i := range sub {
			s := 0.0
			n := c - bins + i
			if n >= 0 && n < len(samples) {
				s = samples[n]
			}

			// Apply Hamming window
			s *= 0.54 - 0.46*math.Cos(float64(i)*math.Pi*2/float64(len(sub)))

			sub[i] = s
		}

		freqs := fft(sub)

		for y := 0; y < bins; y++ {
			r := cmplx.Abs(freqs[y])
			fftResults[x][y] = r
		}
	}

	return fftResults, nil
}

func hfft(samples []float64, freqs []complex128, n, step int) {
	if n == 1 {
		freqs[0] = complex(samples[0], 0)
		return
	}

	half := n / 2

	hfft(samples, freqs, half, 2*step)
	hfft(samples[step:], freqs[half:], half, 2*step)

	for k := 0; k < half; k++ {
		a := -2 * math.Pi * float64(k) / float64(n)
		e := cmplx.Rect(1, a) * freqs[k+half]

		freqs[k], freqs[k+half] = freqs[k]+e, freqs[k]-e
	}
}

func fft(samples []float64) []complex128 {
	n := len(samples)
	freqs := make([]complex128, n)
	hfft(samples, freqs, n, 1)
	return freqs
}

func mapRange(n, srcMin, srcMax, dstMin, dstMax float64) float64 {
	return (n-srcMin)/(srcMax-srcMin)*(dstMax-dstMin) + dstMin
}

// remember that we're not dealing with an normal mp3 file and instead with a stream
// https://github.com/tcolgate/mp3
func audiobytesToFloat64(data []byte) ([]float64, error) {
	// Create a new reader for the MP3 data
	reader := bytes.NewReader(data)

	// Initialize MP3 decoder
	decoder := mp3.NewDecoder(reader)

	var samples []float64
	frame := mp3.Frame{}
	skipped := 0

	for {
		// Decode the next MP3 frame
		err := decoder.Decode(&frame, &skipped)
		if err != nil {
			// Handle EOF gracefully
			if err == io.EOF {
				break
			}
			// Handle unexpected EOF or partial frames more gracefully
			if err.Error() == "unexpected EOF" {
				fmt.Println("Warning: encountered unexpected EOF, skipping frame")
				continue
			}
			return nil, fmt.Errorf("failed to decode MP3 data: %v", err)
		}

		// Get PCM data from the frame
		frameData := frame.Reader()

		// Create a buffer to read the frame's PCM data
		buf := make([]byte, 4096)
		for {
			n, err := frameData.Read(buf)
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, fmt.Errorf("failed to read frame data: %v", err)
			}

			// Convert PCM bytes (assuming 16-bit signed samples) to float64 samples
			for i := 0; i < n; i += 2 {
				if i+1 < n {
					// Combine two bytes into a 16-bit signed integer
					sample := int16(binary.LittleEndian.Uint16(buf[i:]))
					// Normalize to float64 [-1, 1]
					floatSample := float64(sample) / 32768.0
					samples = append(samples, floatSample)
				}
			}
		}
	}

	return samples, nil
}
