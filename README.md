# donutdns

<img align="right" width="240" height="244" src="https://i.imgur.com/1cEeZ3L.png">

Block online ads by intercepting DNS queries

[![License](https://img.shields.io/badge/License-BSD%203--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause)
[![CI Tests](https://github.com/shoenig/donutdns/actions/workflows/tests.yaml/badge.svg)](https://github.com/shoenig/donutdns/actions/workflows/tests.yaml)

## Project Overview

The `github.com/shoenig/donutdns` module provides a [CoreDNS](https://coredns.io) plugin
as well as a standalone executable DNS server that can be used to block DNS queries to
domains used by online advertisers, trackers, scammers, and crypto miners. The project
is meant to be a simpler alternative to the venerable [Pi-Hole](https://pi-hole.net). In
particular, `donutdns` is easy to run as a **non-root** [Docker container](https://hub.docker.com/r/shoenig/donutdns)
with little to no configuration.

#### sample logs

```
[INFO] plugin/donutdns: BLOCK query (A) for www.google-analytics.com.
[INFO] plugin/donutdns: BLOCK query (A) for www-google-analytics.l.google.com.
[INFO] plugin/donutdns: BLOCK query (AAAA) for stats.wp.com.
[INFO] plugin/donutdns: BLOCK query (A) for www.googletagservices.com.
[INFO] plugin/donutdns: BLOCK query (A) for tpc.googlesyndication.com.
[INFO] plugin/donutdns: BLOCK query (A) for c.amazon-adsystem.com.
[INFO] plugin/donutdns: BLOCK query (AAAA) for static.ads-twitter.com.
```

## Domain Block/Allow Lists

The default set of blocked domains are retrieved from the source lists in [sources.json](sources/statics/sources.json).
These lists are compiled and maintained by volunteers; see their respective headers
for more information about terms of use and other metadata. Thank you to those who
contribute to these domain block lists.

The blocking of the default set of domains can be disabled by setting `DONUT_DNS_NO_DEFAULT=1`.

Additional domains can be blocked by `donutdns` by setting any of the `DONUT_DNS_BLOCK`,
`DONUT_DNS_BLOCK_FILE`, `DONUT_DNS_BLOCK_DIR` environment variables.

Likewise, domains can be explicitly allowed by setting the `DONUT_DNS_ALLOW`,
`DONUT_DNS_ALLOW_FILE`, `DONUT_DNS_ALLOW_DIR` environment variables. The allow lists 
take precedense over the block lists.

For nasty companies like Facebook with dynamic subdomains, `donutdns` supports blocking
domains by suffix matching. By setting any of the `DONUT_DNS_SUFFIX`, `DONUT_DNS_SUFFIX_FILE`,
`DONUT_DNS_SUFFIX_DIR` any query matching the given suffix(es) will be blocked.

## Getting Started

`donutdns` can be used as a CoreDNS Plugin or standalone DNS Server.

#### Install 

The `donutdns` standalone DNS Server is written in Go.

Pre-compiled binaries are available for download from the [Releases](https://github.com/shoenig/donutdns/releases) page.

Docker images are available from [Docker Hub](https://hub.docker.com/r/shoenig/donutdns).

With the Go toolchain, `donutdns` standalone can be compiled and installed in one step:

```
go install github.com/shoenig/donutdns@latest
```

#### DNS Server

The `donutdns` executable uses environment variables for configuration.

| Environment Variable | Description | Default |
| -------------------- | ----------- | ------- |
| `DONUT_DNS_PORT` | The port to listen on | `5301` |
| `DONUT_DNS_NO_DEBUG` | Disable CoreDNS debug logging | unset |
| `DONUT_DNS_NO_LOG` | Disable CoreDNS logging | unset |
| `DONUT_DNS_ALLOW` | Comma separated list of domains to NOT block | unset |
| `DONUT_DNS_ALLOW_FILE` | File with list of domains to NOT block | unset |
| `DONUT_DNS_ALLOW_DIR` | Directory with one or more files of list of domains to NOT block | unset |
| `DONUT_DNS_BLOCK` | Comma separated list of domains to block | unset |
| `DONUT_DNS_BLOCK_FILE` | File with list of domains to block | unset |
| `DONUT_DNS_BLOCK_DIR` | Directory with one or more files of list of domains to block | unset |
| `DONUT_DNS_SUFFIX` | Comma separated list of domains to block by suffix | unset |
| `DONUT_DNS_SUFFIX_FILE` | File with list of domains to block by suffix | unset |
| `DONUT_DNS_SUFFIX_DIR` | Directory with one or more files of list of domains to block by suffix | unset |
| `DONUT_DNS_NO_DEFAULTS` | Disable blocking of default domain block lists | unset |
| `DONUT_DNS_UPSTREAM_1` | Fallback DNS Server for non-blocked queries |`1.1.1.1` |
| `DONUT_DNS_UPSTREAM_2` | Fallback DNS Server for non-blocked queries | `1.0.0.1` |
| `DONUT_DNS_UPSTREAM_NAME` | Fallback DNS Server TLS name | `cloudflare-dns.com` |

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

## Subcommands

#### check

Use the `check` command to simulate whether a DNS query would be blocked or allowed.

Usage: `donutdns check [-quiet] [-defaults] <domain>`

`-quiet` will suppress verbose debug logging output

`-defaults` will activate the built-in block lists (which is slow)

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

#### as a systemd unit

The [donutdns.service](donutdns.service) file provides an example Systemd Service Unit file for running
donutdns via systemd.

```
# A minimal unit file, see donutdns.service for more.

[Unit]
Description=Block ads, trackers, and malicioius sites using DonutDNS.

[Service]
ExecStart=/opt/bin/donutdns
Environment=DONUT_DNS_PORT=53

[Install]
WantedBy=multi-user.target
```

Typically this file would be created at `/etc/systemd/system/donutdns.service`.

Configure systemd to run the new service.

```shell
sudo systemctl daemon-reload           # update systemd configurations
sudo systemctl enable donutdns.service # enable donutdns service in systemd
sudo systemctl start donutdns          # start donutdns service in systemd
sudo systemctl status donutdns         # inspect status of donutdns service in systemd
```

#### as a docker container

`donutdns` is available from [Docker Hub](https://hub.docker.com/repository/docker/shoenig/donutdns/general)

This will run the `donutdns` Docker container as the `nobody` user, mapping traffic from port 53. 
```
docker run --rm -p 53:5301 -u nobody shoenig/donutdns:v0.2.0
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
        image = "shoenig/donutdns:v0.1.2"
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

## Troubleshooting

Certain systems (looking at you RHEL/CentOS) make running a useable DNS server particularly
difficult. On my homelab CentOS 9 system I had to disable ipv6 at the kernel level, disable
SELinux, and disable firewalld. You may need to do something similar (ideally updating rules
rather than disabling things) on your system.

## Contributing

The `github.com/shoenig/donutdns` module is always improving with new features
and bug fixes. For contributing such bug fixes and new features please file an issue.

## License

The `github.com/shoenig/donutdns` module is open source under the [BSD-3-Clause](LICENSE) license.
