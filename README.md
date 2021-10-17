# donutdns

Block online ads by intercepting DNS queries

[![Go Report Card](https://goreportcard.com/badge/gophers.dev/cmds/donutdns)](https://goreportcard.com/report/gophers.dev/cmds/donutdns)
[![Build Status](https://app.travis-ci.com/shoenig/donutdns.svg)](https://app.travis-ci.com/github/shoenig/donutdns)
[![GoDoc](https://godoc.org/gophers.dev/cmds/donutdns?status.svg)](https://godoc.org/gophers.dev/cmds/donutdns)
![NetflixOSS Lifecycle](https://img.shields.io/osslifecycle/shoenig/donutdns.svg)
[![License](https://img.shields.io/badge/License-BSD%203--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause)

## Project Overview

The `gophers.dev/cmds/donutdns` module provides a [CoreDNS](https://coredns.io) plugin
as well as a standalone executable DNS server that can be used to block DNS queries to
domains used by online advertisers, trackers, scammers, and crypto miners. The project
is meant to be a simpler alternative to the venerable [Pi-Hole](https://pi-hole.net). In
particular, `donutdns` is easy to run as a **non-root** Docker container with little to
no configuration.

## Domain Block Lists

The default set of blocked domains are retrieved from the source lists in [sources.json](sources/statics/sources.json).
These lists are compiled and maintained by volunteers; see their respective headers
for more information about terms of use and other metadata. Thank you to those who
contribute to these domain block lists.

The blocking of the default set of domains can be disabled by setting `DONUT_DNS_NO_DEFAULT=1`.

Additional domains can be blocked by `donutdns` by setting the `DONUT_DNS_BLOCK` and/or
`DONUT_DNS_BLOCK_FILE` environment variables.

(!) Currently `donutdns` does not support wildcard subdomain blocking. Each subdomain
to be blocked will also need to be added. (e.g. `example.com` and `www.example.com`)

## Getting Started

`donutdns` can be used as a CoreDNS Plugin or standalone DNS Server.

#### DNS Server

The `donutdns` executable uses environment variables for configuration.

| Environment Variable | Description |
| -------------------- | ----------- |
| `DONUT_DNS_PORT` | The port to listen to (default `5301`) |
| `DONUT_DNS_NO_DEBUG` | Disable CoreDNS debug logging (default unset) |
| `DONUT_DNS_NO_LOG` | Disable CoreDNS logging (default unset) |
| `DONUT_DNS_ALLOW` | Comma separated list of domains to NOT block (default unset) |
| `DONUT_DNS_BLOCK` | Comma separated list of domains to block (default unset) |
| `DONUT_DNS_BLOCK_FILE` | File with list of domains to block (default unset) |
| `DONUT_DNS_NO_DEFAULTS` | Disable blocking of default domain block lists (default unset) |
| `DONUT_DNS_UPSTREAM_1` | Fallback DNS Server for non-blocked queries (default `1.1.1.1`) |
| `DONUT_DNS_UPSTREAM_2` | Fallback DNS Server for non-blocked queries (default `1.0.0.1`) |
| `DONUT_DNS_UPSTREAM_NAME` | Fallback DNS Server TLS name (default `cloudflare-dns.com`) |

#### CoreDNS Plugin

The `donutdns` CoreDNS plugin is configured using the `donutdns` block in a standard
[CoreConfig](https://coredns.io/manual/toc/#configuration) configuration file.

Minimal `donutdns` plugin configuration. `defaults` can be set to `true` or `false`
to enable or disable the use of default domain block lists.

```
donutdns {
  defaults true
}
```

This configuration uses `block_file` to explicitly block a set of domains listed
in a file on local disk.

```
donutdns {
  defaults false
  block_file /etc/blocked-domains.txt
}
```

This configuration uses `block` and `allow` to explicitly block and allow certain
domains.

```
donutdns {
  defaults true
  block facebook.com,www.facebook.com,m.facebook.com,fb.com
  allow example.com
}
```

When using `donutdns` as a CoreDNS plugin, the fallthrough behavior must be configured
as desired using one or more other plugins. To recreate the same recursive behavior
as the standalone executable, use the [`forward`](https://coredns.io/plugins/forward/) plugin.

```
forward . 1.1.1.1 1.0.0.1 {
  tls_servername cloudflare-dns.com
}
```

#### Custom block file

The file format for `block_file` or `DONUT_DNS_BLOCK_FILE` is simply a newline
delimited list of domains. Empty lines and lines beginning with `#` are always
ignored. All other lines are scanned with a regular expression to find the first
plausible domain name in the line. [social-media.list](hack/social-media.list)
contains an example file for blocking facebook, instagram, and whatsapp.

```
# An example block list
example.com
www.example.com
```

## Run

#### as an executable

With no configuration, `donutdns` will use the built-in domain block lists
by default.

```
$ donutdns
```

Use the environment variables described above to configure things.

```
$ DONUT_DNS_PORT=5533 DONUT_DNS_NO_DEBUG=1 donutdns
```

#### as a docker container

`donutdns` is available from [Docker Hub](https://hub.docker.com/repository/docker/shoenig/donutdns/general)

This will run the `donutdns` Docker container as the `nobody` user, mapping traffic from port 53. 
```
docker run --rm -p 53:5301 -u nobody shoenig/donutdns:v0.1.0
```

#### as a Nomad job

<details><summary>using docker driver</summary>
  
```hcl
job "donutdns" {
  datacenters = ["dc1"]

  group "donut" {
    network {
      mode = "bridge"
      port "dns" {
        static       = 53
        to           = 5301
        host_network = "public"
      }
    }

    task "dns" {
      driver = "docker"
      user   = "nobody"

      resources {
        cpu    = 120
        memory = 64
        disk   = 128
      }

      env {
        DONUT_DNS_NO_DEBUG   = 1
        DONUT_DNS_BLOCK_FILE = "/local/blocks.txt"
      }

      config {
        image = "shoenig/donutdns:v0.1.0"
      }

      template {
        destination = "local/blocks.txt"
        change_mode = "restart"
        perms       = "644"
        data        = <<EOH
# [example]
example.com
www.example.com
EOH
      }
    }
  }
}
```
  
</details>

## Build

The `donutdns` standalone DNS Server is written in Go. It can be compiled and
installed using the normal Go toolchain in one step.

```
go install gophers.dev/cmds/donutdns@latest
```

## Contributing

The `gophers.dev/cmds/donutdns` module is always improving with new features
and bug fixes. For contributing such bug fixes and new features please file an issue.

## License

The `gophers.dev/cmds/donutdns` module is open source under the [BSD-3-Clause](LICENSE) license.
