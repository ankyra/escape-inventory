---
title: "Configuring an Inventory"
slug: quickstart-configure-inventory 
type: "docs"
toc: true

back: /docs/installation/
backLabel: Installation
next: /docs/quickstart-building-a-package
nextLabel: Building a Package
contributeLink: https://example.com/
---

Before we can truly kick off we need to configure an Escape Inventory. The
Escape Inventory is used by Escape to store, retrieve and answer questions
about packages.  By default Escape is configured to use the central Escape
Inventory hosted by Ankyra, but at the moment of writing this does not provide
write access to members of the public. (Watch this space)

Setting up your own Inventory is easy though. As with Escape we can either use 
the pre-built binaries or run a Docker container. 

### Linux

```
curl -O https://storage.googleapis.com/escape-releases-eu/escape-inventory/0.12.9/escape-inventory-v0.12.9-linux-amd64.tgz
tar -xvzf escape-v0.12.9-linux-amd64.tgz
sudo mv escape-inventory /usr/bin/escape-inventory
escape-inventory
```

### MacOS

```
curl -O https://storage.googleapis.com/escape-releases-eu/escape-inventory/0.12.9/escape-inventory-v0.12.9-darwin-amd64.tgz
tar -xvzf escape-v0.12.9-darwin-amd64.tgz
sudo mv escape-inventory /usr/bin/escape-inventory
escape-inventory
```

### Docker

```
docker run -P -it ankyra/escape-inventory:v0.12.9
```

## Configuring the Inventory

The Inventory can be configured to store packages on Google Cloud Storage
instead of disk and use a Postgres database instead of `ql`, but right now
we'll keep it simple. For more information see [Escape Inventory
Usage](/docs/escape-inventory/).

## Configuring Escape

We have our own instance of Escape Inventory running now, but we need to tell
Escape to use it, because right now it will be configured to go to the central
repository:

```
escape config profile
```

So let's login:

```
escape login --url http://localhost:7770/
```

Oh bim.

