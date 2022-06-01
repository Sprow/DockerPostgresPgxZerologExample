package logger

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"log"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
)

var Log *Logger

type Logger struct {
	*zerolog.Logger
}

func Stack(err error) error {
	//TODO delete first element from stack
	return errors.New(err.Error())
}

func Init() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel) // <--- change level

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	var logsPath = fmt.Sprintf("%v/internal/logger/log.txt", dir)
	logFile, err := os.Create(logsPath)
	if err != nil {
		log.Fatal(err)
	}

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
	multi := zerolog.MultiLevelWriter(consoleWriter, logFile)
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	zlog := zerolog.New(multi).With().Timestamp().Stack().Logger()
	zlog = zlog.Hook(CallerHook{})

	Log = &Logger{
		Logger: &zlog,
	}
}

type CallerHook struct{}

func (h CallerHook) Run(event *zerolog.Event, level zerolog.Level, msg string) {
	if _, file, line, ok := runtime.Caller(3); ok {
		var file = path.Base(file)
		var line = line
		filename := file[strings.LastIndex(file, "/")+1:] + ":" + strconv.Itoa(line)
		event.Str("caller", filename)
	}
}
