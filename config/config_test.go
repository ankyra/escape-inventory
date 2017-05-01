package config

import (
	"testing"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type configSuite struct{}

var _ = Suite(&configSuite{})

func (s *configSuite) Test_NewConfig_uses_Sqlite_and_local_storage_by_default(c *C) {
    conf := NewConfig()
    c.Assert(conf.Database, Equals, "sqlite")
    c.Assert(conf.DatabaseSettings.Path, Equals, "/var/lib/escape/registry.db")
    c.Assert(conf.StorageBackend, Equals, "local")
    c.Assert(conf.StorageSettings.Path, Equals, "/var/lib/escape/releases")
}

func (s *configSuite) Test_LoadConfig_InMemoryDb(c *C) {
    conf, err := LoadConfig("testdata/in_memory_storage_config.json")
    c.Assert(err, IsNil)
    c.Assert(conf.Database, Equals, "memory")
}

func (s *configSuite) Test_LoadConfig_SqliteDb(c *C) {
    conf, err := LoadConfig("testdata/sqlite_storage_config.json")
    c.Assert(err, IsNil)
    c.Assert(conf.Database, Equals, "sqlite")
    c.Assert(conf.DatabaseSettings.Path, Equals, "/var/lib/escape/registry.db")
}

func (s *configSuite) Test_LoadConfig_LocalStorage(c *C) {
    conf, err := LoadConfig("testdata/local_storage_backend.json")
    c.Assert(err, IsNil)
    c.Assert(conf.StorageBackend, Equals, "local")
    c.Assert(conf.StorageSettings.Path, Equals, "/var/lib/escape/releases")
}

func (s *configSuite) Test_LoadConfig_GCS(c *C) {
    conf, err := LoadConfig("testdata/gcs_storage_backend.json")
    c.Assert(err, IsNil)
    c.Assert(conf.StorageBackend, Equals, "gcs")
    c.Assert(conf.StorageSettings.Path, Equals, "")
    c.Assert(conf.StorageSettings.Bucket, Equals, "gs://escape-releases/")
    c.Assert(conf.StorageSettings.Credentials["project-id"], Equals, "test")
}

func (s *configSuite) Test_LoadConfig_fails_if_not_exists(c *C) {
    conf, err := LoadConfig("testdata/doesnt_exist.json")
    c.Assert(conf, IsNil)
    c.Assert(err, Not(IsNil))
}

func (s *configSuite) Test_LoadConfig_fails_if_malformed(c *C) {
    conf, err := LoadConfig("testdata/malformed.json")
    c.Assert(conf, IsNil)
    c.Assert(err, Not(IsNil))
}


func (s *configSuite) Test_LoadConfig_Parses_Yaml(c *C) {
    conf, err := LoadConfig("testdata/yaml_config.yaml")
    c.Assert(err, IsNil)
    c.Assert(conf.Database, Equals, "sqlite")
    c.Assert(conf.DatabaseSettings.Path, Equals, "/var/lib/escape/registry.db")
    c.Assert(conf.StorageBackend, Equals, "gcs")
    c.Assert(conf.StorageSettings.Path, Equals, "")
    c.Assert(conf.StorageSettings.Bucket, Equals, "gs://escape-releases/")
    c.Assert(conf.StorageSettings.Credentials["project-id"], Equals, "test")
}
