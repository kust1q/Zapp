package mocks

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/kust1q/Zapp/backend/internal/config"
)

func GenerateTestRSAKeys() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate RSA key: %w", err)
	}
	publicKey := &privateKey.PublicKey
	return privateKey, publicKey, nil
}

func NewTestAuthConfig() (*config.AuthServiceConfig, error) {
	privateKey, publicKey, err := GenerateTestRSAKeys()
	if err != nil {
		return nil, err
	}

	return &config.AuthServiceConfig{
		PrivateKey:  privateKey,
		PublicKey:   publicKey,
		AccessTTL:   15 * time.Minute,
		RefreshTTL:  24 * time.Hour,
		RecoveryTTL: 1 * time.Hour,
	}, nil
}

func MockAvatarReader() io.Reader {
	return strings.NewReader("fake_avatar_content")
}
