package core

import (
	"math/rand"
	"strconv"
)

type Pack struct {
	Collection string
	Id         string
	Owner      WalletAddress
}

// TODO: Define these a bit more, based on data. For now, keep it simple and random
type PackType struct {
	Name            string
	PackSize        int
	GroundReserved  int
	MonsterReserved int
	LootReserved    int
}

func (pack Pack) GetOwner() WalletAddress {
	return pack.Owner
}

func (packType PackType) Validate() (bool, string) {
	if packType.PackSize < packType.GroundReserved+
		packType.MonsterReserved+
		packType.LootReserved {

		return false, "Invalid pack definition: Reserved cards exceed pack size"
	} else {
		return true, ""
	}
}

func GetPackType(packType string) PackType {
	if packType != "Default" {
		panic("Only the 'Default' pack type is supported for now")
	}

	var packTypes = map[string]PackType{
		"Default": {Name: "Default", PackSize: 10, GroundReserved: 5, MonsterReserved: 3, LootReserved: 2},
	}

	returnPack := packTypes[packType]

	ok, message := returnPack.Validate()

	if ok {
		return returnPack
	} else {
		panic(message)
	}
}

var PACKS_OPENED int = 0
var OPENED_CARDS []*Tile = make([]*Tile, 0, 100)

func NewPack(collection string, owner WalletAddress) Pack {
	PACKS_OPENED++
	return Pack{Collection: collection, Id: collection + "_" + strconv.Itoa(PACKS_OPENED), Owner: owner}
}

func (pack Pack) OpenPack() []*Tile {
	packType := GetPackType(pack.Collection)

	returnPack := make([]*Tile, 0, packType.PackSize)

	// TODO: why can't I append to a slice?
	currentPosition := 0

	for i := 0; i < packType.GroundReserved; i++ {
		returnPack = append(returnPack, (NewTile(Ground, pack.Owner)))
		currentPosition++
	}

	for i := 0; i < packType.MonsterReserved; i++ {
		returnPack = append(returnPack, (NewTile(Monster, pack.Owner)))
		currentPosition++
	}

	for i := 0; i < packType.LootReserved; i++ {
		returnPack = append(returnPack, (NewTile(Loot, pack.Owner)))
		currentPosition++
	}

	// For the remaining items, take a random number to assign the type, favoring the focus
	for i := currentPosition; i < packType.PackSize-1; i++ {
		typePicker := rand.Intn(3)

		var tileType TileType

		switch typePicker {
		case 1:
			tileType = Ground
		case 2:
			tileType = Monster
		case 3:
			tileType = Loot
		default:
			panic("Invalid randomization when generating a pack")
		}

		returnPack = append(returnPack, (NewTile(tileType, pack.Owner)))
	}

	// TEMP: Add new cards into global store
	for _, v := range returnPack {
		OPENED_CARDS = append(OPENED_CARDS, v)
	}

	return returnPack
}
