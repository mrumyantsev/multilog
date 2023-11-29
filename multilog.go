package multilog

import (
	"os"
	"time"
)

type LogWriter interface {
	WriteLog(datetime string, messageType string, message string) error
}

type logMessage struct {
	datetime    string
	messageType *string
	message     *string
}

const (
	errorWord string = ". error: "
)

var (
	writers           *map[string]LogWriter = nil
	writer            LogWriter             = nil
	writersBeforeExit int                   = 0
	isWriterExists    bool                  = false
	isFatalLog        bool                  = false
	errMsg            string
	logMsg            *logMessage = &logMessage{}

	InfoOutputStream        *os.File = os.Stderr
	DebugOutputStream       *os.File = os.Stderr
	ErrorOutputStream       *os.File = os.Stderr
	FatalOutputStream       *os.File = os.Stderr
	ExitStatusCodeWhenFatal int      = 1
	IsEnableDebugLogs       bool     = true
	ItemSeparator           string   = " "
	LineEnding              string   = "\n"
	InfoMessageType         string   = "INF"
	DebugMessageType        string   = "DBG"
	ErrorMessageType        string   = "ERR"
	FatalMessageType        string   = "FTL"
	TimeFormat              string   = "2006-01-02T15:04:05-07:00"
)

func RegisterWriter(name string, w LogWriter) {
	if writers == nil {
		writers = &map[string]LogWriter{}
	}

	(*writers)[name] = w

	writersBeforeExit++
}

func UnregisterWriter(name string) {
	delete((*writers), name)

	writersBeforeExit--
}

func Info(msg string) *logMessage {
	logMsg.datetime = time.Now().Format(TimeFormat)
	logMsg.messageType = &InfoMessageType
	logMsg.message = &msg

	writeToStream(InfoOutputStream)

	return logMsg
}

func Debug(msg string) *logMessage {
	if !IsEnableDebugLogs {
		return nil
	}

	logMsg.datetime = time.Now().Format(TimeFormat)
	logMsg.messageType = &DebugMessageType
	logMsg.message = &msg

	writeToStream(DebugOutputStream)

	return logMsg
}

func Error(desc string, err error) *logMessage {
	errMsg = desc + errorWord + err.Error()

	logMsg.datetime = time.Now().Format(TimeFormat)
	logMsg.messageType = &ErrorMessageType
	logMsg.message = &errMsg

	writeToStream(ErrorOutputStream)

	return logMsg
}

func Fatal(desc string, err error) *logMessage {
	errMsg = desc + errorWord + err.Error()

	logMsg.datetime = time.Now().Format(TimeFormat)
	logMsg.messageType = &FatalMessageType
	logMsg.message = &errMsg

	writeToStream(FatalOutputStream)

	if writersBeforeExit <= 0 {
		os.Exit(ExitStatusCodeWhenFatal)
	}

	isFatalLog = true

	return logMsg
}

func writeToStream(stream *os.File) {
	stream.Write(
		[]byte(logMsg.datetime +
			ItemSeparator +
			*logMsg.messageType +
			ItemSeparator +
			*logMsg.message +
			LineEnding,
		),
	)
}

func (l *logMessage) WriteTo(writerName string) *logMessage {
	if (l == nil) || (writers == nil) {
		return nil
	}

	writer, isWriterExists = (*writers)[writerName]
	if !isWriterExists {
		return l
	}

	writer.WriteLog(l.datetime, *l.messageType, *l.message)

	if isFatalLog {
		writersBeforeExit--

		if writersBeforeExit <= 0 {
			os.Exit(ExitStatusCodeWhenFatal)
		}
	}

	return l
}
