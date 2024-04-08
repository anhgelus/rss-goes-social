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

## Technologies

- Go 1.22
- [mmcdole/gofeed](https://github.com/mmcdole/gofeed)
- [pelletier/go-toml/v2](https://github.com/pelletier/go-toml)
- [redis/go-redis/v9](https://github.com/redis/go-redis)
- and their dependencies
