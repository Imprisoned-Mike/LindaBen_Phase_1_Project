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
	
	RoleID   uint
	UserRole *Role `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` //Belongs to Role


	AvatarID *uint
	Avatar   *File `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
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
func GetUsers(User *[]User) (err error) {
	err = db.Db.Find(User).Error
	if err != nil {
		return err
	}
	return nil
}

// Get user by name
func GetUserByName(name string) (User, error) {
	var user User
	err := db.Db.Where("name=?", name).Find(&user).Error
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
	err := db.Db.Where("id=?", id).Find(&user).Error
	if err != nil {
		return User{}, err
	}
	return user, nil
}

// Get user by id
func GetUser(User *User, id int) (err error) {
	err = db.Db.Where("id = ?", id).First(User).Error
	if err != nil {
		return err
	}
	return nil
}

// Update user
func UpdateUser(User *User) (err error) {
	err = db.Db.Omit("password").Updates(User).Error
	if err != nil {
		return err
	}
	return nil
}
