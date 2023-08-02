# Docker Compose

## Configuring VMClarity

Each configurable service has a <service-name>.env file, in this file set the
attributes required for that service and it will be loaded by the compose file
when started.

## Overriding Parameters in the dockercompose.yml

You can override parameters in the dockercompose.yml by passing a custom env
file into the `docker compose up` command via the `--env-file` flag. An example
overriding all the container images `image_override.env` can be modified or
copied for this.

## Starting VMClarity
```
docker compose --project-name vmclarity --file dockercompose.yml up -d --wait --remove-orphans
```

## Stopping VMClarity
```
docker compose --project-name vmclarity --file dockercompose.yml down --remove-orphans
```
