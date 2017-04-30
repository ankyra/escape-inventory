package local

import (
    "io"
    "os"
    "fmt"
    "path/filepath"
    "github.com/ankyra/escape-registry/shared"
)

type LocalStorageBackend struct {
    localStoragePath string
}

func NewLocalStorageBackend() *LocalStorageBackend {
    return &LocalStorageBackend{}
}

func NewLocalStorageBackendWithStoragePath(localStoragePath string) *LocalStorageBackend {
    return &LocalStorageBackend{
        localStoragePath: localStoragePath,
    }
}

func (ls *LocalStorageBackend) getStoragePath() (string, error) {
    return filepath.Abs(ls.localStoragePath)
}

func (ls *LocalStorageBackend) Upload(releaseId *shared.ReleaseId, pkg io.ReadSeeker) (string, error) {
    storage, err := ls.getStoragePath()
    if err != nil {
        return "", err
    }
    typ := releaseId.Type
    name := releaseId.Name
    targetDir := filepath.Join(storage, typ, name)
    if !PathExists(targetDir) {
        os.MkdirAll(targetDir, 0755)
    }
    if !IsDir(targetDir) {
        return "", fmt.Errorf("Path %s exists, but is not a directory", targetDir)
    }
    target := filepath.Join(targetDir, releaseId.ToString() + ".tgz")
    dst, err := os.Create(target)
    if err != nil {
        return "", err
    }
    if _, err := io.Copy(dst, pkg); err != nil {
        return "", err
    }
    return "file://" + target, nil
}

func (ls *LocalStorageBackend) Download(uri string) (io.ReadSeeker, error) {
    file, err := os.Open(uri[len("file://"):])
    if err != nil {
        return nil, err
    }
    return file, nil
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func IsDir(path string) bool {
	st, err := os.Stat(path)
	if err != nil {
		return false
	}
	return st.IsDir()
}

