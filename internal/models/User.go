package models

import (
	"LindaBen_Phase_1_Project/internal/db"
	"html"
	"strings"

	"gorm.io/gorm"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Model
	Name     string `json:"name"`
	Password string `json:"-"`
	Email    string `gorm:"unique" json:"email"`
	Phone    string `json:"phone"`

	Roles string `json:"roles"`

	AvatarID *uint `json:"avatarId"`
	Avatar   *File `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"avatar,omitempty"`
}

// Save user details
func (user *User) Save() (*User, error) {
	err := db.Db.Create(&user).Error
	if err != nil {
		return &User{}, err
	}
	return user, nil
}

// Generate encrypted password
func (user *User) BeforeSave(*gorm.DB) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(passwordHash)
	user.Name = html.EscapeString(strings.TrimSpace(user.Name))
	return nil
}

// Get all users
func GetUsers(Users *[]User) (err error) {
	err = db.Db.Find(Users).Error
	if err != nil {
		return err
	}
	return nil
}

// Get user by email
func GetUserByEmail(email string) (User, error) {
	var user User
	err := db.Db.Where("email=?", email).Find(&user).Error
	if err != nil {
		return User{}, err
	}
	return user, nil
}

// Validate user password
func (user *User) ValidateUserPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

// Get user by id
func GetUser(Users *User, id int) (err error) {
	err = db.Db.Preload("Avatar").Where("id = ?", id).First(Users).Error
	if err != nil {
		return err
	}
	return nil
}

// Update user
func UpdateUser(Users *User) (err error) {
	err = db.Db.Updates(Users).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete User
func DeleteUser(user *User) (err error) {
	err = db.Db.Delete(user).Error
	if err != nil {
		return err
	}
	return nil
}

// Upload User Avatar
func UploadUserAvatar(user *User, file *File) (err error) {
	user.Avatar = file
	err = db.Db.Save(user).Error
	if err != nil {
		return err
	}
	return nil
}
