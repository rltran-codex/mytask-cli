package test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/rltran-codex/mytask-cli/filehandler"
	"github.com/rltran-codex/mytask-cli/task"
)

var timeFormat = "2006-01-02 15:04:05"
var subjectFormat = "Test %d"

func TestTaskPrint(t *testing.T) {
	now := time.Now().Format(timeFormat)
	testObj := task.Task{
		Id:        1,
		Status:    task.DONE,
		Subject:   "Test the auto printing of this task object",
		CreatedAt: now,
		UpdatedAt: now,
	}

	testStr := fmt.Sprintf("%v", testObj)
	expeStr := fmt.Sprintf("[done] ID:     1 | Created On: %s | Updated On: %s | \n\tTest the auto printing of this task object", now, now)
	if testStr != expeStr {
		t.Error("Task.String() is not correctly formatted")
	}
}

func TestAddTask(t *testing.T) {
	subject := "Test adding task"
	eId := 1
	tasks := task.Tasks{}
	tasks.AddTask(subject)

	if a, ok := tasks.TodoList[eId]; !ok {
		t.Errorf("failed to add to todo list: %v", a)
	} else {
		if a.Id != eId {
			t.Errorf("Expected id: %d. Actual: %d", eId, a.Id)
		}
		if a.Subject != subject {
			t.Errorf("Expected subject: %s. Actual: %s", subject, a.Subject)
		}
		if a.Status != task.TODO {
			t.Errorf("Expected status: %s. Actual: %s", task.TODO, a.Status)
		}
	}
}

func TestGetAll(t *testing.T) {
	tasks := task.Tasks{TodoList: make(map[int]*task.Task)}
	for i := 0; i < 100; i++ {
		var status task.Status
		if i%2 == 0 {
			status = task.DONE
		} else {
			status = task.IN_PROG
		}
		tasks.TodoList[i] = &task.Task{
			Id:        i,
			Status:    status,
			Subject:   fmt.Sprintf(subjectFormat, i),
			CreatedAt: time.Now().Format(timeFormat),
			UpdatedAt: time.Now().Format(timeFormat),
		}
	}

	tList := tasks.FetchAllTasks()
	for i := 0; i < 100; i++ {
		task := tList[i]
		if i != task.Id {
			t.Errorf("FetchAllTasks is not in sequential order")
			t.Fail()
		}
	}
}

func TestGetByStatus(t *testing.T) {
	tasks := task.Tasks{TodoList: make(map[int]*task.Task)}
	for i := 0; i < 100; i++ {
		var status task.Status
		if i%2 == 0 {
			status = task.DONE
		} else {
			status = task.IN_PROG
		}
		tasks.TodoList[i] = &task.Task{
			Id:        i,
			Status:    status,
			Subject:   fmt.Sprintf(subjectFormat, i),
			CreatedAt: time.Now().Format(timeFormat),
			UpdatedAt: time.Now().Format(timeFormat),
		}
	}

	doneTasks := tasks.FetchTaskByStatus(task.DONE)
	inProgTasks := tasks.FetchTaskByStatus(task.IN_PROG)

	if len(doneTasks) != 50 {
		t.Errorf("Done tasks length %d != 50", len(doneTasks))
	}
	if len(inProgTasks) != 50 {
		t.Errorf("In-Progress tasks length %d != 50", len(inProgTasks))
	}
	for i := 0; i < 100; i++ {
		var status task.Status
		if i%2 == 0 {
			status = task.DONE
		} else {
			status = task.IN_PROG
		}
		task := tasks.TodoList[i]
		if task.Status != status {
			t.Errorf("Expected %s. Actual %s.", status, task.Status)
		}
	}
}

func TestDeleteTask(t *testing.T) {
	tasks := task.Tasks{TodoList: make(map[int]*task.Task)}
	for i := 0; i < 100; i++ {
		tasks.TodoList[i] = &task.Task{
			Id:        i,
			Status:    task.TODO,
			Subject:   fmt.Sprintf(subjectFormat, i),
			CreatedAt: time.Now().Format(timeFormat),
			UpdatedAt: time.Now().Format(timeFormat),
		}
	}

	for i := 0; i < 100; i++ {
		// remove task then try to fetch it. if remove is successful, then fetching will result in err
		old, err := tasks.DeleteTaskById(i)
		if old.Id != i {
			t.Errorf("DeleteTaskById returned the wrong object:\n%+v", old)
		}
		if err != nil {
			t.Errorf("unsuccessfully deleted task %d: %v", i, err)
		}
		_, err = tasks.FetchTaskById(i)
		if err == nil {
			t.Errorf("task with id %d still exists", i)
		}
	}
}

func TestWriteJson(t *testing.T) {
	tasks := task.Tasks{TodoList: make(map[int]*task.Task)}
	for i := 1; i <= 5; i++ {
		tasks.TodoList[i] = &task.Task{
			Id:        i,
			Status:    task.DONE,
			Subject:   fmt.Sprintf(subjectFormat, i),
			CreatedAt: time.Now().Format(timeFormat),
			UpdatedAt: time.Now().Format(timeFormat),
		}
	}

	filehandler.Update(tasks)
	if _, err := os.Stat(filehandler.TaskFile); os.IsNotExist(err) {
		t.Errorf("did not successfully create tasks.json")
	}
}

func TestLoadJson(t *testing.T) {
	tasks := task.Tasks{}
	filehandler.ReadJson(&tasks)
}
