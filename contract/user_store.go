package contract

import "github.com/SalehGoML/entity"

type UserStore interface {
	UserReadStore
	UserWriteStore
}

type UserWriteStore interface {
	Save(u entity.User)
}

type UserReadStore interface {
	Load(serializationMode string) []entity.User
}
