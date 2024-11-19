### How to setup Rate Shield?
This document will guide you through the process of setting up Rate Shield. 

#### Prerequisite
* Redis Cluster
* Go Programming Language

#### Setup Redis Cluster
Rate Shield require redis cluster with ReJSON enabled. I'll update this guide soon for steps related setting up Redis Cluster.

#### Setup .env
Inside `rate_shield` folder create .env file and add following variables.
```
SLACK_TOKEN=
SLACK_CHANNEL=

# dev or prod
ENV=

# Redis Ports Config
REDIS_RULES_INSTANCE_URL=

REDIS_CLUSTERS_URLS=
```

Rate Shield store all rules in one Redis instance and all other rate limiting related data in cluster. 

Add complete URL for your redis instance such as `REDIS_RULES_INSTANCE_URL=127.0.0.1:6379`. 

Add comma seperated values of instances for cluster such as `REDIS_CLUSTERS_URLS=localhost:6380,localhost:6381`

#### Run Application
1. Go inside `rate_shield` folder.
2. Run it via `go run main.go`