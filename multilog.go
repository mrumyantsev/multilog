package multilog

import (
	"os"
	"time"
)

type LogWriter interface {
	WriteLog(datetime *string, messageType *string, message *string) error
}

const (
	errorWord string = ". error: "
)

var (
	logWriters   map[string]LogWriter = nil
	logWriterErr error

	InfoOutputStream    *os.File = os.Stderr
	DebugOutputStream   *os.File = os.Stderr
	WarnOutputStream    *os.File = os.Stderr
	ErrorOutputStream   *os.File = os.Stderr
	FatalOutputStream   *os.File = os.Stderr
	FatalExitStatusCode int      = 1
	IsEnableDebugLogs   bool     = true
	IsEnableWarnLogs    bool     = true
	ItemSeparator       string   = " "
	LineEnding          string   = "\n"
	InfoMessageType     string   = "INF"
	DebugMessageType    string   = "DBG"
	WarnMessageType     string   = "WRN"
	ErrorMessageType    string   = "ERR"
	FatalMessageType    string   = "FTL"
	TimeFormat          string   = "2006-01-02T15:04:05-07:00"
)

func RegisterWriter(name string, writer LogWriter) {
	if logWriters == nil {
		logWriters = map[string]LogWriter{}
	}

	logWriters[name] = writer
}

func UnregisterWriter(name string) {
	delete(logWriters, name)
}

func Info(msg string) {
	datetime := time.Now().Format(TimeFormat)

	writeToStream(
		&datetime,
		&InfoMessageType,
		&msg,
		InfoOutputStream,
	)

	if logWriters != nil {
		writeToLogWriters(
			&datetime,
			&FatalMessageType,
			&msg,
		)
	}
}

func Debug(msg string) {
	if !IsEnableDebugLogs {
		return
	}

	datetime := time.Now().Format(TimeFormat)

	writeToStream(
		&datetime,
		&DebugMessageType,
		&msg,
		DebugOutputStream,
	)

	if logWriters != nil {
		writeToLogWriters(
			&datetime,
			&FatalMessageType,
			&msg,
		)
	}
}

func Warn(msg string) {
	if !IsEnableWarnLogs {
		return
	}

	datetime := time.Now().Format(TimeFormat)

	writeToStream(
		&datetime,
		&WarnMessageType,
		&msg,
		WarnOutputStream,
	)

	if logWriters != nil {
		writeToLogWriters(
			&datetime,
			&FatalMessageType,
			&msg,
		)
	}
}

func Error(desc string, err error) {
	datetime := time.Now().Format(TimeFormat)

	desc = desc + errorWord + err.Error()

	writeToStream(
		&datetime,
		&ErrorMessageType,
		&desc,
		ErrorOutputStream,
	)

	if logWriters != nil {
		writeToLogWriters(
			&datetime,
			&FatalMessageType,
			&desc,
		)
	}
}

func Fatal(desc string, err error) {
	datetime := time.Now().Format(TimeFormat)

	desc = desc + errorWord + err.Error()

	writeToStream(
		&datetime,
		&FatalMessageType,
		&desc,
		FatalOutputStream,
	)

	if logWriters != nil {
		writeToLogWriters(
			&datetime,
			&FatalMessageType,
			&desc,
		)
	}

	os.Exit(FatalExitStatusCode)
}

func writeToStream(
	datetime *string,
	messageType *string,
	message *string,
	stream *os.File,
) {
	stream.Write(
		[]byte(
			*datetime +
				ItemSeparator +
				*messageType +
				ItemSeparator +
				*message +
				LineEnding,
		),
	)
}

func writeToLogWriters(
	datetime *string,
	messageType *string,
	message *string,
) {
	for name, writer := range logWriters {
		logWriterErr = writer.WriteLog(
			datetime,
			messageType,
			message,
		)

		if logWriterErr != nil {
			Error("could not write to log writer \""+name+"\"", logWriterErr)
		}
	}
}
