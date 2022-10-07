package pic

import (
	"embed"
	"fmt"
	"os"
	"path"
)

var DefaultCollection Collection = nil

//go:embed assets/img
var defaultImageFiles embed.FS

func init() {
	dc, err := OpenCollection(defaultImageFiles, path.Join("assets", "img"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to load embedded DefaultCollection: %s\n", err)
		return
	}
	DefaultCollection = dc
}
