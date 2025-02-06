package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/rltran-codex/mytask-cli/filehandler"
	"github.com/rltran-codex/mytask-cli/task"
)

var (
	addCmd    *flag.FlagSet
	updateCmd *flag.FlagSet
	deleteCmd *flag.FlagSet
	listCmd   *flag.FlagSet
)

var Tasks *task.Tasks

func main() {
	// make a flag group for add
	addCmd = flag.NewFlagSet("add", flag.ExitOnError)
	subAdd := addCmd.String("subj", "", "REQUIRED. Define the new task's subject.")

	// make a flag group for update
	updateCmd = flag.NewFlagSet("update", flag.ExitOnError)
	uId := updateCmd.Int("id", -1, "REQUIRED. Update task with id.")
	subUpdate := updateCmd.String("subj", "", "OPTIONAL. Update task's subject.")
	statusUpdate := updateCmd.String("stat", "", "OPTIONAL. Update task's status. [todo, in-progress, done]")

	// make a flag group for delete
	deleteCmd = flag.NewFlagSet("delete", flag.ExitOnError)
	dId := deleteCmd.Int("id", -1, "REQUIRED. Delete task with id.")

	// make a flag group for list
	listCmd = flag.NewFlagSet("list", flag.ExitOnError)
	statusList := listCmd.String("stat", "", "OPTIONAL. List all tasks with specified status. [todo, in-progress, done]")

	// init Tasks
	Tasks = &task.Tasks{}
	filehandler.ReadJson(Tasks)
	defer filehandler.Update(*Tasks)

	if len(os.Args) <= 1 {
		printUsage()
	}
	switch os.Args[1] {
	case "add":
		addCmd.Parse(os.Args[2:])
		AddHandler(*subAdd)
	case "update":
		updateCmd.Parse(os.Args[2:])
		UpdateHandler(*uId, *subUpdate, *statusUpdate)
	case "delete":
		deleteCmd.Parse(os.Args[2:])
		DeleteHandler(*dId)
	case "list":
		listCmd.Parse(os.Args[2:])
		ListHandler(*statusList)
	default:
		printUsage()
	}
}

func printUsage() {
	cmds := []struct {
		cmd  string
		desc string
	}{
		{cmd: "add", desc: "Add a new task"},
		{cmd: "update", desc: "Update an existing task"},
		{cmd: "delete", desc: "Delete a task"},
		{cmd: "list", desc: "List tasks based on criteria"},
	}
	fmt.Println("Usage:")
	fmt.Println("  task-cli <command> [options]")
	fmt.Println("")
	fmt.Println("Available commands:")
	for _, c := range cmds {
		fmt.Printf("%-2s%-10s %10s\n", "", c.cmd, c.desc)
	}
	fmt.Println("")
	fmt.Println("Use 'task-cli <command> -h' for more information on a specific command.")
	os.Exit(0)
}

func AddHandler(subject string) {
	if len(subject) == 0 {
		log.Fatalln("cannot add task with empty subject.")
	}
	newId := Tasks.AddTask(subject)
	log.Printf("SUCCESS - added new task (ID: %d)", newId)
}

func UpdateHandler(uId int, subUpdate string, statusUpdate string) {
	var s task.Status
	if len(statusUpdate) > 0 {
		parseStatus(statusUpdate, &s)
	}

	if target, err := Tasks.UpdateTaskById(uId, subUpdate, s); err == nil {
		log.Printf("SUCCESS - updated task (ID: %d)", uId)
		log.Println(target)
	} else {
		log.Fatal(err)
	}
}

func DeleteHandler(dId int) {
	if target, err := Tasks.DeleteTaskById(dId); err == nil {
		log.Printf("SUCCESS - deleted task (ID: %d)", dId)
		log.Println(target)
	} else {
		log.Fatal(err)
	}
}

func ListHandler(status string) {
	var list []task.Task
	if len(status) == 0 {
		list = Tasks.FetchAllTasks()
	} else {
		var s task.Status
		parseStatus(status, &s)
		list = Tasks.FetchTaskByStatus(s)
	}

	for i := range list {
		fmt.Println(list[i])
	}
}

func parseStatus(s string, st *task.Status) {
	switch s {
	case "done":
		*st = task.DONE
	case "todo":
		*st = task.TODO
	case "in-progress":
		*st = task.IN_PROG
	default:
		log.Fatal("invalid status given")
	}
}
