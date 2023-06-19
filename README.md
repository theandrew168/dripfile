# DripFile
File transfers made easy

## Setup
This project depends on the [Go programming language](https://golang.org/dl/) and [NodeJS JavaScript environment](https://nodejs.org/en).
I like to use a [POSIX-compatible Makefile](https://pubs.opengroup.org/onlinepubs/9699919799.2018edition/utilities/make.html) to facilitate the various project operations but traditional [Go commands](https://pkg.go.dev/cmd/go) and [NPM scripts](https://docs.npmjs.com/cli/v9/commands/npm-run-script) will work just as well.

## Building
To build the application into a standalone binary, run:
```
make
```

## Local Development
### Services
This project depends on various services.
To develop locally, you'll need to run these services locally somehow or another.
I find [Docker](https://www.docker.com/) to be a nice tool for this but you can do whatever works best.
* [PostgreSQL](https://www.postgresql.org/) - for persistent storage and task queue

The following command starts the necessary containers:
```
docker compose up -d
```

These containers can be stopped via:
```
docker compose down
```

### Running
To run the backend API and watch for changes on the frontend simultaneously:
```
make -j2 run
```

## Features
* Schedule-based transfers
* Polling-based transfers
* Automatic retries
* Notify via email, text, etc
* Works with S3, FTP, FTPS, SFTP, etc
* Choose from common schedules or build your own

## User Stories
* Guest registers and is logged in
* Guest logs in
* User verifies their email address (required to log in?)
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
