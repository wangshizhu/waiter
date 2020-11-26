package log

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	fileMod = 0666
)

var (
	pathPrefix     string
	fileNameFormat = "2006_01_02"
	logPath        = ""
	fileWriter     writer
)

func Init() bool {
	path := strings.Split(os.Args[0], "/")
	pathPrefix = path[len(path)-1]

	f := fileWriter.createFile(time.Now())
	if f == nil {
		return false
	}

	fileWriter.f = f
	fileWriter.stdout = true

	log.Logger = log.Output(&fileWriter)

	// time format
	// RFC3339 "2006-01-02 15:04:05.000"
	zerolog.TimeFieldFormat = time.RFC3339

	// check
	go fileWriter.checkFileChange()
	go fileWriter.checkFileExists()

	// log
	log.Info().Msg("log init success!")

	return true
}

func Info() *zerolog.Event {
	return log.Info()
}

func Debug() *zerolog.Event {
	return log.Debug()
}

func Error() *zerolog.Event {
	return log.Error()
}

func Warn() *zerolog.Event {
	return log.Warn()
}

func Fatal() *zerolog.Event {
	return log.Fatal()
}

func Print(v ...interface{}) {
	log.Print(v...)
}

func Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func EnableStdOut() {
	fileWriter.stdout = true
}

func DisableStdOut() {
	fileWriter.stdout = false
}

func SetLevel(lvl zerolog.Level) {
	log.Logger = log.Level(zerolog.Level(lvl))
}

// Close 关闭日志系统
func Close() {
	fileWriter.Close()
}

type writer struct {
	mu     sync.Mutex
	f      *os.File
	stdout bool
}

func (w *writer) Write(b []byte) (n int, err error) {
	if w.stdout {
		os.Stderr.Write(b)
	}

	p := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&w.f)))
	f := (*os.File)(p)
	if f == nil {
		return 0, nil
	}

	// w.mu.Lock()
	n, writeErr := f.Write(b)
	// w.mu.Unlock()

	return n, writeErr
}

func (w *writer) Close() {
	p := atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(&w.f)), nil)
	f := (*os.File)(p)
	if f == nil {
		return
	}

	f.Sync()
	f.Close()
}

func (w *writer) createFile(createTime time.Time) *os.File {
	fileSuffix := createTime.Format(fileNameFormat)
	name := fmt.Sprintf("%s%s_%s.log", logPath, pathPrefix, fileSuffix)
	file, _ := os.OpenFile(name, os.O_RDWR|os.O_APPEND|os.O_CREATE, fileMod)
	return file
}

func (w *writer) checkFileChange() {
	for {
		tomorrow := time.Now().Add(time.Hour * 24)
		tomorrow = time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, tomorrow.Location())
		tm := time.NewTimer(tomorrow.Sub(time.Now()))
		select {
		case <-tm.C:
			{
				f := w.createFile(tomorrow)
				if f != nil {
					oldfile := w.f
					atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&w.f)), unsafe.Pointer(f))
					time.Sleep(10 * time.Second)

					oldfile.Sync()
					oldfile.Close()
				}
			}
		}
	}
}

func (w *writer) checkFileExists() {
	for {
		tm := time.NewTimer(time.Second)
		select {
		case <-tm.C:
			{
				if w.existsLogFile() {
					break
				}
				now := time.Now()
				f := w.createFile(now)
				if f != nil {
					atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&w.f)), unsafe.Pointer(f))
				}
			}
		}
	}
}

func (w *writer) existsLogFile() bool {
	now := time.Now()
	t := now.Format(fileNameFormat)
	fileName := fmt.Sprintf("%s_%s.log", pathPrefix, t)
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}
