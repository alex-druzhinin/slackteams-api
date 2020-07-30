package logger

import (
	"bitbucket.org/iwlab-standuply/slackteams-api/config"
	"bitbucket.org/iwlab-standuply/slackteams-api/errors"
	"bitbucket.org/iwlab-standuply/slackteams-api/handler"
	log "github.com/sirupsen/logrus"
)

const (
	LogKeyUserID    string = "userId"
	LogKeyRequestID string = "requestId"
)

func InitLogger(conf config.Config) {
	if conf.Env == config.EnvTypeDevelopment {
		log.SetFormatter(newMyTextFormatter())
		log.SetLevel(log.TraceLevel)
	} else {
		log.SetFormatter(newMyJSONFormatter())
	}
}

type myJSONFormatter struct {
	*log.JSONFormatter
}

func newMyJSONFormatter() *myJSONFormatter {
	return &myJSONFormatter{&log.JSONFormatter{}}
}

func (f *myJSONFormatter) Format(entry *log.Entry) ([]byte, error) {
	if err, has := entry.Data[log.ErrorKey]; has {
		switch err.(type) {
		case error:
			entry.Data[log.ErrorKey] = log.Fields{
				"message": err.(error).Error(),
				"stack":   errors.GetStackTraceString(err.(error)),
			}
		}
	}

	fillWithContext(entry)

	return f.JSONFormatter.Format(entry)
}

type myTextFormatter struct {
	*log.TextFormatter
}

func newMyTextFormatter() *myTextFormatter {
	return &myTextFormatter{&log.TextFormatter{}}
}

func (f *myTextFormatter) Format(entry *log.Entry) ([]byte, error) {
	if err, has := entry.Data[log.ErrorKey]; has {
		switch err.(type) {
		case error:
			entry.Data[log.ErrorKey] = err.(error).Error() + "\n" + errors.GetStackTraceString(err.(error))
		}
	}

	fillWithContext(entry)

	return f.TextFormatter.Format(entry)
}

func fillWithContext(entry *log.Entry) {
	if entry.Context == nil {
		return
	}

	if requestID, ok := entry.Context.Value(handler.CtxKeyRequestID).(string); ok {
		entry.Data[LogKeyRequestID] = requestID
	}
}
