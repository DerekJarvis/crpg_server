package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"jarvispowered.com/color_rpg/logic/src/core"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// Make new players
	player1, _ := core.NewPlayer("harry", "0x0000000000000000000000000000000000000001", "")
	player2, _ := core.NewPlayer("sally", "0x0000000000000000000000000000000000000002", "")

	// Make a new character and assign it to player 1
	newCharacter := core.NewCharacter("Bob", player1.EthAddress)

	// Open 10 packs and assign to player2
	for i := 0; i < 10; i++ {
		newPack := core.NewPack("Default", player2.EthAddress)
		newPack.OpenPack()
	}

	// Create a Dungeon out of the cards (put monsters and treasure at the right)
	newDungeon := core.MakeDumbDungeon(player2)

	// player 1 plays dungeon
	match, _ := newDungeon.PlayDumbDungeon(newCharacter)

	player1.Gold += match.Reward
	player2.Gold -= match.Reward

	// Report out player 1, player 2, dungeon, and dungeon match
	var jsonBytes []byte
	var jsonError error

	jsonBytes, jsonError = json.Marshal(player1)
	if jsonError != nil {
		fmt.Println(jsonError)
	} else {
		fmt.Println(string(jsonBytes))
	}

	jsonBytes, jsonError = json.Marshal(player2)
	if jsonError != nil {
		fmt.Println(jsonError)
	} else {
		fmt.Println(string(jsonBytes))
	}

	jsonBytes, jsonError = json.Marshal(newDungeon)
	if jsonError != nil {
		fmt.Println(jsonError)
	} else {
		fmt.Println(string(jsonBytes))
	}

	jsonBytes, jsonError = json.Marshal(match)
	if jsonError != nil {
		fmt.Println(jsonError)
	} else {
		fmt.Println(string(jsonBytes))
	}
}
