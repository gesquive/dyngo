# digitalocean-ddns

Update a WebFaction DNS entry for a domain.

## Getting started

This project requires Go to be installed. On OS X with Homebrew you can just run `brew install go`.

Running it then should be as simple as:

```console
$ make
$ ./bin/digitalocean-ddns
```

Add your long-running agent logic to `command/agent/command.go`, and any status or action commands you need to `commands.go`.


This app needs a DigitalOcean api key, read more here: https://www.digitalocean.com/community/tutorials/how-to-use-the-digitalocean-api-v2#how-to-generate-a-personal-access-token

### Testing

``make test``

## License

_Fill me in._

## Contributing

See `CONTRIBUTING.md` for more details.
