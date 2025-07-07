package ffmpeg

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

// OpenError ffmpeg 打开文件失败
const OpenError = "Error opening input:"

// mu 任务逐个执行
var mu sync.Mutex

// InspectInfo 检查指定路径文件的元信息
func InspectInfo(path string) (Info, error) {
	if !execOk {
		return Info{}, errors.New("ffmpeg 未初始化")
	}
	mu.Lock()
	defer mu.Unlock()

	cmd := exec.Command(execPath, "-i", path)
	var eb bytes.Buffer
	cmd.Stderr = &eb
	cmd.Run()

	output := eb.String()
	if strings.Contains(output, OpenError) {
		return Info{}, errors.New(output[strings.Index(output, OpenError):])
	}

	i := Info{}
	if durationReg.MatchString(output) {
		i.Duration = resolveDuration(output)
	}

	return i, nil
}

// InspectMusic 检查指定音乐文件的元信息
func InspectMusic(path string) (Music, error) {
	if !execOk {
		return Music{}, errors.New("ffmpeg 未初始化")
	}
	mu.Lock()
	defer mu.Unlock()

	cmd := exec.Command(execPath, "-i", path)
	var eb bytes.Buffer
	cmd.Stderr = &eb
	cmd.Run()

	output := eb.String()
	if strings.Contains(output, OpenError) {
		return Music{}, errors.New(output[strings.Index(output, OpenError):])
	}

	m := Music{}
	wg := sync.WaitGroup{}
	wg.Add(10)

	go func() {
		defer wg.Done()
		if albumReg.MatchString(output) {
			m.Album = albumReg.FindStringSubmatch(output)[1]
		}
	}()

	go func() {
		defer wg.Done()
		if artistReg.MatchString(output) {
			m.Artist = artistReg.FindStringSubmatch(output)[1]
		}
	}()

	go func() {
		defer wg.Done()
		if commentReg.MatchString(output) {
			m.Comment = commentReg.FindStringSubmatch(output)[1]
		}
	}()

	go func() {
		defer wg.Done()
		if dateReg.MatchString(output) {
			m.Date = dateReg.FindStringSubmatch(output)[1]
		}
	}()

	go func() {
		defer wg.Done()
		if tdorReg.MatchString(output) {
			m.Date = tdorReg.FindStringSubmatch(output)[1]
		}
	}()

	go func() {
		defer wg.Done()
		if genreReg.MatchString(output) {
			m.Genre = genreReg.FindStringSubmatch(output)[1]
		}
	}()

	go func() {
		defer wg.Done()
		m.Title = resolveTitle(output)
	}()

	go func() {
		defer wg.Done()
		m.Duration = resolveDuration(output)
	}()

	go func() {
		defer wg.Done()
		m.Track = resolveTrack(output)
	}()

	go func() {
		defer wg.Done()
		m.Lyrics = resolveLyrics(output)
	}()

	wg.Wait()
	return m, nil
}

// ExtractMusicCover 解析音乐海报
func ExtractMusicCover(path string) ([]byte, error) {
	if !execOk {
		return nil, errors.New("ffmpeg 未初始化")
	}
	mu.Lock()
	defer mu.Unlock()

	cmd := exec.Command(execPath, "-i", path, "-an", "-vcodec", "copy", "-f", "image2pipe", "pipe:1")

	var out bytes.Buffer
	var eb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &eb
	cmd.Run()

	if strings.Contains(eb.String(), OpenError) {
		return nil, errors.New(eb.String()[strings.Index(eb.String(), OpenError):])
	}

	return out.Bytes(), nil
}

// GenSilentMP3Bytes 使用 ffmpeg 生成静音 MP3 并返回字节内容
func GenSilentMP3Bytes(durationSec float64) ([]byte, error) {
	args := []string{
		"-f", "lavfi",
		"-i", "anullsrc=r=8000:cl=mono",
		"-t", fmt.Sprintf("%.2f", durationSec),
		"-acodec", "libmp3lame",
		"-b:a", "8k", // 极低比特率
		"-ar", "8000", // 采样率降为 8000Hz
		"-ac", "1", // 单声道
		"-f", "mp3",
		"pipe:1",
	}

	cmd := exec.Command(execPath, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = nil
	cmd.Run()

	return out.Bytes(), nil
}
