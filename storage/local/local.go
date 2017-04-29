package local

import (
    "io"
    "os"
)

type LocalStorageBackend struct {}

func NewLocalStorageBackend() *LocalStorageBackend {
    return &LocalStorageBackend{}
}

func (ls *LocalStorageBackend) Upload(releaseId string, pkg io.ReadSeeker) (string, error) {
//        target_dir = os.path.join(local_storage_path, release.releasetype.name, release.application.name)
//        if not os.path.exists(target_dir):
//            os.makedirs(target_dir)
//        p = os.path.join(target_dir, release.get_release_id() + ".tgz")
//        file.save(p)
//        return 'File saved', 200
    return "", nil
}

func (ls *LocalStorageBackend) Download(uri string) (io.ReadSeeker, error) {
    file, err := os.Open(uri[len("file://"):])
    if err != nil {
        return nil, err
    }
    return file, nil
}

