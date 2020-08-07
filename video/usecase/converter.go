package usecase

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gitlab.com/prosa-backend-test/video_converter/domain"
)

type videoUsecase struct {
}

//NewVideoUsecase ...
func NewVideoUsecase() domain.VideoUsecase {
	return &videoUsecase{}
}

func save(file *multipart.FileHeader) (err error) {
	src, err := file.Open()
	if err != nil {
		log.Println(err)
		return err
	}
	defer src.Close()

	dst, err := os.Create(os.TempDir() + "/" + file.Filename)
	if err != nil {
		log.Println(err)
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		log.Println(err)
		return err
	}
	return
}

func generateOutputFilename(format string, inputFile string) (outputFilename string) {
	toHash := fmt.Sprintf("%s_%s", inputFile, string(time.Now().UnixNano()))
	hashed := sha256.Sum256([]byte(toHash))
	if format == "keep" {
		format = filepath.Ext(inputFile)
	}
	outputFilename = fmt.Sprintf("%s___%x%s", randomString(5), hashed, format)

	return
}

func randomString(n int) (generated string) {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	length := n
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	generated = b.String() // E.g. "e8Yqr"
	return
}

func getAVCodec(format string, inVidCodec string, inAudCodec string) (vidCodec string, audCodec string, err error) {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	file, err := os.Open(path + "/config.json")
	if err != nil {
		log.Println(err)
		return "", "", err
	}
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
		return "", "", err
	}
	var formats domain.Formats

	json.Unmarshal(byteValue, &formats)

	for i := 0; i < len(formats.Formats); i++ {
		if formats.Formats[i].Extension == format {
			if inVidCodec == "default" {
				vidCodec = formats.Formats[i].Video.Default
				break
			} else {
				for j := 0; j < len(formats.Formats[i].Video.Others); j++ {
					if inVidCodec == formats.Formats[i].Video.Others[j] {
						vidCodec = formats.Formats[i].Video.Others[j]
						break
					} else {
						err = errors.New("invalid video codec options")
					}
				}
			}
		}
	}

	for i := 0; i < len(formats.Formats); i++ {
		if formats.Formats[i].Extension == format {
			if inAudCodec == "default" {
				audCodec = formats.Formats[i].Audio.Default
				break
			} else {
				for j := 0; j < len(formats.Formats[i].Audio.Others); j++ {
					if inAudCodec == formats.Formats[i].Audio.Others[j] {
						audCodec = formats.Formats[i].Audio.Others[j]
						break
					} else {
						err = errors.New("invalid audio codec options")
					}
				}
			}
		}
	}

	return
}

func getCompressRate(rate string, filepath string) (bitrate string, err error) {
	cmd := exec.Command("ffprobe",
		"-i", filepath,
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "stream=bit_rate",
		"-of", "default=noprint_wrappers=1")
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Error running command: %v", err)
		return "", err
	}
	trimmed := strings.TrimSuffix(string(output), "\n")
	cmdResult := strings.Split(trimmed, "=")
	srcBitrate, err := strconv.Atoi(cmdResult[1])
	if err != nil {
		log.Println(err)
		return "", err
	}
	fRate, err := strconv.ParseFloat(rate, 64)
	if err != nil {
		log.Println(err)
		return "", err
	}

	bitrate = fmt.Sprintf("%.2f", float64(srcBitrate)*fRate)
	return
}

func convert(inputFile, vidCodec, audCodec, compressRate, outputFile string, chanErr chan error) {
	cmd := exec.Command("ffmpeg",
		"-i", inputFile,
		"-c:v", vidCodec,
		"-c:a", audCodec,
		"-b:v", compressRate,
		outputFile)
	chanErr <- cmd.Start()
	if chanErr != nil {
		chanErr <- cmd.Wait()
	}
}

func (uc *videoUsecase) Convert(file *multipart.FileHeader, op domain.OutputPreset) (downloadPath string, err error) {
	filepath := os.TempDir() + "/" + file.Filename
	if err = save(file); err != nil {
		log.Println(err)
		return "", err
	}
	vidCodec, audCodec, err := getAVCodec(op.Format, op.VideoCodec, op.AudioCodec)
	if err != nil {
		return "", err
	}
	compressRate, err := getCompressRate(op.Compression, filepath)
	if err != nil {
		return "", err
	}

	outputFilename := generateOutputFilename(op.Format, filepath)

	outputFile := "files/" + outputFilename

	chanErr := make(chan error)
	go convert(filepath, vidCodec, audCodec, compressRate, outputFile, chanErr)
	err = <-chanErr
	if err != nil {
		log.Printf("Error running command: %v", <-chanErr)
		return "", err
	}

	downloadPath = "static/" + outputFilename

	return
}
