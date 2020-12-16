package log

import (
	"fmt"
	"os"
	"reflect"
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
	appName        string
	fileNameFormat = "2006_01_02"
	logPath        = ""
	loggerMgr      = loggerManager{logger: make(map[string]*zerolog.Logger)}
)

func Init() bool {
	path := strings.Split(os.Args[0], "/")
	appName = path[len(path)-1]

	if len(appName) == 0 {
		return false
	}

	l := loggerMgr.createLogger(appName)
	log.Logger = *l

	// time format
	// RFC3339 "2006-01-02 15:04:05.000"
	zerolog.TimeFieldFormat = time.RFC3339

	return true
}

func Info() *zerolog.Event {
	return loggerMgr.getSetLogger(appName).Info()
}

func NamedInfo(loggerName string) *zerolog.Event {
	return loggerMgr.getSetLogger(loggerName).Info()
}

func EasyInfo(v ...interface{}) *zerolog.Event {
	l := len(v)
	if l == 0 {
		return loggerMgr.getSetLogger(appName).Info()
	}

	if reflect.TypeOf(v[0]).Kind() != reflect.String {
		loggerMgr.getSetLogger(appName).Fatal().Msg("easy info first param isnt string ")
	}

	name := reflect.ValueOf(v[0]).String()

	return loggerMgr.getSetLogger(name).Info()
}

func Debug() *zerolog.Event {
	return loggerMgr.getSetLogger(appName).Debug()
}

func NamedDebug(loggerName string) *zerolog.Event {
	return loggerMgr.getSetLogger(loggerName).Debug()
}

func EasyDebug(v ...interface{}) *zerolog.Event {
	l := len(v)
	if l == 0 {
		return loggerMgr.getSetLogger(appName).Debug()
	}

	if reflect.TypeOf(v[0]).Kind() != reflect.String {
		loggerMgr.getSetLogger(appName).Fatal().Msg("easy debug first param isnt string ")
	}

	name := reflect.ValueOf(v[0]).String()

	return loggerMgr.getSetLogger(name).Debug()
}

func Error() *zerolog.Event {
	return loggerMgr.getSetLogger(appName).Error()
}

func NamedError(loggerName string) *zerolog.Event {
	return loggerMgr.getSetLogger(loggerName).Error()
}

func EasyError(v ...interface{}) *zerolog.Event {
	l := len(v)
	if l == 0 {
		return loggerMgr.getSetLogger(appName).Error()
	}

	if reflect.TypeOf(v[0]).Kind() != reflect.String {
		loggerMgr.getSetLogger(appName).Fatal().Msg("easy error first param isnt string ")
	}

	name := reflect.ValueOf(v[0]).String()

	return loggerMgr.getSetLogger(name).Error()
}

func Warn() *zerolog.Event {
	return loggerMgr.getSetLogger(appName).Warn()
}

func NamedWarn(loggerName string) *zerolog.Event {
	return loggerMgr.getSetLogger(loggerName).Warn()
}

func EasyWarn(v ...interface{}) *zerolog.Event {
	l := len(v)
	if l == 0 {
		return loggerMgr.getSetLogger(appName).Warn()
	}

	if reflect.TypeOf(v[0]).Kind() != reflect.String {
		loggerMgr.getSetLogger(appName).Fatal().Msg("easy warn first param isnt string ")
	}

	name := reflect.ValueOf(v[0]).String()

	return loggerMgr.getSetLogger(name).Warn()
}

func Fatal() *zerolog.Event {
	return loggerMgr.getSetLogger(appName).Fatal()
}

func NamedFatal(loggerName string) *zerolog.Event {
	return loggerMgr.getSetLogger(loggerName).Fatal()
}

func EasyFatal(v ...interface{}) *zerolog.Event {
	l := len(v)
	if l == 0 {
		return loggerMgr.getSetLogger(appName).Fatal()
	}

	if reflect.TypeOf(v[0]).Kind() != reflect.String {
		loggerMgr.getSetLogger(appName).Fatal().Msg("easy fatal first param isnt string ")
	}

	name := reflect.ValueOf(v[0]).String()

	return loggerMgr.getSetLogger(name).Fatal()
}

func Print(v ...interface{}) {
	loggerMgr.getSetLogger(appName).Print(v...)
}

func EasyPrint(loggerName string, v ...interface{}) {
	l := len(loggerName)
	if l == 0 {
		loggerMgr.getSetLogger(appName).Print(v...)
		return
	}

	loggerMgr.getSetLogger(loggerName).Print(v...)
}

func Printf(format string, v ...interface{}) {
	loggerMgr.getSetLogger(appName).Printf(format, v...)
}

func EasyPrintf(loggerName string, format string, v ...interface{}) {
	l := len(loggerName)
	if l == 0 {
		loggerMgr.getSetLogger(appName).Printf(format, v...)
		return
	}

	loggerMgr.getSetLogger(loggerName).Printf(format, v...)
}

func EnableStdOut() {
	loggerMgr.enableStdOut()
}

func DisableStdOut() {
	loggerMgr.disableStdOut()
}

func SetLevel(lvl zerolog.Level, v ...interface{}) {
	loggerName := appName

	if len(v) > 0 {
		if reflect.TypeOf(v[0]).Kind() != reflect.String {
			loggerMgr.getSetLogger(appName).Error().Msg("SetLevel first param isnt string ")
			return
		}

		loggerName = reflect.ValueOf(v[0]).String()
	}

	logger := loggerMgr.getLogger(loggerName)
	if logger == nil {
		loggerMgr.getSetLogger(appName).Error().Str("loggerName", loggerName).Msg("SetLevel dont find the logger")
		return
	}

	newLogger := logger.Level(zerolog.Level(lvl))

	loggerMgr.setLogger(loggerName, &newLogger)
}

func AddHook(h zerolog.Hook, v ...interface{}) {
	loggerName := appName

	if len(v) > 0 {
		if reflect.TypeOf(v[0]).Kind() != reflect.String {
			loggerMgr.getSetLogger(appName).Error().Msg("AddHook first param isnt string ")
			return
		}

		loggerName = reflect.ValueOf(v[0]).String()
	}

	logger := loggerMgr.getLogger(loggerName)
	if logger == nil {
		loggerMgr.getSetLogger(appName).Error().Str("loggerName", loggerName).Msg("AddHook dont find the logger")
		return
	}

	newLogger := logger.Hook(h)

	loggerMgr.setLogger(loggerName, &newLogger)
}

func Close() {
	// fileWriter.Close()
}

type writer struct {
	mu   sync.Mutex
	f    *os.File
	name string
}

func (w *writer) Write(b []byte) (n int, err error) {
	if loggerMgr.toStdOut() {
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
	name := fmt.Sprintf("%s%s_%s.log", logPath, w.name, fileSuffix)
	file, err := os.OpenFile(name, os.O_RDWR|os.O_APPEND|os.O_CREATE, fileMod)
	if err != nil {
		log.Error().Err(err).Str("logger_name", w.name).Msg("create file faild")
	}

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
	fileName := fmt.Sprintf("%s_%s.log", w.name, t)
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}

type loggerManager struct {
	mu     sync.RWMutex
	stdout bool
	logger map[string]*zerolog.Logger
}

func (lm *loggerManager) enableStdOut() {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	lm.stdout = true
}

func (lm *loggerManager) disableStdOut() {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	lm.stdout = false
}

func (lm *loggerManager) toStdOut() bool {
	lm.mu.RLock()
	defer lm.mu.RUnlock()

	return lm.stdout
}

func (lm *loggerManager) getSetLogger(loggerName string) *zerolog.Logger {
	if len(loggerName) == 0 {
		log.Fatal().Str("logger_name", loggerName).Msg("get logger failed,must point logger name!")
	}

	lm.mu.RLock()

	l, ok := lm.logger[loggerName]
	if !ok {
		lm.mu.RUnlock()
		return lm.createLogger(loggerName)
	}

	lm.mu.RUnlock()
	return l
}

func (lm *loggerManager) getLogger(loggerName string) *zerolog.Logger {
	if len(loggerName) == 0 {
		return nil
	}

	lm.mu.RLock()
	defer lm.mu.RUnlock()

	l, ok := lm.logger[loggerName]
	if !ok {
		return nil
	}

	return l
}

func (lm *loggerManager) setLogger(name string, logger *zerolog.Logger) {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	lm.logger[name] = logger
}

func (lm *loggerManager) createLogger(name string) *zerolog.Logger {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	l, ok := lm.logger[name]
	if ok {
		return l
	}

	w := writer{name: name}

	f := w.createFile(time.Now())
	if f == nil {
		log.Fatal().Str("logger_name", name).Msg("create logger failed")
	}

	w.f = f

	logger := log.Output(&w)

	lm.logger[name] = &logger

	// check
	go w.checkFileChange()
	go w.checkFileExists()

	// log
	log.Info().Str("logger", name).Msg("create logger success!")

	return &logger
}
