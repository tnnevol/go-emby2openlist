package ffmpeg

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	durationReg = regexp.MustCompile(`Duration: (\d+):(\d+):(\d+\.\d+)`)
	albumReg    = regexp.MustCompile(`(?mi)album\s*:\s*(.+?)\s*$`)
	artistReg   = regexp.MustCompile(`(?mi)artist\s*:\s*(.+?)\s*$`)
	commentReg  = regexp.MustCompile(`(?mi)comment\s*:\s*(.+?)\s*$`)
	dateReg     = regexp.MustCompile(`(?mi)date\s*:\s*(.+?)\s*$`)
	titleReg    = regexp.MustCompile(`(?mi)title\s*:\s*(.+?)\s*$`)
	titleTiReg  = regexp.MustCompile(`(?mi):\s*\[ti:(.*?)\]\s*$`)
	trackReg    = regexp.MustCompile(`(?mi)track\s*:\s*(.+?)\s*$`)
	genreReg    = regexp.MustCompile(`(?mi)genre\s*:\s*(.+?)\s*$`)
	tdorReg     = regexp.MustCompile(`(?mi)tdor\s*:\s*(.+?)\s*$`)
	lyricsReg   = regexp.MustCompile(`(?mi):\s*(\[.*?\].*?)\s*$`)
)

// resolveDuration 解析 ffmpeg 的 Duration 参数
func resolveDuration(raw string) time.Duration {
	if !durationReg.MatchString(raw) {
		return 0
	}

	res := durationReg.FindStringSubmatch(raw)
	if len(res) != 4 {
		return 0
	}

	hour, _ := strconv.Atoi(res[1])
	minute, _ := strconv.Atoi(res[2])
	second, _ := strconv.ParseFloat(res[3], 64)
	return time.Hour*time.Duration(hour) +
		time.Minute*time.Duration(minute) +
		time.Duration(float64(time.Second)*second)
}

// resolveLyrics 解析 ffmpeg 的 Lyrics 参数
func resolveLyrics(raw string) string {
	if !lyricsReg.MatchString(raw) {
		return ""
	}

	sb := strings.Builder{}

	results := lyricsReg.FindAllStringSubmatch(raw, -1)
	for i, result := range results {
		sb.WriteString(result[1])
		if i < len(results)-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// resolveTrack 解析 ffmpeg 的 Track 参数
func resolveTrack(raw string) int {
	if !trackReg.MatchString(raw) {
		return 0
	}

	track := trackReg.FindStringSubmatch(raw)[1]
	segIdx := strings.Index(track, "/")
	if segIdx != -1 {
		track = track[:segIdx]
	}

	trackNum, _ := strconv.Atoi(track)
	return trackNum

}

// resolveTitle 解析 ffmpeg 的 title 参数
func resolveTitle(raw string) string {
	// 移除 Duration 字段之后的信息, 防止匹配到干扰字段
	if durationReg.MatchString(raw) {
		loc := durationReg.FindStringIndex(raw)
		if len(loc) > 0 {
			raw = raw[:loc[1]]
		}
	}

	// 优先匹配 title
	if titleReg.MatchString(raw) {
		res := strings.TrimSpace(titleReg.FindStringSubmatch(raw)[1])
		if res != "" {
			return res
		}
	}

	// title 为空则空歌词中的 ti 属性提取
	if titleTiReg.MatchString(raw) {
		return strings.TrimSpace(titleTiReg.FindStringSubmatch(raw)[1])
	}

	return ""
}
