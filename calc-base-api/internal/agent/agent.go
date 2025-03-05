package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Se623/calc-base-api/internal/lib"
)

func Agent(id int) {
	lib.Sugar.Infof("Agent %d: I am started", id)
	ticker := time.NewTicker(time.Duration(lib.TIME_REQUESTING_MS) * time.Millisecond)

	calcsin := make(chan lib.Task)
	calcsout := make(chan lib.Task)
	exprslot := lib.Expr{}
	busy := false

	for i := 0; i < lib.COMPUTING_POWER; i++ {
		go calculator(calcsin, calcsout, i+1, id)
	}

	for {
		select {
		case resTask := <-calcsout:
			if resTask.Status == 2 {
				lib.Sugar.Infof("Agent %d: Got task %d", id, resTask.ID)
				for i, v := range exprslot.Tasks {
					if v.ID == resTask.ID {
						exprslot.Tasks[i] = resTask
						if i <= len(exprslot.Tasks)-2 && exprslot.Tasks[i+1].Links[0] {
							exprslot.Tasks[i+1].Arg1 = exprslot.Tasks[i].Ans
							exprslot.Tasks[i+1].Links[0] = false
							lib.Sugar.Infof("Agent %d: Changed 1st link in task %d to number", id, exprslot.Tasks[i+1].ID)
						} else if i <= len(exprslot.Tasks)-3 && exprslot.Tasks[i+1].Links[1] {
							exprslot.Tasks[i+1].Arg2 = exprslot.Tasks[i].Ans
							exprslot.Tasks[i+1].Links[1] = false
							lib.Sugar.Infof("Agent %d: Changed 2nd link in task %d to number", id, exprslot.Tasks[i+1].ID)
						}
						if i == len(exprslot.Tasks)-1 {
							lib.Sugar.Infof("Agent %d: Task %d - last task, sending to orchestartor", id, exprslot.Tasks[i].ID)
							data, _ := json.Marshal(lib.TaskInc{ID: exprslot.ID, Result: exprslot.Tasks[i].Ans})
							r := bytes.NewReader(data)
							_, err := http.Post("http://localhost:8080/internal/task", "application/json", r)
							if err != nil {
								fmt.Println(err)
							}
							busy = false
						}
					}
				}
			}
		case <-ticker.C:
			if !busy {
				lib.Sugar.Infof("Agent %d: Trying to contact orchestrator for more expressions", id)
				resp, err := http.Get("http://localhost:8080/internal/task")
				if err != nil {
					lib.Sugar.Errorf("Agent %d: Something went wrong, aborting contact. Error: %s", id, err.Error())
					continue
				}
				if resp.StatusCode == 200 {
					decoder := json.NewDecoder(resp.Body)
					var expr lib.Expr
					err = decoder.Decode(&expr)
					if err != nil {
						lib.Sugar.Errorf("Agent %d: Something went wrong, aborting contact. Error: %s", id, err.Error())
						continue
					}
					lib.Sugar.Infof("Agent %d: Got expression %d", id, expr.ID)
					exprslot = expr
					busy = true
				} else {
					lib.Sugar.Infof("Agent %d: Orchestrator request unsuccessful, code: %d", id, resp.StatusCode)
				}
			}
		default:
			if busy {
				for i, v := range exprslot.Tasks {
					if v.Status == 0 && !v.Links[0] && !v.Links[1] {
						lib.Sugar.Infof("Agent %d: Got undestributed task %d, sending to calculators", id, v.ID)
						calcsin <- exprslot.Tasks[i]
						exprslot.Tasks[i].Status = 1
					}
				}
			}
		}
	}
}

func calculator(comm chan lib.Task, result chan lib.Task, id int, agid int) {
	lib.Sugar.Infof("Calc %d-%d: I am started", agid, id)
	for task := range comm {
		lib.Sugar.Infof("Calc %d-%d: Got task %d: %d %s %d", agid, id, task.ID, task.Arg1, task.Operation, task.Arg2)
		timer := time.NewTimer(time.Duration(task.Operation_time) * time.Millisecond)
		<-timer.C
		if task.Operation == "+" {
			task.Ans = task.Arg1 + task.Arg2
		} else if task.Operation == "-" {
			task.Ans = task.Arg1 - task.Arg2
		} else if task.Operation == "*" {
			task.Ans = task.Arg1 * task.Arg2
		} else if task.Operation == "/" {
			task.Ans = task.Arg1 / task.Arg2
		}
		lib.Sugar.Infof("Calc %d-%d: Got answer to the task %d: %d", agid, id, task.ID, task.Ans)
		task.Status = 2
		result <- task
	}
}
