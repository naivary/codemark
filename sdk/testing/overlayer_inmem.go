package testing

import (
	"io/fs"
	"path/filepath"

	"github.com/naivary/codemark/sdk"
	"github.com/spf13/afero"
)

var _ sdk.Overlayer = (*inmemOverlayer)(nil)

type inmemOverlayer struct {
	fs afero.Fs
}

func NewInMemOverlayer(fs afero.Fs) sdk.Overlayer {
	return &inmemOverlayer{fs}
}

func (i *inmemOverlayer) Overlay() (map[string][]byte, error) {
	overlay := make(map[string][]byte)
	err := afero.Walk(i.fs, "codemark/", func(path string, info fs.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".go" {
			data, err := afero.ReadFile(i.fs, path)
			if err != nil {
				return err
			}
			p := filepath.Join("/tmp", filepath.Base(path))
			// NOTE: This must be an ABSOLUTE path
			overlay[p] = data
		}
		return nil
	})
	return overlay, err
}
