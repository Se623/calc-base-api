package orchestrator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Se623/calc-base-api/internal/lib"
	"github.com/Se623/calc-base-api/pkg/rpn"
)

var exprs = lib.NewExprDB()

func Displayer(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		exprArr := lib.DspArr{}
		exprs.Mux.Lock()
		for _, v := range exprs.Exprs {
			if v.Status == 0 {
				exprArr.Expressions = append(exprArr.Expressions, lib.ExprDsp{ID: v.ID, Status: "Queued", Result: -1})
			} else if v.Status == 1 {
				exprArr.Expressions = append(exprArr.Expressions, lib.ExprDsp{ID: v.ID, Status: "Solving", Result: -1})
			} else if v.Status == 2 {
				exprArr.Expressions = append(exprArr.Expressions, lib.ExprDsp{ID: v.ID, Status: "Solved", Result: v.Ans})
			}
		}
		exprArrPack, err := json.Marshal(exprArr)
		if err != nil {
			http.Error(w, "Error: Something invalid", http.StatusInternalServerError)
			exprs.Mux.Unlock()
			return
		}
		fmt.Fprint(w, string(exprArrPack))
		exprs.Mux.Unlock()
	} else {
		var exprPack []byte
		var err error
		idInt, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Error: ID not found", http.StatusNotFound)
			return
		}
		exprs.Mux.Lock()
		for _, v := range exprs.Exprs {
			if v.ID == idInt {
				if v.Status == 0 {
					exprPack, err = json.Marshal(lib.ExprDsp{ID: v.ID, Status: "Queued", Result: -1})
				} else if v.Status == 1 {
					exprPack, err = json.Marshal(lib.ExprDsp{ID: v.ID, Status: "Solving", Result: -1})
				} else if v.Status == 2 {
					exprPack, err = json.Marshal(lib.ExprDsp{ID: v.ID, Status: "Solved", Result: v.Ans})
				}
			}
		}
		if err != nil {
			http.Error(w, "Error: Something invalid", http.StatusInternalServerError)
			exprs.Mux.Unlock()
			return
		}
		fmt.Fprint(w, string(exprPack))
		exprs.Mux.Unlock()
	}
}

func Spliter(w http.ResponseWriter, r *http.Request) {
	pr := []string{}
	res := [][]string{}
	opers := []lib.Task{}

	rpnstack := lib.Newstack()

	decoder := json.NewDecoder(r.Body)
	var resp lib.Raw
	err := decoder.Decode(&resp)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error: Invalid JSON", http.StatusInternalServerError)
		return
	}

	rpnarr, err := rpn.InfixToPostfix(resp.Expression)

	if err != nil {
		http.Error(w, "Error: Invalid Input", http.StatusUnprocessableEntity)
		return
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

		var a float64
		var b float64

		optime := 0
		links := [2]bool{false, false}

		if v[2] == "+" {
			optime = lib.TIME_ADDITION_MS
		} else if v[2] == "-" {
			optime = lib.TIME_SUBTRACTION_MS
		} else if v[2] == "*" {
			optime = lib.TIME_MULTIPLICATIONS_MS
		} else if v[2] == "/" {
			optime = lib.TIME_DIVISIONS_MS
		}

		if v[0] == "L" {
			links[0] = true
			a = -1
		} else {
			a, _ = strconv.ParseFloat(v[0], 64)
		}
		if v[1] == "L" {
			links[1] = true
			b = -1
		} else {
			b, _ = strconv.ParseFloat(v[1], 64)
		}

		if i == len(res)-1 {
			opers = append(opers, lib.Task{ID: num, ProbID: 0, Links: links, Arg1: a, Arg2: b, Operation: v[2], Operation_time: optime, Ans: 0, Status: 0})
		}
		opers = append(opers, lib.Task{ID: num, ProbID: 0, Links: links, Arg1: a, Arg2: b, Operation: v[2], Operation_time: optime, Ans: 0, Status: 0})
		num++
	}

	exprs.Mux.Lock()
	exprs.Exprs = append(exprs.Exprs, lib.Expr{ID: len(exprs.Exprs), Oper: resp.Expression, Tasks: opers, Ans: 0, Status: 0})
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"id": "%d"}`, len(exprs.Exprs)-1)
	exprs.Mux.Unlock()
}

func Distributor(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		exprs.Mux.Lock()
		for i, v := range exprs.Exprs {
			if v.Status == 0 {
				exprPack, err := json.Marshal(v)
				if err != nil {
					http.Error(w, "Error: Something invalid", http.StatusInternalServerError)
					exprs.Mux.Unlock()
					return
				}
				fmt.Fprintf(w, string(exprPack), len(exprs.Exprs)-1)
				exprs.Exprs[i].Status = 1
				exprs.Mux.Unlock()
				return
			}
		}
		exprs.Mux.Unlock()
		http.Error(w, "Error: No expressions", http.StatusNotFound)
	} else if r.Method == http.MethodPost {
		decoder := json.NewDecoder(r.Body)
		var resp lib.TaskInc
		err := decoder.Decode(&resp)
		if err != nil {
			lib.Sugar.Errorf("Orchestrator: Error")
			w.WriteHeader(http.StatusInternalServerError)
		}
		lib.Sugar.Infof("Orchestrator: Got expression %d", resp.ID)

		exprs.Mux.Lock()
		for i, v := range exprs.Exprs {
			if v.ID == resp.ID {
				lib.Sugar.Infof("Orchestrator: Replacing expression %d in database", resp.ID)
				exprs.Exprs[i].Ans = resp.Result
				exprs.Exprs[i].Status = 2
			}
		}
		exprs.Mux.Unlock()
	}
}
