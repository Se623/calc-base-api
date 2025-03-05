package lib

import (
	"sync"

	"go.uber.org/zap"
)

// Константы
var TIME_ADDITION_MS = 1000
var TIME_SUBTRACTION_MS = 1000
var TIME_MULTIPLICATIONS_MS = 1000
var TIME_DIVISIONS_MS = 1000
var TIME_REQUESTING_MS = 5000
var COMPUTING_AGENTS = 2
var COMPUTING_POWER = 2

// Логгер
var Sugar *zap.SugaredLogger

func InitLogger() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	Sugar = logger.Sugar()
}

// "Сырое" выражение (Инфиксное, не разделённое)
type Raw struct {
	Expression string `json:"expression"`
}

// Выражение
type Expr struct {
	ID     int     `json:"id"` // Номер выражения
	Oper   string  // Само выражение
	Tasks  []Task  `json:"tasks"` // Задачи выражения
	Ans    float64 // Ответ
	Status int8    // Статус действия: 0 - не решено, 1 - решается, 2 - решено.
}

// Выражение, которое отображается в API
type ExprDsp struct {
	ID     int     `json:"id"`
	Status string  `json:"status"`
	Result float64 `json:"result"`
}

// Массив выражений, которые отображаются в API
type DspArr struct {
	Expressions []ExprDsp `json:"expressions"`
}

// Задача
type Task struct {
	ID             int     `json:"id"` // ID
	ProbID         int     // Номер выражения действия
	Links          [2]bool `json:"links"`
	Arg1           float64 `json:"arg1"`           // Первое число
	Arg2           float64 `json:"arg2"`           // Второе число
	Operation      string  `json:"operation"`      // Операция
	Operation_time int     `json:"operation_time"` // Время выполнения
	Ans            float64 // Ответ
	Status         int8    // Статус действия: 0 - не решено, 1 - решается, 2 - решено.
}

// Ответ на задачу по ID
type TaskInc struct {
	ID     int     `json:"id"`
	Result float64 `json:"float64"`
}

// База данных выражений (С мьюьексом)
type ExprDB struct {
	Exprs []Expr
	Mux   sync.Mutex
}

// Создать базу данных
func NewExprDB() *ExprDB {
	return &ExprDB{[]Expr{}, sync.Mutex{}}
}

// Стэк
type Stack struct {
	stack []string
}

// Создаёт экземпляр стэка
func Newstack() *Stack {
	return &Stack{stack: []string{}}
}

// Добавляет элемент в стак
func (s *Stack) Push(val string) {
	s.stack = append(s.stack, val)
}

// Просматривает последний элемент в стэке
func (s *Stack) GetTop() string {
	if len(s.stack) != 0 {
		return s.stack[len(s.stack)-1]
	} else {
		return ""
	}
}

// Вынимает последний элемент из стэка
func (s *Stack) Pop() string {
	if len(s.stack) != 0 {
		r := s.stack[len(s.stack)-1]
		s.stack = s.stack[:len(s.stack)-1]
		return r
	} else {
		return ""
	}
}
