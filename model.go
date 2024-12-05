package utils

import "sync"

var (
	_authPool = sync.Pool{
		New: func() any {
			return NewAuth()
		}}
)

// CreateAuth return *Auth
func CreateAuth() *Auth {

	auth := _authPool.Get().(*Auth)

	return auth
}

func (o *Auth) initial() {
	o.UserID = 0
	o.UserRight = 0
}

func (o *Auth) Release() {
	o.initial()
	_authPool.Put(o)
}

// CreateAuthes return *AuthCollection
func CreateAuthes() *AuthCollection {

	authes := &AuthCollection{}

	return authes
}

func (o *AuthCollection) Release() {
	for i := 0; i < len(*o); i++ {
		(*o)[i].Release()
	}
}

// NewAuth return *Auth
func NewAuth() *Auth {

	auth := &Auth{}

	return auth
}

// Auth model
type Auth struct {
	UserID    uint64 `json:"userID"`
	UserRight int64  `json:"userRight"`
}

// NewAuthes return *AuthCollection
func NewAuthes() *AuthCollection {

	authes := &AuthCollection{}

	return authes
}

// AuthCollection Auth list
type AuthCollection []Auth
