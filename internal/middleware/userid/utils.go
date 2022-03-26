package userid

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type storage struct {
	ctx        *fiber.Ctx
	secret     string
	cookieName string
}

func NewStorage(ctx *fiber.Ctx, cookieName string, secret string) *storage {
	return &storage{
		ctx:        ctx,
		secret:     secret,
		cookieName: cookieName,
	}
}

func (s *storage) getHash(data string) string {
	h := hmac.New(sha256.New, []byte(s.secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func (s *storage) checkHash(data string, hash string) bool {
	h := hmac.New(sha256.New, []byte(s.secret))
	h.Write([]byte(data))
	sign, err := hex.DecodeString(hash)
	if err != nil {
		return false
	}
	return hmac.Equal(sign, h.Sum(nil))
}

func (s *storage) Get() (string, error) {
	value := s.ctx.Cookies(s.cookieName)
	if value == "" {
		return "", errors.New("cookie with userID not found")
	}

	parts := strings.Split(value, ":")
	if len(parts) != 2 {
		return "", errors.New("invalid cookie value")
	}

	userID, hash := parts[0], parts[1]
	if s.checkHash(userID, hash) {
		return userID, nil
	}

	return "", errors.New("invalid cookie digest")
}

func (s *storage) Set(userID string) {
	cookie := new(fiber.Cookie)
	cookie.Name = s.cookieName
	cookie.Value = fmt.Sprintf("%s:%s", userID, s.getHash(userID))
	cookie.Expires = time.Now().Add(24 * time.Hour)

	s.ctx.Cookie(cookie)
}
