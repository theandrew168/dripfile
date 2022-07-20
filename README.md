# DripFile
Managed file transfer as a service

## Setup
This project depends on the [Go programming language](https://golang.org/dl/).
I like to use a [POSIX-compatible Makefile](https://pubs.opengroup.org/onlinepubs/9699919799.2018edition/utilities/make.html) to facilitate the various project operations but traditional [Go commands](https://pkg.go.dev/cmd/go) will work just as well.

## Building
To build the application into a standalone binary, run:
```
make
```

## Local Development
### Services
This project uses [PostgreSQL](https://www.postgresql.org/) for persistent storage and [Redis](https://redis.io/) for queuing background tasks.
To develop locally, you'll need to run these services locally somehow or another.
I find [Docker](https://www.docker.com/) to be a nice tool for this but you can do whatever works best.

The following command starts the necessary containers:
```
docker compose up -d
```

These containers can be stopped via:
```
docker compose down
```

### Running
To start the web server:
```
make web
```

To start the periodic task scheduler:
```
make scheduler
```

To start a background task worker:
```
make worker
```

To apply any pending database migrations:
```
make migrate
```

### Testing
Unit and integration tests can be ran after starting the aforementioned services:
```
make test
```

## Innovation Tokens
* [Asynq Task Queue](https://github.com/hibiken/asynq)
* [Tachyons CSS](https://tachyons.io/)

## Features
* Managed file transfer (MFT) as a service
* Schedule-based transfers
* Polling-based transfers
* Automatic retries
* Role-based access controls
* Detailed, customizable notifications
* Notify via email, text, etc
* Automatic compression and encryption
* Works with S3, FTP, FTPS, SFTP, etc
* Choose from common schedules or build your own

## User Stories
* Guest registers and is logged in
* Guest logs in
* User verifies their email address (required to log in?)
* User invites a Guest to their project (by email)
* User deletes another User from their project (owner or admin)
* User updates another User's role (owner or admin)
* User deletes Account (and all ref'd info)
* User CRUDs a Location (no delete if affected Transfers)
* User CRUDs a Transfer (history entries persist but won't link back)
* User CRUDs a Schedule (no delete if affected Transfers)
* User links / unlinks a Transfer and Schedule
* User runs a Transfer adhoc
* User reruns a completed Transfer adhoc
* User views Transfer history
* User cancels an in-progress Transfer
* User switches account auth method (email, github, google, etc)
