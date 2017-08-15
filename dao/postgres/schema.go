// Code generated by go-bindata.
// sources:
// dao/postgres/schemas/.1_initial_schema.up.sql.swp
// dao/postgres/schemas/.3_migrate_existing_projects.up.sql.swp
// dao/postgres/schemas/.4_application_metadata.up.sql.swp
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

var __1_initial_schemaUpSqlSwp = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\xd9\x4f\x6b\x13\x41\x18\x06\xf0\xa7\x7a\xf1\xd2\xd9\xa2\xe0\x79\xd4\x4b\x02\xdb\x6c\x93\xd8\x3f\x22\x1e\xa2\x44\x08\x5a\x95\x34\x8a\xf1\x12\xde\x6c\xa6\x9b\xb1\xc9\xee\x3a\x33\x5b\x1a\x0f\xfa\x05\x14\x3c\x28\x08\x82\x47\x3f\x88\x7e\x0a\xf1\x0b\xf8\x0d\x3c\x2a\xc9\xa6\x96\xd6\xd4\xea\xc5\x22\x7d\x7f\x97\x65\x67\xde\x7d\xe7\x79\xe7\x16\xd2\x5d\x7a\xd0\x58\x97\xab\xa5\xcb\x00\xb0\x00\x6c\x2e\xbc\x6a\x9f\x79\xff\x0e\x5f\x77\x80\xae\x4d\x89\x62\x8b\xa3\x74\xad\x5e\x3d\xb2\x08\xc0\xb3\x69\xc3\xc0\x9a\x30\xa8\x9a\x5e\x4a\xc6\x8d\x82\x28\x49\xc9\xf5\x27\x6b\x91\x76\xfd\xac\x5b\x0a\x93\x61\x40\xf1\xd6\xc8\x50\xa0\x6c\x48\xa9\x5a\x34\x2a\xd2\xd6\x99\x51\xd0\xa3\x24\x48\x13\xeb\x22\xa3\x6c\x60\xc3\xbe\x1a\x92\x0d\xca\x1d\x1d\x6b\xa7\x69\xd0\xc9\x57\x4a\x59\x5a\xb2\x4f\x06\x7f\x12\x89\xb1\x93\x21\x73\x9b\x8b\x6b\xf3\xa8\x56\xca\x4b\xe3\xd7\x4b\x17\x2f\xc8\x73\x67\xef\x1f\x77\x2a\xc6\x18\x63\x8c\x31\xc6\xd8\x3f\xe4\xd2\x39\x3c\x07\x70\x6a\xfa\x7e\x7e\xfa\x9c\x3b\xf0\x64\x8c\x31\xc6\x18\x63\x8c\x31\xc6\xd8\xff\x8b\x7a\xc0\x8b\x79\x00\x22\xff\xff\x7f\xf7\xf7\xff\x17\x0f\xf8\xe4\x01\x1f\x3c\xe0\xa5\x07\x3c\xf5\x00\xf2\x80\x2b\x1e\xb0\xe2\x01\xcb\x1e\xb0\xe0\x01\xdf\x04\xf0\x59\x00\x1f\x05\xf0\x56\x00\x6f\x04\xf0\x5a\x00\x56\x00\x8f\x04\x70\x4d\x00\x05\x01\x9c\x16\xf9\x19\xdf\xe7\x8f\x79\x60\xc6\x18\x63\x8c\x31\xc6\x18\x3b\xa1\x8a\x57\x21\xa5\x94\xf7\x9a\x8d\xf5\x5a\xb3\x2d\x6f\xd5\xdb\x85\xd4\x24\x8f\x55\xe8\x7c\x19\x99\x24\x4b\x3b\x31\x0d\x55\x71\x52\x94\x2a\x33\xd4\xd6\xea\x24\x96\x3a\x76\xbe\x9c\x2c\xee\x15\xc9\x6d\x32\x61\x9f\x4c\xa1\xb2\xbc\x52\xf4\xf3\x2f\xf2\x56\x3f\x77\xaa\x95\xa2\x8f\x1b\xcd\x7a\xad\x55\x97\xad\xda\xf5\xdb\x75\xd9\xb8\x29\xef\xdc\x6d\xc9\xfa\xc3\xc6\x46\x6b\x43\x52\x38\x90\x85\x99\xa1\x8c\x1a\x28\xb2\xaa\xa3\x7b\xbe\xcc\x8c\xf6\x77\x5b\xe7\xc9\x32\xa3\xf7\x9f\x9e\x67\xdb\xfb\x68\xd6\xee\xdf\x86\x4b\x29\xdc\xa2\x48\x1d\x12\x70\x7c\x03\xbe\xdc\x56\x66\x7c\x3f\x07\xe2\xcd\x3c\x69\xbc\x31\x54\x8e\x7a\xe4\x48\x3a\xb5\xe3\xf2\xa5\x69\x87\x5f\x6b\x0f\x1b\x66\xb2\xb9\xef\xfe\xcb\x95\xb5\xf1\x8c\xbf\x19\x65\xda\x4b\x16\xf0\x23\x00\x00\xff\xff\x73\x61\xa0\xd6\x00\x30\x00\x00")

func _1_initial_schemaUpSqlSwpBytes() ([]byte, error) {
	return bindataRead(
		__1_initial_schemaUpSqlSwp,
		".1_initial_schema.up.sql.swp",
	)
}

func _1_initial_schemaUpSqlSwp() (*asset, error) {
	bytes, err := _1_initial_schemaUpSqlSwpBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: ".1_initial_schema.up.sql.swp", size: 12288, mode: os.FileMode(420), modTime: time.Unix(1502815786, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __3_migrate_existing_projectsUpSqlSwp = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\xce\x31\x4a\x04\x51\x0c\x80\xe1\x80\xb5\x20\x7a\x81\x55\xeb\x9d\xac\x8e\xb0\x5e\xc1\xc2\x52\xc1\x6a\xc8\xbc\x8d\x6f\x9e\x3a\xf3\x9e\x49\x06\x77\x1a\x2f\x61\xe7\x95\xbc\x88\x07\xf0\x00\xa2\xd8\xaf\x95\xc2\x92\xaf\x4b\x08\xe1\x6f\x17\xd7\x17\x97\xb3\x65\x75\x06\x00\xb0\x07\xf0\xb4\xf3\x72\xf3\xf1\xf6\x0a\xef\x6b\x80\x56\x0b\xd1\xa0\xb0\x49\xab\x69\xb9\xf1\x08\x00\x9e\x7f\x1e\xa2\x4a\xc0\x5a\x56\x85\xc4\x26\x8c\xb9\x90\x75\xdf\xbb\x98\xac\x1b\xdb\x2a\xe4\x1e\x69\xb8\x9f\x84\x90\x35\x50\xe1\xb9\x70\x4c\x6a\x32\xe1\x8a\x32\x96\xac\x16\x85\x15\x35\x74\xdc\x93\x62\xdd\xf4\x29\x0a\x19\x37\xbc\x4e\x6a\x69\x88\x4d\x91\x7c\xc7\xc1\xb4\x1a\x4b\xa5\x8f\x0f\xbf\xa9\x73\x6e\x8b\x8d\x76\x3b\x3f\xdf\x85\xfa\xf4\x64\xf1\x35\x1e\x1f\x1d\xce\x0e\xf6\xaf\xfe\xbb\xca\x39\xe7\x9c\x73\xce\x39\xf7\x87\x3e\x03\x00\x00\xff\xff\x29\x0d\x0e\x39\x00\x10\x00\x00")

func _3_migrate_existing_projectsUpSqlSwpBytes() ([]byte, error) {
	return bindataRead(
		__3_migrate_existing_projectsUpSqlSwp,
		".3_migrate_existing_projects.up.sql.swp",
	)
}

func _3_migrate_existing_projectsUpSqlSwp() (*asset, error) {
	bytes, err := _3_migrate_existing_projectsUpSqlSwpBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: ".3_migrate_existing_projects.up.sql.swp", size: 4096, mode: os.FileMode(420), modTime: time.Unix(1502815857, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __4_application_metadataUpSqlSwp = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\xda\xb1\x6e\xd3\x40\x1c\x06\xf0\xaf\x4c\x2c\x60\x04\x1b\x12\xd2\x1f\x90\x68\x23\x85\xb8\x4d\x0a\x45\x82\xc5\x14\x17\x22\x5a\x40\xae\x1b\x9a\x29\xba\xd8\x27\xc7\x90\xd8\xc7\xdd\x55\x6a\x16\x98\x78\x02\x46\x98\x79\x0d\x1e\x80\x07\x60\x63\xe6\x31\x50\x9b\x50\x21\x3a\x94\x01\x51\x21\xbe\xdf\x72\xf6\xfd\x3f\x9d\xbf\xd5\xb2\x87\xcb\xbd\xee\x96\xac\xb5\x56\x01\xe0\x02\xf0\xfc\xc6\xbb\xfe\xe7\x0f\xef\xf1\x6d\x1f\x18\x3a\xa3\x54\xe5\x70\x92\xa1\x2b\xd7\x4e\x0c\x01\x78\x3d\x3f\x30\x74\x36\x0b\x3b\x36\x37\xca\xfa\x69\x58\xd4\x46\xf9\xd1\xe1\x5e\x51\xfa\xd1\xde\xb0\x95\xd5\x93\x50\x55\x2f\xa7\x56\x85\xda\x65\xca\xe8\x9b\x56\x17\xa5\xf3\x76\x1a\xe6\xaa\x0e\x4d\xed\x7c\x61\xb5\x0b\x5d\x36\xd2\x13\xe5\xc2\xd5\x81\x32\x66\x5c\x66\xca\x97\x75\x35\x98\x68\xaf\x72\xe5\x55\x6b\xcf\xb4\xdc\xab\xf1\xef\x14\x23\xfa\x2f\x9c\x43\xa7\xbd\xb2\x7c\x70\x75\xfd\xda\x55\xb9\x74\x71\xe7\xb4\x0b\x11\x11\x11\x11\x11\xd1\x5f\xe4\xcd\x02\xde\x00\x38\x33\xbf\x3f\x3b\x5f\x17\x7e\x59\x89\x88\x88\x88\x88\x88\x88\xe8\xdf\xa5\x72\xe0\xeb\x79\xe0\x72\x30\xfb\xfe\xff\xe3\xfd\xff\x4b\x00\x7c\x0a\x80\x8f\x01\xf0\x36\x00\x7a\x01\x70\x2f\x00\xae\x04\xb3\xec\xc3\xe0\x94\x8b\x13\x11\x11\x11\x11\x11\x11\xfd\x11\x8d\xbb\x10\x11\x79\x96\x74\xb7\xa2\xa4\x2f\x8f\xe3\xfe\x52\xa5\x26\xba\x29\xc6\xd6\x2f\x74\xe6\x1b\x87\xe3\x71\x5d\xd4\xe2\xf5\xbe\x97\x07\xf1\x46\xb4\xb3\x99\xca\xe2\x62\x73\x36\x51\x5e\x3b\x3f\xb0\x7a\xac\x95\xd3\x83\x32\x97\x5e\x94\xac\x3f\x8a\x92\xa5\xf6\xad\xdb\x8d\x63\xf1\x5c\xbb\xcc\x96\xc6\x97\x75\x25\x69\xbc\x9b\x1e\x0b\xcc\x1f\x7b\x74\x4a\xa7\xdd\x98\x0d\x0e\x5a\x1d\xed\xae\xb4\xef\x34\x9a\x82\xf5\x24\x8e\xd2\x58\xd2\xe8\xfe\x66\x2c\xdd\x0d\x79\xf2\x34\x95\x78\xb7\xbb\x9d\x6e\xcb\x4f\xff\xca\xcb\x12\xbe\x07\x00\x00\xff\xff\xd0\x38\x67\xe1\x00\x30\x00\x00")

func _4_application_metadataUpSqlSwpBytes() ([]byte, error) {
	return bindataRead(
		__4_application_metadataUpSqlSwp,
		".4_application_metadata.up.sql.swp",
	)
}

func _4_application_metadataUpSqlSwp() (*asset, error) {
	bytes, err := _4_application_metadataUpSqlSwpBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: ".4_application_metadata.up.sql.swp", size: 12288, mode: os.FileMode(384), modTime: time.Unix(1502815841, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
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
	".1_initial_schema.up.sql.swp": _1_initial_schemaUpSqlSwp,
	".3_migrate_existing_projects.up.sql.swp": _3_migrate_existing_projectsUpSqlSwp,
	".4_application_metadata.up.sql.swp": _4_application_metadataUpSqlSwp,
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
	".1_initial_schema.up.sql.swp": &bintree{_1_initial_schemaUpSqlSwp, map[string]*bintree{}},
	".3_migrate_existing_projects.up.sql.swp": &bintree{_3_migrate_existing_projectsUpSqlSwp, map[string]*bintree{}},
	".4_application_metadata.up.sql.swp": &bintree{_4_application_metadataUpSqlSwp, map[string]*bintree{}},
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

