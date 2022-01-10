# dripfile
Managed file transfer as a service

## Design
This project is broken up into four sub-programs:
* `dripfile-migrate` - compares and applies database migrations
* `dripfile-web` - primary CRUD application web server
* `dripfile-worker` - watches the queue and performs file transfer jobs
* `dripfile-scheduler` - manages job schedules and publishes them to the queue

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
tailwindcss --watch -m -i tailwind.input.css -o static/css/tailwind.min.css
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
* User CRUDs a Location (no delete if affected Jobs)
* User CRUDs a Job (history entries persist but won't link back)
* User CRUDs a Schedule (no delete if affected Jobs)
* User links / unlinks a Job and Schedule
* User runs a Job adhoc
* User reruns a completed Job adhoc
* User views Job history
* User cancels an in-progress Job
* User updates their payment info
* User views old invoices
* User switches account auth method (email, github, google, etc)
