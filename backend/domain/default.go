package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DefaultFields contain struct fields common to all domain models.
// Meant to be used via struct composition only.
type DefaultFields struct {
	DefaultFieldId
	DefaultFieldCreated
	DefaultFieldUpdated
}

func NewDefaultFields(id string, created, updated time.Time) DefaultFields {
	return DefaultFields{
		DefaultFieldId:      DefaultFieldId{ID: id},
		DefaultFieldCreated: DefaultFieldCreated{CreatedAt: created},
		DefaultFieldUpdated: DefaultFieldUpdated{UpdatedAt: updated},
	}
}

func (d *DefaultFields) BeforeCreate(tx *gorm.DB) error {
	err := d.DefaultFieldId.BeforeCreate(tx)
	if err != nil {
		return err
	}
	err = d.DefaultFieldCreated.BeforeCreate(tx)
	if err != nil {
		return err
	}
	return nil
}

type DefaultFieldId struct {
	// ID is the technical primary identifier.
	ID string `gorm:"primaryKey;type:uuid;column:id"`
}

func (d *DefaultFieldId) BeforeCreate(tx *gorm.DB) error {
	if d.ID == "" { // do not overwrite a uuid
		d.ID = uuid.NewString()
	}
	return nil
}

type DefaultFieldCreated struct {
	// CreatedAt is the date the object was created in database.
	CreatedAt time.Time `gorm:"not null; column:created_at"`
}

func (d *DefaultFieldCreated) BeforeCreate(_ *gorm.DB) error {
	if d.CreatedAt.IsZero() {
		d.CreatedAt = time.Now()
	}
	return nil
}

type DefaultFieldUpdated struct {
	// UpdatedAt is the last time the object was updated in database.
	UpdatedAt time.Time `gorm:"not null; column:updated_at"`
}

func (d *DefaultFieldUpdated) BeforeUpdate(db *gorm.DB) error {
	d.UpdatedAt = time.Now()
	return nil
}
