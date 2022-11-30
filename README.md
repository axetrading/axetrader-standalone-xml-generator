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
