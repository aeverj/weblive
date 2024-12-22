package utils

import (
	"fmt"
	"sync"
	"time"
)

type Spinner struct {
	done      chan bool
	counter   int
	message   string
	mu        sync.Mutex
	isRunning bool
}

// Option 定义 Spinner 的可选配置项
type Option func(*Spinner)

// NewSpinner 创建新的进度显示器
func NewSpinner(options ...Option) *Spinner {
	s := &Spinner{
		done:      make(chan bool),
		message:   "Processing tasks",
		isRunning: false,
	}

	for _, option := range options {
		option(s)
	}

	return s
}

// WithMessage 设置显示消息
func WithMessage(message string) Option {
	return func(s *Spinner) {
		s.message = message
	}
}

// Increment 增加计数
func (s *Spinner) Increment() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counter++
}

// GetCounter 获取当前计数
func (s *Spinner) GetCounter() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.counter
}

// Start 开始显示进度
func (s *Spinner) Start() {
	if s.isRunning {
		return
	}
	s.isRunning = true

	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

	s.moveToBottom() // 移动到终端底部

	go func() {
		idx := 0
		for {
			select {
			case <-s.done:
				return
			default:
				s.mu.Lock()
				counter := s.counter
				s.mu.Unlock()

				s.clearLine()    // 清除当前行
				s.moveToBottom() // 移动到底部

				fmt.Printf("\r%s %s... (%d completed)", frames[idx], s.message, counter)

				idx = (idx + 1) % len(frames)
				time.Sleep(time.Millisecond * 80)
			}
		}
	}()
}

// Stop 停止进度显示
func (s *Spinner) Stop() {
	if !s.isRunning {
		return
	}
	s.isRunning = false
	close(s.done)

	s.mu.Lock()
	counter := s.counter
	s.mu.Unlock()

	s.clearLine()    // 清除当前行
	s.moveToBottom() // 移动到底部

	fmt.Printf("\r✓ %s completed! (%d total)\n", s.message, counter)
}

// clearLine 清除当前行
func (s *Spinner) clearLine() {
	fmt.Print("\033[K") // 清除当前行
}

// moveToBottom 移动到终端最下方
func (s *Spinner) moveToBottom() {
	fmt.Print("\033[1B") // 移动到很远的下方
}
