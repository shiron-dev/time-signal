package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/shiron-dev/time-signal/audio"
)

const oshirase = "お知らせします"

func TimeToWav(time time.Time) []WavData {
	voiceTxt := timeToVoiceTxt(time)
	fmt.Println(voiceTxt)
	voice := make([]WavData, len(voiceTxt))
	for i, v := range voiceTxt {
		filePath, err := voiceTxtToFilePath(v)
		fmt.Println(filePath)
		if err != nil {
			continue
		}
		content, err := audio.Himari.ReadFile("voice/himari/" + filePath + ".wav")
		if err != nil {
			panic(err)
		}
		wav, err := ReadWavBytes(content)
		if err != nil {
			panic(err)
		}
		voice[i] = wav
	}

	return voice
}

func voiceTxtToFilePath(voice string) (string, error) {
	txt := string(audio.VoiceTxt)
	lines := strings.Split(txt, "\n")
	for i, l := range lines {
		if l == voice {
			return fmt.Sprintf("%03d", i+1), nil
		}
	}
	return "", fmt.Errorf("voice not found: %s", voice)
}

func timeToVoiceTxt(time time.Time) []string {
	if time.Hour() == 12 && time.Minute() == 0 {
		return []string{
			"正午を",
			oshirase,
		}
	}

	kubun := "午前"
	if time.Hour() > 12 {
		kubun = "午後"
	}

	nji := strconv.Itoa(time.Hour()%12) + "時"

	zyuu := strconv.Itoa(time.Minute()/10) + "0"
	if zyuu == "0" {
		zyuu = ""
	}

	nfun := strconv.Itoa(time.Minute()%10) + "分"

	if time.Minute()%10 == 0 {
		zyuu = ""
		nfun = strconv.Itoa(time.Minute()) + "分"
	}

	nbyou := strconv.Itoa(time.Second()) + "秒を"
	if time.Second() == 0 {
		nbyou = "ちょうどを"
	}

	return lo.Filter([]string{
		kubun,
		nji,
		zyuu,
		nfun,
		nbyou,
		oshirase,
	}, func(s string, i int) bool {
		return s != ""
	})
}
