package core

import "math/rand"

type Tile struct {
	TileType
	ColorAttributes
	Owner WalletAddress
}

type TileType int

const (
	Ground TileType = iota
	Monster
	Loot
)

func (tile Tile) GetOwner() WalletAddress {
	return tile.Owner
}

// type TileDefinition struct {
// 	Collection string
// 	Id         int
// 	TileType
// }

func NewTile(tileType TileType, owner WalletAddress) *Tile {
	t := new(Tile)
	t.TileType = tileType
	t.Magenta = rand.Intn(25) + 10
	t.Yellow = rand.Intn(25) + 10
	t.Cyan = rand.Intn(25) + 10
	t.Owner = owner

	return t
}

func (tile Tile) GetPower() (output int) {
	tileColors := ConvertColorsToArray(&tile.ColorAttributes)

	for _, v := range tileColors {
		output += *v
	}

	return
}
