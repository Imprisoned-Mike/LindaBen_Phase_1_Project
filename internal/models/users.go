package models

import (
	"LindaBen_Phase_1_Project/internal/db"
	"html"
	"strings"

	"gorm.io/gorm"

	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	gorm.Model
	Name     string `json:"name"`
	Password string `json:"-"`
	Email    string `gorm:"unique" json:"email"`
	Phone    string `json:"phone"`

	Roles string `json:"roles"`

	AvatarID *uint `json:"avatarId"`
	Avatar   *File `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"avatar,omitempty"`
}

// Save user details
func (user *Users) Save() (*Users, error) {
	err := db.Db.Create(&user).Error
	if err != nil {
		return &Users{}, err
	}
	return user, nil
}

// Generate encrypted password
func (user *Users) BeforeSave(*gorm.DB) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(passwordHash)
	user.Name = html.EscapeString(strings.TrimSpace(user.Name))
	return nil
}

// Get all users
func GetUsers(Users *[]Users) (err error) {
	err = db.Db.Find(Users).Error
	if err != nil {
		return err
	}
	return nil
}

// Get user by email
func GetUserByEmail(email string) (Users, error) {
	var user Users
	err := db.Db.Where("email=?", email).Find(&user).Error
	if err != nil {
		return Users{}, err
	}
	return user, nil
}

// Validate user password
func (user *Users) ValidateUserPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

// Get user by id
func GetUser(Users *Users, id int) (err error) {
	err = db.Db.Where("id = ?", id).First(Users).Error
	if err != nil {
		return err
	}
	return nil
}

// Update user
func UpdateUser(Users *Users) (err error) {
	err = db.Db.Omit("password").Updates(Users).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete User
func DeleteUser(user *Users) (err error) {
	err = db.Db.Delete(user).Error
	if err != nil {
		return err
	}
	return nil
}

// Upload User Avatar
func UploadUserAvatar(user *Users, file *File) (err error) {
	user.Avatar = file
	err = db.Db.Save(user).Error
	if err != nil {
		return err
	}
	return nil
}
