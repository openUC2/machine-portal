// Package templates provides a path for webpage templates to be lazily loaded from the filesystem.
package templates

import (
	"io/fs"
	"os"
)

type Client struct {
	Config Config
}

func NewClient(c Config) *Client {
	return &Client{
		Config: c,
	}
}

func (c *Client) GetFS() fs.FS {
	paths := c.Config.Paths
	switch len(paths) {
	case 0:
		return nil
	case 1:
		return os.DirFS(paths[0])
	}

	overlay := &OverlayFS{}
	suboverlay := overlay
	for i, path := range paths {
		suboverlay.Upper = os.DirFS(path)
		if i+1 == len(paths)-1 {
			suboverlay.Lower = os.DirFS(paths[i+1])
			break
		}
		nextSuboverlay := &OverlayFS{}
		suboverlay.Lower = nextSuboverlay
		suboverlay = nextSuboverlay
	}
	return overlay
}
