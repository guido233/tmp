package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"
	"sort"
	"strings"
	"sync/atomic"
	"time"
)

// LogLevel is the log level type.
type LogLevel int

const (
	// DEBUG represents debug log level.
	DEBUG LogLevel = iota
	// INFO represents info log level.
	INFO
	// WARN represents warn log level.
	WARN
	// ERROR represents error log level.
	ERROR
	// FATAL represents fatal log level.
	FATAL
)

/*日志间隔周期*/
const (
	LOG_INTERVAL_YEAR   = "year"
	LOG_INTERVAL_MONTH  = "month"
	LOG_INTERVAL_DAY    = "day"
	LOG_INTERVAL_HOUR   = "hour"
	LOG_INTERVAL_MINUTE = "minute"
)

var (
	started        int32
	loggerInstance Logger
	tagName        = map[LogLevel]string{
		DEBUG: "DEBUG",
		INFO:  "INFO",
		WARN:  "WARN",
		ERROR: "ERROR",
		FATAL: "FATAL",
	}
)

// Logger is the logger type.
type Logger struct {
	logger     *log.Logger
	level      LogLevel
	segment    *logSegment
	stopped    int32
	logPath    string
	isStdout   bool
	printStack bool
	/*是否发送到远程服务器*/
	remoteEnabled bool
	/*网络方式 udp/tcp*/
	network string
	/*远程服务器地址*/
	logTarget string
	/*最近一次日志时间*/
	lastLogTime *time.Time
	/*是否按时间输出日志*/
	enableTimeCheck bool
	/*
	 * 日志输出间隔
	 * year - 按年
	 * month - 按月
	 * day - 按天
	 * hour - 按小时
	 * minute - 按分钟
	 */
	logInterval string
	/*是否按大小输出日志*/
	enableSizeCheck bool
	/*日志大小(MB)*/
	logSize int64
	/*日志最多保存个数*/
	logMaxCount int
}

// logSegment implements io.Writer
type logSegment struct {
	logger  *Logger
	logPath string
	logFile *os.File
}

// Start returns a decorated innerLogger.
func Start(decorators ...func(Logger) Logger) Logger {
	if atomic.CompareAndSwapInt32(&started, 0, 1) {
		loggerInstance = Logger{}
		for _, decorator := range decorators {
			loggerInstance = decorator(loggerInstance)
		}
		var logger *log.Logger
		var segment *logSegment
		if loggerInstance.logPath != "" {
			segment = newLogSegment(&loggerInstance, loggerInstance.logPath)
		}
		if segment != nil {
			logger = log.New(segment, "", log.Ldate|log.Ltime|log.Lmicroseconds)
		} else if loggerInstance.isStdout {
			logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds)
		} else {
			logger = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lmicroseconds)
		}
		loggerInstance.segment = segment
		loggerInstance.logger = logger
		now := time.Now()
		loggerInstance.lastLogTime = &now
		return loggerInstance
	}
	panic("Start() already called")
}

// Stop stops the logger.
func (l Logger) Stop() {
	if atomic.CompareAndSwapInt32(&l.stopped, 0, 1) {
		if l.printStack {
			traceInfo := make([]byte, 1<<16)
			n := runtime.Stack(traceInfo, true)
			l.logger.Printf("%s", traceInfo[:n])
			if l.isStdout {
				log.Printf("%s", traceInfo[:n])
			}
		}
		if l.segment != nil {
			l.segment.Close()
		}
		l.segment = nil
		l.logger = nil
		atomic.StoreInt32(&started, 0)
	}
}

func newLogSegment(logger *Logger, logPath string) *logSegment {
	if logPath != "" {
		err := os.MkdirAll(logPath, os.ModePerm)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return nil
		}
		name := getLogFileName(time.Now())
		logFile, err := os.OpenFile(path.Join(logPath, name), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			if os.IsNotExist(err) {
				logFile, err = os.Create(path.Join(logPath, name))
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					return nil
				}
			} else {
				fmt.Fprintln(os.Stderr, err)
				return nil
			}
		}

		return &logSegment{
			logger:  logger,
			logPath: logPath,
			logFile: logFile,
		}
	}
	return nil
}

func (ls *logSegment) Write(p []byte) (n int, err error) {
	if ls.logFile != os.Stdout && ls.logFile != os.Stderr {

		bCheck := false
		// 优先按照大小切割日志
		if ls.logger.enableSizeCheck {
			bCheck = ls.logger.CheckLogSize()
		} else if ls.logger.enableTimeCheck {
			bCheck = ls.logger.CheckLogInterval()
		}

		if bCheck {

			/* 先关闭文件 */
			ls.logFile.Close()
			ls.logFile = nil

			/* 修改上一个文件名 */
			if err := os.Rename(ls.logPath+getLogFileName(time.Now()), ls.logPath+getLogFileNameWithTime(time.Now())); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}

			/* 判断文件个数来限制日志数量 */
			if ls.logger.logMaxCount != 0 {
				if err := ls.logger.LimitLogMaxCount(); err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
			}

			now := time.Now()
			ls.logger.lastLogTime = &now

			/* 生成新的日志文件 */
			name := getLogFileName(time.Now())
			ls.logFile, err = os.Create(path.Join(ls.logPath, name))
			if err != nil {
				// log inlastLogTime if we can't create new file
				fmt.Fprintln(os.Stderr, err)
				ls.logFile = os.Stderr
			}

		}
	}
	return ls.logFile.Write(p)
}

// Close closes the log file.
func (ls *logSegment) Close() {
	ls.logFile.Close()
}

// getLogFileName returns the log file name.
func getLogFileName(t time.Time) string {
	proc := path.Base(os.Args[0])
	return fmt.Sprintf("%s.log", proc)
}

// getLogFileNameWithTime returns the log file name with time.
func getLogFileNameWithTime(t time.Time) string {
	proc := path.Base(os.Args[0])
	now := time.Now()
	year := now.Year()
	month := now.Month()
	day := now.Day()
	hour := now.Hour()
	minute := now.Minute()
	second := now.Second()
	return fmt.Sprintf("%s.%04d-%02d-%02d-%02d-%02d-%02d.log",
		proc, year, month, day, hour, minute, second)
}

const (
	// 是否隐藏日志细节
	hideDetail = false
)

// doPrintln prints log.
func (l Logger) doPrintln(level LogLevel, v ...interface{}) {
	if l.logger == nil {
		return
	}
	if level >= l.level {
		var prefix string
		if !hideDetail {
			funcName, fileName, lineNum := getRuntimeInfo()
			prefix = fmt.Sprintf("%5s [%s] (%s:%d) - ", tagName[level], path.Base(funcName), path.Base(fileName), lineNum)
		} else {
			funcName, _, _ := getRuntimeInfo()
			prefix = fmt.Sprintf("%5s [%s] - ", tagName[level], path.Base(funcName))
		}

		value := fmt.Sprintf("%s%s", prefix, fmt.Sprintln(v...))
		l.logger.Print(value)
		if l.isStdout {
			log.Print(value)
		}
		if level == FATAL {
			os.Exit(1)
		}
	}
}

// doPrintf prints formatted log.
func (l Logger) doPrintf(level LogLevel, format string, v ...interface{}) {
	if l.logger == nil {
		return
	}
	if level >= l.level {
		if !hideDetail {
			funcName, fileName, lineNum := getRuntimeInfo()
			format = fmt.Sprintf("%5s [%s] (%s:%d) - %s", tagName[level], path.Base(funcName), path.Base(fileName), lineNum, format)
		} else {
			funcName, _, _ := getRuntimeInfo()
			format = fmt.Sprintf("%5s [%s] - %s", tagName[level], path.Base(funcName), format)
		}
		l.logger.Printf(format, v...)
		if l.isStdout {
			log.Printf(format, v...)
		}
		if level == FATAL {
			os.Exit(1)
		}
	}

	if l.remoteEnabled {
		l.sendLog(level, format, v)
	}
}

func (l Logger) Name() string {
	return "logger"
}

// Configure configures the provider
func (l Logger) Configure(config map[string]interface{}) error {
	return nil
}

func (l Logger) Printf(format string, v ...interface{}) {
	l.doPrintf(INFO, format, v)
}

func (l Logger) sendLog(level LogLevel, format string, v ...interface{}) error {
	if l.logTarget == "" {
		return nil
	}

	content := fmt.Sprintf(format, v)

	reqBody, err := json.Marshal(map[string]interface{}{
		"level":     level,
		"content":   content,
		"timestamp": time.Now().Unix(),
	})
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	req, err := http.NewRequest(http.MethodPost, l.logTarget, bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}

	// Asynchronous call
	go l.makeRequest(client, req)

	return nil
}

func (l Logger) makeRequest(client *http.Client, request *http.Request) {
	resp, err := client.Do(request)
	if err == nil {
		defer resp.Body.Close()
		resp.Close = true
	} else {
		fmt.Println(err.Error())
	}
}

// getRuntimeInfo returns the function name, file name and line number of the caller.
func getRuntimeInfo() (string, string, int) {
	pc, fn, ln, ok := runtime.Caller(3) // 3 steps up the stack frame
	if !ok {
		fn = "???"
		ln = 0
	}
	function := "???"
	caller := runtime.FuncForPC(pc)
	if caller != nil {
		function = caller.Name()
	}
	return function, fn, ln
}

/* 日志的输出类型 */

// Debugf prints formatted debug log.
func Debugf(format string, v ...interface{}) {
	loggerInstance.doPrintf(DEBUG, format, v...)
}

// Infof prints formatted info log.
func Infof(format string, v ...interface{}) {
	loggerInstance.doPrintf(INFO, format, v...)
}

// Warnf prints formatted warn log.
func Warnf(format string, v ...interface{}) {
	loggerInstance.doPrintf(WARN, format, v...)
}

// Errorf prints formatted error log.
func Errorf(format string, v ...interface{}) {
	loggerInstance.doPrintf(ERROR, format, v...)
}

// Fatalf prints formatted fatal log and exits.
func Fatalf(format string, v ...interface{}) {
	loggerInstance.doPrintf(FATAL, format, v...)
	os.Exit(1)
}

// Debugln prints debug log.
func Debugln(v ...interface{}) {
	loggerInstance.doPrintln(DEBUG, v...)
}

// Infoln prints info log.
func Infoln(v ...interface{}) {
	loggerInstance.doPrintln(INFO, v...)
}

// Warnln prints warn log.
func Warnln(v ...interface{}) {
	loggerInstance.doPrintln(WARN, v...)
}

// Errorln prints error log.
func Errorln(v ...interface{}) {
	loggerInstance.doPrintln(ERROR, v...)
}

// Fatalln prints fatal log and exits.
func Fatalln(v ...interface{}) {
	loggerInstance.doPrintln(FATAL, v...)
	os.Exit(1)
}

/* 日志选项 */

// DebugLevel sets log level to debug.
func DebugLevel() func(Logger) Logger {
	return func(l Logger) Logger {
		l.level = DEBUG
		return l
	}
}

// InfoLevel sets log level to info.
func InfoLevel() func(Logger) Logger {
	return func(l Logger) Logger {
		l.level = INFO
		return l
	}
}

// WarnLevel sets log level to warn.
func WarnLevel() func(Logger) Logger {
	return func(l Logger) Logger {
		l.level = WARN
		return l
	}
}

// ErrorLevel sets log level to error.
func ErrorLevel() func(Logger) Logger {
	return func(l Logger) Logger {
		l.level = ERROR
		return l
	}
}

// FatalLevel sets log level to fatal.
func FatalLevel() func(Logger) Logger {
	return func(l Logger) Logger {
		l.level = FATAL
		return l
	}
}

// GetLogger returns a function to set the log level.
func GetLogger(level string) func(Logger) Logger {
	if level == "debug" {
		return DebugLevel()
	} else if level == "info" {
		return InfoLevel()
	} else if level == "warning" {
		return WarnLevel()
	} else if level == "error" {
		return ErrorLevel()
	} else {
		return WarnLevel()
	}
}

// LogFilePath returns a function to set the log file path.
func LogFilePath(p string) func(Logger) Logger {
	return func(l Logger) Logger {
		l.logPath = p
		return l
	}
}

// EnableRemote 打开logTarget
func EnableRemote(target string) func(Logger) Logger {
	return func(l Logger) Logger {
		l.remoteEnabled = true
		l.logTarget = target

		return l
	}
}

// EveryYear sets new log file created every year.
func EveryYear() func(Logger) Logger {
	return func(l Logger) Logger {
		l.enableTimeCheck = true
		l.logInterval = LOG_INTERVAL_YEAR
		return l
	}
}

// EveryMonth sets new log file created every month.
func EveryMonth() func(Logger) Logger {
	return func(l Logger) Logger {
		l.enableTimeCheck = true
		l.logInterval = LOG_INTERVAL_MONTH
		return l
	}
}

// EveryDay sets new log file created every day.
func EveryDay() func(Logger) Logger {
	return func(l Logger) Logger {
		l.enableTimeCheck = true
		l.logInterval = LOG_INTERVAL_DAY
		return l
	}
}

// EveryHour sets new log file created every hour.
func EveryHour() func(Logger) Logger {
	return func(l Logger) Logger {
		l.enableTimeCheck = true
		l.logInterval = LOG_INTERVAL_HOUR
		return l
	}
}

// EveryMinute sets new log file created every minute.
func EveryMinute() func(Logger) Logger {
	return func(l Logger) Logger {
		l.enableTimeCheck = true
		l.logInterval = LOG_INTERVAL_MINUTE
		return l
	}
}

// LogSize sets the log file size(MB).
func LogSize(sizeMB int64) func(Logger) Logger {
	return func(l Logger) Logger {
		l.enableSizeCheck = true
		l.logSize = sizeMB
		return l
	}
}

// LogMaxCount sets the max count of log files.
// 只有日志切割时才有用
func LogMaxCount(count int) func(Logger) Logger {
	return func(l Logger) Logger {
		l.logMaxCount = count
		return l
	}
}

// AlsoStdout sets log also output to stdio.
func AlsoStdout() func(Logger) Logger {
	return func(l Logger) Logger {
		l.isStdout = true
		return l
	}
}

// PrintStack sets log output the stack trace info.
func PrintStack() func(Logger) Logger {
	return func(l Logger) Logger {
		l.printStack = true
		return l
	}
}

/* 日志切割 */

// CheckLogInterval 按时间切割日志
func (l Logger) CheckLogInterval() bool {

	now := time.Now()
	if l.lastLogTime == nil {
		l.lastLogTime = &now
	}

	year_last, month_last, day_last := l.lastLogTime.Date()
	hour_last := l.lastLogTime.Hour()
	minute_last := l.lastLogTime.Minute()

	year_now, month_now, day_now := now.Date()
	hour_now := now.Hour()
	minute_now := now.Minute()

	if l.logInterval == LOG_INTERVAL_YEAR {
		if year_last != year_now {
			return true
		}
	} else if l.logInterval == LOG_INTERVAL_MONTH {
		if year_last != year_now || int(month_last) != int(month_now) {
			return true
		}
	} else if l.logInterval == LOG_INTERVAL_DAY {
		if year_last != year_now || int(month_last) != int(month_now) || day_last != day_now {
			return true
		}
	} else if l.logInterval == LOG_INTERVAL_HOUR {
		if year_last != year_now || int(month_last) != int(month_now) || day_last != day_now || hour_last != hour_now {
			return true
		}
	} else if l.logInterval == LOG_INTERVAL_MINUTE {
		if year_last != year_now || int(month_last) != int(month_now) || day_last != day_now || hour_last != hour_now || minute_last != minute_now {
			return true
		}
	} else {
		return false
	}

	return false
}

// CheckLogSize 按大小切割日志
func (l Logger) CheckLogSize() bool {

	if l.segment == nil {
		return false
	}

	stat, err := l.segment.logFile.Stat()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return false
	}
	size := stat.Size()

	if size >= l.logSize*1024*1024 {
		return true
	}

	return false
}

// LimitLogMaxCount 限制日志最大个数
func (l Logger) LimitLogMaxCount() error {

	var logTime []string
	proc := path.Base(os.Args[0])

	// 读取文件
	files, err := ioutil.ReadDir(l.logPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	// 匹配文件
	for _, f := range files {
		split := strings.Split(f.Name(), ".")
		if len(split) < 3 {
			continue
		}
		if split[0] != proc {
			continue
		}
		if split[2] != "log" {
			continue
		}

		logTime = append(logTime, split[1])
	}

	if len(logTime) <= l.logMaxCount {
		return nil
	}

	sort.Strings(logTime)

	for i := 0; i < len(logTime)-l.logMaxCount; i++ {
		// 删除文件
		if err := os.Remove(l.logPath + proc + "." + logTime[i] + ".log"); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}

	return nil
}
