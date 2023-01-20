# testdb

## Overview

A repo for creating test databases in PostgreSQL and/or MySQL, with Go clients for testing and debugging code.

## Requirements

- [dbmate](https://github.com/amacneil/dbmate) (for migrations)
- [sqlc](https://github.com/kyleconroy/sqlc) (for generating Go code from SQL queries)

## Usage

### Setup

```bash
# contains postgre_url, mysql_url
cp .env.example .env
```

### Migration

```bash
# migrates database to latest version
make pgdbup or make mydbup

# migrates database to previous version
make pgdbdown or make mydbdown

# create new migration file
make pgdbnew name={migration_name} or make mydbnew name={migration_name}
```

### PostgresSQL

```bash
# start docker container
make pgup or make myup

# stop docker container
make pgdown or make mydown

# destory docker container
make pgstop or make mystop
```

### Development

```bash
# generate queries Go code for both PostgreSQL and MySQL
make sqlc
```
