package storage

import (
	"errors"
	"time"

	"github.com/Stolkerve/kappa/pgk/utils"
	"gorm.io/gorm"
)

type Call struct {
	ID          string    `gorm:"primaryKey"`
	CreatedAt   time.Time `gorm:"not null"`
	Stdout      string    `gorm:"not null"`
	Stderr      string    `gorm:"not null"`
	ErrorMsg    string
	Duration    time.Duration `gorm:"not null"`
	MemoryUsage uint32        `gorm:"not null"`
	Fail        bool          `gorm:"not null"`
	FunctionID  string        `gorm:"not null"`
}

func (c *Call) BeforeCreate(tx *gorm.DB) (err error) {
	id, err := utils.GenerateUid()
	if err != nil {
		err = errors.New("Can't save invalid id")
		return
	}
	c.ID = id
	return
}
