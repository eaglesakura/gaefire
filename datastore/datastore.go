package gaefire

type KindInfo struct {
	Name    string
	Version int
}

type ManagedStoreData interface {
	GetId() string
}