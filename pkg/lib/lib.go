package lib

// Задача выражения
type Task struct {
	ID     int      // Номер действия
	ProbID int      // Номер выражения действия
	Oper   []string // Само действие
	Ans    float64  // Ответ
	Status int8     // Статус действия: 0 - не решено, 1 - решается, 2 - решено, 3 - ошибка.
}

// Выражение
type Expr struct {
	ID     int     // Номер выражения
	Oper   string  // Само выражение
	Tasks  []Task  // Задачи выражения
	Ans    float64 // Ответ
	Status int8    // Статус действия: 0 - не решено, 1 - решается, 2 - решено, 3 - ошибка.
}

// Выдаёт значение из конфига по полю
func GetConfig(field string) int {
	var Config = map[string]int{
		"TIME_ADDITION_MS":        1000,
		"TIME_SUBTRACTION_MS":     1000,
		"TIME_MULTIPLICATIONS_MS": 1000,
		"TIME_DIVISIONS_MS":       1000,
		"COMPUTING_AGENTS":        2,
		"COMPUTING_POWER":         5,
	}
	return Config[field]
}

// Стэк
type Stack struct {
	stack []string
}

// Создаёт экземпляр стэка
func Newstack() Stack {
	return Stack{stack: []string{}}
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
