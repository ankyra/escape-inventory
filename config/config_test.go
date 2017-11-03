/*
Copyright 2017 Ankyra

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
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
	c.Assert(conf.DatabaseSettings.Path, Equals, "/var/lib/escape/inventory.db")
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
	c.Assert(conf.DatabaseSettings.Path, Equals, "/var/lib/escape/inventory.db")
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
	c.Assert(conf.DatabaseSettings.Path, Equals, "/var/lib/escape/inventory.db")
	c.Assert(conf.StorageBackend, Equals, "gcs")
	c.Assert(conf.StorageSettings.Path, Equals, "")
	c.Assert(conf.StorageSettings.Bucket, Equals, "gs://escape-releases/")
	c.Assert(conf.StorageSettings.Credentials, Equals, "test")
	c.Assert(conf.UserServiceURL, Equals, "http://user-service.com")

}

func (s *configSuite) Test_LoadConfig_Parses_Yml(c *C) {
	env := []string{}
	conf, err := LoadConfig("testdata/yml_config.yml", env)
	c.Assert(err, IsNil)
	c.Assert(conf.Database, Equals, "sqlite")
	c.Assert(conf.DatabaseSettings.Path, Equals, "/var/lib/escape/inventory.db")
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

func (s *configSuite) Test_NewConfig_Uses_EnvironmentVariables(c *C) {
	env := []string{
		"PORT=9876",
		"DATABASE=memory",
		"DATABASE_SETTINGS_PATH=",
		"STORAGE_BACKEND=gcs",
		"STORAGE_SETTINGS_PATH=",
		"STORAGE_SETTINGS_BUCKET=gs://escape-releases/",
		"STORAGE_SETTINGS_CREDENTIALS=test",
		"USER_SERVICE_URL=http://example.com",
		"DEV=true",
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
	c.Assert(conf.UserServiceURL, Equals, "http://example.com")
	c.Assert(conf.Dev, Equals, true)
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
		"USER_SERVICE_URL=http://example.com",
		"DEV=true",
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
	c.Assert(conf.UserServiceURL, Equals, "http://example.com")
	c.Assert(conf.Dev, Equals, true)
}
