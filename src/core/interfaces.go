package core

type HasOwner interface {
	GetOwner() WalletAddress
}
