# object-store
This is an object store using redis as a data store. It is built for the purpose of an interview at endor labs. 

### Prerequisite
- Redis 7.0

### Assumptions
- If there are more than one objects exist with the same name then `GetOBjectByName` will return the most recently saved object.

### Setup
**Want to skip setup?:** [Run using docker](#run-using-docker) or just run redis on `localhost:6397` with no auth
- `export DB_Host=<DB_HOST>`
- `export DB_PORT=<DB_PORT>`
- `export DB_USER=<DB_USER>`
- `export DB_PASSWORD=<DB_PASSWORD>`

### Building object store
It is as simple as

``make build``


### Running object store

It is so simple :-)

``make run``

#### Run using docker

``docker-compose run object-store``

### Testing
Grab a cup of coffee and run

``make test``

Happy Judging!!

