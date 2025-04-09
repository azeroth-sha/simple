package simple

// Logger 定义了日志记录器的接口，支持不同级别的日志记录和格式化输出
type Logger interface {
	// Debug 记录调试级别的日志，参数为任意类型的值
	Debug(v ...interface{})

	// Info 记录信息级别的日志，参数为任意类型的值
	Info(v ...interface{})

	// Warn 记录警告级别的日志，参数为任意类型的值
	Warn(v ...interface{})

	// Error 记录错误级别的日志，参数为任意类型的值
	Error(v ...interface{})

	// Fatal 记录致命错误级别的日志，参数为任意类型的值，通常会导致程序退出
	Fatal(v ...interface{})

	// Debugf 记录调试级别的日志，支持格式化字符串和参数
	Debugf(format string, v ...interface{})

	// Infof 记录信息级别的日志，支持格式化字符串和参数
	Infof(format string, v ...interface{})

	// Warnf 记录警告级别的日志，支持格式化字符串和参数
	Warnf(format string, v ...interface{})

	// Errorf 记录错误级别的日志，支持格式化字符串和参数
	Errorf(format string, v ...interface{})

	// Fatalf 记录致命错误级别的日志，支持格式化字符串和参数，通常会导致程序退出
	Fatalf(format string, v ...interface{})
}
