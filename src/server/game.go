//game
package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var dbGame *sql.DB

type Player struct {
	Id        int
	UserId    int
	ResRed    int
	ResBlue   int
	ResYellow int
}

func InitDbGame() {
	dbGame, _ = sql.Open("sqlite3", "data/game.sqlitedb?_bulsy_timeout=5000")
	//dbGame.Exec("PRAGMA busy_timeout = 50000;")
	dbGame.SetMaxOpenConns(1)

	create := `CREATE TABLE 'player' (
    'Id' 		INTEGER PRIMARY KEY AUTOINCREMENT,
    'UserId' 	INTEGER NOT NULL,
    'ResRed' 	INTEGER NOT NULL DEFAULT 0,
    'ResBlue' 	INTEGER NOT NULL DEFAULT 0, 
	'ResYellow' INTEGER NOT NULL DEFAULT 0
	);`

	res, err := dbGame.Exec(create)
	log.Println(res, err)

	res, err = dbGame.Exec("CREATE UNIQUE INDEX index_UserId on player (UserId);")
	log.Println(res, err)
}

func CreatePlayer(UserId int) {
	stmt, _ := dbGame.Prepare("INSERT INTO player (UserId) VALUES (?)")
	defer stmt.Close()

	res, err := stmt.Exec(&UserId)
	log.Println(res, err)
}

func GetPlayerByUserId(UserId int) *Player {

	stmt, _ := dbGame.Prepare("SELECT * FROM player WHERE UserId=?")
	defer stmt.Close()

	player := &Player{}
	err := stmt.QueryRow(UserId).Scan(&player.Id, &player.UserId, &player.ResRed, &player.ResBlue, &player.ResYellow)
	if err != nil {
		log.Println(err)
		return nil
	}
	fmt.Println(player)
	return player
}

func (p *Player) SaveToDatabase() {
	stmt, err := dbGame.Prepare("UPDATE player SET ResRed=?, ResBlue=?, ResYellow=? WHERE Id=?")
	defer stmt.Close()
	res, err := stmt.Exec(&p.ResRed, &p.ResBlue, &p.ResYellow, &p.Id)
	log.Println("SaveToDatabase()", res, err)
}
