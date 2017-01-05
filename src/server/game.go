//game
package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const GAMEDATAVERSION = 1

var dbGame *sql.DB

type Player struct {
	Id        int
	UserId    int
	ResRed    int
	ResBlue   int
	ResYellow int
}

type GameData struct {
	Version int
	Players []*Player
}

var gameData GameData
var User2PlayerId map[int]int
var PlayerMutex sync.RWMutex

func InitDbGame() {
	LoadDbGame()

	go func() {
		for {
			time.Sleep(time.Second * 5 * 60)
			SaveDbGame()
		}
	}()
}

func SaveDbGame() {
	gameData.Version = GAMEDATAVERSION
	PlayerMutex.RLock()
	b, err := json.MarshalIndent(&gameData, "", "\t")
	PlayerMutex.RUnlock()
	if err != nil {
		log.Println(err)
		return
	}

	ioutil.WriteFile("data/game.json_temp", b, 0644)
	os.Remove("data/game.json_old")
	os.Rename("data/game.json", "data/game.json_old")
	os.Rename("data/game.json_temp", "data/game.json")
	os.Remove("data/game.json_old")
	log.Println("Game data saved successfully")
}

func LoadDbGame() {

	User2PlayerId = make(map[int]int)

	b, err := ioutil.ReadFile("data/game.json")
	if err != nil {
		os.Rename("data/game.json_old", "data/game.json")
		log.Println(err)
		b, err = ioutil.ReadFile("data/game.json")
		if err != nil {
			log.Println("load failed")
			return
		}
	}

	PlayerMutex.Lock()
	err = json.Unmarshal(b, &gameData)

	for _, v := range gameData.Players {
		User2PlayerId[v.UserId] = v.Id
	}
	PlayerMutex.Unlock()

	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Game data loaded successfully")
}

func CreatePlayer(UserId int) {
	PlayerMutex.Lock()
	defer PlayerMutex.Unlock()

	_, ok := User2PlayerId[UserId]
	if ok {
		log.Println("already existing player")
		return
	}
	p := &Player{
		UserId: UserId,
		Id:     len(gameData.Players),
	}
	gameData.Players = append(gameData.Players, p)
	User2PlayerId[UserId] = p.Id
	log.Println("player created", UserId, p.Id, len(gameData.Players))
	return
}

func GetPlayerByUserId(UserId int) *Player {

	PlayerMutex.RLock()
	defer PlayerMutex.RUnlock()
	i, ok := User2PlayerId[UserId]
	if !ok {
		return nil
	}
	log.Println("player get", UserId, i, ok, len(gameData.Players))
	return gameData.Players[i]
}
