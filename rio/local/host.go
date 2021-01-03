package local

import (
	"sync"

	"github.com/desops/khan/rio"
)

type Host struct {
	// cache
	usersmu   sync.Mutex
	users     map[string]*rio.User
	groups    map[string]*rio.Group
	passwords map[string]*rio.Password
}

func (host *Host) String() string {
	return "local"
}

func New() *Host {
	return &Host{}
}
