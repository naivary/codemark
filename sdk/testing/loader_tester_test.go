package testing

import (
	"io/fs"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"golang.org/x/tools/go/packages"
)

func TestLoaderTester(t *testing.T) {
	tester, err := NewLoaderTester()
	if err != nil {
		t.Errorf("err occured: %v\n", err)
	}
	fs, err := tester.NewFS()
	if err != nil {
		t.Errorf("err occured: %v\n", err)
	}
	overlay, err := createOverlay(fs)
	if err != nil {
		t.Errorf("err occured: %v\n", err)
	}
	cfg := &packages.Config{
		Mode:    packages.NeedTypesInfo | packages.NeedSyntax | packages.NeedTypes,
		Overlay: overlay,
		Dir:     "codemark",
	}
}

func createOverlay(afs afero.Fs) (map[string][]byte, error) {
	overlay := make(map[string][]byte)
	err := afero.Walk(afs, "/", func(path string, info fs.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".go" {
			data, err := afero.ReadFile(afs, path)
			if err != nil {
				return err
			}
			// NOTE: This must be an ABSOLUTE path
			overlay[path] = data
		}
		return nil
	})
	return overlay, err
}
