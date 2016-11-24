package main

import (
	"log"
	"net/http"
	"strings"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	url := r.URL.Path[1:]
	if strings.HasPrefix(url, "public/") {
		http.ServeFile(w, r, "www/"+url)
		return
	}

	user := getUserFromRequest(r)
	if user == nil {
		http.Redirect(w, r, "/public/login.html", 307)
		return
	}

	//fmt.Fprintf(w, "Hi there, I love %s %v!", r.URL.Path[1:], a)

	http.ServeFile(w, r, "www/"+url)
}

func handlerRpc(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	url := r.URL.Path[5:]

	if url == "login" {
		name := r.Form.Get("Name")
		password := r.Form.Get("Password")

		user := GetUserByName(name)
		if user == nil {
			log.Println("user nil")
			http.Redirect(w, r, "/public/login.html#invalid_login", 307)
			return
		}

		if user.Password != password {
			http.Redirect(w, r, "/public/login.html#invalid_login", 307)
			log.Println("pass wrong")
			return
		}

		if len(user.Session) <= 0 {
			user.GenerateCookie()
		}

		expiration := time.Now().Add(365 * 24 * time.Hour * 100)

		nCookie := http.Cookie{Name: "username", Value: user.Name, Expires: expiration, Path: "/"}
		sCookie := http.Cookie{Name: "session", Value: user.Session, Expires: expiration, Path: "/"}
		http.SetCookie(w, &nCookie)
		http.SetCookie(w, &sCookie)

		http.Redirect(w, r, "/", 307)
		return
	}

	user := getUserFromRequest(r)
	if user == nil {
		w.WriteHeader(401)
		return
	}

	if url == "logout" {
		user.ClearCookie()
		return
	}

	log.Println(url)
}

func main() {

	InitDB()

	http.HandleFunc("/", handler)
	http.HandleFunc("/rpc/", handlerRpc)
	http.ListenAndServe(":80", nil)
}

func getUserFromRequest(r *http.Request) *User {
	username, err := r.Cookie("username")
	if err != nil {
		return nil
	}

	sessioncookie, err := r.Cookie("session")
	if err != nil {
		return nil
	}

	user := GetUserByName(username.Value)
	if user == nil {
		return nil
	}

	if user.Session != sessioncookie.Value {
		return nil
	}

	return user
}

/**/
