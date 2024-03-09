# go-rest-template

## Overview

TODO

## Running the app locally

**Requires Docker installed and running**

This repository includes a `Makefile` with some helpful commands.

These commands are driven from the values set in the `.env` file.

| Command            | Description                                                                                                                   |
|--------------------|-------------------------------------------------------------------------------------------------------------------------------|
| `run`              | Runs the application (requires DB to be running and DB migrations applied                                                     |
| `db/up`            | Bring up a local postgres container (configured in the `docker-compose.yml`) that the app can connect to when running locally |
| `db/down`          | Terminate the postgres container                                                                                              |
| `db/migrations/up` | Runs the database schema migrations.  If the database is not already up, the database will be started first.                  |
