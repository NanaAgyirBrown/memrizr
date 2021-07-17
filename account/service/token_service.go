package service

import (
	"context"
	"crypto/rsa"
	"github.com/nanaagyirbrown/memrizr/handler/model"
	"github.com/nanaagyirbrown/memrizr/handler/model/apperrors"
	"log"
)

type TokenService struct {
	//We will use JSON Web Tokens nas our authorization mechanism
	PrivKey *rsa.PrivateKey
	PubKey *rsa.PublicKey
	RefreshSecret string
}

// TSConfig will hold repositories that will eventually be injected into this
// this service layer
type TSConfig struct {
	// TokenRepository model.TokenRepository
	PrivKey *rsa.PrivateKey
	PubKey *rsa.PublicKey
	RefreshSecret string
}

// NewTokenService is a factory function for
// initializing a UserService with its repository layer dependencies
func NewTokenService(c *TSConfig) model.TokenService{
	return &TokenService{
		PrivKey: c.PrivKey,
		PubKey: c.PubKey,
		RefreshSecret: c.RefreshSecret,
	}
}

// NewPairFromUser creates fresh id and refresh tokens for the current user
// If a previous token is included, the previous token is removed from
// the tokens repository
func (t TokenService) NewPairFromUser(ctx context.Context, u *model.User, prevTokenID string) (*model.TokenPair, error) {
	idToken, err := generateIDToken(u, t.PrivKey)

	if err != nil {
		log.Printf("Error generating idToken for uid: %v. Error: %v\n", u.UID, err.Error())
		return nil, apperrors.NewInternal()
	}

	refreshToken, err := generateRefreshToken(u.UID, t.RefreshSecret)

	if err != nil {
		log.Printf("Error generating refreshToken for uid: %v. Error: %v\n", u.UID, err.Error())
		return nil, apperrors.NewInternal()
	}

	// TODO: store refresh tokens by calling TokenRepository methods

	return &model.TokenPair{
		IDToken:      idToken,
		RefreshToken: refreshToken.SS,
	}, nil
}