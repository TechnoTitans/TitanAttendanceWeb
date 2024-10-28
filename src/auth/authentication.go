package auth

import (
	"TitanAttendance/src/utils"
	"errors"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type Authentication struct {
	Pin string
}

type UserAccess struct {
	UserLevel int
}

const (
	none = iota
	plain
	admin
)

const PasswordExpiration = 1 * time.Hour

var PasswordCreationTime = time.UnixMicro(0)
var userPin = ""

func (a *Authentication) CheckPin() (UserAccess, error) {
	if a.Pin == utils.GetAdminPassword() {
		return UserAccess{UserLevel: admin}, nil
	}
	if !PasswordCreationTime.IsZero() && a.Pin == userPin && !HasPasswordExpired() {
		return UserAccess{UserLevel: plain}, nil
	}

	return UserAccess{UserLevel: none}, errors.New("invalid pin")
}

func CheckWithCookie(cookie http.Cookie) (UserAccess, error) {
	pin := Authentication{Pin: cookie.Value}
	return pin.CheckPin()
}

func (a *UserAccess) IsAdmin() bool {
	return a.UserLevel == admin
}

func (a *UserAccess) IsPlain() bool {
	return a.UserLevel == plain
}

func (a *UserAccess) IsAllowed() bool {
	return (a.IsAdmin() || a.IsPlain()) && a.UserLevel != none
}

func CreateUserPin() string {
	if HasPasswordExpired() {
		userPin = strconv.Itoa(rand.Intn(900000) + 100000)
		PasswordCreationTime = time.Now()
	}
	return userPin
}

func HasPasswordExpired() bool {
	return time.Since(PasswordCreationTime) > PasswordExpiration
}

func PlainUser() UserAccess {
	return UserAccess{UserLevel: plain}
}

func SavePinCookie(w http.ResponseWriter, userAuth Authentication) {
	http.SetCookie(w, &http.Cookie{
		Name:     "TitanAttendancePin",
		Value:    userAuth.Pin,
		Path:     "/",
		MaxAge:   int(PasswordExpiration.Seconds()),
		HttpOnly: true,
	})
}

func ClearPinCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "TitanAttendancePin",
		Value:    "",
		Path:     "/",
		MaxAge:   0,
		HttpOnly: true,
	})
}
