### How to setup Rate Shield?

#### Prerequisite
* Docker & Docker Compose
* Redis Instance (single instance for storing rules)
* Redis Cluster (6 nodes: 3 masters + 3 replicas with ReJSON Module Enabled for distributed rate limiting)

* Environment variables for the application:
    * RATE_SHIELD_PORT: The port number the application will listen on.
    * REDIS_RULES_INSTANCE_URL: URL for the Redis instance storing rate limit rules.
    * REDIS_RULES_INSTANCE_USERNAME: Username for Redis authentication (optional).
    * REDIS_RULES_INSTANCE_PASSWORD: Password for Redis authentication (optional).
    * REDIS_CLUSTERS_URLS: Comma separated cluster node URLs. Ex - `redis-node-1:7000,redis-node-2:7001,redis-node-3:7002,redis-node-4:7003,redis-node-5:7004,redis-node-6:7005`
    * REDIS_CLUSTER_USERNAME: Username for Redis Cluster authentication (optional).
    * REDIS_CLUSTER_PASSWORD: Password for Redis Cluster authentication (optional).
    * SLACK_TOKEN: Slack bot token for error notifications.
    * SLACK_CHANNEL: Slack channel ID for notifications.

---

### Setup with Docker Compose (Recommended)

The easiest way to set up RateShield with Redis Cluster is using Docker Compose.

1. Go to `rate_shield` subfolder.

2. Copy the example environment file and configure it:
```bash
cp .env.example .env
```

3. Edit the `.env` file and set your Slack credentials:
```bash
SLACK_TOKEN=your-slack-token-here
SLACK_CHANNEL=your-slack-channel-id-here
```

4. Start all services (including Redis Cluster with 6 nodes and RedisInsight):
```bash
docker-compose up -d
```

This will start:
- **redis-rules**: Single Redis instance for storing rate limit rules (port 6379)
- **redis-node-1 to redis-node-6**: Redis Cluster with 3 masters and 3 replicas (ports 7000-7005)
- **redis-cluster-init**: Initialization service that creates the cluster
- **redisinsight**: Web UI for monitoring Redis cluster (port 8001)
- **app**: RateShield application (port 8080)

5. Verify the cluster is running:
```bash
docker exec -it redis-node-1 redis-cli --cluster check redis-node-1:7000
```

6. Access the services:
- **RateShield API**: http://localhost:8080
- **RedisInsight Dashboard**: http://localhost:8001

7. To stop all services:
```bash
docker-compose down
```

To remove all data volumes:
```bash
docker-compose down -v
```

---

### Setup with Docker (Manual)

If you prefer to set up Redis manually, follow these steps:

1. Set up a Redis instance for rules storage (must have ReJSON module).

2. Set up a Redis Cluster with 6 nodes (3 masters + 3 replicas) with ReJSON module enabled.

3. Go to `rate_shield` subfolder and build the Docker image:
```bash
docker build -t rate-shield-backend .
```

4. Run the container with appropriate environment variables:
```bash
docker run -d \
  -p 8080:8080 \
  -e RATE_SHIELD_PORT=8080 \
  -e REDIS_RULES_INSTANCE_URL=redis://your-redis-host:6379 \
  -e REDIS_RULES_INSTANCE_PASSWORD=your-password \
  -e REDIS_CLUSTERS_URLS=node1:7000,node2:7001,node3:7002,node4:7003,node5:7004,node6:7005 \
  -e REDIS_CLUSTER_PASSWORD=your-cluster-password \
  -e SLACK_TOKEN=your-slack-token \
  -e SLACK_CHANNEL=your-slack-channel \
  rate-shield-backend
```

<b>Important: - </b> Value for -p and RATE_SHIELD_PORT must be same.

Now you can access rate shield via localhost:8080 (value passed in RATE_SHIELD_PORT).

---

### Development Setup

For development, use the dev compose file which includes only the Redis infrastructure:

```bash
cd rate_shield
docker-compose -f docker-compose-dev.yml up -d
```

Then run the application locally:
```bash
go run main.go
```

---

### Troubleshooting

**Cluster not forming:**
- Ensure all 6 Redis nodes are running: `docker ps`
- Check cluster initialization logs: `docker logs redis-cluster-init`
- Manually verify cluster: `docker exec -it redis-node-1 redis-cli --cluster check redis-node-1:7000`

**Application can't connect to cluster:**
- Verify REDIS_CLUSTERS_URLS contains all 6 node addresses
- Check network connectivity: `docker network inspect rate-shield-network`
- Ensure cluster has ReJSON module: `docker exec -it redis-node-1 redis-cli MODULE LIST`

**RedisInsight not accessible:**
- Check if the container is running: `docker ps | grep redisinsight`
- Access at http://localhost:8001
- Add cluster manually: Use host `redis-node-1` and port `7000`

---

Follow [this](https://github.com/x-sushant-x/Rate-Shield/tree/main/rate_shield/documentation) usage guide to know how to use Rate Shield.

<br>

### Setup Frontend Client
1. Go to `web` subfolder and build docker image using `docker build -t rate-shield-frontend .`
2. Once docker image is built you can run it using below command.

```
docker run -p 6012:6012 \
-e VITE_RATE_SHIELD_BACKEND_BASE_URL=http://localhost:9081 \
-e PORT=6012 \
rate_shield_frontend
```

<b>Important: -</b> VITE_RATE_SHIELD_BACKEND_BASE_URL should have the URL of your backend docker container and -p and -e PORT value should match.
