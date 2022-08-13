package main

import (
	"net/http"
	"text/template"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	ID       uint `gorm:"autoincrement"`
	TaskName string
	Assignee string
	Deadline string
	Status   string `gorm:"default:unfinished"`
}

type Data struct {
	Title     string
	Body      string
	All_Tasks []Task
}

var templates = template.Must(template.ParseGlob("views/*"))

func renderTemplate(w http.ResponseWriter, tmpl string, page *Data) {
	err := templates.ExecuteTemplate(w, tmpl, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {

	var all_tasks []Task

	db, err := gorm.Open(sqlite.Open("tasks.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.Find(&all_tasks)

	data := &Data{
		Title:     "All Tasks",
		Body:      "",
		All_Tasks: all_tasks,
	}

	templates.Execute(w, data)
	renderTemplate(w, "index", data)
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {

	data := &Data{}

	templates.Execute(w, data)
	renderTemplate(w, "create", data)
}

func SaveHandler(w http.ResponseWriter, r *http.Request) {

	var newTask Task

	db, err := gorm.Open(sqlite.Open("tasks.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if r.Method == "POST" {

		newTask.TaskName = r.FormValue("taskname")
		newTask.Assignee = r.FormValue("assignee")
		newTask.Deadline = r.FormValue("deadline")

		db.Create(&newTask)

		if err != nil {
			panic(err.Error())
		}

	}

	http.Redirect(w, r, "/", 301)
}

func EditHandler(w http.ResponseWriter, r *http.Request) {

	var editData Task
	var arrayData []Task

	id := r.URL.Query().Get("id")

	db, err := gorm.Open(sqlite.Open("tasks.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.First(&editData, id)

	arrayData = append(arrayData, editData)

	data := &Data{
		Title:     "Edit Task",
		Body:      "",
		All_Tasks: arrayData,
	}

	templates.Execute(w, data)
	renderTemplate(w, "edit", data)
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	var editData Task

	id := r.URL.Query().Get("id")

	db, err := gorm.Open(sqlite.Open("tasks.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if r.Method == "POST" {
		db.First(&editData, id)
		db.Model(&editData).Updates(Task{TaskName: r.FormValue("taskname"), Assignee: r.FormValue("assignee"), Deadline: r.FormValue("deadline")})

		if err != nil {
			panic(err.Error())
		}

	}

	http.Redirect(w, r, "/", 301)
}

func MarkAsDoneHandler(w http.ResponseWriter, r *http.Request) {
	var editData Task

	id := r.URL.Query().Get("id")

	db, err := gorm.Open(sqlite.Open("tasks.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.First(&editData, id)
	db.Model(&editData).Updates(Task{Status: "finished"})

	if err != nil {
		panic(err.Error())
	}

	http.Redirect(w, r, "/", 301)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	var editData Task

	id := r.URL.Query().Get("id")

	db, err := gorm.Open(sqlite.Open("tasks.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.First(&editData, id)
	db.Delete(&editData)

	if err != nil {
		panic(err.Error())
	}

	http.Redirect(w, r, "/", 301)
}

func main() {
	db, err := gorm.Open(sqlite.Open("tasks.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Task{})

	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/create", CreateHandler)
	http.HandleFunc("/save", SaveHandler)
	http.HandleFunc("/edit", EditHandler)
	http.HandleFunc("/update", UpdateHandler)
	http.HandleFunc("/markasdone", MarkAsDoneHandler)
	http.HandleFunc("/delete", DeleteHandler)

	http.ListenAndServe(":8080", nil)

}
