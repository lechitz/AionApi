package testdata

import (
	"github.com/lechitz/AionApi/internal/core/domain"
)

var TestPerfectToken = domain.TokenDomain{
	UserID: 1,
	Token:  "token_abc123",
	//CreatedAt: time.Now(),
	//ExpiresAt: time.Now().Add(24 * time.Hour),
}
