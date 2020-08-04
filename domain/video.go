package domain

import "mime/multipart"

//OutputPreset ...
type OutputPreset struct {
	Format      string
	Compression string
	VideoCodec  string
	AudioCodec  string
}

//Formats ...
type Formats struct {
	Formats []Format
}

//Format ...
type Format struct {
	Extension string
	Video     Video
	Audio     Audio
}

//Video ...
type Video struct {
	Default string
	Others  []string
}

//Audio ...
type Audio struct {
	Default string
	Others  []string
}

//VideoUsecase ...
type VideoUsecase interface {
	Convert(file *multipart.FileHeader, op OutputPreset) (downloadPath string, err error)
}
