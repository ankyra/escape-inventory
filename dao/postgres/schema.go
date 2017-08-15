// Code generated by go-bindata.
// sources:
// dao/postgres/schemas/1_initial_schema.down.sql
// dao/postgres/schemas/1_initial_schema.up.sql
// dao/postgres/schemas/2_project_metadata.down.sql
// dao/postgres/schemas/2_project_metadata.up.sql
// dao/postgres/schemas/3_migrate_existing_projects.up.sql
// dao/postgres/schemas/4_application_metadata.down.sql
// dao/postgres/schemas/4_application_metadata.up.sql
// DO NOT EDIT!

package postgres

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var __1_initial_schemaDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x09\xf2\x0f\x50\x08\x71\x74\xf2\x71\x55\x28\x4a\xcd\x49\x4d\x2c\x4e\xb5\xe6\x42\x12\x2b\x48\x4c\xce\x4e\x4c\x47\x15\x4b\x4c\xce\xb1\xe6\x02\x04\x00\x00\xff\xff\xd7\x9e\x4f\xa4\x38\x00\x00\x00")

func _1_initial_schemaDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__1_initial_schemaDownSql,
		"1_initial_schema.down.sql",
	)
}

func _1_initial_schemaDownSql() (*asset, error) {
	bytes, err := _1_initial_schemaDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "1_initial_schema.down.sql", size: 56, mode: os.FileMode(436), modTime: time.Unix(1502613619, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __1_initial_schemaUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x91\xcd\x4e\xc3\x30\x10\x84\xef\x7e\x8a\x3d\xc6\x92\x2f\x04\x81\x90\x38\x05\x64\xa4\x88\x5f\xa5\x3e\xd0\x53\xb5\x72\x57\xc5\xd0\x24\xd6\xc6\xad\x78\x7c\x54\xe7\x17\x91\x00\xd7\x99\xc9\xe4\x9b\xf5\x6d\xa1\x33\xa3\xc1\x64\x37\x0f\x1a\xf2\x3b\x78\x7a\x36\xa0\x5f\xf3\x95\x59\x01\xd3\x9e\xb0\x21\x48\x04\x00\x40\x85\x25\xc1\x11\xd9\xbe\x21\x27\x67\xe9\x95\x54\x10\xf5\x2e\xb5\x71\xdb\xc1\x4d\x2f\x2e\xa5\x8a\xe6\x91\xb8\x71\x75\x35\x38\xe7\x69\x67\x94\x14\x70\x8b\x01\x21\xd0\x67\x68\x25\xcf\xf5\x3b\xd9\xf0\x33\xfb\x52\xe4\x8f\x59\xb1\x86\x7b\xbd\x4e\x4e\x14\xaa\xaf\x55\xfd\x37\x52\xc8\x6b\x21\x7e\x99\xe2\xd1\x7e\xe0\xae\x9f\xb2\xf8\xa7\xa5\x2d\xed\xd2\x03\xbb\x39\x79\x8a\x37\x16\xa8\x53\xfc\xff\x80\x68\xf7\x7f\xc1\xed\xb8\x3e\xf8\xcd\xb7\x67\x18\x0f\xed\x89\x4b\xd7\xc4\x5b\xbb\x2a\xcc\x90\x75\xb5\x6a\x52\x13\xa1\xbe\x02\x00\x00\xff\xff\x53\x97\x8f\x29\x00\x02\x00\x00")

func _1_initial_schemaUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__1_initial_schemaUpSql,
		"1_initial_schema.up.sql",
	)
}

func _1_initial_schemaUpSql() (*asset, error) {
	bytes, err := _1_initial_schemaUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "1_initial_schema.up.sql", size: 512, mode: os.FileMode(436), modTime: time.Unix(1502613606, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __2_project_metadataDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x09\xf2\x0f\x50\x08\x71\x74\xf2\x71\x55\x28\x28\xca\xcf\x4a\x4d\x2e\xb1\xe6\x02\x04\x00\x00\xff\xff\xa5\x8e\xd4\xaa\x14\x00\x00\x00")

func _2_project_metadataDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__2_project_metadataDownSql,
		"2_project_metadata.down.sql",
	)
}

func _2_project_metadataDownSql() (*asset, error) {
	bytes, err := _2_project_metadataDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "2_project_metadata.down.sql", size: 20, mode: os.FileMode(436), modTime: time.Unix(1502797856, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __2_project_metadataUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x0e\x72\x75\x0c\x71\x55\x08\x71\x74\xf2\x71\x55\x28\x28\xca\xcf\x4a\x4d\x2e\x51\xd0\xe0\x52\x50\x50\x50\xc8\x4b\xcc\x4d\x55\x08\x73\x0c\x72\xf6\x70\x0c\xd2\x30\x36\xd2\xd4\x01\x8b\xa6\xa4\x16\x27\x17\x65\x16\x94\x64\xe6\xe7\x29\x94\xa4\x56\x94\x40\x44\xf3\x8b\xd2\x43\x83\x7c\xe0\xaa\x8d\x4c\xcd\xa0\xca\x73\xf2\xd3\xf3\x91\xd4\x05\x04\x79\xfa\x3a\x06\x45\x2a\x78\xbb\x46\x6a\x80\xcc\xd7\xe4\xd2\xb4\xe6\x02\x04\x00\x00\xff\xff\x46\xec\xb6\x00\x84\x00\x00\x00")

func _2_project_metadataUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__2_project_metadataUpSql,
		"2_project_metadata.up.sql",
	)
}

func _2_project_metadataUpSql() (*asset, error) {
	bytes, err := _2_project_metadataUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "2_project_metadata.up.sql", size: 132, mode: os.FileMode(436), modTime: time.Unix(1502815779, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __3_migrate_existing_projectsUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xf2\xf4\x0b\x76\x0d\x0a\x51\xf0\xf4\x0b\xf1\x57\x28\x28\xca\xcf\x4a\x4d\x2e\xd1\xc8\x4b\xcc\x4d\xd5\x51\x48\x49\x2d\x4e\x2e\xca\x2c\x28\xc9\xcc\xcf\xd3\x51\xc8\x2f\x4a\x0f\x0d\xf2\xd1\x51\xc8\xc9\x4f\xcf\xd7\xe4\x0a\x76\xf5\x71\x75\x0e\x51\x48\xc9\x2c\x2e\xc9\xcc\x4b\x2e\xd1\x80\x6a\xd4\xd4\x51\x50\x57\x87\x61\x2e\xb7\x20\x7f\x5f\x85\xa2\xd4\x9c\xd4\xc4\xe2\x54\x6b\x2e\x40\x00\x00\x00\xff\xff\x5b\xed\x91\x00\x68\x00\x00\x00")

func _3_migrate_existing_projectsUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__3_migrate_existing_projectsUpSql,
		"3_migrate_existing_projects.up.sql",
	)
}

func _3_migrate_existing_projectsUpSql() (*asset, error) {
	bytes, err := _3_migrate_existing_projectsUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "3_migrate_existing_projects.up.sql", size: 104, mode: os.FileMode(436), modTime: time.Unix(1502806903, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __4_application_metadataDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x09\xf2\x0f\x50\x08\x71\x74\xf2\x71\x55\x48\x2c\x28\xc8\xc9\x4c\x4e\x2c\xc9\xcc\xcf\xb3\xe6\x02\x04\x00\x00\xff\xff\xc6\x19\x92\xd8\x18\x00\x00\x00")

func _4_application_metadataDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__4_application_metadataDownSql,
		"4_application_metadata.down.sql",
	)
}

func _4_application_metadataDownSql() (*asset, error) {
	bytes, err := _4_application_metadataDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "4_application_metadata.down.sql", size: 24, mode: os.FileMode(436), modTime: time.Unix(1502815842, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __4_application_metadataUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x64\x8e\xcd\x8a\x83\x30\x18\x45\xf7\x3e\xc5\xdd\x69\xc0\xcd\x38\xcc\x30\x30\xab\x8c\x13\xa9\xd4\xfe\x10\xd3\xa2\x2b\x09\xfa\x51\x52\xac\x09\x9a\x45\x1f\xbf\xd4\x52\x37\x6e\xef\xb9\x1c\x4e\x2a\x05\x57\x02\x8a\xff\x15\x02\x79\x86\xfd\x41\x41\x54\x79\xa9\x4a\x68\xe7\x7a\xd3\x6a\x6f\xec\x80\x28\x00\x80\x41\xdf\x08\x67\x2e\xd3\x0d\x97\xd1\x47\xf2\xc3\x62\xcc\xbb\x1b\xed\x95\x5a\xbf\xa0\xcf\x84\xc5\x33\xe8\x68\x6a\x47\xe3\x66\x85\x12\x95\xc2\xbf\xc8\xf8\xa9\x50\x08\xc3\xd7\xa1\xd7\x9e\x26\xdf\x8c\xd4\x93\x9e\xa8\x31\xdd\xe2\x48\xbe\xbe\xd9\xfa\x6e\x2f\x16\x9e\xee\x7e\x45\x8e\x32\xdf\x71\x59\x63\x2b\xea\xe8\x99\x19\xbf\xa3\x58\xc0\x7e\x83\x47\x00\x00\x00\xff\xff\x15\xdb\xf7\x53\xe6\x00\x00\x00")

func _4_application_metadataUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__4_application_metadataUpSql,
		"4_application_metadata.up.sql",
	)
}

func _4_application_metadataUpSql() (*asset, error) {
	bytes, err := _4_application_metadataUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "4_application_metadata.up.sql", size: 230, mode: os.FileMode(436), modTime: time.Unix(1502815831, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"1_initial_schema.down.sql": _1_initial_schemaDownSql,
	"1_initial_schema.up.sql": _1_initial_schemaUpSql,
	"2_project_metadata.down.sql": _2_project_metadataDownSql,
	"2_project_metadata.up.sql": _2_project_metadataUpSql,
	"3_migrate_existing_projects.up.sql": _3_migrate_existing_projectsUpSql,
	"4_application_metadata.down.sql": _4_application_metadataDownSql,
	"4_application_metadata.up.sql": _4_application_metadataUpSql,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"1_initial_schema.down.sql": &bintree{_1_initial_schemaDownSql, map[string]*bintree{}},
	"1_initial_schema.up.sql": &bintree{_1_initial_schemaUpSql, map[string]*bintree{}},
	"2_project_metadata.down.sql": &bintree{_2_project_metadataDownSql, map[string]*bintree{}},
	"2_project_metadata.up.sql": &bintree{_2_project_metadataUpSql, map[string]*bintree{}},
	"3_migrate_existing_projects.up.sql": &bintree{_3_migrate_existing_projectsUpSql, map[string]*bintree{}},
	"4_application_metadata.down.sql": &bintree{_4_application_metadataDownSql, map[string]*bintree{}},
	"4_application_metadata.up.sql": &bintree{_4_application_metadataUpSql, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

