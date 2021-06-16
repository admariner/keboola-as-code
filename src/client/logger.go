package client

import (
	"fmt"
	"go.uber.org/zap"
	"regexp"
	"strings"
)

const LoggerPrefix = "HTTP%s\t"

type Logger struct {
	logger *zap.SugaredLogger
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.logWithoutSecrets("", format, v...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.logWithoutSecrets("-WARN", format, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.logWithoutSecrets("-ERROR", format, v...)
}

func (l *Logger) logWithoutSecrets(level string, format string, v ...interface{}) {
	v = append([]interface{}{level}, v...)
	msg := fmt.Sprintf(LoggerPrefix+format, v...)
	msg = removeSecrets(msg)
	msg = strings.TrimSuffix(msg, "\n")
	l.logger.Debug(msg)
}

func removeSecrets(str string) string {
	return regexp.MustCompile(`(?i)(token[^\w/,]\s*)\d[^\s/]*`).ReplaceAllString(str, "$1*****")
}