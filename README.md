# dripfile
Managed file transfer as a service

## Design
This project is broken up into four sub-programs:
* `dripfile-migrate` - compares and applies database migrations
* `dripfile-web` - primary CRUD application web server
* `dripfile-worker` - watches the queue and performs file transfers
* `dripfile-scheduler` - manages transfer schedules and publishes them to the queue

## Setup
This project depends on the [Go programming language](https://golang.org/dl/) and the [TailwindCSS CLI](https://tailwindcss.com/blog/standalone-cli).

## Database
This project uses [PostgreSQL](https://www.postgresql.org/) for persistent storage.
To develop locally, you'll an instance of the database running somehow or another.
I find [Docker](https://www.docker.com/) to be a nice tool for this but you can do whatever works best.

The following command starts the necessary containers:
```bash
docker compose up -d
```

These containers can be stopped via:
```bash
docker compose down
```

## Running
If actively working on frontend templates, set `ENV=dev` to tell the server to reload templates from the filesystem on every page load.
Apply database migrations, run whichever components you need (in background processes), and let Tailwind watch for CSS changes:
```bash
# make run
ENV=dev go run cmd/migrate/main.go
ENV=dev go run cmd/web/main.go &
ENV=dev go run cmd/worker/main.go &
ENV=dev go run cmd/scheduler/main.go &
tailwindcss --watch -m -i tailwind.input.css -o internal/static/static/css/tailwind.min.css
```

## Testing
Tests can be ran after starting the necessary containers and applying database migrations:
```bash
# make test
ENV=dev go run cmd/migrate/main.go
go test -v ./...
```

Note that the tests will leave random test in the database so feel free to flush it out by restarting the containers:
```bash
docker compose down
docker compose up -d
```

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

## Test Cases
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
