package lobby

import (
	"crypto/subtle"
)

const AdminTokenHeader = "X-Admin-Token"

// ValidateAdminToken reports whether token matches the lobby admin token.
func ValidateAdminToken(lobby *Lobby, token string) bool {
	if lobby == nil || lobby.AdminToken == "" || token == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(lobby.AdminToken), []byte(token)) == 1
}
