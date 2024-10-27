package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-vgo/robotgo"
)

const (
	apiURL       = "https://127.0.0.1:2999"
	pollInterval = 1 * time.Second
)

type Player struct {
	SummonerName string `json:"summonerName"`
	IsDead       bool   `json:"isDead"`
}

func main() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	var lastDeathState bool

	for {
		players, err := fetchPlayerList()
		if err != nil {
			fmt.Println("Error fetching player list:", err)
			time.Sleep(pollInterval)
			continue
		}

		activePlayer, found := findActivePlayer(players)
		if !found {
			fmt.Println("Active player not found")
			time.Sleep(pollInterval)
			continue
		}

		if activePlayer.IsDead && !lastDeathState {
			sendChatMessage("JG GAP")
		}

		lastDeathState = activePlayer.IsDead
		fmt.Println("lastDeathState", lastDeathState)
		time.Sleep(pollInterval)
	}
}

func fetchPlayerList() ([]Player, error) {
	res, err := http.Get(apiURL + "/liveclientdata/playerlist")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var players []Player
	err = json.Unmarshal(body, &players)
	if err != nil {
		return nil, err
	}

	return players, nil
}

func findActivePlayer(players []Player) (Player, bool) {
	if len(players) > 0 {
		return players[0], true
	}
	return Player{}, false
}

func sendChatMessage(message string) {
	robotgo.KeyTap("enter")
	time.Sleep(100 * time.Millisecond)

	robotgo.TypeStr(message)
	time.Sleep(100 * time.Millisecond)

	robotgo.KeyTap("enter")
}

func jgGapArt() string {
	return `
■■■■■   ■■■■      ■■■■   ■■■■■   ■■■■■
    ■  ■         ■       ■    ■  ■    ■ 
    ■  ■   ■■    ■  ■■   ■■■■■   ■■■■■
■   ■  ■    ■    ■   ■   ■    ■   ■
 ■■■    ■■■■      ■■■■   ■    ■   ■
`
}
