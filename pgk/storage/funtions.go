package storage

import (
	"errors"
	"time"

	"github.com/Stolkerve/kappa/pgk/utils"
	"gorm.io/gorm"
)

type Function struct {
	ID        string    `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"not null"`
	Wasm      []byte    `gorm:"type:blob;not null"`
}

func (f *Function) BeforeCreate(tx *gorm.DB) (err error) {
	id, err := utils.GenerateUid()
	if err != nil {
		err = errors.New("Can't save invalid id")
		return
	}
	f.ID = id
	return
}
