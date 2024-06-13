# Database

Database operations

## Requirement
The following tool need to be install locally :

 - [liquibase](https://docs.liquibase.com/start/install/home.html)
 - `make`

## Deployment on staging and prod environment

To deploy on `staging`/`production` you need to set some variables (`host`, `port`, `username` and `password`) on a liquibase.properties file.

#### update DB
Then set the LB_PROPERTIES as follow 
```shell
export LB_PROPERTIES=<liquibase.properties> && make db-migrate-up
```
or
```shell
LB_PROPERTIES=<liquibase.properties> make db-migrate-up
```
#### roolback DB by 1 revision
```shell
export LB_PROPERTIES=<liquibase.properties> && make db-migrate-down
```
or
```shell
LB_PROPERTIES=<liquibase.properties> make db-migrate-down
```

## For local dev

### update
```shell
cd migration/liquibase
liquibase update
```

### roolback by n
```shell
cd migration/liquibase
liquidbase rollback-count n
```

### Help

For more help try:
```shell
make help
```