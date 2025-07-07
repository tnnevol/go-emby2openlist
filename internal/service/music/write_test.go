package music_test

import (
	"log"
	"testing"

	"github.com/AmbitiousJun/go-emby2openlist/v2/internal/service/music"
)

func TestWriteNFO(t *testing.T) {
	rmt := "http://0.0.0.0:123456/d/%E9%9F%B3%E4%B9%90/%E9%99%88%E6%A5%9A%E7%94%9F%20-%20%E7%88%B1%E7%9A%84%E5%B0%BD%E5%A4%B4.flac?sign=x3Vazh5gSUYn6fuuYCZoey6J3AViM0XPpGvPUsOalYA=:0"
	meta, err := music.ExtractRemoteTag(rmt)
	if err != nil {
		t.Fatal(err)
		return
	}

	log.Println(meta)
}
