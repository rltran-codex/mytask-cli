# Task Tracker

- [Intro](#overview)
- [Installation](#installation)
- [Usage](#usage)
  - [Add](#add)
  - [Update](#update)
  - [Delete](#delete)
  - [List](#list)

# Overview

Project from [roadmap.sh](https://roadmap.sh/projects/task-tracker) built using Go.

Create a simple TODO application that allows users to manage tasks efficiently.
Each task includes an ID, subject, creation date, and last updated timestamp.
The application stores tasks in a .json file, ensuring persistence across sessions.

# Installation

1. Clone this repo `git clone git@github.com:rltran-codex/mytask-cli.git`
2. `cd mytask-cli`
3. Build the application

### Windows 10/11

```powershell
go build -o task-cli.exe
```

### Linux

```bash
go build -o task-cli
```

# Usage

On the first usage or absence of appdata folder, the application will
initialize a new `tasks.json` file.
Run the application from the command line.

```bash
> task-cli

Usage:
  task-cli <command> [options]

Available commands:
  add        Add a new task
  update     Update an existing task
  delete     Delete a task
  list       List tasks based on criteria

Use 'task-cli <command> -h' for more information on a specific command.
```

## Add

```bash
> task-cli add -h

Usage of add:
  -subj string
        REQUIRED. Define the new task's subject.
```

Example

```bash
> task-cli add -subj "Buy groceries"
```

## Update

```bash
> task-cli update -h

Usage of update:
  -id int
        REQUIRED. Update task with id. (default -1)
  -stat string
        OPTIONAL. Update task's status.
  -subj string
        OPTIONAL. Update task's subject.
```

Example

```bash
> task-cli update -id 1 -stat "in-progress"
> task-cli update -id 1 -stat "done"
> task-cli update -id 1 -subj "Buy groceries and cook dinner"
```

## Delete

```bash
> task-cli delete -h

Usage of delete:
  -id int
        REQUIRED. Delete task with id. (default -1)
```

Example

```bash
> task-cli delete -id 1
```

## List

```bash
> task-cli list -h

Usage of list:
  -stat string
        OPTIONAL. List all tasks with specified status.
```

Example

```bash
> task-cli list
> task-cli list -stat "done"
> task-cli list -stat "todo"
> task-cli list -stat "in-progress"
```
