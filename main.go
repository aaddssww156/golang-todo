package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Task struct {
	Name   string `json:"name"`
	Done   bool   `json:"done"`
	Number string `json:"number"`
}

func init() {
	tasksFile, err := os.Open("tasks.txt")
	if err != nil {
		log.Fatal(err)
	}

	tasks = loadTasks(tasksFile)

	defer tasksFile.Close()
}

var (
	tasks []Task
)

func main() {
	for {
		printTasks()

		var input string
		fmt.Printf("\nWhat do you want to do?\n\n1 - Add a new task\n2 - Mark task as done\n")
		fmt.Scanln(&input)

		if input == "1" {
			addTask()
		} else if input == "2" {
			markAsDone()
		} else {
			log.Println("Wrong input!")
		}
	}
}

func markAsDone() {
	panic("unimplemented")
}

func addTask() {
	fmt.Printf("\n\nEnter a task name - ")

	reader := bufio.NewReader(os.Stdin)
	taskName, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	newTask := Task{
		Name:   taskName,
		Done:   false,
		Number: strconv.Itoa(len(tasks) + 1),
	}

	apendTaskToAFile(newTask)

	tasks = append(tasks, newTask)
}

func apendTaskToAFile(task Task) {
	file, err := os.OpenFile("tasks.txt", os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	var done string

	if task.Done {
		done = "[+]"
	} else {
		done = "[]"
	}

	newString := fmt.Sprintf("\n%s. %s %s", task.Number, done, task.Name)

	_, err = file.WriteAt([]byte(task.Number), 0)
	if err != nil {
		log.Fatal(err)
	}

	file.Close()

	file, err = os.OpenFile("tasks.txt", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatal(err)
	}
	_, err = file.WriteString(newString)
	if err != nil {
		log.Fatal(err)
	}
}

func loadTasks(file *os.File) []Task {
	scanner := bufio.NewScanner(file)
	// tasks := make([]Task, )
	var tasks []Task

	i := 0
	for scanner.Scan() {
		if i == 0 {
			number, err := strconv.Atoi(scanner.Text())
			if err != nil {
				log.Fatal(err)
			}
			tasks = make([]Task, number)
		} else {
			line := scanner.Text()
			tasks[i-1].Number = string(line[0])
			tasks[i-1].Name = strings.Split(line, "]")[1]
			tasks[i-1].Done = IsDone(&line)
		}
		i++
	}
	return tasks
}

func printTasks() {
	for _, task := range tasks {
		var done string

		if task.Done {
			done = "[+]"
		} else {
			done = "[]"
		}

		fmt.Printf("#%s %s %s\n", task.Number, done, task.Name)
	}
}

func IsDone(line *string) bool {
	done := strings.Index(*line, "+")

	if done > -1 {
		return true
	} else {
		return false
	}
}
