package task

import (
	"fmt"
	"slices"
	"time"
)

type Status string

const (
	DONE    = "done"
	IN_PROG = "in-progress"
	TODO    = "todo"
)

type Tasks struct {
	TodoList map[int]*Task `json:"todo_list"`
}

type Task struct {
	Id        int    `json:"id"`
	Subject   string `json:"subject"`
	Status    Status `json:"status"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func (t *Tasks) FetchTaskById(id int) (*Task, error) {
	if task, ok := t.TodoList[id]; ok {
		return task, nil
	}
	return nil, fmt.Errorf("NOT FOUND: unable to fetch task (ID: %d)", id)
}

func (t *Tasks) FetchAllTasks() []Task {
	size := len(t.TodoList)
	keys := make([]int, size)
	if size > 0 {
		idx := 0
		for i := range t.TodoList {
			keys[idx] = i
			idx++
		}
		slices.Sort(keys)
	}

	tasks := make([]Task, size)
	for i, k := range keys {
		tasks[i] = *(t.TodoList[k])
	}

	return tasks
}

func (t *Tasks) FetchTaskByStatus(status Status) []Task {
	taskList := []Task{}
	for _, task := range t.TodoList {
		if task.Status != status {
			continue
		}
		taskList = append(taskList, *task)
	}

	return taskList
}

func (t *Tasks) AddTask(subject string) int {
	if t.TodoList == nil {
		t.TodoList = make(map[int]*Task)
	}

	max := 0
	size := len(t.TodoList)
	if size > 0 {
		keys := make([]int, size)
		idx := 0
		for i := range t.TodoList {
			keys[idx] = i
			idx++
		}
		slices.Sort(keys)

		// grab the next sequential number in the list
		max = keys[size-1]
	}
	max++

	now := getTimestamp()
	t.TodoList[max] = &Task{
		Id:        max,
		Subject:   subject,
		Status:    TODO,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return max
}

func (t *Tasks) UpdateTaskById(id int, val string, status Status) (Task, error) {
	if task, err := t.FetchTaskById(id); err == nil {
		task.UpdatedAt = getTimestamp()
		if len(val) > 0 {
			task.Subject = val
		}
		if len(status) > 0 {
			task.Status = status
		}
		return *task, nil
	}

	return Task{}, fmt.Errorf("NOT FOUND: unable to update task (ID: %d)", id)
}

func (t *Tasks) DeleteTaskById(id int) (Task, error) {
	if task, err := t.FetchTaskById(id); err == nil {
		delete(t.TodoList, id)
		return *task, nil
	}

	return Task{}, fmt.Errorf("NOT FOUND: unable to delete task (ID: %d)", id)
}

func (t Task) String() string {
	return fmt.Sprintf("[%s] ID: %5d | Created On: %s | Updated On: %s | \n\t%s", t.Status, t.Id, t.CreatedAt, t.UpdatedAt, t.Subject)
}

func getTimestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
