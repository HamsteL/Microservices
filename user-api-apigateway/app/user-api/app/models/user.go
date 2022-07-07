package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type User struct {
	ID           uint      `gorm:"primary_key" json:"id"`
	FirstName    string    `gorm:"size:255;" json:"first_name"`
	LastName     string    `gorm:"size:255;" json:"last_name"`
	Email        string    `gorm:"size:255; not null;" json:"email"`
	PasswordHash string    `gorm:"size:1000; not null;" json:"password_hash"`
	UpdateDate   time.Time `json:"update_date"`
}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {
	u.UpdateDate = time.Now()
	var err error
	err = db.Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) FindAllUsers(db *gorm.DB) (*[]User, error) {
	var err error
	var users []User
	err = db.Debug().Model(&User{}).Limit(10).Find(&users).Error
	if err != nil {
		return &[]User{}, err
	}
	return &users, err
}

func (u *User) FindUserByID(db *gorm.DB, uid uint32) (*User, error) {
	var err error
	err = db.Select([]string{"first_name", "last_name", "email"}).Where("id = ?", uid).Find(&u).Error
	if err != nil {
		return &User{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("User Not Found")
	}
	return u, err
}

func (u *User) FindUserByEmail(db *gorm.DB, email string) (*User, error) {
	var err error
	err = db.Debug().Model(User{}).Where("email = ?", email).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("User Not Found")
	}
	return u, err
}

func (u *User) UpdateUser(db *gorm.DB, uid uint32) (*User, error) {
	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"first_name":  u.FirstName,
			"last_name":   u.LastName,
			"update_date": time.Now(),
		},
	)
	if db.Error != nil {
		return &User{}, db.Error
	}
	err := db.Debug().Model(&User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) DeleteUser(db *gorm.DB, uid uint32) (int64, error) {
	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).Delete(&User{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
