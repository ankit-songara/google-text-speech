package cmd

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ankit-songara/google-text-speech/speech"
)

const (
	ExitCodeOK              = 0
	ExitCodeParseFlagsError = 1
	ExitCodeValidateError   = 2
	ExitCodeInternalError   = 3
	ExitCodeOutputFileError = 4
)

type CLI struct {
	ErrStream io.Writer
}

func (cli *CLI) Run(args []string) int {
	flags := flag.NewFlagSet("google-text-to-speech", flag.ContinueOnError)
	var (
		text, voice, out string
		rate, pitch      float64
	)

	flags.StringVar(&text, "text", "", "Text to synthesize")
	flags.StringVar(&voice, "voice", "ja-JP-Standard-A", "Voice name")
	flags.Float64Var(&rate, "rate", 1.00, "speech rate (0.25 ~ 4.00)")
	flags.Float64Var(&pitch, "pitch", 0.00, "speech pitch (-20.00 ~ 20.00)")
	flags.StringVar(&out, "o", "", "output audio file (support format of the audio: :LINEAR16, MP3, OGG_OPUS)")

	if err := flags.Parse(args[1:]); err != nil {
		fmt.Fprintf(cli.ErrStream, "%v", err)
		return ExitCodeParseFlagsError
	}

	opt, err := makeSpeechOpt(text, voice, out, rate, pitch)
	if err != nil {
		fmt.Fprintf(cli.ErrStream, "%v", err)
		return ExitCodeValidateError
	}

	ctx := context.Background()
	speaker, err := speech.NewSpeechClient(ctx)
	if err != nil {
		fmt.Fprintf(cli.ErrStream, "failed to create speech client: %v\n", err)
		return ExitCodeInternalError
	}
	b, err := speaker.Run(ctx, speech.NewRequest(text, opt))
	if err != nil {
		fmt.Fprintf(cli.ErrStream, "failed to synthesize speech: %v\n", err)
		return ExitCodeInternalError
	}

	if err = os.WriteFile(out, b, 0644); err != nil {
		fmt.Fprintf(cli.ErrStream, "failed to write output file: %v\n", err)
		return ExitCodeOutputFileError
	}
	fmt.Printf("mp3 file created successfully: %s\n", out)
	return ExitCodeOK
}

func makeSpeechOpt(text, voice, out string, rate, pitch float64) (*speech.SpeechOption, error) {
	if text == "" {
		return nil, fmt.Errorf("text is required")
	}

	var voiceName string
	switch v := strings.ToLower(voice); v {
	case "standard-a":
		voiceName = speech.VoiceStandardA
	case "standard-b":
		voiceName = speech.VoiceStandardB
	case "standard-c":
		voiceName = speech.VoiceStandardC
	case "standard-d":
		voiceName = speech.VoiceStandardD
	case "wavenet-a":
		voiceName = speech.VoiceWavenetA
	case "wavenet-b":
		voiceName = speech.VoiceWavenetB
	case "wavenet-c":
		voiceName = speech.VoiceWavenetC
	case "wavenet-d":
		voiceName = speech.VoiceWavenetD
	default:
		return nil, fmt.Errorf("invalid voice name: %s", voice)
	}

	if rate < 0.25 || rate > 4.00 {
		return nil, fmt.Errorf("speech rate must be between 0.25 and 4.00, rate: %f", rate)
	}

	if pitch < -20.00 || pitch > 20.00 {
		return nil, fmt.Errorf("speech pitch must be between -20.00 and 20.00, pitch: %f", pitch)
	}

	switch ext := strings.ToLower(filepath.Ext(out)); ext {
	case ".wav":
		return &speech.SpeechOption{
			LanguageCode:      "ja-JP",
			VoiceName:         voiceName,
			AudioEncoding:     speech.AudioEncoding_LINEAR16,
			AudioSpeakingRate: rate,
			AudioPitch:        pitch,
		}, nil
	case ".mp3":
		return &speech.SpeechOption{
			LanguageCode:      "ja-JP",
			VoiceName:         voiceName,
			AudioEncoding:     speech.AudioEncoding_MP3,
			AudioSpeakingRate: rate,
			AudioPitch:        pitch,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported output file format: %s", ext)
	}

}
