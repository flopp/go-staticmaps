package sm

import (
	"fmt"
	"os"

	"github.com/Wessie/appdirs"
)

type TileCache interface {
	// Root path to store cached tiles in with no trailing slash.
	Path() string
	// Permission to set when creating missing cache directories.
	Perm() os.FileMode
}

type tileCacheStaticPath struct {
	path string
	perm os.FileMode
}

func (c *tileCacheStaticPath) Path() string {
	return c.path
}

func (c *tileCacheStaticPath) Perm() os.FileMode {
	return c.perm
}

// Stores cache files in a static path.
func NewTileCache(rootPath string, perm os.FileMode) *tileCacheStaticPath {
	return &tileCacheStaticPath{
		path: rootPath,
		perm: perm,
	}
}

// Stores cache files in a user-specific cache directory.
func NewTileCacheFromUserCache(name string, perm os.FileMode) *tileCacheStaticPath {
	app := appdirs.New("go-staticmaps", "flopp.net", "0.1")
	return &tileCacheStaticPath{
		path: fmt.Sprintf("%s/%s", app.UserCache(), name),
		perm: perm,
	}
}
