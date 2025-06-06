package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type status int

func (s status) String() string {
	switch s {
	case Todo:
		return "todo"
	case InProgress:
		return "in-progress"
	case Done:
		return "done"
	default:
		return "unknown"
	}
}

func (s status) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *status) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	switch str {
	case "todo":
		*s = Todo
	case "in-progress":
		*s = InProgress
	case "done":
		*s = Done
	default:
		*s = Todo
	}

	return nil
}

const (
	Todo       status = iota
	InProgress status = iota
	Done       status = iota
)

// Task class for reading
type task struct {
	Description string
	Status      status
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// task holder
var TodoList []task = make([]task, 0)

// main programm function
func main() {
	const NAME = "to-do_list.txt"

	var err error
	TodoList, err = loadData(NAME)

	if err != nil {
		panic(err)
	}

	todoCmd := flag.NewFlagSet("todo", flag.ExitOnError)

	todoAdd := todoCmd.String("add", "", "Add a task with the given name")
	todoUpdate := todoCmd.Int("update", -1, "Update a task with the given id")
	todoDelete := todoCmd.Int("delete", -1, "Delete a task with the given id")
	todoMarkDone := todoCmd.Int("mark-done", -1, "Mark a task with the given id as done")
	todoMarkProgress := todoCmd.Int("mark-in-progress", -1, "Mark a task with the given id as in progress")
	todoList := todoCmd.String("list", "none", "List all tasks or with specified status (todo, done, in-progress)")

	if len(os.Args) < 2 {
		fmt.Println("Expected a 'todo' command")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "todo":
		todoCmd.Parse(os.Args[2:])
		HandleTodo(todoAdd, todoUpdate, todoDelete, todoMarkDone, todoMarkProgress, todoList, todoCmd.Args())
	default:
		fmt.Println("Expected a 'todo' command")
		os.Exit(1)
	}

	saveData(NAME)
}

func loadData(location string) ([]task, error) {

	fi, err := os.Open(location)
	if err != nil {
		file, err := os.Create(location)

		if err != nil {
			return []task{}, err
		}

		defer file.Close()
		return []task{}, nil
	}

	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()

	// read our opened jsonFile as a byte array.
	byteValue, _ := io.ReadAll(fi)

	// we initialize our Users array
	var tasks []task
	json.Unmarshal(byteValue, &tasks)

	return tasks, nil
}

func saveData(location string) {
	// open output file
	fo, err := os.Create(location)
	if err != nil {
		panic(err)
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	textJson, err := json.Marshal(TodoList)
	if err != nil {
		panic(err)
	}

	fo.Write(textJson)
}

func HandleTodo(add *string, update *int, delete *int, markDone *int, markProgress *int, list *string, extraArgs []string) {

	// if *add == "" && *update == -1 && *delete == -1 && *markDone == -1 && *markProgress == -1 && *list == "none" {
	// 	fmt.Print("Expected an argument after 'todo' command")
	// 	os.Exit(1)
	// }

	if *add != "" {
		var newTask = task{
			Description: *add,
			Status:      Todo,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		TodoList = append(TodoList, newTask)
	}

	if *update != -1 {

		if len(extraArgs) == 0 {
			fmt.Println("Missing description for task update.")
			return
		} else if *update < 0 || *update >= len(TodoList) {
			fmt.Println("Invalid task ID.")
			return
		}
		description := strings.Join(extraArgs, " ")

		TodoList[*update].Description = description
		TodoList[*update].UpdatedAt = time.Now()
	}

	if *delete != -1 {
		if *delete < 0 || *delete >= len(TodoList) {
			fmt.Println("Invalid task ID.")
			return
		}
		TodoList = append(TodoList[:*delete], TodoList[*delete+1:]...)
	}

	if *markDone != -1 {
		if *markDone < 0 || *markDone >= len(TodoList) {
			fmt.Println("Invalid task ID.")
			return
		}
		TodoList[*markDone].Status = Done
	}

	if *markProgress != -1 {
		if *markProgress < 0 || *markProgress >= len(TodoList) {
			fmt.Println("Invalid task ID.")
			return
		}
		TodoList[*markProgress].Status = InProgress
	}

	switch *list {
	case "done":
		for _, v := range TodoList {
			if v.Status == Done {
				fmt.Printf("%s; Created at: %s, Last Update: %s \n", v.Description, v.CreatedAt, v.UpdatedAt)
			}
		}
	case "in-progress":
		for _, v := range TodoList {
			if v.Status == InProgress {
				fmt.Printf("%s; Created at: %s, Last Update: %s \n", v.Description, v.CreatedAt, v.UpdatedAt)
			}
		}
	case "todo":
		for _, v := range TodoList {
			if v.Status == Todo {
				fmt.Printf("%s; Created at: %s, Last Update: %s \n", v.Description, v.CreatedAt, v.UpdatedAt)
			}
		}
	case "":
		for _, v := range TodoList {
			var itemStatus string
			switch v.Status {
			case Done:
				itemStatus = "done"
			case Todo:
				itemStatus = "to do"
			case InProgress:
				itemStatus = "in progress"

			}
			fmt.Printf("%s; Status: %s Created at: %s; Last Update: %s \n", v.Description, itemStatus, v.CreatedAt, v.UpdatedAt)
		}
	case "none":
		break
	default:
		fmt.Println("Wrong argument after list command")
		return
	}
}
