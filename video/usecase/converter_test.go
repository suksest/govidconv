package usecase

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/prosa-backend-test/video_converter/domain"
)

func TestNewVideoUsecase(t *testing.T) {
	tests := []struct {
		name string
		want domain.VideoUsecase
	}{
		// TODO: Add test cases.
		{
			name: "new video usecase",
			want: &videoUsecase{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewVideoUsecase(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewVideoUsecase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getSampleRequest() (req *http.Request, err error) {
	file, err := os.Open("../../testdata/001.mp4")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var reqBody bytes.Buffer
	multipartWriter := multipart.NewWriter(&reqBody)

	fileWriter, err := multipartWriter.CreateFormFile("file", "001.mp4")
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return nil, err
	}
	multipartWriter.Close()

	req = httptest.NewRequest("POST", "http://localhost:5000/video/convert", &reqBody)
	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())
	return
}

func Test_save(t *testing.T) {
	req, err := getSampleRequest()
	if err != nil {
		return
	}
	_, file, err := req.FormFile("file")
	if err != nil {
		return
	}

	type args struct {
		file *multipart.FileHeader
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "valid file header",
			args: args{
				file: file,
			},
			wantErr: false,
		},
		{
			name: "invalid file header",
			args: args{
				file: &multipart.FileHeader{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := save(tt.args.file); (err != nil) != tt.wantErr {
				t.Errorf("save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_generateOutputFilename(t *testing.T) {
	type args struct {
		format    string
		inputFile string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "test generated output filename with different extension",
			args: args{
				format:    ".mp4",
				inputFile: "Golang.mkv",
			},
		},
		{
			name: "test generated output filename with keep extension",
			args: args{
				format:    "keep",
				inputFile: "Golang.mkv",
			},
		},
		{
			name: "test generated output filename with same extension",
			args: args{
				format:    ".mp4",
				inputFile: "Golang.mp4",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOutputFilename := generateOutputFilename(tt.args.format, tt.args.inputFile)
			ok := assert.Regexp(t, regexp.MustCompile("^[a-zA-Z0-9]{5}___[a-zA-Z0-9]{64}.(mkv|mp4|mov)$"), gotOutputFilename)
			if !ok {
				t.Error("Filename pattern not matched")
			}
		})
	}
}

func Test_randomString(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "test string length",
			args: args{5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotGenerated := randomString(tt.args.n)
			ok := assert.Regexp(t, regexp.MustCompile("^([a-zA-Z0-9]){5}$"), gotGenerated)
			if !ok {
				t.Error("String pattern not match")
			}
		})
	}
}

func Test_getAVCodec(t *testing.T) {
	type args struct {
		configFilePath string
		format         string
		inVidCodec     string
		inAudCodec     string
	}
	tests := []struct {
		name         string
		args         args
		wantVidCodec string
		wantAudCodec string
		wantErr      bool
	}{
		// TODO: Add test cases.
		{
			name: "config file not exists",
			args: args{
				configFilePath: "notexist.json",
				format:         ".mp4",
				inVidCodec:     "theora",
				inAudCodec:     "mp3",
			},
			wantVidCodec: "",
			wantAudCodec: "",
			wantErr:      true,
		},
		{
			name: "invalid video codec options",
			args: args{
				configFilePath: "config.json",
				format:         ".mp4",
				inVidCodec:     "theora",
				inAudCodec:     "mp3",
			},
			wantVidCodec: "",
			wantAudCodec: "mp3",
			wantErr:      true,
		},
		{
			name: "invalid audio codec options",
			args: args{
				configFilePath: "config.json",
				format:         ".mp4",
				inVidCodec:     "h265",
				inAudCodec:     "flac",
			},
			wantVidCodec: "h265",
			wantAudCodec: "",
			wantErr:      true,
		},
		{
			name: "valid options",
			args: args{
				configFilePath: "config.json",
				format:         ".mkv",
				inVidCodec:     "theora",
				inAudCodec:     "mp3",
			},
			wantVidCodec: "theora",
			wantAudCodec: "mp3",
			wantErr:      false,
		},
		{
			name: "default options",
			args: args{
				configFilePath: "config.json",
				format:         ".mp4",
				inVidCodec:     "default",
				inAudCodec:     "default",
			},
			wantVidCodec: "h264",
			wantAudCodec: "aac",
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVidCodec, gotAudCodec, err := getAVCodec(tt.args.configFilePath, tt.args.format, tt.args.inVidCodec, tt.args.inAudCodec)
			if gotVidCodec != tt.wantVidCodec {
				t.Errorf("getAVCodec() gotVidCodec = %v, want %v", gotVidCodec, tt.wantVidCodec)
			}
			if gotAudCodec != tt.wantAudCodec {
				t.Errorf("getAVCodec() gotAudCodec = %v, want %v", gotAudCodec, tt.wantAudCodec)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("getAVCodec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_getCompressRate(t *testing.T) {
	type args struct {
		rate     string
		filepath string
	}
	tests := []struct {
		name        string
		args        args
		wantBitrate string
		wantErr     bool
	}{
		// TODO: Add test cases.
		{
			name: "file not exists",
			args: args{
				rate:     "0.5",
				filepath: "../../testdata/notexist.mp4",
			},
			wantBitrate: "",
			wantErr:     true,
		},
		{
			name: "correct bit rate",
			args: args{
				rate:     "0.5",
				filepath: "../../testdata/001.mp4",
			},
			wantBitrate: "129380.50",
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBitrate, err := getCompressRate(tt.args.rate, tt.args.filepath)
			if gotBitrate != tt.wantBitrate {
				t.Errorf("getCompressRate() = %v, want %v", gotBitrate, tt.wantBitrate)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("getCompressRate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_convert(t *testing.T) {
	type args struct {
		inputFile    string
		vidCodec     string
		audCodec     string
		compressRate string
		outputFile   string
		chanErr      chan error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		/*{
			name: "valid file",
			args: args{
				inputFile:    "../../testdata/001.mp4",
				vidCodec:     "h264",
				audCodec:     "aac",
				compressRate: "129380.50",
				outputFile:   "../../testdata/output.mp4",
			},
			wantErr: false,
		},
		{
			name: "invalid file",
			args: args{
				inputFile:    "../../testdata/notexist.mp4",
				vidCodec:     "h264",
				audCodec:     "aac",
				compressRate: "129380.50",
				outputFile:   "../../testdata/output.mp4",
			},
			wantErr: true,
		},*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			convert(tt.args.inputFile, tt.args.vidCodec, tt.args.audCodec, tt.args.compressRate, tt.args.outputFile, tt.args.chanErr)
		})
	}
}

func Test_Convert(t *testing.T) {
	req, err := getSampleRequest()
	if err != nil {
		return
	}
	_, file, err := req.FormFile("file")
	if err != nil {
		return
	}
	type args struct {
		file *multipart.FileHeader
		op   domain.OutputPreset
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "compress file",
			args: args{
				file: file,
				op: domain.OutputPreset{
					Format:      "keep",
					Compression: "0.9",
					VideoCodec:  "default",
					AudioCodec:  "default",
				},
			},
			wantErr: false,
		},
		{
			name: "different format with default compression and codec",
			args: args{
				file: file,
				op: domain.OutputPreset{
					Format:      ".mkv",
					Compression: "1.0",
					VideoCodec:  "default",
					AudioCodec:  "default",
				},
			},
			wantErr: false,
		},
		{
			name: "different format, default compression, different video codec, default audio codec",
			args: args{
				file: file,
				op: domain.OutputPreset{
					Format:      ".mkv",
					Compression: "1.0",
					VideoCodec:  "theora",
					AudioCodec:  "default",
				},
			},
			wantErr: false,
		},
		{
			name: "different format, default compression, different video codec, different audio codec",
			args: args{
				file: file,
				op: domain.OutputPreset{
					Format:      ".mkv",
					Compression: "1.0",
					VideoCodec:  "theora",
					AudioCodec:  "flac",
				},
			},
			wantErr: false,
		},
		{
			name: "different format, different compression, different video codec, different audio codec",
			args: args{
				file: file,
				op: domain.OutputPreset{
					Format:      ".mkv",
					Compression: "0.8",
					VideoCodec:  "theora",
					AudioCodec:  "flac",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &videoUsecase{}
			gotDownloadPath, err := uc.Convert(tt.args.file, tt.args.op)
			ok := assert.Regexp(t, regexp.MustCompile("^static/[a-zA-Z0-9]{5}___[a-zA-Z0-9]{64}.(mkv|mp4|mov)$"), gotDownloadPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("videoUsecase.Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !ok {
				t.Error("dowloadPath pattern not matched")
			}
		})
	}
}
