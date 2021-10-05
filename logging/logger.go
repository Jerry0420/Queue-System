package logging

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"time"
)

type LoggerTool interface {
	parseContextValue(context context.Context) []interface{}
	generateLogContent(level string, message []interface{}) []interface{}
	INFOf(message ...interface{})
	WARNf(message ...interface{})
	ERRORf(message ...interface{})
	FATALf(message ...interface{})
}

type loggerTool struct {
	logger *log.Logger
	contextKeys []string
}

func getLogFilePath() string{
	err := os.MkdirAll("/app/logs", 0777)
    if err != nil {
		log.Fatalf("Fail to getLogFilePath :%v", err)
    }
	return fmt.Sprintf("/app/logs/%s.log", time.Now().Format("20060102"))
}

func NewLogger(contextKeys []string, disable bool) *loggerTool {
	var lt *loggerTool
	logFilePath := getLogFilePath()
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Fail to NewLogger :%v", err)
    }
	lt = &loggerTool{
		logger: log.New(io.MultiWriter(os.Stdout, logFile), "", log.Ldate|log.Ltime),
		contextKeys: contextKeys,
	}

	if disable == true {
		lt.logger = log.New(ioutil.Discard, "", log.Ldate)
	}
	return lt
}

func (lt *loggerTool) parseContextValue(context context.Context) []interface{} {
	var values []interface{}
	for _, key := range lt.contextKeys {
		if value := context.Value(key); value != nil {
			values = append(values, value)
		}
	}
	return values
}

func (lt *loggerTool) generateLogContent(level string, message []interface{}) []interface{} {
	var content []interface{}

	_, file, line, _ := runtime.Caller(2)
	content = append(content, fmt.Sprint(level, file, ":", line))

	if len(message) == 0 {
		return content
	}

	if firstMessage, ok := message[0].(context.Context); ok {
		contextValues := lt.parseContextValue(firstMessage)
		content = append(content, contextValues...)
		if len(message) > 1 {
			message = message[1:]
		} else {
			return content
		}
	}

	content = append(content, "|")

	if formatString, ok := message[0].(string); ok && len(message) > 1 {
		content = append(content, fmt.Sprintf(formatString, message[1:]...))
	} else {
		content = append(content, message...)
	}

	return content
}

func (lt *loggerTool) INFOf(message ...interface{}) {
	content := lt.generateLogContent("[INFO] ", message)
	lt.logger.Println(content...)
}

func (lt *loggerTool) WARNf(message ...interface{}) {
	content := lt.generateLogContent("[WARN] ", message)
	lt.logger.Println(content...)
}

func (lt *loggerTool) ERRORf(message ...interface{}) {
	content := lt.generateLogContent("[ERROR] ", message)
	lt.logger.Println(content...)
}

func (lt *loggerTool) FATALf(message ...interface{}) {
	content := lt.generateLogContent("[FATAL] ", message)
	lt.logger.Fatalln(content...)
}