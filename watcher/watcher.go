package watcher

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/mitchellh/go-ps"
)

const (
	PROCESSLEN = 16
	Change     = 0
	Pid        = 1
)

type Event struct {
	Name string
}

type EventPid struct {
	Pid int
}

type Watcher struct {
	files      map[string]time.Time
	Event      chan Event
	EventPid   chan EventPid
	Errors     chan error
	mu         sync.Mutex
	paths      string
	extension  string
	executable string
	pid        int
}

func NewWatcher() *Watcher {
	w := &Watcher{
		Event:    make(chan Event),
		EventPid: make(chan EventPid),
		Errors:   make(chan error),
	}
	go w.handlerChange()
	go w.handlePid()
	return w
}

func (w *Watcher) NewPath(path string) *Watcher {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.paths = path
	return (w)
}

func (w *Watcher) NewExtension(ext string) *Watcher {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.extension = ext
	return (w)
}

func (w *Watcher) NewExecutable(exec string) *Watcher {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.executable = exec
	return (w)
}

func (w *Watcher) handlePid() {
	for {
		if w.executable != "" && w.pid == 0 {
			process, err := ps.Processes()
			if err != nil {
				w.Errors <- err
			}
			for _, proc := range process {
				tmp := ""
				if len(w.executable) > PROCESSLEN {
					tmp = w.executable[:PROCESSLEN]
				} else {
					tmp = w.executable
				}
				if proc.Executable() == tmp {
					w.mu.Lock()
					w.pid = proc.Pid()
					w.mu.Unlock()
					w.EventPid <- EventPid{Pid: proc.Pid()}
					break
				}
			}
		} else {
			w.mu.Lock()
			w.pid = 0
			w.mu.Unlock()
		}
		time.Sleep(time.Duration(500 * time.Millisecond))
	}
}

func (w *Watcher) handlerChange() {
	for {
		if w.paths != "" && w.extension != "" {
			files := make(map[string]time.Time)

			err := filepath.Walk(w.paths, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if !info.IsDir() && filepath.Ext(path) == w.extension {
					files[path] = info.ModTime()
					time1, ok1 := files[path]
					w.mu.Lock()
					time2, ok2 := w.files[path]
					w.mu.Unlock()
					if ok1 && ok2 && !time1.Equal(time2) {
						w.Event <- Event{Name: path}
					}
				}
				return nil
			})
			w.mu.Lock()
			w.files = files
			w.mu.Unlock()
			if err != nil {
				w.Errors <- err
			}
			time.Sleep(time.Duration(500 * time.Millisecond))
		}
	}
}
