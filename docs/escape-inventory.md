---
title: "Escape Inventory"
slug: escape-inventory 
type: "docs"
toc: true
wip: false
---

* [Usage](#usage)
* [Configuration Options](#configuration-options)
* [Storage Backends](#storage-backends)
* [Databases](#databases)

# Usage

```
escape-inventory [CONFIG_FILE]
```

The Escape Inventory can be configured using a simple JSON or YAML file (default
`/etc/escape-inventory/config.json`), and/or environment variables. If the
provided configuration file does not exist the program falls back to the
following default configuration: 

```
{
  "port": "7770",
  "database": "ql",
  "database_settings": {
    "path": "/var/lib/escape/inventory.db"
  },
  "storage_backend": "local",
  "storage_settings": {
    "path": "/var/lib/escape/releases"
  }
}
```

This configures the Inventory with an <a href='https://github.com/cznic/ql' target='_blank'>ql</a>
database and a local file system storage backend, which is the recommended way
to run the Inventory for development purposes. For production use we currently
support a Google Cloud Storage backend and a Postgresql database.


# Configuration Options

JSON/YAML field | Environment Variable | Default |Description
----------------|----------------------|---------|-----------
|`port`|`PORT`|`7770`|The port to listen on. 
|`database`|`DATABASE`|`ql`|The database to use (one of: `ql`, `postgres`).
|`database_settings . path`|`DATABASE_SETTINGS_PATH`|`/var/lib/escape/inventory.db`|The path to the database. Only relevant for the `ql` backend.
|`database_settings . postgres_url`|`DATABASE_SETTINGS_POSTGRES_URL`||The URL to a postgres database. For more information see the documentation for the postgres backend.
|`storage_backend`|`STORAGE_BACKEND`|`local`|The storage backend to use (one of: `local`, `gcs`).
|`storage_settings . path`|`STORAGE_SETTINGS_PATH`|`/var/lib/escape/releases/`|Where packages will be stored. Only relevant for the the `local` storage backend.
|`storage_settings . bucket`|`STORAGE_SETTINGS_BUCKET`||The bucket where packages will be stored. Only relevant for the `gcs` storage backend. 
|`storage_settings . credentials`|`STORAGE_SETTINGS_CREDENTIALS`||This path points to the credentials for the GCS bucket. For more information see the documentation for the GCS storage backend. 


# Storage Backends

The storage backends are used to store and retrieve uploaded packages.

## Local File Storage

The default way to use the Inventory is to have it store uploaded packages on
the local file system.  The destination path can be configured using the
`database_settings.path` variable:

```
{
  "storage_backend": "local",
  "storage_settings": {
    "path": "/var/lib/escape/releases"
  }
}
```

Using environment variables:

```
STORAGE_BACKEND=local
STORAGE_SETTINGS_PATH=/var/lib/escape/releases
```


## Google Cloud Storage

Stores uploaded packages in Google Cloud Storage.  The `credentials` variable
is optional, but should point to an existing service account json file if
provided. The service account should have the "Storage -> Storage Admin" role.
If no credentials are provided the Inventory is assumed to be running in GCP 
under the `storage-rw` scope.

```json
{
  "storage_backend": "gcs",
  "storage_settings": {
    "bucket": "my-bucket",
    "credentials": "/my/secret/service/credentials.json",
  }
}
```

Or using environment variables:

```bash
export STORAGE_BACKEND=gcs
export STORAGE_SETTINGS_BUCKET=my-bucket
export STORAGE_SETTINGS_CREDENTIALS=/my/secret/service/credentials.json
```

# Databases

## QL

This is the default database, which will work out of the box (provided the path
is accessible) and will store the data on the local disk. The path is
configurable using the `database_settings.path` variable:

```json
{
  "database": "ql",
  "database_settings": {
    "path": "/var/lib/escape/inventory.db"
  }
}
```

Or using environment variables:

```bash
export DATABASE=ql
export DATABASE_SETTINGS_PATH=/var/lib/escape/inventory.db
```


## Postgresql

Postgresql can be configured using the `postgres_url` variable. Please see 
<a href='https://godoc.org/github.com/lib/pq' target='_blank'>pq</a> for the full connection string parameters.

```json
{
  "database": "postgres",
  "database_settings": {
    "postgres_url": "postgres://user:pass@localhost/database?sslmode=disable"
  }
}
```

Or using environment variables:

```bash
export DATABASE=postgres
export DATABASE_SETTINGS_POSTGRES_URL=postgres://user:pass@localhost/database?sslmode=disable
```
