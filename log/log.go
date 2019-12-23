package log

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"

	"github.com/kitabisa/perkakas/v2/httputil"
	"github.com/kitabisa/perkakas/v2/token/jwt"
)

type Level uint32

const (
	FieldLogID           = "log_id"
	FieldHTTPStatus      = "http_status"
	FieldEndpoint        = "endpoint"
	FieldMethod          = "method"
	FieldServiceName     = "service"
	FieldUserID          = "user_id"
	FieldRequestBody     = "request_body"
	FieldRequestHeaders  = "request_headers"
	FieldResponseBody    = "response_body"
	FieldResponseHeaders = "response_headers"
)

type message struct {
	Message  interface{} `json:"message"`
	Level    Level       `json:"level"`
	File     string      `json:"file"`
	FuncName string      `json:"func"`
	Line     int         `json:"line"`
}

const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

const (
	// DefaultIndex specifies the default index of runtime caller.
	DefaultIndex = 2
)

func (level Level) MarshalText() ([]byte, error) {
	switch level {
	case TraceLevel:
		return []byte("trace"), nil
	case DebugLevel:
		return []byte("debug"), nil
	case InfoLevel:
		return []byte("info"), nil
	case WarnLevel:
		return []byte("warning"), nil
	case ErrorLevel:
		return []byte("error"), nil
	case FatalLevel:
		return []byte("fatal"), nil
	case PanicLevel:
		return []byte("panic"), nil
	}

	return nil, fmt.Errorf("not a valid logrus level %d", level)
}

type Logger struct {
	logger  *log.Logger
	fields  sync.Map
	id      string
	service string
}

func (l *Logger) NewChildLogger() (logger *Logger) {
	logger = newLogger(l.service, l.id)
	return
}

func (l *Logger) SetRequest(req interface{}) {
	switch v := req.(type) {
	case *http.Request:
		token, ok := v.Context().Value("token").(*jwt.UserClaim)
		if ok {
			l.fields.Store(FieldUserID, token.UserID)
		}

		header := httputil.ExcludeSensitiveHeader(v.Header)

		l.fields.Store(FieldEndpoint, v.URL.String())
		l.fields.Store(FieldMethod, v.Method)
		l.fields.Store(FieldRequestHeaders, header)

		switch v.Method {
		case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
			body := httputil.ReadRequestBody(v)
			httputil.ExcludeSensitiveRequestBody(&body)

			l.fields.Store(FieldRequestBody, body)
		}
	default:
		l.fields.Store(FieldRequestBody, req)
	}
}

func (l *Logger) SetResponse(res interface{}, body []byte) {
	switch v := res.(type) {
	case http.ResponseWriter:
		l.fields.Store(FieldResponseHeaders, v.Header())
		l.fields.Store(FieldResponseBody, string(body))
	case *http.Response:
		l.fields.Store(FieldResponseHeaders, v.Header)
		l.fields.Store(FieldResponseBody, string(body))
	}
}

func (l *Logger) AddMessage(level Level, message ...interface{}) *Logger {
	l.setCaller(level, message...)
	return l
}

func (l *Logger) Print(directMsg ...interface{}) {
	if len(directMsg) > 0 {
		l.AddMessage(DebugLevel, directMsg...)
	}

	stackVal, _ := l.fields.Load("stack")
	messages := ensureStackType(stackVal)

	if len(messages) > 0 {
		entry := l.logger.WithFields(l.syncMapToLogFields())
		entry.Tracef("%+v", messages[0].Message)
	}

	l.clear()
}

func (l *Logger) clear() {
	l.fields.Range(func(key interface{}, value interface{}) bool {
		k, ok := key.(string)
		if ok {
			if k != FieldLogID {
				l.fields.Delete(key)
			}

			return true
		}

		return false
	})

	id := uuid.NewV1().String()

	l.fields.Store(FieldLogID, id)
	l.fields.Store("stack", []message{})
}

func (l *Logger) syncMapToLogFields() (fields log.Fields) {
	fields = make(log.Fields)

	l.fields.Range(func(key interface{}, value interface{}) bool {
		k, ok := key.(string)
		if ok {
			fields[k] = value
			return true
		}

		return false
	})

	return
}

var (
	regexPerkakas = regexp.MustCompile("perkakas/log")
	regexTest     = regexp.MustCompile("_test.go")
)

func (l *Logger) setCaller(level Level, msgs ...interface{}) {
	if msgs == nil || len(msgs) == 0 {
		return
	}

	var (
		fileName      []byte
		matchPerkakas bool
		matchTest     bool
	)

	for _, val := range msgs {
		var index int
		if val == "" {
			continue
		}

		index = DefaultIndex

		for ; ; index++ {
			_, file, _, ok := runtime.Caller(index)
			if !ok {
				break
			}
			fileName = append(fileName[:0], []byte(strings.ToLower(file))...)
			matchPerkakas = regexPerkakas.Match(fileName)
			matchTest = regexTest.Match(fileName)
			if !matchPerkakas {
				break
			}
			if matchPerkakas && matchTest {
				break
			}
		}

		if pc, file, line, ok := runtime.Caller(index); ok {
			fName := runtime.FuncForPC(pc).Name()

			err, ok := val.(error)
			if ok && err != nil {
				val = err.Error()
			}

			vmsg := message{
				Message:  val,
				Level:    level,
				File:     file,
				FuncName: fName,
				Line:     line,
			}

			l.addMessageStack(vmsg)
		}
	}
}

func (l *Logger) addMessageStack(msg ...message) {
	stackVal, _ := l.fields.Load("stack")
	stack := ensureStackType(stackVal)

	for _, val := range msg {
		stack = append(stack, val)
	}

	l.fields.Store("stack", stack)
}

func ensureStackType(stack interface{}) (val []message) {
	val, ok := stack.([]message)
	if !ok {
		panic("perkakas/log: stack trace is expecting a message but found other types")
	}

	return
}

func newLog(formatter log.Formatter, out io.Writer, level log.Level, reportCaller bool) (l *log.Logger) {
	l = log.New()
	l.SetFormatter(formatter)
	l.SetOutput(out)
	l.SetLevel(level)

	return
}

func newLogger(serviceName string, logID string) (logger *Logger) {
	formatter := &log.JSONFormatter{
		TimestampFormat: time.RFC3339,
		// PrettyPrint:     true,
		FieldMap: log.FieldMap{
			log.FieldKeyMsg: "log_message",
		},
	}

	newLogger := newLog(formatter, os.Stdout, log.TraceLevel, false)

	logger = new(Logger)
	logger.logger = newLogger
	logger.fields.Range(func(key interface{}, value interface{}) bool {
		logger.fields.Delete(key)
		return true
	})
	logger.fields.Store("stack", []message{})

	id := uuid.NewV1().String()
	if logID != "" {
		id = logID
	}

	logger.fields.Store(FieldLogID, id)
	logger.fields.Store(FieldServiceName, serviceName)
	logger.id = id
	logger.service = serviceName
	return
}

func NewLogger(serviceName string) (logger *Logger) {
	logger = newLogger(serviceName, "")
	return
}
