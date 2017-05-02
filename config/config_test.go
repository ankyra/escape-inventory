package config

import (
    "os"
	"testing"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type configSuite struct{}

var _ = Suite(&configSuite{})

func (s *configSuite) Test_NewConfig_uses_Sqlite_and_local_storage_by_default(c *C) {
    env := []string{}
    conf, err := NewConfig(env)
    c.Assert(err, IsNil)
    c.Assert(conf.Port, Equals, "7770")
    c.Assert(conf.Database, Equals, "sqlite")
    c.Assert(conf.DatabaseSettings.Path, Equals, "/var/lib/escape/registry.db")
    c.Assert(conf.StorageBackend, Equals, "local")
    c.Assert(conf.StorageSettings.Path, Equals, "/var/lib/escape/releases")
}

func (s *configSuite) Test_LoadConfig_InMemoryDb(c *C) {
    env := []string{}
    conf, err := LoadConfig("testdata/in_memory_db.json", env)
    c.Assert(err, IsNil)
    c.Assert(conf.Database, Equals, "memory")
}

func (s *configSuite) Test_LoadConfig_Uses_Default_Storage_Backend_If_Not_Configured(c *C) {
    env := []string{}
    conf, err := LoadConfig("testdata/in_memory_db.json", env)
    c.Assert(err, IsNil)
    c.Assert(conf.Database, Equals, "memory")
    c.Assert(conf.StorageBackend, Equals, "local")
    c.Assert(conf.StorageSettings.Path, Equals, "/var/lib/escape/releases")
}

func (s *configSuite) Test_LoadConfig_SqliteDb(c *C) {
    env := []string{}
    conf, err := LoadConfig("testdata/sqlite_db.json", env)
    c.Assert(err, IsNil)
    c.Assert(conf.Database, Equals, "sqlite")
    c.Assert(conf.DatabaseSettings.Path, Equals, "/var/lib/escape/registry.db")
}

func (s *configSuite) Test_LoadConfig_PostgresDb(c *C) {
    env := []string{}
    conf, err := LoadConfig("testdata/postgres_db.json", env)
    c.Assert(err, IsNil)
    c.Assert(conf.Database, Equals, "postgres")
    c.Assert(conf.DatabaseSettings.Path, Equals, "")
    c.Assert(conf.DatabaseSettings.PostgresUrl, Equals, "postgres://")
}

func (s *configSuite) Test_LoadConfig_LocalStorage(c *C) {
    env := []string{}
    conf, err := LoadConfig("testdata/local_storage_backend.json", env)
    c.Assert(err, IsNil)
    c.Assert(conf.StorageBackend, Equals, "local")
    c.Assert(conf.StorageSettings.Path, Equals, "/var/lib/escape/releases")
}

func (s *configSuite) Test_LoadConfig_GCS(c *C) {
    env := []string{}
    conf, err := LoadConfig("testdata/gcs_storage_backend.json", env)
    c.Assert(err, IsNil)
    c.Assert(conf.StorageBackend, Equals, "gcs")
    c.Assert(conf.StorageSettings.Path, Equals, "")
    c.Assert(conf.StorageSettings.Bucket, Equals, "gs://escape-releases/")
    c.Assert(conf.StorageSettings.Credentials, Equals, "test")
}

func (s *configSuite) Test_LoadConfig_fails_if_not_exists(c *C) {
    env := []string{}
    conf, err := LoadConfig("testdata/doesnt_exist.json", env)
    c.Assert(conf, IsNil)
    c.Assert(err, Not(IsNil))
}

func (s *configSuite) Test_LoadConfig_fails_if_malformed(c *C) {
    env := []string{}
    conf, err := LoadConfig("testdata/malformed.json", env)
    c.Assert(conf, IsNil)
    c.Assert(err, Not(IsNil))
}


func (s *configSuite) Test_LoadConfig_Parses_Yaml(c *C) {
    env := []string{}
    conf, err := LoadConfig("testdata/yaml_config.yaml", env)
    c.Assert(err, IsNil)
    c.Assert(conf.Port, Equals, "9876")
    c.Assert(conf.Database, Equals, "sqlite")
    c.Assert(conf.DatabaseSettings.Path, Equals, "/var/lib/escape/registry.db")
    c.Assert(conf.StorageBackend, Equals, "gcs")
    c.Assert(conf.StorageSettings.Path, Equals, "")
    c.Assert(conf.StorageSettings.Bucket, Equals, "gs://escape-releases/")
    c.Assert(conf.StorageSettings.Credentials, Equals, "test")
}

func (s *configSuite) Test_LoadConfig_Parses_Yml(c *C) {
    env := []string{}
    conf, err := LoadConfig("testdata/yml_config.yml", env)
    c.Assert(err, IsNil)
    c.Assert(conf.Database, Equals, "sqlite")
    c.Assert(conf.DatabaseSettings.Path, Equals, "/var/lib/escape/registry.db")
    c.Assert(conf.StorageBackend, Equals, "gcs")
    c.Assert(conf.StorageSettings.Path, Equals, "")
    c.Assert(conf.StorageSettings.Bucket, Equals, "gs://escape-releases/")
    c.Assert(conf.StorageSettings.Credentials, Equals, "test")
}

func (s *configSuite) Test_LoadConfig_fails_if_yaml_malformed(c *C) {
    env := []string{}
    conf, err := LoadConfig("testdata/malformed.yaml", env)
    c.Assert(conf, IsNil)
    c.Assert(err, Not(IsNil))
}

func (s *configSuite) Test_LoadConfig_fails_if_file_cant_be_read(c *C) {
    env := []string{}
    os.Chmod("testdata/cant_read.json", 0)
    conf, err := LoadConfig("testdata/cant_read.json", env)
    c.Assert(conf, IsNil)
    c.Assert(err, Not(IsNil))
    os.Chmod("testdata/cant_read.json", 0666)
}

func (s *configSuite) Test_NewConfig_Uses_EnvironmentVariables(c *C) {
    env := []string{
        "PORT=9876",
        "DATABASE=memory",
        "DATABASE_SETTINGS_PATH=",
        "STORAGE_BACKEND=gcs",
        "STORAGE_SETTINGS_PATH=",
        "STORAGE_SETTINGS_BUCKET=gs://escape-releases/",
        "STORAGE_SETTINGS_CREDENTIALS=test",
    }
    conf, err := NewConfig(env)
    c.Assert(err, IsNil)
    c.Assert(conf.Port, Equals, "9876")
    c.Assert(conf.Database, Equals, "memory")
    c.Assert(conf.DatabaseSettings.Path, Equals, "")
    c.Assert(conf.StorageBackend, Equals, "gcs")
    c.Assert(conf.StorageSettings.Path, Equals, "")
    c.Assert(conf.StorageSettings.Bucket, Equals, "gs://escape-releases/")
    c.Assert(conf.StorageSettings.Credentials, Equals, "test")
}

func (s *configSuite) Test_LoadConfig_Uses_EnvironmentVariables(c *C) {
    env := []string{
        "PORT=9876",
        "DATABASE=memory",
        "DATABASE_SETTINGS_PATH=",
        "DATABASE_SETTINGS_POSTGRES_URL=postgres://",
        "STORAGE_BACKEND=local",
        "STORAGE_SETTINGS_PATH=/test/",
        "STORAGE_SETTINGS_BUCKET=",
        "STORAGE_SETTINGS_CREDENTIALS=",
    }
    conf, err := LoadConfig("testdata/yml_config.yml", env)
    c.Assert(err, IsNil)
    c.Assert(conf.Port, Equals, "9876")
    c.Assert(conf.Database, Equals, "memory")
    c.Assert(conf.DatabaseSettings.Path, Equals, "")
    c.Assert(conf.DatabaseSettings.PostgresUrl, Equals, "postgres://")
    c.Assert(conf.StorageBackend, Equals, "local")
    c.Assert(conf.StorageSettings.Path, Equals, "/test/")
    c.Assert(conf.StorageSettings.Bucket, Equals, "")
    c.Assert(conf.StorageSettings.Credentials, Equals, "")
}
