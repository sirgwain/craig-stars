package cs

import "time"

// APIToken is used to store a hash of an api token after authenticating with discord and redirecting to a tool
type APIToken struct {
	DBObject
	UserID      int64      `json:"userId"`
	Name        string     `json:"name"`
	TokenPrefix string     `json:"tokenPrefix"`
	Scope       string     `json:"scope"`
	ExpiresAt   time.Time  `json:"expiresAt"`
	LastUsedAt  *time.Time `json:"lastUsedAt,omitempty"`
	RevokedAt   *time.Time `json:"revokedAt,omitempty"`
}
