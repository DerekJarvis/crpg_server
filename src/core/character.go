package core

import (
	"math/rand"
)

type Character struct {
	Name  string
	Owner WalletAddress
	ColorAttributes
	Power float32
}

func (character Character) GetOwner() WalletAddress {
	return character.Owner
}

const maxCharacterStartPoints int = 100

func NewCharacter(name string, owner WalletAddress) *Character {
	c := new(Character)

	c.Name = name
	c.Owner = owner
	c.Red = rand.Intn(25) + 10
	c.Green = rand.Intn(25) + 10
	c.Blue = rand.Intn(25) + 10

	highestStat := GetHighestStat(&c.ColorAttributes)

	remainingPoints := maxCharacterStartPoints - c.Red - c.Green - c.Blue

	*highestStat += remainingPoints

	c.Power = 100

	return c
}

func (character Character) GetPower() (output int) {
	characterColors := ConvertColorsToArray(&character.ColorAttributes)

	for _, v := range characterColors {
		output += *v
	}

	return
}
