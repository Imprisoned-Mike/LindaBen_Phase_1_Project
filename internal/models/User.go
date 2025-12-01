package models

import (
	"LindaBen_Phase_1_Project/internal/db"
	"html"
	"strings"

	"gorm.io/gorm"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Name     string
	Password string
	Email    string `gorm:"unique"`
	Phone    string
	Role     string // e.g. "admin", "school_admin:34", "vendor_admin:12"

	AvatarID *uint
	Avatar   *File `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// Save user details
func (user *User) Save() (*User, error) {
	err := db.Create(&user).Error
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
	user.Username = html.EscapeString(strings.TrimSpace(user.Username))
	return nil
}

// Get all users
func GetUsers(User *[]User) (err error) {
	err = database.Db.Find(User).Error
	if err != nil {
		return err
	}
	return nil
}

// Get user by username
func GetUserByUsername(username string) (User, error) {
	var user User
	err := database.Db.Where("username=?", username).Find(&user).Error
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
func GetUserById(id uint) (User, error) {
	var user User
	err := database.Db.Where("id=?", id).Find(&user).Error
	if err != nil {
		return User{}, err
	}
	return user, nil
}

// Get user by id
func GetUser(User *User, id int) (err error) {
	err = database.Db.Where("id = ?", id).First(User).Error
	if err != nil {
		return err
	}
	return nil
}

// Update user
func UpdateUser(User *User) (err error) {
	err = database.Db.Omit("password").Updates(User).Error
	if err != nil {
		return err
	}
	return nil
}
