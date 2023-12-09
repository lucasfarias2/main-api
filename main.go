package main

import (
	"log"
	"main-api/db"
	"net/http"
	"strings"
)

import (
	"github.com/fsnotify/fsnotify"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"html/template"
	"main-api/routes"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	r := chi.NewRouter()

	db.ConnectDatabase()

	fileServer(r, "/static", http.Dir("static"))

	r.HandleFunc("/ws", handleConnections)

	r.Mount("/api", routes.ApiRouter())
	//r.Mount("/", routes.WebRouter())

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, r, "index", map[string]interface{}{
			"Title":   "Packlify",
			"Pattern": r.URL.Path,
		})
	})

	r.Get("/projects", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, r, "projects", map[string]interface{}{
			"Title":   "Packlify - Projects",
			"Pattern": r.URL.Path,
		})
	})

	log.Println("Running on localhost:8080")
	http.ListenAndServe(":8080", r)
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	})
}

func renderTemplate(w http.ResponseWriter, r *http.Request, name string, data map[string]interface{}) {
	tmpl, err := template.ParseFiles("templates/layout.html", "templates/"+name+".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.ExecuteTemplate(w, "layout.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	// Watch for changes in the templates folder
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
					ws.WriteMessage(websocket.TextMessage, []byte("reload"))
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("./templates")
	if err != nil {
		log.Fatal(err)
	}

	// Keep the connection open
	for {
		if _, _, err := ws.NextReader(); err != nil {
			ws.Close()
			break
		}
	}
}
