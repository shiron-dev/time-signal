package main

import (
	"time"

	"github.com/shiron-dev/time-signal/cmd"
)

const (
	silenceSeconds = 0.05
)

func main() {

	cmd.Timer(OnSecondChange, SpeakTime)
}

func OnSecondChange(seconds int) {
	cmd.WriteWavStdout(cmd.PlayBeep(seconds%10 == 0))
}

func SpeakTime(time time.Time) {
	wavs := cmd.TimeToWav(time)
	wav, err := cmd.CombineWavData(wavs, silenceSeconds)
	if err != nil {
		panic(err)
	}
	cmd.WriteWavStdout(wav)
}
