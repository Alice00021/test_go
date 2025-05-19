package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	SecretKey          = "какой-то-интересный-секретный-ключ"
	AccessTokenExpire  = 15 * time.Minute
	RefreshTokenExpire = 7 * 24 * time.Hour
	TokenExpiration    = 24 * time.Hour
)

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userID uint, username string)(string, error){
	return generateToken(userID, username, AccessTokenExpire)
}

func GenerateRefreshToken(userID uint, username string)(string, error){
	return generateToken(userID, username, RefreshTokenExpire)
}
func generateToken(userID uint, username string, expire time.Duration) (string, error) {
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SecretKey))
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrInvalidKey
}
