package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "这是home处理器")
}
func handleDetail(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "这是detail处理器")
}
func handleIndex(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title string
		Items []string
	}{
		Title: "梦回故里",
		Items: []string{
			"fly",
			"run",
		},
	}

	templates.ExecuteTemplate(w, "index.html", data)
}

func handleDb(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test_new")
	if err != nil {
		panic(err)
	}

	rows, err := db.Query("select id,title from news")

	if err != nil {
		panic(err)
	}

	var (
		id    string
		title string
	)
	var news []string

	for rows.Next() {
		rows.Scan(&id, &title)
		news = append(news, title)
	}

	templates.ExecuteTemplate(w, "db.html", news)
}

func handleForm(w http.ResponseWriter, r *http.Request) {
	//数据库获取数据
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test_new")
	if err != nil {
		panic(err)
	}

	rows, err := db.Query("select id,title from news")

	if err != nil {
		panic(err)
	}

	var (
		id    string
		title string
	)
	var news []string

	for rows.Next() {
		rows.Scan(&id, &title)
		news = append(news, title)
	}

	templates.ExecuteTemplate(w, "form.html", news)
}

func handleFormAdd(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	title := r.PostForm.Get("title")
	//TODO:将数据保存到数据库，然后再显示到页面上
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test_new")
	if err != nil {
		panic(err)
	}

	db.Exec("insert into news set title = ?", title)

}

var templates *template.Template

var db *sql.DB

var store = sessions.NewCookieStore([]byte("test"))

func main() {
	templates = template.Must(template.ParseGlob("templates/*.html"))
	r := mux.NewRouter()
	r.HandleFunc("/", handleIndex)
	r.HandleFunc("/home", handleHome)
	r.HandleFunc("/detail", handleDetail)
	r.HandleFunc("/db", handleDb)
	r.HandleFunc("/form", handleForm)
	r.HandleFunc("/formadd", handleFormAdd)
	r.HandleFunc("/setsession", handleSetSession)
	r.HandleFunc("/getsession", handleGetSession)
	fs := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static", fs))
	http.Handle("/", r)
	http.ListenAndServe(":8888", nil)
}

func handleSetSession(w http.ResponseWriter, r *http.Request) {
	//templates.ExecuteTemplate(w,"form.html",nil)
	session, err := store.Get(r, "menghuiguli")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["name"] = "xiaoming"
	session.Values["age"] = 19
	session.Save(r, w)
}

func handleGetSession(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "menghuiguli")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(session.Values["name"])
	fmt.Println(session.Values["height"])
	//templates.ExecuteTemplate(w,"getsession.html",nil)
}
