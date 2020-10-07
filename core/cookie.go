package core

import (
	"github.com/space307/go-lmax-api/model"
	"github.com/space307/go-lmax-api/version"
)

const (
	sessionID = "Set-cookie"
)

type (
	cookie struct {
		sessionID string
	}
)

func createUserAgent(userID string) string {
	var result string
	result += "LMAX Java API v"
	result += version.ApiLibraryVersion
	result += " p"
	result += version.ProtocolVersion
	result += " jdk11.0.8"
	return result
}

func extractCookie(h model.Header) (cookie cookie) {
	cookie.sessionID = h.Get(sessionID)
	return
}
