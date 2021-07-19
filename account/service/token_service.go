package service

import (
	"context"
	"crypto/rsa"
	"github.com/nanaagyirbrown/memrizr/handler/model"
	"github.com/nanaagyirbrown/memrizr/handler/model/apperrors"
	"log"
)

type tokenService struct {
	TokenRepository model.TokenRepository
	//We will use JSON Web Tokens nas our authorization mechanism
	PrivKey *rsa.PrivateKey
	PubKey *rsa.PublicKey
	RefreshSecret string
	IDExpirationSecs	int64
	RefreshExpirationSecs	int64
}

// TSConfig will hold repositories that will eventually be injected into this
// this service layer
type TSConfig struct {
	TokenRepository model.TokenRepository
	// TokenRepository model.TokenRepository
	PrivKey *rsa.PrivateKey
	PubKey *rsa.PublicKey
	RefreshSecret string
	IDExpirationSecs	int64
	RefreshExpirationSecs	int64
}

// NewTokenService is a factory function for
// initializing a UserService with its repository layer dependencies
func NewTokenService(c *TSConfig) model.TokenService{
	return &tokenService{
		TokenRepository: c.TokenRepository,
		PrivKey: c.PrivKey,
		PubKey: c.PubKey,
		RefreshSecret: c.RefreshSecret,
		IDExpirationSecs: c.IDExpirationSecs,
		RefreshExpirationSecs: c.RefreshExpirationSecs,
	}
}

// NewPairFromUser creates fresh id and refresh tokens for the current user
// If a previous token is included, the previous token is removed from
// the tokens repository
func (t tokenService) NewPairFromUser(ctx context.Context, u *model.User, prevTokenID string) (*model.TokenPair, error) {
	idToken, err := generateIDToken(u, t.PrivKey, t.IDExpirationSecs)

	if err != nil {
		log.Printf("Error generating idToken for uid: %v. Error: %v\n", u.UID, err.Error())
		return nil, apperrors.NewInternal()
	}

	refreshToken, err := generateRefreshToken(u.UID, t.RefreshSecret, t.RefreshExpirationSecs)

	if err != nil {
		log.Printf("Error generating refreshToken for uid: %v. Error: %v\n", u.UID, err.Error())
		return nil, apperrors.NewInternal()
	}

	// TODO: store refresh tokens by calling TokenRepository methods
	// set freshly minted refresh token to valid list
	if err := t.TokenRepository.SetRefreshToken(ctx, u.UID.String(), refreshToken.ID, refreshToken.ExpiresIn); err != nil {
		log.Printf("Error storing tokenID for uid: %v. Error: %v\n", u.UID, err.Error())
		return nil, apperrors.NewInternal()
	}

	log.Printf("saved in redis ", u.UID)

	// delete user's current refresh token (used when refreshing idToken)
	if prevTokenID != "" {
		if err := t.TokenRepository.DeleteRefreshToken(ctx, u.UID.String(), prevTokenID); err != nil {
			log.Printf("Could not delete previous refreshToken for uid: %v, tokenID: %v\n", u.UID.String(), prevTokenID)
		}
	}

	return &model.TokenPair{
		IDToken:      idToken,
		RefreshToken: refreshToken.SS,
	}, nil
}