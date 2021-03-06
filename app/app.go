package app

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/endofcake/go-todo-rest-api-example/app/model"
	"github.com/endofcake/go-todo-rest-api-example/config"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/mingrammer/go-todo-rest-api-example/app/handler"
	"gopkg.in/matryer/try.v1"
)

// App has router and db instances
type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

// Initialize  the app with predefined configuration
func (a *App) Initialize(config *config.Config) {
	dbURI := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		config.DB.Host,
		config.DB.Port,
		config.DB.DbName,
		config.DB.Username,
		config.DB.Password,
		config.DB.SslMode)

	log.Print(dbURI)

	log.Print("Trying connection to the database...")

	var db *gorm.DB
	err := try.Do(func(attempt int) (bool, error) {
		var err error
		db, err = gorm.Open(config.DB.Dialect, dbURI)
		if err != nil {
			log.Print("Connection failed, retrying...")
			time.Sleep(5 * time.Second)
		}
		return attempt < 5, err // try 5 times
	})
	if err != nil {
		log.Fatalf("Could not connect to the database. %s", err)
	}

	// db, err := gorm.Open(config.DB.Dialect, dbURI)
	// if err != nil {
	// 	print(err)
	// }

	a.DB = model.DBMigrate(db)
	a.Router = mux.NewRouter()
	a.setRouters()
}

// setRouters sets the all required routers
func (a *App) setRouters() {
	// Routing for handling the projects
	a.Get("/projects", a.GetAllProjects)
	a.Post("/projects", a.CreateProject)
	a.Get("/projects/{title}", a.GetProject)
	a.Put("/projects/{title}", a.UpdateProject)
	a.Delete("/projects/{title}", a.DeleteProject)
	a.Put("/projects/{title}/archive", a.ArchiveProject)
	a.Delete("/projects/{title}/archive", a.RestoreProject)

	// Routing for handling the tasks
	a.Get("/projects/{title}/tasks", a.GetAllTasks)
	a.Post("/projects/{title}/tasks", a.CreateTask)
	a.Get("/projects/{title}/tasks/{id:[0-9]+}", a.GetTask)
	a.Put("/projects/{title}/tasks/{id:[0-9]+}", a.UpdateTask)
	a.Delete("/projects/{title}/tasks/{id:[0-9]+}", a.DeleteTask)
	a.Put("/projects/{title}/tasks/{id:[0-9]+}/complete", a.CompleteTask)
	a.Delete("/projects/{title}/tasks/{id:[0-9]+}/complete", a.UndoTask)
}

// Get wraps the router for GET method
func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

// Post wraps the router for POST method
func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

// Put wraps the router for PUT method
func (a *App) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("PUT")
}

// Delete wraps the router for DELETE method
func (a *App) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("DELETE")
}

/*
** Projects Handlers
 */
func (a *App) GetAllProjects(w http.ResponseWriter, r *http.Request) {
	handler.GetAllProjects(a.DB, w, r)
}

func (a *App) CreateProject(w http.ResponseWriter, r *http.Request) {
	handler.CreateProject(a.DB, w, r)
}

func (a *App) GetProject(w http.ResponseWriter, r *http.Request) {
	handler.GetProject(a.DB, w, r)
}

func (a *App) UpdateProject(w http.ResponseWriter, r *http.Request) {
	handler.UpdateProject(a.DB, w, r)
}

func (a *App) DeleteProject(w http.ResponseWriter, r *http.Request) {
	handler.DeleteProject(a.DB, w, r)
}

func (a *App) ArchiveProject(w http.ResponseWriter, r *http.Request) {
	handler.ArchiveProject(a.DB, w, r)
}

func (a *App) RestoreProject(w http.ResponseWriter, r *http.Request) {
	handler.RestoreProject(a.DB, w, r)
}

/*
** Tasks Handlers
 */
func (a *App) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	handler.GetAllTasks(a.DB, w, r)
}

func (a *App) CreateTask(w http.ResponseWriter, r *http.Request) {
	handler.CreateTask(a.DB, w, r)
}

func (a *App) GetTask(w http.ResponseWriter, r *http.Request) {
	handler.GetTask(a.DB, w, r)
}

func (a *App) UpdateTask(w http.ResponseWriter, r *http.Request) {
	handler.UpdateTask(a.DB, w, r)
}

func (a *App) DeleteTask(w http.ResponseWriter, r *http.Request) {
	handler.DeleteTask(a.DB, w, r)
}

func (a *App) CompleteTask(w http.ResponseWriter, r *http.Request) {
	handler.CompleteTask(a.DB, w, r)
}

func (a *App) UndoTask(w http.ResponseWriter, r *http.Request) {
	handler.UndoTask(a.DB, w, r)
}

// Run the app on it's router
func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}
