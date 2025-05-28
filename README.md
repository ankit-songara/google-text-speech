# google-text-speech

A command-line tool to synthesize Japanese text to speech using Google Cloud Text-to-Speech API.

## Features

- Supports multiple Japanese voices (Standard and Wavenet)
- Output audio in MP3 or WAV format
- Adjustable speech rate and pitch

## Requirements

- Go 1.24+
- Google Cloud Platform account with Text-to-Speech API enabled
- Service account key JSON file (see below)

## Setup

1. **Clone the repository:**
   ```sh
   git clone https://github.com/ankit-songara/google-text-speech.git
   cd google-text-speech
   ```

2. **Set up Google Cloud credentials:**
   - Create a service account with the "Text-to-Speech Admin" role in Google Cloud Console.
   - Download the service account key as `service-account-key.json`.
   - Set the environment variable:
     ```sh
     export GOOGLE_APPLICATION_CREDENTIALS="path/to/service-account-key.json"
     ```

3. **Install dependencies:**
   ```sh
   go mod tidy
   ```

4. **Build the project:**
   ```sh
   go build -o google-text-speech
   ```

## Usage

```sh
./google-text-speech -text "こんにちは、世界" -voice wavenet-a -rate 1.0 -pitch 0.0 -o output.mp3
```

### Command-line Options

| Flag      | Description                                                                                  | Default                |
|-----------|----------------------------------------------------------------------------------------------|------------------------|
| `-text`   | Text to synthesize (required)                                                                |                        |
| `-voice`  | Voice name: standard-a, standard-b, standard-c, standard-d, wavenet-a, wavenet-b, ...        | ja-JP-Standard-A       |
| `-rate`   | Speech rate (0.25 ~ 4.00)                                                                    | 1.00                   |
| `-pitch`  | Speech pitch (-20.00 ~ 20.00)                                                                | 0.00                   |
| `-o`      | Output audio file (extension determines format: .mp3 or .wav)                                |                        |

### Example

```sh
./google-text-speech -text "おはようございます" -voice wavenet-b -rate 1.2 -pitch 2.0 -o greeting.mp3
```

## Supported Voices

- standard-a
- standard-b
- standard-c
- standard-d
- wavenet-a
- wavenet-b
- wavenet-c
- wavenet-d

## Supported Output Formats

- `.mp3` (MP3)
- `.wav` (LINEAR16)

## License

MIT

---

**Note:**  
This project requires a valid Google Cloud service account key with access to the Text-to-Speech API. Do not commit your credentials to version control.
