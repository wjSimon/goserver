package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var database *sql.DB

type User struct {
	Id       int
	Name     string
	Password string
	Session  string
}

func InitDB() {
	database, _ = sql.Open("sqlite3", "data/user.sqlitedb")
	create := `CREATE TABLE 'userinfo' (
    'Id' 		INTEGER PRIMARY KEY AUTOINCREMENT,
    'Name' 		TEXT    NOT NULL,
    'Password' 	TEXT    NOT NULL,
    'Session' 	TEXT    NOT NULL
	);`

	res, err := database.Exec(create)
	log.Println(res, err)

	res, err = database.Exec("CREATE UNIQUE INDEX index_Name on userinfo (Name);")
	log.Println(res, err)

	//GetUserByName("dummy2")
}

func GetUserByName(name string) *User {

	stmt, _ := database.Prepare("SELECT * FROM userinfo WHERE Name=?")
	defer stmt.Close()

	user := &User{}
	err := stmt.QueryRow(name).Scan(&user.Id, &user.Name, &user.Password, &user.Session)
	if err != nil {
		log.Println(err)
		return nil
	}
	fmt.Println(user)
	return user
}

func (u *User) ClearCookie() { //Logout
	u.Session = ""
	u.SaveToDatabase()
}

func (u *User) GenerateCookie() {
	u.Session = "session"
	u.SaveToDatabase()
}

func (u *User) SaveToDatabase() {
	stmt, err := database.Prepare("UPDATE userinfo SET Password=?, Session=? WHERE Id=?")
	log.Println(err)
	defer stmt.Close()
	stmt.Exec(&u.Password, &u.Session, &u.Id)
}
