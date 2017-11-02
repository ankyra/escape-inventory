---
title: "Escape Inventory"
slug: escape-inventory 
type: "docs"
toc: true
---

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
  "database": "sqlite",
  "database_settings": {
    "path": "/var/lib/escape/inventory.db"
  },
  "storage_backend": "local",
  "storage_settings": {
    "path": "/var/lib/escape/releases"
  }
}
```

This configures the Inventory with an SQLite database and a local file system
storage backend.

# Environment Variables

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

# Storage Backends

The storage backends are used to store and retrieve uploaded packages.

## Local file storage

Stores uploaded packages on the local file system.
The `path` variable points to a directory in which the releases will be stored.

### JSON

```
{
  "storage_backend": "local",
  "storage_settings": {
    "path": "/var/lib/escape/releases"
  }
}
```

### Environment Variables

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

### JSON

```
{
  "storage_backend": "gcs",
  "storage_settings": {
    "bucket": "my-bucket",
    "credentials": "/my/secret/service/credentials.json",
  }
}
```

### Environment Variables

```
STORAGE_BACKEND=gcs
STORAGE_SETTINGS_BUCKET=my-bucket
STORAGE_SETTINGS_CREDENTIALS=/my/secret/service/credentials.json
```

# Databases

## SQLite

This is the default database, which will work out of the box (provided the path
is accessible). 

### JSON

```
{
  "database": "sqlite",
  "database_settings": {
    "path": "/var/lib/escape/inventory.db"
  }
}
```

### Environment Variables

```
DATABASE=sqlite
DATABASE_SETTINGS_PATH=/var/lib/escape/inventory.db
```


## Postgresql

Postgresql can be configured using the `postgres_url` variable. Please see 
https://godoc.org/github.com/lib/pq for the full connection string parameters.

### JSON

```
{
  "database": "postgres",
  "database_settings": {
    "postgres_url": "postgres://user:pass@localhost/database?sslmode=disable"
  }
}
```

### Environment Variables

```
DATABASE=postgres
DATABASE_SETTINGS_POSTGRES_URL=postgres://user:pass@localhost/database?sslmode=disable
```

## In Memory 

This database is only provided for unit testing purposes and shouldn't be used,
ever.

### JSON

```
{
  "database": "memory"
}
```

### Environment Variables

```
DATABASE=memory
```
