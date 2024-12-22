package pkg

import (
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

type Logrus struct {
	*logrus.Logger
}

var Logger *logrus.Logger
var llogrus *Logrus

func InitLogrusConfig() Logrus {
	jsonFormatter := &logrus.JSONFormatter{
		PrettyPrint:     false,
		TimestampFormat: "2006-01-02 15:04:05",
	}
	Logger := logrus.New()
	Logger.SetOutput(os.Stdout)
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(jsonFormatter)
	logrus.SetLevel(logrus.InfoLevel)
	llogrus = &Logrus{Logger}
	return *llogrus
}

func (l Logrus) Level() log.Lvl {
	switch l.Logger.Level {
	case logrus.DebugLevel:
		return log.DEBUG
	case logrus.WarnLevel:
		return log.WARN
	case logrus.ErrorLevel:
		return log.ERROR
	case logrus.InfoLevel:
		return log.INFO
	default:
		l.Panic("Invalid level")
	}

	return log.OFF
}

func (l Logrus) SetHeader(_ string) {}

func (l Logrus) SetPrefix(s string) {}

func (l Logrus) Prefix() string {
	return ""
}

func (l Logrus) SetLevel(lvl log.Lvl) {
	switch lvl {
	case log.DEBUG:
		Logger.SetLevel(logrus.DebugLevel)
	case log.WARN:
		Logger.SetLevel(logrus.WarnLevel)
	case log.ERROR:
		Logger.SetLevel(logrus.ErrorLevel)
	case log.INFO:
		Logger.SetLevel(logrus.InfoLevel)
	default:
		l.Panic("Invalid level")
	}
}

func (l Logrus) Output() io.Writer {
	return l.Out
}

func (l Logrus) SetOutput(w io.Writer) {
	Logger.SetOutput(w)
}

func (l Logrus) Printj(j log.JSON) {
	Logger.WithFields(logrus.Fields(j)).Print()
}

func (l Logrus) Debugj(j log.JSON) {
	Logger.WithFields(logrus.Fields(j)).Debug()
}

func (l Logrus) Infoj(j log.JSON) {
	Logger.WithFields(logrus.Fields(j)).Info()
}

func (l Logrus) Warnj(j log.JSON) {
	Logger.WithFields(logrus.Fields(j)).Warn()
}

func (l Logrus) Errorj(j log.JSON) {
	Logger.WithFields(logrus.Fields(j)).Error()
}

func (l Logrus) Fatalj(j log.JSON) {
	Logger.WithFields(logrus.Fields(j)).Fatal()
}

func (l Logrus) Panicj(j log.JSON) {
	Logger.WithFields(logrus.Fields(j)).Panic()
}

func (l Logrus) Print(i ...interface{}) {
	Logger.Print(i[0].(string))
}

func (l Logrus) Debug(i ...interface{}) {
	Logger.Debug(i[0].(string))
}

func (l Logrus) Info(i ...interface{}) {
	Logger.Info(i[0].(string))
}

func (l Logrus) Warn(i ...interface{}) {
	Logger.Warn(i[0].(string))
}

func (l Logrus) Error(i ...interface{}) {
	Logger.Error(i[0].(string))
}

func (l Logrus) Fatal(i ...interface{}) {
	Logger.Fatal(i[0].(string))
}

func (l Logrus) Panic(i ...interface{}) {
	Logger.Panic(i[0].(string))
}
