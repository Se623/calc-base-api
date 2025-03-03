package calc

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Se623/calc-base-api/internal/rpn"
	"github.com/Se623/calc-base-api/pkg/lib"
	"go.uber.org/zap"
)

func Orchestrator(entry chan string, sugar *zap.SugaredLogger) {
	rpnstack := lib.Newstack()
	agentsin := make(chan lib.Expr)
	agentsout := make(chan lib.Expr)
	exprs := []lib.Expr{}

	sugar.Info("Orchestrator is started")

	for i := 0; i < lib.GetConfig("COMPUTING_AGENTS"); i++ {
		go Agent(agentsin, agentsout, sugar, i+1)
	}

	for {
		select {
		case infoper := <-entry:
			sugar.Infof("Orchestrator: Got the raw (%s)", infoper)
			pr := []string{}
			res := [][]string{}
			opers := []lib.Task{}

			rpnarr, err := rpn.InfixToPostfix(infoper)

			if err != nil {
				sugar.Errorf("Orchestrator: Raw (%s) error", infoper)
				exprs = append(exprs, lib.Expr{ID: len(exprs), Oper: infoper, Tasks: opers, Ans: 0, Status: 3})
				break
			}

			// Разделение выражений на задания
			for _, v := range rpnarr {
				if _, err := strconv.ParseFloat(v, 64); err == nil {
					rpnstack.Push(v)
				} else {
					pr = append(pr, rpnstack.Pop())
					pr = append(pr, rpnstack.Pop())

					pr = append(pr, v)
					rpnstack.Push("L")

					res = append(res, pr)
					pr = []string{}
				}
			}

			var num int

			for i, v := range res {
				v[0], v[1] = v[1], v[0]
				if i == len(res)-1 {
					opers = append(opers, lib.Task{ID: num, ProbID: 0, Oper: v, Ans: 0, Status: 0})
					break
				}
				opers = append(opers, lib.Task{ID: num, ProbID: 0, Oper: v, Ans: 0, Status: 0})
				num++
			}

			exprs = append(exprs, lib.Expr{ID: len(exprs), Oper: infoper, Tasks: opers, Ans: 0, Status: 0})
			fmt.Println(opers)
			sugar.Infof("Orchestrator: Expresssion %d is splited into tasks and loaded", len(exprs)-1)

		case expr := <-agentsout:
			if expr.Status == 2 {
				sugar.Infof("Orchestrator: Got the solved expression %d, updating expression array", expr.ID)
				fmt.Println(expr)
				for i, v := range exprs {
					if v.ID == expr.ID {
						exprs[i] = expr
					}
				}
			}
		default:
			for i, v := range exprs {
				if v.Status == 0 {
					sugar.Infof("Orchestrator: Got the undestributed expression %d, sending to agents", v.ID)
					agentsin <- exprs[i]
					exprs[i].Status = 1
				}
			}
		}
	}
}

func Agent(comm chan lib.Expr, result chan lib.Expr, sugar *zap.SugaredLogger, id int) {
	calcsin := make(chan lib.Task)
	calcsout := make(chan lib.Task)
	exprslot := lib.Expr{}
	busy := false

	sugar.Infof("Agent %d is started", id)

	for i := 0; i < lib.GetConfig("COMPUTING_AGENTS"); i++ {
		go Calculator(calcsin, calcsout, sugar, i+1, id)
	}
	for {
		select {
		case expr := <-comm:
			sugar.Infof("Agent %d: Got the expression %d", id, expr.ID)
			if !busy {
				exprslot = expr
				busy = true
			} else {
				comm <- expr
			}
		case resTask := <-calcsout:
			if resTask.Status == 2 {
				sugar.Infof("Agent %d: Got the result to the task %d", id, resTask.ID)
				for i, v := range exprslot.Tasks {
					if v.ID == resTask.ID {
						exprslot.Tasks[i] = resTask
						if i <= len(exprslot.Tasks)-2 && exprslot.Tasks[i+1].Oper[0] == "L" {
							exprslot.Tasks[i+1].Oper[0] = fmt.Sprint(exprslot.Tasks[i].Ans)
						} else if i <= len(exprslot.Tasks)-3 && exprslot.Tasks[i+1].Oper[1] == "L" {
							exprslot.Tasks[i+1].Oper[1] = fmt.Sprint(exprslot.Tasks[i].Ans)
						}
						if i == len(exprslot.Tasks)-1 {
							sugar.Infof("Agent %d: Task %d is the last one, sending the expression to orchestrator", id, resTask.ID)
							exprslot.Ans = exprslot.Tasks[i].Ans
							exprslot.Status = 2
							result <- exprslot
							busy = false
						}
					}
				}
			}
		default:
			for i, v := range exprslot.Tasks {
				if v.Status == 0 && v.Oper[0] != "L" && v.Oper[1] != "L" {
					sugar.Infof("Agent %d: Got the undestributed task %d, sending to calculators", id, v.ID)
					calcsin <- exprslot.Tasks[i]
					exprslot.Tasks[i].Status = 1
				}
			}
		}
	}
}

func Calculator(comm chan lib.Task, result chan lib.Task, sugar *zap.SugaredLogger, id int, agid int) {
	sugar.Infof("Calculator %d-%d is started", agid, id)
	for {
		select {
		case task := <-comm:
			sugar.Infof("Calculator %d-%d: Got the task %d", agid, id, task.ID)
			if task.Oper[2] == "+" {
				sugar.Infof("Calculator %d-%d: Task %d - addition, starting timer", agid, id, task.ID)
				timer := time.NewTimer(time.Duration(lib.GetConfig("TIME_ADDITION_MS")))
				<-timer.C
				a, _ := strconv.ParseFloat(task.Oper[0], 64)
				b, _ := strconv.ParseFloat(task.Oper[1], 64)
				task.Ans = a + b
			} else if task.Oper[2] == "-" {
				sugar.Infof("Calculator %d-%d: Task %d - substraction, starting timer", agid, id, task.ID)
				timer := time.NewTimer(time.Duration(lib.GetConfig("TIME_SUBTRACTION_MS")))
				<-timer.C
				a, _ := strconv.ParseFloat(task.Oper[0], 64)
				b, _ := strconv.ParseFloat(task.Oper[1], 64)
				task.Ans = a - b
			} else if task.Oper[2] == "*" {
				sugar.Infof("Calculator %d-%d: Task %d - multiplication, starting timer", agid, id, task.ID)
				timer := time.NewTimer(time.Duration(lib.GetConfig("TIME_MULTIPLICATIONS_MS")))
				<-timer.C
				a, _ := strconv.ParseFloat(task.Oper[0], 64)
				b, _ := strconv.ParseFloat(task.Oper[1], 64)
				task.Ans = a * b
			} else if task.Oper[2] == "/" {
				sugar.Infof("Calculator %d-%d: Task %d - division, starting timer", agid, id, task.ID)
				timer := time.NewTimer(time.Duration(lib.GetConfig("TIME_DIVISIONS_MS")))
				<-timer.C
				a, _ := strconv.ParseFloat(task.Oper[0], 64)
				b, _ := strconv.ParseFloat(task.Oper[1], 64)
				task.Ans = a / b
			}
			sugar.Infof("Calculator %d-%d: Task %d - timer ended, result: %g", agid, id, task.ID, task.Ans)
			task.Status = 2
			result <- task
		}
	}
}
