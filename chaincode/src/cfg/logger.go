package cfg

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"
)

var (
	LOGFILE_MAXSIZE_DEFAULT int64 = 50 << 20
	_logger_map             map[string]*log.Logger
	_logfile_maxsize        int64 = LOGFILE_MAXSIZE_DEFAULT
	_logfile_base           string
)

func init() {
	_logger_map = make(map[string]*log.Logger)
}

func GetLogger(typ string) *log.Logger {
	if logger, ok := _logger_map[typ]; ok {
		return logger
	}
	file, err := OpenLogFile(_logfile_base+"."+typ, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("error opening file %v\n", err)
		return nil
	}

	logger := log.New(file, "", log.LstdFlags)
	_logger_map[typ] = logger

	return logger
}

func InitLogger(logfile string, maxSize int64) {

	var err error
	var fullpath string

	if filepath.IsAbs(logfile) { // start with slash, just open
		fullpath = logfile
	} else {
		fullpath = path.Join(_base_path, "", logfile)
	}
	dir, filename := path.Split(logfile)
	if filename == "" {
		fullpath = path.Join(fullpath, "log")
	}
	_logfile_base = fullpath

	err = os.MkdirAll(dir, 0777)
	if err != nil {
		LogFatalf("MkdirAll err:", err)
		return
	}

	_logfile_maxsize = int64(maxSize) << 10 //单位是k
	if _logfile_maxsize < (1 << 16) {       //日志文件最小64k
		_logfile_maxsize = LOGFILE_MAXSIZE_DEFAULT
	}
	log.Println("logfile size:", _logfile_maxsize)

	startLogger(fullpath)
}
func startLogger(logfile string) {
	f, err := OpenLogFile(logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		log.Printf("cannot open logfile %v\n", err)
		os.Exit(-1)
	}

	log.SetOutput(f)
}

func tmpLog(p *[]byte, format string, v ...interface{}) {
	*p = append([]byte(fmt.Sprintf(format, v...)), (*p)...)
}

type LogFile struct {
	*os.File
}

func OpenLogFile(name string, flag int, perm os.FileMode) (file *LogFile, err error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}
	lf := LogFile{}
	lf.File = f
	return &lf, nil
}

func (f *LogFile) Write(p []byte) (int, error) {
	fi, err := f.Stat()
	if err != nil {
		tmpLog(&p, "file.Stat err:%v.", err)
	}

	if fi.Size() >= _logfile_maxsize {
		now := int64(time.Now().UnixNano() / 1000000)
		curFileName := f.Name()
		newFileName := fmt.Sprintf("%s.%d", f.Name(), now)

		err = os.Rename(curFileName, newFileName)
		if err != nil {
			tmpLog(&p, "[RAW] rename [%s] to [%s] err:%v\n",
				curFileName, newFileName, err)
		}

		newFile, err := os.OpenFile(curFileName,
			os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			tmpLog(&p, "[RAW] open file %s err:%v", curFileName, err)
		} else {
			f.File.Close()
			f.File = newFile
		}
	}

	return f.File.Write(p)
}
