# Escape Registry

Please see http://escape.ankyra.io for the full documentation.

## Features

* Centralized Escape release metadata
  * Postgres database
  * SQLite database
* Upload and download packages
  * Google Cloud Storage back-end
  * Local file-system storage back-end
* Importing and Exporting

## Usage

```
escape-registry [CONFIG_FILE]
```

The Escape Registry can be configured using a simple JSON or YAML file (default
`/etc/escape-registry/config.json`), and/or environment variables. If the
provided configuration file does not exist the program falls back to the
following default configuration: 

```
{
  "port": "7770",
  "database": "sqlite",
  "database_settings": {
    "path": "/var/lib/escape/registry.db"
  },
  "storage_backend": "local"
  "storage_settings": {
    "path": "/var/lib/escape/releases"
  }
}
```

This configures the Registry with an SQLite database and a local file system
storage backend.

### Environment Variables

Configuration variables can be overridden using environment variables:

```
PORT
DATABASE
DATABASE_SETTINGS_PATH
DATABASE_SETTINGS_POSTGRES_URL
STORAGE_BACKEND
STORAGE_SETTINGS_PATH
STORAGE_SETTINGS_BUCKET
STORAGE_SETTINGS_CREDENTIALS
```

## Storage Backends

The storage backends are used to store and retrieve uploaded packages.

### Local file storage

```
{
  "storage_backend": "local",
  "storage_settings": {
    "path": "/var/lib/escape/releases"
  }
}
```

Stores uploaded packages on the local file system.
The `path` variable points to a directory in which the releases will be stored.

### Google Cloud Storage

```
{
  "storage_backend": "gcs",
  "storage_settings": {
    "bucket": "my-bucket",
    "credentials": "/my/secret/service/credentials.json",
  }
}
```

Stores uploaded packages in Google Cloud Storage. 
The `credentials` variable is optional, but should point to an existing service
account json file if provided.

## Databases

### SQLite

```
{
  "database": "sqlite",
  "database_settings": {
    "path": "/var/lib/escape/registry.db"
  }
}
```

This is the default database, which will work out of the box (provided the path
is accessible).


### Postgresql

```
{
  "database": "postgres",
  "database_settings": {
    "postgres_url": "postgres://user:pass@localhost/database?sslmode=disable"
  }
}
```

Postgresql can be configured using the `postgres_url` variable. Please see 
https://godoc.org/github.com/lib/pq for the full connection string parameters.

### In Memory 

```
{
  "database": "memory"
}
```

This database is only provided for testing purposes and shouldn't be used,
ever.

## License

```
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
```
