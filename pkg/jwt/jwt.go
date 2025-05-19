package jwt

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
)

/* func init_secret_key(cfg *config.Config) string {
	SecretKey := fmt.Sprintf(
		"SecretKey=%s",
		cfg.Jwt.SecretKey,
	)
	return SecretKey
} */

type JWTManager struct {
	secretKey string
	AccessTokenExpire time.Duration
	RefreshTokenExpire time.Duration
}

func NewJWTManager(secretKey string) *JWTManager {
	return &JWTManager{
		secretKey: secretKey,
		AccessTokenExpire:  15 * time.Minute,
		RefreshTokenExpire: 7 * 24 * time.Hour,
	}
}

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}


func (jm *JWTManager) GenerateAccessToken(userID uint, username string)(string, error){
	return jm.generateToken(userID, username, jm.AccessTokenExpire)
}

func (jm  *JWTManager) GenerateRefreshToken(userID uint, username string)(string, error){
	return jm.generateToken(userID, username, jm.RefreshTokenExpire)
}
func (jm  *JWTManager) generateToken(userID uint, username string, expire time.Duration) (string, error) {
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jm.secretKey))
}

func (jm  *JWTManager) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jm.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrInvalidKey
}
