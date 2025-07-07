package ffmpeg

import "time"

// Info 记录文件元信息
type Info struct {
	Duration time.Duration
}

// Music 记录音乐元信息
type Music struct {
	Info
	Album   string // 专辑
	Artist  string // 艺术家
	Comment string // 备注
	Date    string // 发布日期
	Lyrics  string // 歌词
	Title   string // 标题
	Track   int    // 轨道
	Genre   string // 流派
}
