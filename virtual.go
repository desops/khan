package khan

import (
	"os"
)

// Virtual is an in-memory model for all changes we are capable of making on a
// server, including service status and a filesystem. This way a dry run can
// be pixel-perfect.
//
// If you start managing files with contents too large for RAM, this will need
// to be improved. (Currently all managed files, including contents, are in
// kept in memory in their entirety.)
type Virtual struct {
	// Host metadata extracted from uname command
	Uname    string
	Hostname string
	Kernel   string
	OS       string
	Arch     string

	// File system model
	Files   map[string]os.FileInfo
	Content map[string]string

	// User and group model
	Users  map[string]*User
	Groups map[string]*Group

	cacheUsers  map[string]*User
	cacheGroups map[string]*Group
	withshadow  bool
}

func NewVirtual() *Virtual {
	v := &Virtual{
		Files:   make(map[string]os.FileInfo),
		Content: make(map[string]string),
		Users:   make(map[string]*User),
		Groups:  make(map[string]*Group),
	}
	return v
}
