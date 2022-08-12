package core

import "errors"

type WalletAddress string

type Player struct {
	EthAddress WalletAddress
	Name       string
	Gold       int
}

func GetSignatureMessage(address string, name string) string {
	return ""
}

func CheckSignatureMessage(signedMessage string) bool {
	return true
}

func NewPlayer(name string, ethAddress string, signature string) (*Player, error) {
	player := new(Player)

	if !CheckSignatureMessage(signature) {
		return nil, errors.New("Signature is not valid")
	}

	player.Name = name
	player.EthAddress = WalletAddress(ethAddress)
	player.Gold = 1000

	return player, nil
}
