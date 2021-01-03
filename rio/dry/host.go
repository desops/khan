package dry

import (
	"fmt"
	"runtime"
	"sync"

	hhost "github.com/desops/khan/host"
)

type Host struct {
	// the user/group to emulate operations as
	uid uint32
	gid uint32

	cascade hhost.Host

	fsmu sync.Mutex
	fs   map[string]*File

	usersmu   sync.Mutex
	users     map[string]*hhost.User
	groups    map[string]*hhost.Group
	passwords map[string]*hhost.Password
}

func (host *Host) String() string {
	return fmt.Sprintf("virtual host (uid %d gid %d) cascade → %s", host.uid, host.gid, host.cascade)
}

func New(uid, gid uint32, cascade hhost.Host) *Host {
	return &Host{
		uid:       uid,
		gid:       gid,
		cascade:   cascade,
		fs:        map[string]*File{},
		users:     map[string]*hhost.User{},
		groups:    map[string]*hhost.Group{},
		passwords: map[string]*hhost.Password{},
	}
}

func (host *Host) Info() (*hhost.Info, error) {
	if host.cascade != nil {
		return host.cascade.Info()
	}

	// TODO let you customize this
	return &hhost.Info{
		Hostname: "dry",
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
	}, nil
}