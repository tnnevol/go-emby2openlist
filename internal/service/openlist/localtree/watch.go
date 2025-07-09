package localtree

import (
	"fmt"
	"time"

	"github.com/radovskyb/watcher"
)

// WatcherInterval 本地目录树监控轮询间隔
const WatcherInterval = time.Second

// watchDirChange 监听目录变更, 触发回调函数
func watchDirChange(dirPath string, handler func(e watcher.Event)) error {
	w := watcher.New()

	w.SetMaxEvents(1)
	w.IgnoreHiddenFiles(true)
	w.FilterOps(watcher.Create, watcher.Remove, watcher.Rename, watcher.Write, watcher.Move)

	go func() {
		for e := range w.Event {
			handler(e)
		}
	}()

	if err := w.AddRecursive(dirPath); err != nil {
		return fmt.Errorf("目录监听失败 [%s]: %w", dirPath, err)
	}

	if err := w.Start(WatcherInterval); err != nil {
		return fmt.Errorf("无法开启监听 [%s]: %w", dirPath, err)
	}
	return nil
}
