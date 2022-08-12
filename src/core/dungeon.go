package core

import (
	"sort"
)

type DungeonTile struct {
	X    int
	Y    int
	Tile *Tile
}

type Dungeon struct {
	Name       string
	Owner      WalletAddress
	Width      int
	Height     int
	Power      int
	Difficulty int
	Tiles      []DungeonTile
}

type DungeonMatch struct {
	*Dungeon
	*Character
	DungeonMultiplier   float32
	CharacterMultiplier float32
	StartingPower       float32
	EndingPower         float32
	Reward              int
	PlayedTiles         []int
}

func (dungeon Dungeon) GetOwner() WalletAddress {
	return dungeon.Owner
}

func GetCoordinates(position int, width int, height int) (x int, y int) {
	x = position % width
	y = position / width

	if y > height {
		y = height
	}

	return
}

func GetPosition(x, y, width, height int) (position int) {
	position = y*width + x

	return
}

func MakeDumbDungeon(player *Player) (newDungeon Dungeon) {
	playerTiles := make([]*Tile, 0, 100)

	var groundTiles int

	for _, t := range OPENED_CARDS {
		if t.GetOwner() == player.EthAddress {
			playerTiles = append(playerTiles, t)

			if t.TileType == Ground {
				groundTiles++
			}
		}
	}

	newDungeon.Name = player.Name + "_dungeon_001"
	newDungeon.Owner = player.EthAddress
	newDungeon.Power = 100
	newDungeon.Height = 5
	newDungeon.Width = (groundTiles / 5)

	if groundTiles%5 != 0 {
		newDungeon.Width++ // Round up to fit them all
	}

	totalTiles := newDungeon.Height * newDungeon.Width

	groundPosition := 0
	monsterPosition, lootPosition := totalTiles-1, totalTiles-1

	for _, t := range playerTiles {
		var currentPosition, currentX, currentY int

		switch t.TileType {
		case Ground:
			currentPosition = groundPosition
			groundPosition++

			if groundPosition > totalTiles {
				panic("Ground tiles has exceeded total tile limit when creating a dungeon")
			}
		case Monster:
			currentPosition = monsterPosition
			monsterPosition--

			if monsterPosition < 0 {
				monsterPosition = totalTiles - 1
			}

		case Loot:
			currentPosition = lootPosition
			lootPosition--

			if lootPosition < 0 {
				lootPosition = totalTiles - 1
			}
		default:
			panic("Unknown tile type when creating dungeon") // TODO - Find out how to convert the "enum" to a string or int
		}

		currentX, currentY = GetCoordinates(currentPosition, newDungeon.Width, newDungeon.Height)

		newDungeon.Tiles = append(newDungeon.Tiles, DungeonTile{X: currentX, Y: currentY, Tile: t})

	}

	// Sort the tiles so they can be read in sequence to iterate through the map
	sort.Slice(newDungeon.Tiles, func(l, r int) bool {
		if newDungeon.Tiles[l].Y == newDungeon.Tiles[r].Y {
			if newDungeon.Tiles[l].X == newDungeon.Tiles[r].X {
				return newDungeon.Tiles[l].Tile.TileType < newDungeon.Tiles[r].Tile.TileType
			} else {
				return newDungeon.Tiles[l].X < newDungeon.Tiles[r].X
			}
		} else {
			return newDungeon.Tiles[l].Y < newDungeon.Tiles[r].Y
		}
		//
		//  // Ground, then monsters, then loot
	})

	return
}

func CalculateTile(character *Character, dungeonTile *DungeonTile, characterMultiplier float32, dungeonMultiplier float32) (neededPower float32, reward int) {
	colorDifference := float32(SmashColors(&character.ColorAttributes, &dungeonTile.Tile.ColorAttributes)) *
		dungeonMultiplier * characterMultiplier

	// Calculate Reward
	switch dungeonTile.Tile.TileType {
	case Ground:
		neededPower = colorDifference / 255
		reward = 10
	case Monster:
		neededPower = colorDifference / 25.5
		reward = 10 * int(neededPower)
	case Loot:
		neededPower = colorDifference / 25.5
		reward = 10 * int(neededPower)
	}

	return
}

func (dungeon *Dungeon) PlayDumbDungeon(character *Character) (newMatch DungeonMatch, ok bool) {
	if character.Power == 0 {
		return newMatch, false
	}

	// make a new match
	newMatch.Dungeon = dungeon
	newMatch.Character = character
	newMatch.DungeonMultiplier = 1.0
	newMatch.CharacterMultiplier = 1.0
	newMatch.StartingPower = character.Power
	// If this works, it really wasn't any easier than just remembering to put this at the end of the function body
	// probably not the best use case for defer
	defer func(newMatch *DungeonMatch, character *Character) {
		newMatch.EndingPower = character.Power
	}(&newMatch, character)

	// Tiles are already sorted in a way that allows us to row scan
	// to iterate through the dungeon until we die
	for n, t := range newMatch.Dungeon.Tiles {
		neededPower, reward := CalculateTile(character, &t, newMatch.DungeonMultiplier, newMatch.CharacterMultiplier)

		if character.Power-neededPower < 0 {
			// Deduct power, but tile was not defeated
			character.Power = 0
			break
		} else {
			newMatch.Reward += reward

			// Update Power
			character.Power -= neededPower
			dungeon.Power -= int(neededPower)

			// Add tile to list of played tiles
			newMatch.PlayedTiles = append(newMatch.PlayedTiles, n) // Tile needs ID and maybe dungeon tile? And maybe reference DT here?
		}
	}

	return newMatch, true
}
