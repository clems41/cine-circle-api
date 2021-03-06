package repositoryModel

import (
	"gorm.io/gorm"
	"strconv"
	"time"
)

type Metadata gorm.Model

// GetID Return the ID
func (mtd *Metadata) GetID() uint {
	return mtd.ID
}

// GetIDAsString Return the ID in string format
func (mtd *Metadata) GetIDAsString() string {
	return strconv.FormatUint(uint64(mtd.ID), 10)
}

// SetID Set the ID like a string
func (mtd *Metadata) SetID(id uint) {
	mtd.ID = id
}

func (mtd *Metadata) CreatedAtNow() {
	mtd.CreatedAt = time.Now()
}

func (mtd *Metadata) UpdatedAtNow() {
	mtd.CreatedAt = time.Now()
}

func (mtd *Metadata) HasBeenModified() bool {

	if mtd == nil {
		return false
	}

	return !mtd.UpdatedAt.IsZero() && !mtd.UpdatedAt.Equal(mtd.CreatedAt)
}

