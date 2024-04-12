# RSS Goes Social

RSS Goes Social is an application that adds RSS feeds to the Fediverse through a Mastodon-compatible API
([Mastodon](https://joinmastodon.org/), [GoToSocial](https://gotosocial.org), 
[Firefish](https://joinfirefish.org/)...).

## Usage

I recommend using Docker, but you can also build this application yourself with Go 1.22.

The official Docker image is [anhgelus/rss-goes-social](https://hub.docker.com/r/anhgelus/rss-goes-social).

1. Download the `docker-compose.yml` at the root of the repository.
It contains the application and Redis.
2. Then start it with `docker compose up -d`. It will create the config file `config/config.toml`
3. Customize the config file and restart.
4. Enjoy !

### Config file

The default config file is
```toml
version = '1'
fetch_every_X_minutes = 5

[redis]
host = 'localhost'
port = 6379
password = ''

[[feed]]
rss_feed_url = 'https://blog.example.org/rss'
server_url = 'https://gts.example.org'
token = 'account_token'
enabled = false
language = 'language of the feed (e.g. en, fr, de)'
```

Do not modify `version`.

`fetch_every_X_minutes` is the time (in minute) between two fetches.

`[redis].host` is the Redis' host. If you are using the `docker-compose.yml` provided, it's `redis`.

`[redis].port` is the port of Redis. If you are using the `docker-compose.yml` provided, it's 6379.

`[redis].passowrd` is the Redis' password. If you are using the `docker-compose.yml` provided, there is no password.

`[[feed]].rss_feed_url` is the url of the RSS feed.

`[[feed]].server_url` is the url of the Mastodon-like API.

`[[feed]].token` is the bot's token.

`[[feed]].enabled` is set to true when the feed is enabled.

`[[feed]].language` is the language of the feed (ISO 639)

To add a new feed, just copy the `[[feed]]` part and edit every variable. e.g.
```toml
version = '1'
fetch_every_X_minutes = 5

[redis]
host = 'redis'
port = 6379
password = ''

[[feed]]
rss_feed_url = 'https://blog.example.org/rss'
server_url = 'https://gts.example.org'
token = 'account_token'
enabled = false
language = 'en'

[[feed]]
rss_feed_url = 'https://blog2.example.org/rss'
server_url = 'https://mastodon.example.org'
token = 'account_token2'
enabled = true
language = 'fr'
```

### CLI

- `rss-goes-social` and `rss-goes-social` show the help
- `rss-goes-social help {command}` shows the help for the command
- `rss-goes-social run` runs the application
- `rss-goes-social setup {url}` setup the Mastodon application for the given server (url could be `https://mastodon.social`)

#### Setup

`setup` can be used without any flag or with two flags.
Without any flags, you will create a new application, and you will get a new `client_id` and a new `client_secret`.
With the two flags, you will use an existing application.

If you never registered RSS Goes Social on `https://example.org`, you must use `rss-goes-social https//example.org`.

If you already registered RSS Goes Social on `https://example.org`, you must use 
`rss-goes-social setup https://example.org -id client_id -secret client_secret` where `client_id` is the client's id of
the previous registered application and `client_secret` is the client's secret of the previous registered application.

During the setup, RSS Goes Social will ask you to log in to an account and to copy a token given by the instance.
You must log in to the account for which you wish to obtain the token.
After the login, you have to click on the "Allow" button to allow RSS Goes Social to use this account.
Then, the instance will give you the token. 
You have to give this token to RSS Goes Social to finish the process and to finally get the token to put in the config.

This can also be done graphically if the instance supports it.
For example, Mastodon supports it but GoToSocial does not.

## Technologies

- Go 1.22
- [mmcdole/gofeed](https://github.com/mmcdole/gofeed)
- [pelletier/go-toml/v2](https://github.com/pelletier/go-toml)
- [redis/go-redis/v9](https://github.com/redis/go-redis)
- and their dependencies
