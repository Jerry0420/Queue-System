package logging

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
	"path/filepath"
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
	sync.Mutex
	logger *log.Logger
	logFile *os.File
	contextKeys []string
}

func getLogFilePathOfToday() string{
	// save log file inside "logs" dir.
	executableFilePath, _ := os.Executable()
	executableFileDir := filepath.Dir(executableFilePath)
	logBaseDir := filepath.Join(executableFileDir, "logs")
	
	err := os.MkdirAll(logBaseDir, 0777)
    if err != nil {
		log.Fatalf("Fail to getLogFilePath :%v", err)
    }
	return fmt.Sprintf("%s/%s.log", logBaseDir, time.Now().Format("20060102"))
}

func NewLogger(contextKeys []string, disable bool) *loggerTool {
	var lt *loggerTool

	logFilePath := getLogFilePathOfToday()
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Fail to open log file :%v", err)
    }

	lt = &loggerTool{
		logger: log.New(io.MultiWriter(os.Stdout, logFile), "", log.Ldate|log.Ltime),
		logFile: logFile,
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
		// add pipe symbol for seperation.
		if key == "sep" {
			values = append(values, "|")
			continue
		}
		if value := context.Value(key); value != nil {
			values = append(values, value)
		}
	}
	return values
}

func (lt *loggerTool) generateLogContent(level string, message []interface{}) []interface{} {
	var content []interface{}

	_, file, line, _ := runtime.Caller(2)
	content = append(content, fmt.Sprint(level, file, ":", line, " |"))

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

// daily rotate the log file.
func (lt *loggerTool) writeContentOut(content []interface{}) {
	today := fmt.Sprint(time.Now().Format("20060102"))
	if !strings.Contains(lt.logFile.Name(), today) {
		logFilePath := getLogFilePathOfToday()
		lt.Lock()
		newlogFile, err := os.OpenFile(logFilePath, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Fail to open log file :%v", err)
		}

		lt.logger.SetOutput(io.MultiWriter(os.Stdout, newlogFile))
		
		err = lt.logFile.Close()
		if err != nil {
			log.Fatalf("Fail to close old log file :%v", err)
		}
		lt.logFile = newlogFile
		defer lt.Unlock()
	}
	lt.logger.Println(content...)
}

func (lt *loggerTool) INFOf(message ...interface{}) {
	content := lt.generateLogContent("[INFO] ", message)
	lt.writeContentOut(content)
}

func (lt *loggerTool) WARNf(message ...interface{}) {
	content := lt.generateLogContent("[WARN] ", message)
	lt.writeContentOut(content)
}

func (lt *loggerTool) ERRORf(message ...interface{}) {
	content := lt.generateLogContent("[ERROR] ", message)
	lt.writeContentOut(content)
}

func (lt *loggerTool) FATALf(message ...interface{}) {
	content := lt.generateLogContent("[FATAL] ", message)
	lt.writeContentOut(content)
	os.Exit(1)
}