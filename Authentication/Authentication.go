package Authentication

import (
	"crypto/rand"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"main/Database"
	"time"
)

type Authentication struct {
	db *Database.DB
	// jwtSecretKey is the JWT secret key. Each time the server starts, new key is generated.
	jwtUserSecretKey      []byte
	jwtExpirationDuration time.Duration
	logger                *logrus.Logger
}
type claims struct {
	jwt.MapClaims
	Email string `json:"username"`
}
type Credentials struct {
	Email    string
	Password string
}

func CreateAuthentication(db *Database.DB, jwtExpirationInMinutes int64, logger *logrus.Logger) (*Authentication, error) {
	UserSecretKey, err := generateRandomKey()
	if err != nil {
		return nil, err
	}

	if db == nil {
		return nil, errors.New("the database is essential for authentication")
	}

	return &Authentication{
		db:               db,
		jwtUserSecretKey: UserSecretKey,

		jwtExpirationDuration: time.Duration(int64(time.Minute) * jwtExpirationInMinutes),
		logger:                logger,
	}, nil
}

func (a *Authentication) AuthenticateUserWithCredentials(cred Credentials) error {

	// check user exist
	user, err := a.db.GetUserByEmail(cred.Email)
	if err != nil {
		if err.Error() == "record not found" {
			return errors.New("user not exist")
		}
		return err
	}

	//check user password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(cred.Password)); err != nil {
		return errors.New("the password is not correct")
	}
	return nil
}

func (a *Authentication) GenerateJwtToken(email string) (token *string, err error) {
	expirationTime := time.Now().Add(a.jwtExpirationDuration)
	tokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
		Email: email,
		MapClaims: jwt.MapClaims{
			"expired_at": expirationTime.Unix(),
		},
	})

	// Calculate the signed account string format of JWT key
	tokenString := ""

	tokenString, err = tokenJWT.SignedString(a.jwtUserSecretKey)

	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

func (a *Authentication) CheckToken(token string) (*claims, error) {
	c := &claims{}
	tkn, err := jwt.ParseWithClaims(token, c, func(token *jwt.Token) (interface{}, error) {
		return a.jwtUserSecretKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.New("invalid token")
		}
		if err.Error() == "token signature is invalid: signature is invalid" {
			return nil, errors.New("invalid token")
		}
		a.logger.WithError(err).Warn("can not validate the token of the user")
		return nil, errors.New("something went wrong")
	}

	if !tkn.Valid {
		return nil, errors.New("unauthorized")
	}

	return c, nil
}

// generateRandomKey
// Each time that Auth is initialized, generateRandomKey is called to
// generate another key
func generateRandomKey() ([]byte, error) {
	jwtKey := make([]byte, 32)
	if _, err := rand.Read(jwtKey); err != nil {
		return nil, err
	}

	return jwtKey, nil
}
