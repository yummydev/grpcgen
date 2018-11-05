package phdcid

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// contextKey type setup to serve as context key type for the mw package
type contextKey string

func (c contextKey) String() string {
	return strings.Join([]string{"mw contextKey:", string(c)}, " ")
}

const cidContextKeyName = contextKey("ConversationID")

// GetCIDKeyName provides the key value used to store Conversation ID
//	in a context as a string
func GetCIDKeyName() string {
	return string(cidContextKeyName)
}

// GetCIDContextKeyName is an alternative to GetCIDKeyName()
func GetCIDContextKeyName() string {
	return cidContextKeyName.String()
}

// CreateCID - Create a new conversation ID
func CreateCID() string {
	uuidStruct, _ := uuid.NewRandom()
	return uuidStruct.String()
}

// GetCIDFromContext - Get ConversationID from context
func GetCIDFromContext(ctx context.Context) string {
	var cid string

	val := ctx.Value(cidContextKeyName)
	if val != nil {
		cid = val.(string)
	}
	return cid
}

// NewContextWithCID can be used when you have no context
//	and want one with a Conversation ID
func NewContextWithCID() (ctx context.Context, conversationID string) {
	conversationID = CreateCID()
	return SetCIDInContext(context.Background(), conversationID), conversationID
}

// SetCIDInContext - Set ConversationID in context
func SetCIDInContext(ctx context.Context, cid string) context.Context {
	return context.WithValue(ctx, cidContextKeyName, cid)
}

// GetCIDFromHTTP - Get ConversationID from http header
func GetCIDFromHTTP(header http.Header) string {
	return header.Get(GetCIDKeyName())
}

// SetCIDInHTTP - Set ConversationID in context
func SetCIDInHTTP(header http.Header, cid string) {
	header.Set(GetCIDKeyName(), cid)
}
