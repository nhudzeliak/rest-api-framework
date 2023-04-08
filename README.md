# Web application task

## Quick Start
For those who don't want to get into details, there are dockerized versions of all services this projects provides. 
Simply execute `make docker_up` in the project's root directory to build, compose, and run the produced images.

## Running the App
To run the app outside of Docker, execute `make run` from the project's root directory. This will re-build and start a 
new instance of the API on host and port provided in the config files (more on that in `Config` section below) or ones 
provided in `OVERRIDE_HOST` and `OVERRIDE_PORT` flags. Data source name for the database is always fetched from configs
and can't be overridden by cmd flags.

## Database
The database uses PostgreSQL v14.1 driver. Schema structure is managed using migrations, which you can create with 
`make new_migration --name="migration name here"` (instantiating `.up.sql` and `down.slq` files in `./db/migrations/`), 
and apply using `make migrate`. The latter will compare all migration ids found in `./db/migrations/` folder with 
`versions` table in the target database and run the missing ones, registering them in the aforementioned table.

There are two users created in the initial migration of the database (`./db/migrations/_init.sql`) -- a readonly one, 
and one with full privileges. This is done to allow for later separation of the database into a readonly replica and a 
common one, without touching the code.

## Testing
The project implements both unit and integration tests.
* To run unit tests, execute `make unit_test`. This will traverse all the directories and subdirectories of the project,
    and run all the test files found in it (except for those located under `./test/` folder).
* To run integration tests, execute `make integration_test`. This will re-build the app and start a new instance of the
    api on port specified in `Makefile` running in the background, and then run the unit tests located in `test` package. 
    The underlying API instance is killed when tests are done running.

## Config
There are four section in the `app/config/app.conf` file, each responsible for a separate instance of the app. This 
files main usage is supposed to be storing credentials/secrets that would otherwise be set in the environment. 