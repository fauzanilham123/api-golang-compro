package models

import (
	"crypto/rand"
	"encoding/hex"
	"html"
	"strings"
	"time"

	"api-golang-compro/utils/token"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type (
    // User
    User struct {
        ID        uint      `json:"id" gorm:"primary_key"`
        Username  string    `gorm:"not null;unique" json:"username"`
        Email     string    `gorm:"not null;unique" json:"email"`
        Password  string    `gorm:"not null;" json:"password"`
        Salt      string    `gorm:"not null" json:"-"`
        CreatedAt time.Time `json:"created_at"`
        UpdatedAt time.Time `json:"updated_at"`
        LogActivity 	  []LogActivity				`json:"-"`
    }
)

func generateRandomSalt() ([]byte, error) {
    salt := make([]byte, 16) // Panjang salt yang dihasilkan (16 byte)
    _, err := rand.Read(salt)
    if err != nil {
        return nil, err
    }
    return salt, nil
}

func VerifyPassword(password, hashedPassword string) error {
    return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoginCheck(username string, password string, db *gorm.DB) (string, error) {

    var err error

    u := User{}

    err = db.Model(User{}).Where("username = ?", username).Take(&u).Error
    if err != nil {
        return "", err
    }

    // Decode salt dari hex
    salt, err := hex.DecodeString(u.Salt)
    if err != nil {
        return "", err
    }

    // Gabungkan salt dengan kata sandi yang dimasukkan saat login
    saltedPassword := append([]byte(password), salt...)

    // Decode hash dari hex
    hashedPassword, err := hex.DecodeString(u.Password)
    if err != nil {
        return "", err
    }

    err = bcrypt.CompareHashAndPassword(hashedPassword, saltedPassword)
    if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
        return "", err
    }

    token, err := token.GenerateToken(u.ID)
    if err != nil {
        return "", err
    }

    return token, nil

}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {

    // Membuat salt acak
    salt, errs := generateRandomSalt()
    if errs != nil {
        return nil, errs
    }

     // Menggabungkan salt dengan kata sandi
    saltedPassword := append([]byte(u.Password), salt...)

    //turn password into hash
    hashedPassword, errPassword := bcrypt.GenerateFromPassword(saltedPassword, bcrypt.DefaultCost)
    if errPassword != nil {
        return &User{}, errPassword
    }
    // Simpan salt dan hash dalam basis data
    u.Salt = hex.EncodeToString(salt)
    u.Password = hex.EncodeToString(hashedPassword)
    //remove spaces in username
    u.Username = html.EscapeString(strings.TrimSpace(u.Username))

    var err error = db.Create(&u).Error
    if err != nil {
        return &User{}, err
    }
    return u, nil
}