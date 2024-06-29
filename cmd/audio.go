package cmd

import (
	"bytes"
	"encoding/binary"
	"math"
	"os"
)

type WavData struct {
	Data        []byte
	SampleRate  int
	NumChannels int
}

func CombineWavData(wavs []WavData, silenceSeconds float64) (WavData, error) {
	var combinedData []byte
	var sampleRate, numChannels int

	for _, wav := range wavs {
		if sampleRate == 0 {
			sampleRate = wav.SampleRate
			numChannels = wav.NumChannels
		}

		combinedData = append(combinedData, wav.Data...)

		silenceData := generateSilenceData(sampleRate, numChannels, silenceSeconds)
		combinedData = append(combinedData, silenceData...)
	}

	combinedData = combinedData[:len(combinedData)-len(generateSilenceData(sampleRate, numChannels, silenceSeconds))]

	return WavData{
		Data:        combinedData,
		SampleRate:  sampleRate,
		NumChannels: numChannels,
	}, nil
}

func ReadWavBytes(data []byte) (WavData, error) {
	header := data[:44]
	data = data[44:]

	sampleRate := int(binary.LittleEndian.Uint32(header[24:28]))
	numChannels := int(binary.LittleEndian.Uint16(header[22:24]))

	return WavData{
		Data:        data,
		SampleRate:  sampleRate,
		NumChannels: numChannels,
	}, nil
}

func generateSilenceData(sampleRate, numChannels int, durationSeconds float64) []byte {
	samples := int(float64(sampleRate) * durationSeconds)
	sampleSize := 2 // 16-bit audio
	totalBytes := samples * numChannels * sampleSize
	return make([]byte, totalBytes)
}

func WriteWavFile(filePath string, wav WavData) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	header := createWavHeader(len(wav.Data), wav.SampleRate, wav.NumChannels)
	if _, err := file.Write(header); err != nil {
		return err
	}

	if _, err := file.Write(wav.Data); err != nil {
		return err
	}

	return nil
}

func WriteWavStdout(wav WavData) error {
	header := createWavHeader(len(wav.Data), wav.SampleRate, wav.NumChannels)
	_, err := os.Stdout.Write(header)
	if err != nil {
		return err
	}

	_, err = os.Stdout.Write(wav.Data)
	return err
}

func createWavHeader(dataLength, sampleRate, numChannels int) []byte {
	totalLength := dataLength + 36
	header := &bytes.Buffer{}

	// ChunkID
	header.WriteString("RIFF")
	// ChunkSize
	binary.Write(header, binary.LittleEndian, uint32(totalLength))
	// Format
	header.WriteString("WAVE")
	// Subchunk1ID
	header.WriteString("fmt ")
	// Subchunk1Size
	binary.Write(header, binary.LittleEndian, uint32(16))
	// AudioFormat
	binary.Write(header, binary.LittleEndian, uint16(1))
	// NumChannels
	binary.Write(header, binary.LittleEndian, uint16(numChannels))
	// SampleRate
	binary.Write(header, binary.LittleEndian, uint32(sampleRate))
	// ByteRate
	byteRate := sampleRate * numChannels * 2
	binary.Write(header, binary.LittleEndian, uint32(byteRate))
	// BlockAlign
	blockAlign := numChannels * 2
	binary.Write(header, binary.LittleEndian, uint16(blockAlign))
	// BitsPerSample
	binary.Write(header, binary.LittleEndian, uint16(16))
	// Subchunk2ID
	header.WriteString("data")
	// Subchunk2Size
	binary.Write(header, binary.LittleEndian, uint32(dataLength))

	return header.Bytes()
}

// sin(A3) or sin(A4)
func PlayBeep(isPon bool) WavData {
	duration := 0.2
	if isPon {
		duration = 0.4
	}
	frequency := 440.0
	if isPon {
		frequency = 880.0
	}
	const (
		sampleRate  = 44100
		numChannels = 1
	)

	samples := int(sampleRate * duration)
	sampleSize := 2 // 16-bit audio
	totalBytes := samples * numChannels * sampleSize
	data := make([]byte, totalBytes)

	for i := 0; i < samples; i++ {
		value := 0.5 * float64(1<<15-1) * math.Sin(2*math.Pi*frequency*float64(i)/sampleRate)
		binary.LittleEndian.PutUint16(data[i*2:], uint16(value))
	}

	return WavData{
		Data:        data,
		SampleRate:  sampleRate,
		NumChannels: numChannels,
	}
}
