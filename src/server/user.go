package main

import (
	"database/sql"
	"fmt"
	"log"
	"util"

	_ "github.com/mattn/go-sqlite3"
)

var dbUser *sql.DB

type User struct {
	Id        int
	Name      string
	Password  string
	Session   string
	UserLevel int
}

func InitDbUser() {
	dbUser, _ = sql.Open("sqlite3", "data/user.sqlitedb")
	create := `CREATE TABLE 'userinfo' (
    'Id' 		INTEGER PRIMARY KEY AUTOINCREMENT,
    'Name' 		TEXT    NOT NULL,
    'Password' 	TEXT    NOT NULL,
    'Session' 	TEXT    NOT NULL DEFAULT '',
	'UserLevel' INTEGER NOT NULL DEFAULT 0
	);`

	res, err := dbUser.Exec(create)
	log.Println(res, err)

	res, err = dbUser.Exec("CREATE UNIQUE INDEX index_Name on userinfo (Name);")
	log.Println(res, err)

	//GetUserByName("dummy2")
}

func GetUserByName(name string) *User {

	stmt, _ := dbUser.Prepare("SELECT * FROM userinfo WHERE Name=?")
	defer stmt.Close()

	user := &User{}
	err := stmt.QueryRow(name).Scan(&user.Id, &user.Name, &user.Password, &user.Session, &user.UserLevel)
	if err != nil {
		log.Println(err)
		return nil
	}
	fmt.Println(user)
	return user
}

func CreateUser(name string, password string) {
	stmt, _ := dbUser.Prepare("INSERT INTO userinfo (Name, Password) VALUES (?,?)")
	defer stmt.Close()

	res, err := stmt.Exec(&name, &password)
	log.Println(res, err)
}

func (u *User) ClearCookie() { //Logout
	u.Session = ""
	u.SaveToDatabase()
}

func (u *User) GenerateCookie() {
	u.Session, _ = util.GenerateRandomString(30)
	u.SaveToDatabase()
}

func (u *User) SaveToDatabase() {
	stmt, err := dbUser.Prepare("UPDATE userinfo SET Password=?, Session=? WHERE Id=?")
	defer stmt.Close()
	res, err := stmt.Exec(&u.Password, &u.Session, &u.Id)
	log.Println("SaveToDatabase()", res, err)
}
