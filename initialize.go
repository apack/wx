package wx

import (
	"bytes"
	"embed"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

//go:embed example/app
var exampleApp embed.FS

func Initialize(dir string) error {
	if _, err := os.Stat(dir); err == nil {
		return fmt.Errorf("directory %s already exists", dir)
	} else if !os.IsNotExist(err) {
		return err
	}
	return copyDir(exampleApp, "example/app", dir)
}

func copyDir(files embed.FS, src, dst string) error {
	entries, err := files.ReadDir(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			err := copyDir(files, path.Join(src, entry.Name()), path.Join(dst, entry.Name()))
			if err != nil {
				return err
			}
		} else {
			err := copyFile(files, path.Join(src, entry.Name()), path.Join(dst, entry.Name()))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func copyFile(files embed.FS, src, dst string) error {
	fmt.Printf("creating %s\n", dst)
	data, err := files.ReadFile(src)
	if err != nil {
		return err
	}
	if dst == "gen.go" {
		return nil
	}
	data = bytes.Replace(data,
		[]byte("//go:generate go run ../../cmd/wx build"),
		[]byte("//go:generate wx build"), -1)
	return ioutil.WriteFile(dst, data, 0644)
}
