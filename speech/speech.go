package speech

import (
	"context"
	"fmt"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	texttospeechpb "cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
)

// Voice name constants for Japanese language
const (
	VoiceStandardA = "ja-JP-Standard-A"
	VoiceStandardB = "ja-JP-Standard-B"
	VoiceStandardC = "ja-JP-Standard-C"
	VoiceStandardD = "ja-JP-Standard-D"
	VoiceWavenetA  = "ja-JP-Wavenet-A"
	VoiceWavenetB  = "ja-JP-Wavenet-B"
	VoiceWavenetC  = "ja-JP-Wavenet-C"
	VoiceWavenetD  = "ja-JP-Wavenet-D"
)

// Audio encoding constants
const (
	AudioEncoding_LINEAR16 = "LINEAR16"
	AudioEncoding_MP3      = "MP3"
	AudioEncoding_OGG_OPUS = "OGG_OPUS"
)

// SpeechOption contains configuration options for text-to-speech synthesis
type SpeechOption struct {
	LanguageCode      string
	VoiceName         string
	AudioEncoding     string
	AudioSpeakingRate float64
	AudioPitch        float64
}

// SpeechClient is a wrapper for the Google Cloud Text-to-Speech client
type SpeechClient struct {
	client *texttospeech.Client
}

// SpeechRequest represents a request to the Text-to-Speech API
type SpeechRequest struct {
	Text string
	Opt  *SpeechOption
}

// NewSpeechClient creates a new SpeechClient
func NewSpeechClient(ctx context.Context) (*SpeechClient, error) {
	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create texttospeech client: %v", err)
	}
	return &SpeechClient{client: client}, nil
}

// NewRequest creates a new SpeechRequest
func NewRequest(text string, opt *SpeechOption) *SpeechRequest {
	return &SpeechRequest{
		Text: text,
		Opt:  opt,
	}
}

// Run sends a synthesis request to the Text-to-Speech API and returns the audio content
func (s *SpeechClient) Run(ctx context.Context, req *SpeechRequest) ([]byte, error) {
	// Create the text-to-speech request
	request := &texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{
				Text: req.Text,
			},
		},
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: req.Opt.LanguageCode,
			Name:         req.Opt.VoiceName,
			SsmlGender:   texttospeechpb.SsmlVoiceGender_NEUTRAL,
		},
		AudioConfig: &texttospeechpb.AudioConfig{
			SpeakingRate: req.Opt.AudioSpeakingRate,
			Pitch:        req.Opt.AudioPitch,
		},
	}

	// Set the audio encoding based on the option
	switch req.Opt.AudioEncoding {
	case AudioEncoding_LINEAR16:
		request.AudioConfig.AudioEncoding = texttospeechpb.AudioEncoding_LINEAR16
	case AudioEncoding_MP3:
		request.AudioConfig.AudioEncoding = texttospeechpb.AudioEncoding_MP3
	case AudioEncoding_OGG_OPUS:
		request.AudioConfig.AudioEncoding = texttospeechpb.AudioEncoding_OGG_OPUS
	default:
		return nil, fmt.Errorf("unsupported audio encoding: %s", req.Opt.AudioEncoding)
	}

	// Call the API
	response, err := s.client.SynthesizeSpeech(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to synthesize speech: %v", err)
	}

	return response.AudioContent, nil
}

