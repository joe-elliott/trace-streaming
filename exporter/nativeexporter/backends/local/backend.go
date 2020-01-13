package local

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/joe-elliott/trace-streaming/exporter/nativeexporter/backends"
)

type localBackend struct {
	path string
}

func NewBackend(path string) backends.Writer {
	return &localBackend{
		path: path,
	}
}

func (b *localBackend) Write(bIndex []byte, bTraces []byte, bBloom []byte) error {
	folderName := fmt.Sprintf("%d", time.Now().Unix())
	rootPath := path.Join(b.path, folderName)

	err := os.Mkdir(rootPath, os.ModePerm)
	if err != nil {
		return err
	}

	pathIndex := path.Join(rootPath, "index")
	err = ioutil.WriteFile(pathIndex, bIndex, os.ModePerm)
	if err != nil {
		return err
	}

	pathTraces := path.Join(rootPath, "traces")
	err = ioutil.WriteFile(pathTraces, bTraces, os.ModePerm)
	if err != nil {
		return err
	}

	pathBloom := path.Join(rootPath, "bloom")
	err = ioutil.WriteFile(pathBloom, bBloom, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
