package models

import (
	databases "SnappVotingBack/app"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

type User struct {
	Id          int64     `json:"id"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	IsSuperUser bool      `json:"isSuperUser"`
	CreatedAt   time.Time `json:"createdAt"`
}

func (u *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	u.Password = string(bytes)
	return err
}
func (u *User) Get() error {
	return databases.PostgresDB.QueryRow("SELECT id,password,is_superuser FROM users WHERE email = $1", u.Email).Scan(&u.Id, &u.Password, &u.IsSuperUser)
}
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
func (u *User) Create() (err error) {
	u.CreatedAt = time.Now().UTC()
	_, err = databases.PostgresDB.Query("INSERT INTO users (email,password,is_superuser,created_at) VALUES ($1,$2,$3,$4)", u.Email, u.Password, true, u.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
func (u *User) UpdatePassword() error {
	_, err := databases.PostgresDB.Exec("UPDATE users SET password = $1 WHERE id = $2", u.Password, u.Id)
	return err
}
func (u *User) Auth() (access string, err error) {
	godotenv.Load("local.env")
	atClaims := jwt.MapClaims{}
	now := time.Now().UTC()

	atClaims["user_id"] = u.Id
	atClaims["is_superuser"] = u.IsSuperUser
	atClaims["exp"] = now.Add(time.Hour * 12).Unix()
	atClaims["iat"] = now.Unix() // The time at which the token was issued.
	atClaims["nbf"] = now.Unix()

	access, err = jwt.NewWithClaims(jwt.SigningMethodHS512, atClaims).SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return access, err
	}

	return access, err
}
