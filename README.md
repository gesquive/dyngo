# digitalocean-ddns
[![Travis CI](https://img.shields.io/travis/gesquive/digitalocean-ddns/master.svg?style=flat-square)](https://travis-ci.org/gesquive/digitalocean-ddns)
[![Software License](https://img.shields.io/badge/License-MIT-orange.svg?style=flat-square)](https://github.com/gesquive/digitalocean-ddns/blob/master/LICENSE)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/gesquive/digitalocean-ddns)

Sync a DigitalOcean Domain Record entry with your public IP.

### Why?
I created this because the domain I own was being managed on the DigitalOcean nameservers and I didn't want to pay for a DDNS service for another domain. Using this app, I can host my website at `mydomain.com` but also have a subdomain of my choosing (ie. `dev.mydomain.com`) point to my dev network hosted elsewhere behind a dynamic IP.

## Installing

### Compile
This project has only been tested with go1.11+. To compile just run `go get -u github.com/gesquive/digitalocean-ddns` and the executable should be built for you automatically in your `$GOPATH`.

Optionally you can run `make install` to build and copy the executable to `/usr/local/bin/` with correct permissions.

### Download
You could also download the latest release for your platform from [github](https://github.com/gesquive/digitalocean-ddns/releases).

Once you have an executable, make sure to copy it somewhere on your path like `/usr/local/bin` or `C:/Program Files/`.
If on a \*nix/mac system, make sure to run `chmod +x /path/to/digitalocean-ddns`.

## Configuration

Before configuring, make sure that the domain exists in your DigitalOcean account. DigitalOcean provides excellent [documentation](https://www.digitalocean.com/docs/networking/dns/how-to/add-domains/) on this subject.
Also, when generating your DigitalOcean [personal access token](https://www.digitalocean.com/docs/api/create-personal-access-token/), make sure the token has read/write permissions.


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
 - `~/.config/digitalocean-ddns/config.yml`
 - `/etc/digitalocean-ddns/config.yml`

Copy `config.example.yml` to one of these locations and populate the values with your own. Since the config contains a writable API token, make sure to set permissions on the config file appropriately so others cannot read it. A good suggestion is `chmod 600 /path/to/config.yml`.

If you are planning to run this app as a service/cronjob, it is recommended that you place the config in `/etc/digitalocean-ddns/config.yml`. Otherwise, if running from the command line, place the config in `~/.config/digitalocean-ddns/config.yml` and make sure to set `run_once: true`.

### Environment Variables
Optionally, instead of using a config file you can specify config entries as environment variables. Use the prefix `DODDNS_` in front of the uppercased variable name. For example, the config variable `sync-interval` would be the environment variable `DODDNS_SYNC_INTERVAL`.

## Usage

```console
A service application that watches your external IP for changes and updates a DigitalOcean domain record when a change is detected

Usage:
  digitalocean-ddns [flags]

Flags:
      --config string          Path to a specific config file (default "./config.yaml")
  -d, --domain string          The DigitalOcean domain record to update
  -h, --help                   help for digitalocean-ddns
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
To run as a cronjob on an Ubuntu system create a cronjob entry under the user the app is run with. If running as root, you can copy `services/digitalocean-ddns.cron` to `/etc/cron.d/digitalocean-ddns` or copy the following into you preferred crontab:
```shell
  0  *  *  *  * /usr/local/bin/digitalocean-ddns --run-once
```

Add any flags/env vars needed to make sure the job runs as intended. If not using arguments, then make sure the config file path is specified with a flag or can be found at one of the expected locations.

### Service
By default, the process is setup to run as a service. Feel free to use upstart, systemd, runit or any other service manager to run the `digitalocean-ddns` executable.

Example systemd & upstart scripts can be found in the `services` directory.

## Documentation

This documentation can be found at github.com/gesquive/digitalocean-ddns

## License

This package is made available under an MIT-style license. See LICENSE.

## Contributing

PRs are always welcome!
