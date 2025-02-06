package filehandler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/rltran-codex/mytask-cli/task"
)

var TaskFile string

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	// check if appdata directory exists
	baseDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	appDir := filepath.Join(baseDir, "appdata")
	TaskFile = filepath.Join(appDir, "tasks.json")

	os.Mkdir(appDir, 0744)
	if _, err := os.Stat(TaskFile); os.IsNotExist(err) {
		f, err := os.Create(TaskFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
	}
}

func ReadJson(obj *task.Tasks) {
	f, err := os.Open(TaskFile)
	handleError(err)
	defer f.Close()

	// load .json into a temp struct
	temp := struct {
		TodoList map[string]*task.Task `json:"todo_list"`
	}{}
	obj.TodoList = make(map[int]*task.Task)
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&temp)
	if err != nil {
		log.Println("WARNING: task.json was empty or corrupted, starting fresh")
		return
	}

	// filter out any keys that are not int
	for k, v := range temp.TodoList {
		var key int
		if i, err := fmt.Sscanf(k, "%d", &key); i == 1 && err == nil {
			obj.TodoList[key] = v
		}
		handleError(err)
	}
}

func Update(data task.Tasks) {
	// use a temp file to try and write
	tempDir, err := os.MkdirTemp("", "task-cli-")
	handleError(err)
	defer os.RemoveAll(tempDir)

	tempFile, err := os.CreateTemp(tempDir, "task-tmp-*.json")
	handleError(err)
	defer tempFile.Close()

	encoder := json.NewEncoder(tempFile)
	err = encoder.Encode(data)
	handleError(err)

	out, err := os.OpenFile(TaskFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
	handleError(err)
	defer out.Close()

	// copy task-tmp-*.json to task.json
	tempFile.Seek(0, io.SeekStart)
	_, err = io.Copy(out, tempFile)
	handleError(err)
}
