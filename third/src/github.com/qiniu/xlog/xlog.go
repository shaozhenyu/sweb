package xlog

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/qiniu/log"
)

const logKey = "X-Log"
const reqidKey = "X-Reqid"

const (
	Ldate         = log.Ldate
	Ltime         = log.Ltime
	Lmicroseconds = log.Lmicroseconds
	Llongfile     = log.Llongfile
	Lshortfile    = log.Lshortfile
	Lmodule       = log.Lmodule
	Llevel        = log.Llevel
	LstdFlags     = log.LstdFlags
	Ldefault      = log.Ldefault
)

const (
	Ldebug = log.Ldebug
	Linfo  = log.Linfo
	Lwarn  = log.Lwarn
	Lerror = log.Lerror
	Lpanic = log.Lpanic
	Lfatal = log.Lfatal
)

// ============================================================================

type reqIder interface {
	ReqId() string
}

type header interface {
	Header() http.Header
}

// ============================================================================
// type *Logger

type Logger struct {
	h     http.Header
	reqId string
}

var pid = uint32(os.Getpid())

func genReqId() string {
	var b [12]byte
	binary.LittleEndian.PutUint32(b[:], pid)
	binary.LittleEndian.PutUint64(b[4:], uint64(time.Now().UnixNano()))
	return base64.URLEncoding.EncodeToString(b[:])
}

func New(w http.ResponseWriter, req *http.Request) *Logger {
	reqId := req.Header.Get(reqidKey)
	if reqId == "" {
		reqId = genReqId()
		req.Header.Set(reqidKey, reqId)
	}
	h := w.Header()
	h.Set(reqidKey, reqId)
	return &Logger{h, reqId}
}

func NewWith(a interface{}) *Logger {

	var h http.Header
	var reqId string
	if a == nil {
		reqId = genReqId()
	} else {
		l, ok := a.(*Logger)
		if ok {
			return l
		}
		reqId, ok = a.(string)
		if !ok {
			if g, ok := a.(reqIder); ok {
				reqId = g.ReqId()
			} else {
				panic("xlog.NewWith: unknown param")
			}
			if g, ok := a.(header); ok {
				h = g.Header()
			}
		}
	}
	if h == nil {
		h = http.Header{reqidKey: []string{reqId}}
	}
	return &Logger{h, reqId}
}

func NewDummy() *Logger {
	return NewWith(genReqId())
}

func (xlog *Logger) Spawn() *Logger {
	return NewWith(xlog.reqId)
}

// ============================================================================

func (xlog *Logger) Xget() []string {
	return xlog.h[logKey]
}

func (xlog *Logger) Xput(logs []string) {
	xlog.h[logKey] = append(xlog.h[logKey], logs...)
}

func (xlog *Logger) Xlog(v ...interface{}) {
	s := fmt.Sprint(v...)
	xlog.h[logKey] = append(xlog.h[logKey], s)
}

func (xlog *Logger) Xlogf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	xlog.h[logKey] = append(xlog.h[logKey], s)
}

/*
* 用法示意：

	func Foo(log xlog.*Logger) {
		...
		now := time.Now()
		err := longtimeOperation()
		log.Xprof("longtimeOperation", now, err)
		...
	}
*/
func (xlog *Logger) Xprof(modFn string, start time.Time, err error) {
	const maxErrorLen = 32
	durMs := time.Since(start).Nanoseconds() / 1e6
	if durMs > 0 {
		modFn += ":" + strconv.FormatInt(durMs, 10)
	}
	if err != nil {
		msg := err.Error()
		if len(msg) > maxErrorLen {
			msg = msg[:maxErrorLen]
		}
		modFn += "/" + msg
	}
	xlog.h[logKey] = append(xlog.h[logKey], modFn)
}

/*
* 用法示意：

	func Foo(log xlog.*Logger) (err error) {
		defer log.Xtrack("Foo", time.Now(), &err)
		...
	}

	func Bar(log xlog.*Logger) {
		defer log.Xtrack("Bar", time.Now(), nil)
		...
	}
*/
func (xlog *Logger) Xtrack(modFn string, start time.Time, errTrack *error) {
	var err error
	if errTrack != nil {
		err = *errTrack
	}
	xlog.Xprof(modFn, start, err)
}

// ============================================================================

func (xlog *Logger) ReqId() string {
	return xlog.reqId
}

func (xlog *Logger) Header() http.Header {
	return xlog.h
}

// Print calls Output to print to the standard Logger.
// Arguments are handled in the manner of fmt.Print.
func (xlog *Logger) Print(v ...interface{}) {
	log.Std.Output(xlog.reqId, log.Linfo, 2, fmt.Sprint(v...))
}

// Printf calls Output to print to the standard Logger.
// Arguments are handled in the manner of fmt.Printf.
func (xlog *Logger) Printf(format string, v ...interface{}) {
	log.Std.Output(xlog.reqId, log.Linfo, 2, fmt.Sprintf(format, v...))
}

// Println calls Output to print to the standard Logger.
// Arguments are handled in the manner of fmt.Println.
func (xlog *Logger) Println(v ...interface{}) {
	log.Std.Output(xlog.reqId, log.Linfo, 2, fmt.Sprintln(v...))
}

// -----------------------------------------

func (xlog *Logger) Debugf(format string, v ...interface{}) {
	if log.Ldebug < log.Std.Level {
		return
	}
	log.Std.Output(xlog.reqId, log.Ldebug, 2, fmt.Sprintf(format, v...))
}

func (xlog *Logger) Debug(v ...interface{}) {
	if log.Ldebug < log.Std.Level {
		return
	}
	log.Std.Output(xlog.reqId, log.Ldebug, 2, fmt.Sprintln(v...))
}

// -----------------------------------------

func (xlog *Logger) Infof(format string, v ...interface{}) {
	if log.Linfo < log.Std.Level {
		return
	}
	log.Std.Output(xlog.reqId, log.Linfo, 2, fmt.Sprintf(format, v...))
}

func (xlog *Logger) Info(v ...interface{}) {
	if log.Linfo < log.Std.Level {
		return
	}
	log.Std.Output(xlog.reqId, log.Linfo, 2, fmt.Sprintln(v...))
}

// -----------------------------------------

func (xlog *Logger) Warnf(format string, v ...interface{}) {
	log.Std.Output(xlog.reqId, log.Lwarn, 2, fmt.Sprintf(format, v...))
}

func (xlog *Logger) Warn(v ...interface{}) {
	log.Std.Output(xlog.reqId, log.Lwarn, 2, fmt.Sprintln(v...))
}

// -----------------------------------------

func (xlog *Logger) Errorf(format string, v ...interface{}) {
	log.Std.Output(xlog.reqId, log.Lerror, 2, fmt.Sprintf(format, v...))
}

func (xlog *Logger) Error(v ...interface{}) {
	log.Std.Output(xlog.reqId, log.Lerror, 2, fmt.Sprintln(v...))
}

// -----------------------------------------

// Fatal is equivalent to Print() followed by a call to os.Exit(1).
func (xlog *Logger) Fatal(v ...interface{}) {
	log.Std.Output(xlog.reqId, log.Lfatal, 2, fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalf is equivalent to Printf() followed by a call to os.Exit(1).
func (xlog *Logger) Fatalf(format string, v ...interface{}) {
	log.Std.Output(xlog.reqId, log.Lfatal, 2, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Fatalln is equivalent to Println() followed by a call to os.Exit(1).
func (xlog *Logger) Fatalln(v ...interface{}) {
	log.Std.Output(xlog.reqId, log.Lfatal, 2, fmt.Sprintln(v...))
	os.Exit(1)
}

// -----------------------------------------

// Panic is equivalent to Print() followed by a call to panic().
func (xlog *Logger) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	log.Std.Output(xlog.reqId, log.Lpanic, 2, s)
	panic(s)
}

// Panicf is equivalent to Printf() followed by a call to panic().
func (xlog *Logger) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	log.Std.Output(xlog.reqId, log.Lpanic, 2, s)
	panic(s)
}

// Panicln is equivalent to Println() followed by a call to panic().
func (xlog *Logger) Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	log.Std.Output(xlog.reqId, log.Lpanic, 2, s)
	panic(s)
}

func (xlog *Logger) Stack(v ...interface{}) {
	s := fmt.Sprint(v...)
	s += "\n"
	buf := make([]byte, 1024*1024)
	n := runtime.Stack(buf, true)
	s += string(buf[:n])
	s += "\n"
	log.Std.Output(xlog.reqId, log.Lerror, 2, s)
}

func (xlog *Logger) SingleStack(v ...interface{}) {
	s := fmt.Sprint(v...)
	s += "\n"
	buf := make([]byte, 1024*1024)
	n := runtime.Stack(buf, false)
	s += string(buf[:n])
	s += "\n"
	log.Std.Output(xlog.reqId, log.Lerror, 2, s)
}

// ============================================================================

func Debugf(reqId string, format string, v ...interface{}) {
	if log.Ldebug < log.Std.Level {
		return
	}
	log.Std.Output(reqId, log.Ldebug, 2, fmt.Sprintf(format, v...))
}

func Debug(reqId string, v ...interface{}) {
	if log.Ldebug < log.Std.Level {
		return
	}
	log.Std.Output(reqId, log.Ldebug, 2, fmt.Sprintln(v...))
}

// -----------------------------------------

func Infof(reqId string, format string, v ...interface{}) {
	if log.Linfo < log.Std.Level {
		return
	}
	log.Std.Output(reqId, log.Linfo, 2, fmt.Sprintf(format, v...))
}

func Info(reqId string, v ...interface{}) {
	if log.Linfo < log.Std.Level {
		return
	}
	log.Std.Output(reqId, log.Linfo, 2, fmt.Sprintln(v...))
}

// -----------------------------------------

func Warnf(reqId string, format string, v ...interface{}) {
	log.Std.Output(reqId, log.Lwarn, 2, fmt.Sprintf(format, v...))
}

func Warn(reqId string, v ...interface{}) {
	log.Std.Output(reqId, log.Lwarn, 2, fmt.Sprintln(v...))
}

// -----------------------------------------

func Errorf(reqId string, format string, v ...interface{}) {
	log.Std.Output(reqId, log.Lerror, 2, fmt.Sprintf(format, v...))
}

func Error(reqId string, v ...interface{}) {
	log.Std.Output(reqId, log.Lerror, 2, fmt.Sprintln(v...))
}

// -----------------------------------------

// Fatal is equivalent to Print() followed by a call to os.Exit(1).
func Fatal(reqId string, v ...interface{}) {
	log.Std.Output(reqId, log.Lfatal, 2, fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalf is equivalent to Printf() followed by a call to os.Exit(1).
func Fatalf(reqId string, format string, v ...interface{}) {
	log.Std.Output(reqId, log.Lfatal, 2, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Fatalln is equivalent to Println() followed by a call to os.Exit(1).
func Fatalln(reqId string, v ...interface{}) {
	log.Std.Output(reqId, log.Lfatal, 2, fmt.Sprintln(v...))
	os.Exit(1)
}

// -----------------------------------------

// Panic is equivalent to Print() followed by a call to panic().
func Panic(reqId string, v ...interface{}) {
	s := fmt.Sprint(v...)
	log.Std.Output(reqId, log.Lpanic, 2, s)
	panic(s)
}

// Panicf is equivalent to Printf() followed by a call to panic().
func Panicf(reqId string, format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	log.Std.Output(reqId, log.Lpanic, 2, s)
	panic(s)
}

// Panicln is equivalent to Println() followed by a call to panic().
func Panicln(reqId string, v ...interface{}) {
	s := fmt.Sprintln(v...)
	log.Std.Output(reqId, log.Lpanic, 2, s)
	panic(s)
}

func Stack(reqId string, v ...interface{}) {
	s := fmt.Sprint(v...)
	s += "\n"
	buf := make([]byte, 1024*1024)
	n := runtime.Stack(buf, true)
	s += string(buf[:n])
	s += "\n"
	log.Std.Output(reqId, log.Lerror, 2, s)
}

func SingleStack(reqId string, v ...interface{}) {
	s := fmt.Sprint(v...)
	s += "\n"
	buf := make([]byte, 1024*1024)
	n := runtime.Stack(buf, false)
	s += string(buf[:n])
	s += "\n"
	log.Std.Output(reqId, log.Lerror, 2, s)
}

func SetOutput(w io.Writer) {
	log.SetOutput(w)
}

func SetFlags(flag int) {
	log.SetFlags(flag)
}

func SetOutputLevel(lvl int) {
	log.SetOutputLevel(lvl)
}
