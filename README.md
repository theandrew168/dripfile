# dripfile
Managed file transfer as a service

## Setup
This project depends on the [Go programming language](https://golang.org/dl/).

## Building
To build the application into a standalone binary, run:
```bash
make
```

## Local Development
### Services
This project uses [PostgreSQL](https://www.postgresql.org/) for persistent storage, [MinIO](https://min.io/) for object storage, [Asynq](https://github.com/hibiken/asynq) for background jobs, and [Redis](https://redis.io/) for caching.
To develop locally, you'll need to run these services locally somehow or another.
I find [Docker](https://www.docker.com/) to be a nice tool for this but you can do whatever works best.

The following command starts the necessary containers:
```bash
docker compose up -d
```

These containers can be stopped via:
```bash
docker compose down
```

### Running
To apply any pending database migrations:
```bash
make migrate
```

To start the web server:
```bash
make web
```

To start the periodic task scheduler:
```bash
make scheduler
```

To start a background task worker:
```bash
make worker
```

### Testing
Unit and integration tests can be ran after starting the aforementioned services:
```bash
make test
```

## Innovation Tokens
* Asynq (task queue)

## Features
* Managed file transfer (MFT) as a service
* Schedule-based transfers
* Polling-based transfers
* Automatic retries
* Pay-as-you-go pricing ($0.10 per GB)
* Generous free tier (first 10GB free)
* Seamless integration with isolated networks
* Unlimited team members per project
* Role-based access controls
* Detailed, customizable notifications
* Notify via email, text, etc
* Automatic compression and encryption
* Works with S3, FTP, FTPS, SFTP, etc
* Choose from common schedules or build your own
* Database backups (postgresql, mysql, mongodb)

## User Stories
* Guest registers and is logged in
* Guest logs in
* User verifies their email address (required to log in?)
* User switches current project
* User invites a Guest to their project (by email)
* User deletes another User from their project (owner or admin)
* User updates another User's role (owner or admin)
* User updates payment info (require upon reg?)
* User deletes Account (and all ref'd info)
* User CRUDs a Location (no delete if affected Transfers)
* User CRUDs a Transfer (history entries persist but won't link back)
* User CRUDs a Schedule (no delete if affected Transfers)
* User links / unlinks a Transfer and Schedule
* User runs a Transfer adhoc
* User reruns a completed Transfer adhoc
* User views Transfer history
* User cancels an in-progress Transfer
* User updates their payment info
* User views old invoices
* User switches account auth method (email, github, google, etc)
