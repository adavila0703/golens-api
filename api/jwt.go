package api

import (
	"golens-api/config"
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

type JWTClaim struct {
	DeviceUUID string `json:"deviceUUID"`
	Salt       string
	Pepper     string
	jwt.Claims
}

var (
	// mocks
	GenerateAccessTokenF  = GenerateAccessToken
	GenerateRefreshTokenF = GenerateRefreshToken
)
var accessTokenSecret = []byte(config.Cfg.AccessTokenSecret)
var refreshTokenSecret = []byte(config.Cfg.RefreshTokenSecret)

func GenerateRefreshToken(deviceUUID string) (string, error) {
	expirationTime := time.Now().Add(72 * time.Hour)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &JWTClaim{
		DeviceUUID: deviceUUID,
		Salt:       generateRandomString(25),
		Pepper:     generateRandomString(256),
		Claims: jwt.MapClaims{
			"exp": expirationTime.Unix(),
		},
	})
	tokenString, err := token.SignedString(refreshTokenSecret)

	return tokenString, err
}

func GenerateAccessToken(deviceUUID string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &JWTClaim{
		DeviceUUID: deviceUUID,
		Pepper:     generateRandomString(256),
		Claims: jwt.MapClaims{
			"exp": expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(accessTokenSecret)

	return tokenString, err
}

func ValidateAccessToken(signedToken string) error {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(accessTokenSecret), nil
		},
	)

	if err != nil {
		return errors.WithStack(err)
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		return errors.WithStack(ErrCannotParseClaims)
	}

	currentTime := time.Now().Local().Unix()
	if expirationTime, err := claims.Claims.GetExpirationTime(); expirationTime.Unix() < currentTime || err != nil {
		return errors.WithStack(ErrTokenExpired)
	}

	return nil
}

func ExtractClaimsFromToken(signedToken string, secretToken string) (*JWTClaim, error) {
	secretTokenBytes := []byte(secretToken)

	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secretTokenBytes), nil
		},
	)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		return nil, errors.WithStack(ErrCannotParseClaims)
	}

	return claims, nil
}

var characters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func generateRandomString(length int) string {
	randomString := make([]rune, length)
	for i := range randomString {
		randomString[i] = characters[rand.Intn(len(characters))]
	}
	return string(randomString)
}
