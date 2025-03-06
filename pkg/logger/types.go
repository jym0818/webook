package logger

type LoggerV1 interface {
	Debug(msg string, args ...Field)
	Info(msg string, args ...Field)
	Warn(msg string, args ...Field)
	Error(msg string, args ...Field)
}

type Field struct {
	Key   string
	Value any
}

func LoggerExample() {
	var l LoggerV1
	phone := "152xxxx1234"
	l.Info("用户未注册", Field{
		Key:   "phone",
		Value: phone,
	})
}
