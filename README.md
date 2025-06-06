
# task_tracker

CLI app to task your tasks.
Made as part of the Go backend roadmap on roadmap.sh
task: https://roadmap.sh/projects/task-tracker



## Installation

```bash
  gh repo clone chaosbooy/task_tracker
  go build
```
    
## Documentation

to use the task tracker after building the execution file use the todo flag and one the arguments listed below

```bash
Usage of todo:
  -add string
        Add a task with the given name
  -delete int
        Delete a task with the given id (default -1)
  -list string
        List all tasks or with specified status (todo, done, in-progress) (default "none")
  -mark-done int
        Mark a task with the given id as done (default -1)
  -mark-in-progress int
        Mark a task with the given id as in progress (default -1)
  -update int
        Update a task with the given id (default -1)
```
## Demo

```bash
    ./task_tracker todo --add do_chores
    ./task_tracker todo --list todo
    do_chores; Created at: 2025-06-06 21:21:34.6962502 +0200 CEST, Last Update: 2025-06-06 21:21:34.6962502 +0200 CES
```
