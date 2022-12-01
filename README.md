Config generator for running the AxeTrader server in Wildfly.

Available as a standalone binary and docker image.

# Use within Dockerfile

```Dockerfile
FROM axetrader-base-image AS base

FROM ghcr.io/axetrading/axetrader-standalone-xml-generator AS config

COPY config.json config.json
RUN /generate < config.json > standalone.xml

FROM base

COPY --from=config /standalone.xml /home/axe/axetrader/wildfly/standalone/configuration/standalone.xml
```

# Use standalone with docker

```shell
docker run -i ghcr.io/axetrading/axetrader-standalone-xml-generator < config.json > standalone.xml
```

# Use a downloaded binary

```shell
axetrader-standalone-xml-generator < config.json > standalone.xml
```

# Interface

This tool provides an opinionated interface to the config. Any values in the incoming JSON that it doesn't recognise will cause the tool to fail fast. The only options that are supported are those that are directly reflected in the standalone.xml - this tool has one responsibility and one only: generate standalone.xml.

In a couple of cases (e.g. database hostname and password) we have dropped support for configuration - this tool is designed to run at the start of the pipeline to generate the common config shared across environments. For valuse like these that need to be different between different environments we reference system properties to allow these values to be injected in.

## Config structure

### `database`

This top level key contains details for connecting to the database. It includes the following fields:

### `dialect`

This must be `{ "Psql": "postgresql" }` (the default if not present), or `{ "Mssql": "mssqlserver" }`

### `name`

The name of the database we connect to. Defaults to `axetrader` (default recommended).

### `port`

The port, defaults to 5432 for postgres or 1433 for mssql (defaults recommended).

### `User`

The username to connect to the database. Defaults to `axetrader` (default recommended).

## `wildfly`

Top-level key with details for wildfly config.

### `port_client` (integer)

The port to listen for the client on.

### `db_jndi_name` (strong)

TODO

### `db_max_pool_size` (integer)

TODO

### `db_pool_name` (string)

TODO

### `metrics` (boolean)

TODO

### `statistics` (boolean)

TODO

### systemProperties (map[string]string)

TODO
