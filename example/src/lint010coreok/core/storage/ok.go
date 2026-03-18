package storage

type Service interface {
	Read()
}

type Store interface {
	Write()
}
