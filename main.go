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
	db, err = sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		log.Fatal(err)
	}

	q, err := db.Prepare("CREATE TABLE IF NOT EXISTS tasks (id INTEGER PRIMARY KEY AUTOINCREMENT, name VARCHAR(255), done BOOLEAN)")
	if err != nil {
		log.Fatal(err)
	}
	q.Exec()

	q, err = db.Prepare(`INSERT INTO tasks (id, name, done) VALUES (1, "Do a homework", false), (2, "Wash the dishes", true)`)
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
		printAllTasks()

		fmt.Printf("\nWhat do you want to do?\n1 - Add new task\n2 - Mark task as done\nInput:")

		var input string
		fmt.Scan(&input)

		if input == "1" {
			addNewTask()
		} else if input == "2" {
			markTaskAsDone()
		} else {
			fmt.Println("Wrong input!")
		}
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

	_, err = q.Exec(taskName, false)
	if err != nil {
		log.Fatal(err)
	}
}

func printAllTasks() {
	data, err := db.Query("SELECT COUNT(*) FROM tasks")
	if err != nil {
		log.Fatal(err)
	}

	var count int
	for data.Next() {
		data.Scan(&count)
	}

	rows, err := db.Query("SELECT * FROM tasks")
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

	for _, task := range tasks {
		var done string

		if task.Done {
			done = "[+]"
		} else {
			done = "[]"
		}

		fmt.Printf("#%d %s %s\n", task.Number, done, task.Name)
	}

}
