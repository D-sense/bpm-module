package audio

import (
	"bytes"
	"errors"
	"github.com/h2non/filetype"
	errors2 "github.com/pkg/errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

type BpmAudio struct {
   Filepath string
}

func (b BpmAudio) ExtractBpm(url string) (float64, bool, error) {
	//download track
	filePath, err := b.downloadTrack(url)
	if err != nil {
		if err != io.EOF{
			return 0, false,  errors2.WithMessage(err, "could not download track")
		}
	}

	file, err := os.Open(filePath)
	if err != nil {
		err = clearDir(filePath)
		if err != nil {
			return 0, false, errors2.WithMessage(err, "could not delete file: %v")
		}
		return 0, false, errors2.WithMessage(err, "could not open a downloaded file")
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		err = clearDir(filePath)
		if err != nil {
			return 0, false, errors2.WithMessage(err, "could not delete file: %v")
		}
		return 0, false, errors2.WithMessage(err, "could not convert file to bytes")
	}

	// Supported file types: images
	if !filetype.IsAudio(fileBytes) {
		err = clearDir(filePath)
		if err != nil {
			return 0, false, errors2.WithMessage(err, "could not delete file: %v")
		}
		return 0, false, errors.New("unsupported file type (not an audio)")
	}

	// extract BPM
	bpmResult, err := b.ffmpegBpm(filePath)
	if err != nil {
		err = clearDir(filePath)
		if err != nil {
			return 0, false, errors2.WithMessage(err, "could not delete file: %v")
		}
		return 0, false, errors2.WithMessage(err, "could not extract BPM: %v",)
	}

	// delete file
	err = clearDir(filePath)
	if err != nil {
		return 0, false, errors2.WithMessage(err, "could not delete file: %v")
	}

	return bpmResult, true, nil
}

func (b BpmAudio) ffmpegBpm(file string) (float64, error){
	cmdArguments := "ffmpeg -v quiet -i "+ file +" -f f32le -ac 1 -c:a pcm_f32le -ar 44100 pipe:1 | bpm"
	cmd := exec.Command("bash","-c", cmdArguments)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return 0, err
	}

	result, err  := strconv.ParseFloat(strings.Replace(out.String(),"\n", "", -1 ),32)
	if err != nil {
		return 0, err
	}

	if result == 146 || result == 0 {
		return 0, errors.New("invalid bpm number (file may be corrupt)")
	}

	return result, nil
}

func (b BpmAudio) downloadTrack(url string) (string, error) {
	err := createDir(setPath(b.Filepath))
	if err != nil {
		return "", errors2.WithMessage(err, "error while creating directory")
	}

	fileName := getFileName(url)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	out, err := os.Create(b.Filepath+fileName)
	if err != nil {
		return "", errors2.WithMessage(err, "error while putting file in a directory")
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)

	return b.Filepath+fileName, err
}

func clearDir(dir string) error {
	fileName := getFileName(dir)
	folder := setPath(dir)

	names, err := ioutil.ReadDir(folder)
	if err != nil {
		return err
	}

	for _, entry := range names {
		if entry.Name() == fileName{
			_ = os.RemoveAll(path.Join([]string{folder, entry.Name()}...))
		}
	}

	return nil
}

func createDir(name string) error {
	dir := setPath(name)

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	return nil
}

func setPath(path string) string {
	return  path[:strings.LastIndexByte(path, '/')+1]
}

func getFileName(path string) string {
	return  path[strings.LastIndexByte(path, '/')+1:]
}

func TimeSpent(start time.Time, name string, success, failure int64) {
	elapsed := time.Since(start)
	log.Printf("%s took %s | %v passed | %v failed", name, elapsed, success, failure)
}
