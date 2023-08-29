package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Task struct {
	Name   string `json:"name"`
	Done   bool   `json:"done"`
	Number int    `json:"id"`
}

func init() {
	var err error
	db, err = sql.Open("sqlite3", "/home/aaddssww/tasks.db")
	if err != nil {
		log.Fatal(err)
	}

	q, err := db.Prepare("CREATE TABLE IF NOT EXISTS tasks (id INTEGER PRIMARY KEY AUTOINCREMENT, name VARCHAR(255), done BOOLEAN)")
	if err != nil {
		log.Fatal(err)
	}
	q.Exec()
}

var (
	db    *sql.DB
	tasks []Task
)

func main() {
	for {
		fmt.Printf(`
What do you want to do?
1 - Add new task
2 - Mark task as done
3 - Print all tasks
4 - Print going tasks
5 - Delete all tasks
Input:`)

		var input string
		fmt.Scan(&input)

		if input == "1" {
			addNewTask()
		} else if input == "2" {
			markTaskAsDone()
		} else if input == "3" {
			printAllTasks(true)
		} else if input == "4" {
			printAllTasks(false)
		} else if input == "5" {
			deleteAll()
		} else {
			fmt.Println("Wrong input!")
		}
	}
}

func deleteAll() {
	_, err := db.Exec("DELETE FROM tasks")
	if err != nil {
		log.Fatal(err)
	}
}

func markTaskAsDone() {
	fmt.Printf("\nEnter task number: ")
	var number string
	fmt.Scanln(&number)

	q, err := db.Prepare("UPDATE tasks SET done = true WHERE id = $1")
	if err != nil {
		log.Fatal(err)
	}

	q.Exec(number)
}

func addNewTask() {
	fmt.Printf("\nEnter a name: ")

	reader := bufio.NewReader(os.Stdin)
	taskName, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	q, err := db.Prepare(`INSERT INTO tasks (name, done) VALUES ($1, $2)`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = q.Exec(taskName[:len(taskName)-1], false)
	if err != nil {
		log.Fatal(err)
	}
}

func printAllTasks(all bool) {
	var query string
	if all {
		query = "SELECT * FROM tasks"
	} else {
		query = "SELECT * FROM tasks where done = false"
	}
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}

	tasks = nil

	var id int
	var name string
	var done bool
	for rows.Next() {
		rows.Scan(&id, &name, &done)

		task := Task{
			Name:   name,
			Done:   done,
			Number: id,
		}

		tasks = append(tasks, task)
	}

	fmt.Println("\n\n----------------------")

	for _, task := range tasks {
		var done string

		if task.Done {
			done = "[+]"
		} else {
			done = "[]"
		}

		fmt.Printf("#%d %s %s\n", task.Number, done, task.Name)
	}

	fmt.Println("----------------------")
}
