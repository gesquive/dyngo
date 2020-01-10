# dyngo
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/gesquive/dyngo)
[![Software License](https://img.shields.io/badge/License-MIT-orange.svg?style=flat-square)](https://github.com/gesquive/dyngo/blob/master/LICENSE)
[![Build Status](https://img.shields.io/gitlab/pipeline/gesquive/dyngo?style=flat-square)](https://gitlab.com/gesquive/dyngo/pipelines)
[![Coverage Report](https://gitlab.com/gesquive/dyngo/badges/master/coverage.svg?style=flat-square)](https://gesquive.gitlab.io/dyngo/coverage.html)
[![Docker Pulls](https://img.shields.io/docker/pulls/gesquive/dyngo?style=flat-square)](https://hub.docker.com/r/gesquive/dyngo)

Sync a DigitalOcean/Cloudflare/Custom DNS entry with your public IP.

### Why?
I created this because the domain I own was being managed in some cloud nameservers and I didn't want to pay for a DDNS service for another domain. Using this app, I can host my website at `mydomain.com` but also have a subdomain of my choosing (ie. `dev.mydomain.com`) point to my dev network hosted elsewhere behind a dynamic IP.

## Installing

### Compile
This project has only been tested with go1.11+. To compile just run `go get -u github.com/gesquive/dyngo` and the executable should be built for you automatically in your `$GOPATH`. This project uses go mods, so you might need to set `GO111MODULE=on` in order for `go get` to complete properly.

Optionally you can run `make install` to build and copy the executable to `/usr/local/bin/` with correct permissions.

### Download
You could download the latest release for your platform from [github](https://github.com/gesquive/dyngo/releases).

Once you have an executable, make sure to copy it somewhere on your path like `/usr/local/bin` or `C:/Program Files/`.
If on a \*nix/mac system, make sure to run `chmod +x /path/to/dyngo`.

### Docker
You can also run dyngo from the provided [Docker image](https://hub.docker.com/r/gesquive/dyngo) with the sample configuration file:

```shell
mkdir -p dyngo && cp pkg/config.example.yml dyngo/config.yml
docker run -d -v $PWD/dyngo:/config dyngo:latest
```

While dyngo will technically run with the sample config, it won't have any domain to sync too. You will need to edit the config to add dns providers.

For more details read the [Docker image documentation](https://hub.docker.com/r/gesquive/dyngo).

### Precedence Order
The application looks for variables in the following order:
 - command line flag
 - environment variable
 - config file variable
 - default

So any variable specified on the command line would override values set in the environment or config file.

### Config File
The application looks for a configuration file at the following locations in order:
 - `./config.yml`
 - `~/.config/dyngo/config.yml`
 - `/etc/dyngo/config.yml`

Copy `config.example.yml` to one of these locations and populate the values with your own. Since the config contains a writable API token, make sure to set permissions on the config file appropriately so others cannot read it. A good suggestion is `chmod 600 /path/to/config.yml`.

If you are planning to run this app as a service/cronjob, it is recommended that you place the config in `/etc/dyngo/config.yml`. Otherwise, if running from the command line, place the config in `~/.config/dyngo/config.yml` and make sure to set `run_once: true`.

### Environment Variables
Optionally, instead of using a config file you can specify config entries as environment variables. Use the prefix `DYNGO_` in front of the uppercased variable name. For example, the config variable `sync-interval` would be the environment variable `DYNGO_SYNC_INTERVAL`.

## Usage

```console
A service application that watches your external IP for changes and updates a DigitalOcean domain record when a change is detected

Usage:
  dyngo [flags]

Flags:
      --config string          Path to a specific config file (default "./config.yaml")
  -d, --domain string          The DigitalOcean domain record to update
  -h, --help                   help for dyngo
      --log-file string        Path to log file (default "-")
  -o, --run-once               Only run once and exit
  -i, --sync-interval string   The duration between DNS updates (default "60m")
  -t, --token string           The DigitalOcean API token to authenticate with
      --version                Display the version number and exit
```

It is helpful to use the `--run-once` when first setting up to find any misconfigurations.

Optionally, a hidden debug flag is available in case you need additional output.
```console
Hidden Flags:
  -D, --debug                  Include debug statements in log output
```


### Cronjob
To run as a cronjob on an Ubuntu system create a cronjob entry under the user the app is run with. If running as root, you can copy `services/dyngo.cron` to `/etc/cron.d/dyngo` or copy the following into you preferred crontab:
```shell
  0  *  *  *  * /usr/local/bin/dyngo --run-once
```

Add any flags/env vars needed to make sure the job runs as intended. If not using arguments, then make sure the config file path is specified with a flag or can be found at one of the expected locations.

### Service
By default, the process is setup to run as a service. Feel free to use upstart, systemd, runit or any other service manager to run the `dyngo` executable.

Example systemd & upstart scripts can be found in the `services` directory.


## DNS Provider Configuration

Before configuring and running dyngo, make sure that the domain exists in your cloud account. Specifics can be found below.

### `digitalocean`
DigitalOcean DNS provides excellent [documentation](https://www.digitalocean.com/docs/networking/dns/how-to/add-domains/) on this adding domains to DNS.

 - `record`: The record to set the IP on (ie. `ddns.mydomain.com`)
 - `token`: The DigitalOcean [personal access token](https://www.digitalocean.com/docs/api/create-personal-access-token/). Token requires read/write permissions.


### `cloudflare`
You will need to [add your domain](https://support.cloudflare.com/hc/en-us/articles/201720164-Creating-a-Cloudflare-account-and-adding-a-website) to Cloudflare before you can [manage any records](https://support.cloudflare.com/hc/en-us/articles/360019093151-Managing-DNS-records-in-Cloudflare).
- `record`: The record to set the IP on (ie. `ddns.mydomain.com`)
- `token`: The [Cloudflare API Token](https://support.cloudflare.com/hc/en-us/articles/200167836-Managing-API-Tokens-and-Keys). API token requires at least the `Zone.Zone:Read, Zone.DNS:Edit` permissions.

### `custom`
If your provider is not found above, it is possible to run a custom script as well. The `custom` DNS provider supports the following config options:

- `record`: The record to set the IP on (ie. `ddns.mydomain.com`)
- `path`: The relative path to the script
- `args`: Arguments to pass to the script when executing

Scripts should expect the following arguments when being executed:
- Record - the record value from the config example: `custom.domain.com`
- Record Type - example: `A`, `AAAA`
- IP Address - example: `192.168.1.100`, `2001:cdba::3257:9652`
- args -  The args value from the config

For example, a config with the following value:

```yaml
dns_providers:
  name: custom
  # The domain record to pass to the script
  record: custom.domain.com
  # The relative path to the script
  path: path/to/custom_script.sh
  # The arguments to pass when executing the script
  args: "-D"
```
would possible result in the following call: 

`path/to/custom_script.sh custom.domain.com A 192.168.10.10 -D`


## Documentation

This documentation can be found at github.com/gesquive/dyngo

## License

This package is made available under an MIT-style license. See LICENSE.

## Contributing

PRs are always welcome!
