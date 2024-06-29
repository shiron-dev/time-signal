package audio

import (
	"embed"
	_ "embed"
)

//go:embed voice/voice.txt
var VoiceTxt []byte

//go:embed voice/himari/*
var Himari embed.FS
