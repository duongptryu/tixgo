# Database Migrations

This folder contains the database migration scripts for the TixGo project. Migrations are managed using [golang-migrate/migrate](https://github.com/golang-migrate/migrate), a database migration tool for Go projects.

## Folder Structure

- `000001_init_schema.up.sql`   — Creates the initial database schema (tables, indexes, constraints, etc.)
- `000001_init_schema.down.sql` — Drops all tables and objects created by the corresponding `up` migration
- Additional migration files should follow the naming convention: `<version>_<description>.up.sql` and `<version>_<description>.down.sql`

## Migration Tool: golang-migrate

We use [golang-migrate/migrate](https://github.com/golang-migrate/migrate) to manage schema changes. See the [official PostgreSQL tutorial](https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md) for more details.

### Installation

You can install the CLI tool via Homebrew (macOS):

```sh
brew install golang-migrate
```

Or download a binary from the [releases page](https://github.com/golang-migrate/migrate/releases).

### Creating a New Migration

To create a new migration, run:

```sh
migrate create -ext sql -dir ./migrations -seq <migration_name>
```

This will generate two files:
- `<version>_<migration_name>.up.sql`   — for applying the migration
- `<version>_<migration_name>.down.sql` — for rolling back the migration

Edit these files to add your SQL changes.

### Running Migrations

To apply all up migrations:

```sh
migrate -path ./migrations -database "postgres://<user>:<password>@<host>:<port>/<database>?sslmode=disable" up
```

To rollback the last migration:

```sh
migrate -path ./migrations -database "postgres://<user>:<password>@<host>:<port>/<database>?sslmode=disable" down 1
```

You can also step up/down by N migrations, or force a version. See the [CLI documentation](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) for more options.

### Example Connection String

```
postgres://postgres:postgres@localhost:5432/tixgo?sslmode=disable
```

- Replace `postgres:postgres` with your username and password
- Replace `localhost:5432` with your host and port
- Replace `tixgo` with your database name

### Best Practices

- Always create both `up` and `down` migration scripts
- Test migrations on a local/dev database before applying to production
- Keep migration files in version control
- Use sequential version numbers for easy tracking

---

For more details, see the [golang-migrate PostgreSQL tutorial](https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md).
