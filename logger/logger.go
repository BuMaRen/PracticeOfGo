package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	SIZEMAX = iota
	GETINFOFAIL
	TERMINATE
	MAXLOGSIZE = 1024 * 1024 * 100 //100mb
	DATE       = 0
	TIME       = 1
)

type Logger interface {
	Write(str string)
	WriteS(format string, a ...interface{})
	Run()
}

type logger struct {
	file        *os.File
	size        int
	writeMtx    chan int
	path        string
	buffer      strings.Builder
	callerDepth int
}

var lgr *logger = nil
var running bool = false

// GetLogger 获取唯一的日志实例
func GetLogger(path string) Logger {
	if lgr != nil {
		return lgr
	}
	lgr = new(logger)
	if lgr.init(path) != nil {
		return nil
	}
	return lgr
}

func getTimeString(t time.Time) string {
	sb := strings.Builder{}
	buffer := strings.Split(t.String(), " ")
	for _, v := range strings.Split(buffer[DATE], "-") {
		sb.WriteString(v)
	}
	times := strings.Split(buffer[TIME], ".")
	for _, v := range strings.Split(times[0], ":") {
		sb.WriteString(v)
	}
	sb.WriteString(times[1][:3])
	fmt.Println(times[1][:3])
	return sb.String()
}

// setOutPut 定向输出到日志文件
func (lgr *logger) setOutPut(file string) error {
	fd, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	fd.Seek(0, os.SEEK_END)
	lgr.file = fd
	log.SetOutput(lgr.file)
	return nil
}

// Init 初始化日志组件，入参为日志文件的目录路径
func (lgr *logger) init(path string) error {
	lgr.path = path
	// 0大小的chan必须有一端听才能写的进去
	lgr.writeMtx = make(chan int, 1)
	log.SetFlags(log.Flags() | log.Lshortfile)
	lgr.buffer = strings.Builder{}
	if err := lgr.setOutPut(path + "/server_log.log"); err != nil {
		return err
	}
	lgr.callerDepth = 1
	return nil
}

func (lgr *logger) Write(str string) {
	lgr.writeMtx <- 1
	defer func() {
		<-lgr.writeMtx
	}()

	//[文件:行->方法]
	callerPtr, filePath, line, _ := runtime.Caller(lgr.callerDepth)
	funcName := runtime.FuncForPC(callerPtr).Name()
	filePathSlice := strings.Split(filePath, "/")
	fileName := filePathSlice[len(filePathSlice)-1]
	timeStr := time.Now().Format(time.ANSIC)
	caller := fmt.Sprintf("%s [%s:%d->%s] ", timeStr, fileName, line, funcName)
	lgr.buffer.WriteString(string(caller))
	lgr.buffer.WriteString(str)
	lgr.buffer.WriteRune('\n')
	sz, err := lgr.file.WriteString(lgr.buffer.String())
	lgr.buffer.Reset()
	if err != nil {
		fmt.Println("writing log fail, because", err)
		return
	}
	lgr.size += sz
}

func (lgr *logger) WriteS(format string, a ...interface{}) {
	lgr.callerDepth = 2
	logString := fmt.Sprintf(format, a...)
	lgr.Write(logString)
	lgr.callerDepth = 1
}

// Run 启动日志检测进程
func (lgr *logger) Run() {
	if running {
		return
	}
	running = true
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			lgr.writeMtx <- 1
			lgr.file.Close()
			running = false
			<-lgr.writeMtx
		}
	}()
	lgr.fileGuardian()
	lgr.Write("logger ready")
}

// fileGuardian 监控文件大小
func (lgr *logger) fileGuardian() {
	for {
		if lgr.size > MAXLOGSIZE {
			lgr.writeMtx <- 1
			lgr.updateLogFile()
			<-lgr.writeMtx
		}
		time.Sleep(time.Second)
	}
}

// updateLogFile 超过限制的文件进行压缩或者分开
func (lgr *logger) updateLogFile() error {
	oldName := lgr.path + "/" + lgr.file.Name()
	NewName := lgr.path + "/server_log" + getTimeString(time.Now()) + ".log"
	if err := lgr.file.Close(); err != nil {
		fmt.Println("close file failed, because", err)
		return err
	}
	if err := os.Rename(oldName, NewName); err != nil {
		fmt.Println("rename file failed, because", err)
		return err
	}
	if err := lgr.setOutPut(oldName); err != nil {
		fmt.Println("open new file failed, because", err)
		return err
	}
	return nil
}
