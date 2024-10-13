package common

import (
	"fmt"
	"time"

	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

type TokenPayload struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiryAt  time.Time `json:"expiryAt"`
}

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPayload(id string, duration time.Duration) *TokenPayload {
	return &TokenPayload{
		ID:        id,
		CreatedAt: time.Now().UTC(),
		ExpiryAt:  time.Now().UTC().Add(duration),
	}
}

func (p *TokenPayload) Valid() error {
	if time.Now().UTC().Before(p.ExpiryAt) {
		return nil
	}

	return fmt.Errorf("token has expired")
}

func NewPaseto(symmetricKey string) (*PasetoMaker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("SymmetricKey too short should be: %v", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return maker, nil
}

func (p *PasetoMaker) CreateToken(id string, duration time.Duration) (string, error) {
	payload := NewPayload(id, duration)

	return p.paseto.Encrypt(p.symmetricKey, payload, nil)
}

func (p *PasetoMaker) VerifyToken(token string) (*TokenPayload, error) {
	payload := &TokenPayload{}

	err := p.paseto.Decrypt(token, p.symmetricKey, payload, nil)
	if err != nil {
		return nil, err
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
