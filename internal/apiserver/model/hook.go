package model

import (
	"gorm.io/gorm"

	"github.com/onexstack/onexstack/pkg/authn"

	"github.com/onexstack/fastgo/internal/pkg/rid"
)

func (m *Post) AfterCreate(tx *gorm.DB) error {
	m.PostID = rid.PostID.New(uint64(m.ID))

	return tx.Save(m).Error
}

func (m *User) AfterCreate(tx *gorm.DB) error {
	m.UserID = rid.UserID.New(uint64(m.ID))

	return tx.Save(m).Error
}

func (m *User) BeforeCreate(tx *gorm.DB) error {
	var err error
	m.Password, err = authn.Encrypt(m.Password)
	if err != nil {
		return err
	}
	return nil
}
