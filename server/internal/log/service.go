package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// LogLevel 日志级别
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// String 返回日志级别字符串
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// LogEntry 日志条目
type LogEntry struct {
	Timestamp time.Time              `json:"timestamp"`
	Level     LogLevel               `json:"level"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
	RequestID string                 `json:"request_id,omitempty"`
	UserID    int64                  `json:"user_id,omitempty"`
}

// LogService 日志服务
type LogService struct {
	mu       sync.Mutex
	file     *os.File
	logger   *log.Logger
	logDir   string
	logLevel LogLevel
}

// NewLogService 创建日志服务
func NewLogService(logDir string, level LogLevel) (*LogService, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("create log directory: %w", err)
	}

	service := &LogService{
		logDir:   logDir,
		logLevel: level,
	}

	if err := service.rotateLog(); err != nil {
		return nil, fmt.Errorf("rotate log: %w", err)
	}

	return service, nil
}

// rotateLog 轮转日志文件
func (s *LogService) rotateLog() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 关闭现有文件
	if s.file != nil {
		s.file.Close()
	}

	// 生成日志文件名
	date := time.Now().Format("2006-01-02")
	logFile := filepath.Join(s.logDir, fmt.Sprintf("app-%s.log", date))

	// 打开日志文件（追加模式）
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open log file: %w", err)
	}

	s.file = file
	s.logger = log.New(io.MultiWriter(os.Stdout, file), "", 0)

	return nil
}

// Log 记录日志
func (s *LogService) Log(entry LogEntry) {
	if entry.Level < s.logLevel {
		return
	}

	// 检查是否需要轮转（每天）
	if s.file != nil {
		stat, err := s.file.Stat()
		if err == nil {
			fileDate := stat.ModTime().Format("2006-01-02")
			today := time.Now().Format("2006-01-02")
			if fileDate != today {
				s.rotateLog()
			}
		}
	}

	// 格式化日志消息
	timestamp := entry.Timestamp.Format("2006-01-02 15:04:05")
	level := entry.Level.String()

	// 构建日志消息
	msg := fmt.Sprintf("[%s] [%s] %s", timestamp, level, entry.Message)

	// 添加字段
	if len(entry.Fields) > 0 {
		for k, v := range entry.Fields {
			msg += fmt.Sprintf(" %s=%v", k, v)
		}
	}

	// 添加请求ID和用户ID
	if entry.RequestID != "" {
		msg += fmt.Sprintf(" request_id=%s", entry.RequestID)
	}
	if entry.UserID > 0 {
		msg += fmt.Sprintf(" user_id=%d", entry.UserID)
	}

	s.logger.Println(msg)
}

// Debug 记录调试日志
func (s *LogService) Debug(msg string, fields ...map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     DEBUG,
		Message:   msg,
	}
	if len(fields) > 0 {
		entry.Fields = fields[0]
	}
	s.Log(entry)
}

// Info 记录信息日志
func (s *LogService) Info(msg string, fields ...map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     INFO,
		Message:   msg,
	}
	if len(fields) > 0 {
		entry.Fields = fields[0]
	}
	s.Log(entry)
}

// Warn 记录警告日志
func (s *LogService) Warn(msg string, fields ...map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     WARN,
		Message:   msg,
	}
	if len(fields) > 0 {
		entry.Fields = fields[0]
	}
	s.Log(entry)
}

// Error 记录错误日志
func (s *LogService) Error(msg string, fields ...map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     ERROR,
		Message:   msg,
	}
	if len(fields) > 0 {
		entry.Fields = fields[0]
	}
	s.Log(entry)
}

// WithRequest 添加请求上下文
func (s *LogService) WithRequest(requestID string, userID int64) *RequestLogger {
	return &RequestLogger{
		service:   s,
		requestID: requestID,
		userID:    userID,
	}
}

// Close 关闭日志服务
func (s *LogService) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.file != nil {
		return s.file.Close()
	}
	return nil
}

// RequestLogger 请求日志记录器
type RequestLogger struct {
	service   *LogService
	requestID string
	userID    int64
}

// Debug 记录调试日志
func (r *RequestLogger) Debug(msg string, fields ...map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     DEBUG,
		Message:   msg,
		RequestID: r.requestID,
		UserID:    r.userID,
	}
	if len(fields) > 0 {
		entry.Fields = fields[0]
	}
	r.service.Log(entry)
}

// Info 记录信息日志
func (r *RequestLogger) Info(msg string, fields ...map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     INFO,
		Message:   msg,
		RequestID: r.requestID,
		UserID:    r.userID,
	}
	if len(fields) > 0 {
		entry.Fields = fields[0]
	}
	r.service.Log(entry)
}

// Warn 记录警告日志
func (r *RequestLogger) Warn(msg string, fields ...map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     WARN,
		Message:   msg,
		RequestID: r.requestID,
		UserID:    r.userID,
	}
	if len(fields) > 0 {
		entry.Fields = fields[0]
	}
	r.service.Log(entry)
}

// Error 记录错误日志
func (r *RequestLogger) Error(msg string, fields ...map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     ERROR,
		Message:   msg,
		RequestID: r.requestID,
		UserID:    r.userID,
	}
	if len(fields) > 0 {
		entry.Fields = fields[0]
	}
	r.service.Log(entry)
}