package authctx

import (
	"context"
	"errors"
)

// IdentityType discriminates the type of a resource identity.
type IdentityType string

const (
	IdentityTypeUser      IdentityType = "USER"
	IdentityTypeBot       IdentityType = "BOT"
	IdentityTypeConnector IdentityType = "CONNECTOR"
	IdentityTypeAgent     IdentityType = "AGENT"
	IdentityTypeDevice    IdentityType = "DEVICE"
	IdentityTypeService   IdentityType = "SERVICE"
	IdentityTypeSystem    IdentityType = "SYSTEM"
)

// Identity is the unified resource identity used for JWT claims and context.
type Identity struct {
	Type      IdentityType
	ID        string
	Name      string
	OwnerID   string
	OwnerType IdentityType
	Provider  string
}

type ctxKey string

const identityKey ctxKey = "identity"

// WithIdentity stores the identity in the context.
func WithIdentity(ctx context.Context, identity *Identity) context.Context {
	return context.WithValue(ctx, identityKey, identity)
}

// GetIdentity retrieves the identity from the context.
func GetIdentity(ctx context.Context) (*Identity, error) {
	identity, ok := ctx.Value(identityKey).(*Identity)
	if !ok || identity == nil {
		return nil, errors.New("identity not found in context")
	}
	return identity, nil
}
