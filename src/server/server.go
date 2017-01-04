package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
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
	timeNow := time.Now().UnixNano()

	r.ParseForm()
	mParams := r.URL.Query()
	url := r.URL.Path[5:]

	if url == "reg" {
		name := r.Form.Get("Name")
		password := r.Form.Get("Password")
		password2 := r.Form.Get("Password2")
		email := r.Form.Get("Email")
		validMail := regexp.MustCompile(`(^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$)`)
		//(^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$)
		user := GetUserByName(name)
		if user != nil {
			log.Println("user exists")
			http.Redirect(w, r, "/public/register.html#invalid_register", 307)
			return
		}

		if len(name) < 3 {
			log.Println("name too short")
			http.Redirect(w, r, "/public/register.html#name_too_short", 307)
			return
		}
		if len(password) < 3 {
			log.Println("pass too short")
			http.Redirect(w, r, "/public/register.html#password_too_short", 307)
			return
		}
		if password != password2 {
			log.Println("pass not equal")
			http.Redirect(w, r, "/public/register.html#passwords_different", 307)
			return
		}

		if !validMail.MatchString(email) {
			log.Println("email not an email")
			http.Redirect(w, r, "/public/register.html#email_invalid", 307)
			return
		}

		CreateUser(name, password)
		http.Redirect(w, r, "/public/login.html#register_succesful", 307)
	}
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

		var RawURLEncoding = base64.URLEncoding.WithPadding(base64.NoPadding)
		nCookie := http.Cookie{Name: "username", Value: RawURLEncoding.EncodeToString([]byte(user.Name)), Expires: expiration, Path: "/"}
		sCookie := http.Cookie{Name: "session", Value: user.Session, Expires: expiration, Path: "/"}
		http.SetCookie(w, &nCookie)
		http.SetCookie(w, &sCookie)

		http.Redirect(w, r, "/", 307)
		log.Println("login successful")
		log.Println(r.UserAgent())

		return
	}
	if url == "loginunity" {
		name := r.Form.Get("Name")
		password := r.Form.Get("Password")

		user := GetUserByName(name)
		if user == nil {
			log.Println("user nil")
			w.WriteHeader(401)
			return
		}

		if user.Password != password {
			log.Println("pass wrong")
			w.WriteHeader(401)
			return
		}

		if len(user.Session) <= 0 {
			user.GenerateCookie()
		}

		var RawURLEncoding = base64.URLEncoding.WithPadding(base64.NoPadding)
		w.Header().Set("hUsername", RawURLEncoding.EncodeToString([]byte(user.Name)))
		w.Header().Set("hSession", user.Session)

		w.WriteHeader(200)
		log.Println("login successful")
		log.Println(r.UserAgent())

		return
	}
	//BEYOND THIS -> LOGIN REQUIRED
	user := getUserFromRequest(r)
	if user == nil {
		ClearHTTPCookie(w)
		w.WriteHeader(401)
		return
	}

	if url == "logout" {
		log.Println("logout")
		ClearHTTPCookie(w)
		user.ClearCookie()
		return
	}

	player := GetPlayerByUserId(user.Id)
	if player == nil {
		CreatePlayer(user.Id)
		player = GetPlayerByUserId(user.Id)
	}

	//timeNow := time.Now().UnixNano()
	//requestCount++

	var err error
	var data interface{}

	defer func() {
		timeDone := time.Now().UnixNano()
		//requestDuration += timeDone - timeNow
		log.Println("Dur:", (timeDone-timeNow)/(1000000))

		var result = map[string]interface{}{}
		result["ts"] = ((timeDone + timeNow) / 2) / 1000000
		result["tc"] = mParams.Get("tc")

		if err == nil {
			result["result"] = true
		} else {
			//requestErr++
			log.Println("rpc error:", url, err, data)
			result["result"] = false
			result["error"] = err.Error()
		}
		if data != nil {
			result["data"] = data
		}
		msg, _ := json.Marshal(result)
		w.Write(msg)
	}()

	if url == "mineYellow" {
		player.ResYellow += 10
		//player.SaveToDatabase()

		data = map[string]interface{}{"player": player}

		return
	}

	log.Println(url)
}

func ClearHTTPCookie(w http.ResponseWriter) {
	expiration := time.Now().Add(-24 * time.Hour)
	nCookie := http.Cookie{Name: "username", Value: "", Expires: expiration, Path: "/"}
	sCookie := http.Cookie{Name: "session", Value: "", Expires: expiration, Path: "/"}
	http.SetCookie(w, &nCookie)
	http.SetCookie(w, &sCookie)
}

func main() {

	InitDbUser()
	InitDbGame()

	http.HandleFunc("/", handler)
	http.HandleFunc("/rpc/", handlerRpc)
	http.ListenAndServe(":80", nil)
}

func getUserFromRequest(r *http.Request) *User {
	username, err := r.Cookie("username")
	if err != nil {
		return nil
	}

	var RawURLEncoding = base64.URLEncoding.WithPadding(base64.NoPadding)
	usernameByte, _ := RawURLEncoding.DecodeString(username.Value)
	username.Value = string(usernameByte)
	sessioncookie, err := r.Cookie("session")
	if err != nil {
		return nil
	}

	user := GetUserByName(username.Value)
	log.Println(username.Value, sessioncookie.Value, user)
	if user == nil {
		return nil
	}

	if user.Session != sessioncookie.Value || len(sessioncookie.Value) <= 0 {
		return nil
	}

	return user
}

/**/
