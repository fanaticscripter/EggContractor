# EggContractor

[![Build status](https://github.com/fanaticscripter/EggContractor/workflows/build/badge.svg)](https://github.com/fanaticscripter/EggContractor/actions)
[![Docker pulls](https://img.shields.io/docker/pulls/fanaticscripter/eggcontractor)](https://hub.docker.com/r/fanaticscripter/eggcontractor)
[![Gallery](https://shields.io/badge/-gallery-blueviolet)](https://github.com/fanaticscripter/EggContractor/wiki/Gallery)

EggContractor is a self-hosted contract monitoring web app + CLI client for [Egg, Inc.](https://en.wikipedia.org/wiki/Egg,_Inc.). It allows you to easily monitor all your contract progress, as well as peeking into prospective coops you may want to join.

Note that reverse engineered API protobufs are independently available at [api/egginc.proto](api/egginc.proto), so you may also find this repo useful for building your own client.

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [Demo](#demo)
- [Comparison to egginc.mioi.io](#comparison-to-eggincmioiio)
- [Installation & deployment](#installation--deployment)
  - [Nginx reverse proxying](#nginx-reverse-proxying)
- [CLI](#cli)
- [Known issues](#known-issues)
- [License](#license)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Demo

Below are some semi-interactive demo pages (pre-rendered, with in-page interactivity, but inter-page interactivity is likely broken). Note that the demo pages were captured from the MVP version of EggContractor, hence somewhat outdated; since then a number of major features have been implemented; for instance, offline timer for each coop member, and projection of eggs laid based on that.

- Home page: [/](https://egg-contractor-static-demo.netlify.app/)
- Home page with time travel (state at a certain timestamp in the past): [/?by=1608802500](https://egg-contractor-static-demo.netlify.app/?by=1608802500)
- Peeking a coop: [/peek/space-bots-2020/XXXXXXXXX/](https://egg-contractor-static-demo.netlify.app/peek/space-bots-2020/XXXXXXXXX/)
- List of recently peeked coops: [/peeked/](https://egg-contractor-static-demo.netlify.app/peeked/)

You can also find screenshots in the [gallery](https://github.com/fanaticscripter/EggContractor/wiki/Gallery).

## Comparison to egginc.mioi.io

Advantages:

- **No rate limiting**, obviously.

- **All the info about all your active contracts (solos & coops) is right on the home page**. Even the coop "peeker" widget is right there. **No clicking around**.

- **Sort coop members by EB, egg laying rate, etc.** Very handy.

- View a list of coops you recently "peeked". [Demo page](https://egg-contractor-static-demo.netlify.app/peeked/).

- **Stats are routinely retrieved in the background** (frequency easily configurable) and **stored in a database.** So you can project the actual number of eggs laid by taking into account how long each coop member has been offline (implemented), travel back in time to view your contract statuses in the past (implemented), or plot every player's progress.

Disadvantages:

- Frontend fanciness in general. I didn't bother to invest time into the frontend, so no pretty little pictures, no progress bars, and no dark theme (kind of a shame, not hard to add though). Hopefully you can still easily pick out whatever info you need from my UI.

- No home farm info. I don't need a separate web app to learn about my home farm, so not much of a disadvantage actually.

- No "Contract Calculator". I never used that feature so not sure how useful.

## Installation & deployment

docker-compose is the recommended method of deployment. Sorry k8s fans.

`docker-compose.yml`:

```yaml
version: "3"
services:
  app:
    image: fanaticscripter/eggcontractor:latest
    container_name: EggContractor
    restart: always
    environment:
      # TZ should be set to the local timezone in .env; e.g. TZ=America/New_York
      - TZ=${TZ}
    ports:
      # Use 0.0.0.0 only if you want to access the web app on the local network
      # directly, without a reverse proxy layer. You may change the host port
      # (the first port number) to another value.
      - "0.0.0.0:8080:8080"
    volumes:
      # config.toml should set database.path to /data/data.db
      - ./config.toml:/config.toml
      - ./data:/data
    labels:
      ofelia.enabled: "true"
      # Scheduling refreshes.
      #
      # Schedule syntax is documented at
      # - https://github.com/mcuadros/ofelia
      # - https://pkg.go.dev/github.com/robfig/cron@v1.2.0
      # Note that if you use cron syntax, the syntax has been extended to add
      # a second field at the beginning, so every minute would be "0 * * * * *"
      # instead of "* * * * *", and so on.
      ofelia.job-exec.refresh.schedule: "@every 2m"
      ofelia.job-exec.refresh.user: 0
      ofelia.job-exec.refresh.command: "/EggContractor refresh"
      ofelia.job-exec.refresh.save-folder: /logs
      ofelia.job-exec.refresh.save-only-on-error: true
      # Scheduling daily database backups.
      ofelia.job-exec.db-backup.schedule: "0 0 0 * * *"
      ofelia.job-exec.db-backup.user: 0
      ofelia.job-exec.db-backup.command: "/EggContractor backup"
      ofelia.job-exec.db-backup.save-folder: /logs
      ofelia.job-exec.db-backup.save-only-on-error: true

  ofelia:
    image: mcuadros/ofelia:latest
    container_name: EggContractor_sched
    restart: always
    depends_on:
      - app
    command: daemon --docker
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
      - ./logs:/logs
```

`config.toml`:

```toml
[player]
# player.id, required.
#
# Your unique player ID. To view your player ID, go to Main Menu -> Settings ->
# Privacy & Data, and the ID should be in the bottom left corner. On iOS (at
# least when signed in via Game Center) this would be of the format G:1234567890.
# Copy the string verbatim for this field.
id = "G:1234567890"

[database]
path = "/data/data.db"
```

`.env`:

```sh
# Set your local timezone here.
TZ=America/New_York
```

With these files in place,

```console
$ docker-compose up
```

### Nginx reverse proxying

In case you need help putting nginx in front for SSL termination and stuff, here's my nginx config:

```nginx
server {
    server_name egg.my.domain;
    root /var/www/html;

    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    ssl_certificate /etc/letsencrypt/live/my.domain/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/my.domain/privkey.pem;
    include /etc/letsencrypt/options-ssl-nginx.conf;
    ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem;

    access_log /var/log/nginx/egg.access.log;
    error_log /var/log/nginx/egg.error.log;

    add_header Strict-Transport-Security "max-age=31536000" always;

    location / {
        proxy_pass         http://127.0.0.1:8080;

        proxy_http_version 1.1;
        proxy_buffering    off;
        proxy_set_header   Host                 $host;
        proxy_set_header   Origin               http://$host;
        proxy_set_header   X-Real-IP            $remote_addr;
        proxy_set_header   X-Forwarded-For      $proxy_add_x_forwarded_for;
        proxy_set_header   X-Forwarded-Proto    $scheme;
    }
}

server {
    server_name egg.my.domain;

    listen 80;
    listen [::]:80;

    return 301 https://egg.my.domain$request_uri;
}
```

## CLI

```console
$ ./EggContractor help
Usage:
  EggContractor [command]

Available Commands:
  config          Print current configurations
  config-template Print a config file template
  contracts       Print a list of current and past contracts
  events          Print current and past events
  help            Help about any command
  peek            Peek at a coop
  peeked          Print list of recently peeked coops
  refresh         Refresh game state and print statuses of active solo contracts & coops
  serve           Run web server
  status          Print statuses of active solo contracts & coops from last refresh
  units           Print a table of units (order of magnitudes)

Flags:
      --config string    config file, could also be set through env var EGGCONTRACTOR_CONFIG_FILE (default ~/.config/EggContractor/config.toml)
      --debug            enable debug logging
  -h, --help             help for EggContractor
  -s, --sort criterion   sort coop members by one of the following criteria: 'eggs_laid' (aliases: 'contribution', 'total', 'laid'), 'laying_rate' (alias: 'rate'), or 'earning_bonus' (alias: 'eb') (default eggs_laid)
  -v, --verbose          enable verbose logging

Use "EggContractor [command] --help" for more information about a command.
```

Use `help` on individual subcommands to learn more about them.

If you're running the docker-compose setup, you need to use

    docker exec EggContractor /EggContractor [command]

or

    docker-compose exec app /EggContractor [command]

to use the CLI. You may want to set up a shell alias.

## Known issues

- When running the Docker image via Docker for Mac, with a mounted data directory, everything would seem fine and dandy until one accesses the SQLite database from the host system (e.g. do a SELECT on it), at which point all subsequent attempts to open the database from within the container would fail with "unable to open database file" (the error comes from sqlite itself, so doesn't matter if one uses the `sqlite3` CLI or the golang driver; the `go-sqlite3` driver does however misleadingly add a file/directory does not exist error message, which is irrelevant since the file can be open(2)'ed alright).

  This might be an issue in SQLite / Docker for Mac filesystem driver interactions. I've yet to isolate it.

  Workaround: run the Docker image on a Linux host. Probably already doing that outside of development anyway.

## License

The MIT license.
