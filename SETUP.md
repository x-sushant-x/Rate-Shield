### How to setup Rate Shield?

#### Prerequisite
* Docker
* Redis Instance (for rules) - Optional if fallback is enabled
* Redis Cluster (with ReJSON Module Enabled) - Optional if fallback is enabled

* Environment variables for the application:
    * RATE_SHIELD_PORT: The port number the application will listen on.
    * REDIS_RULES_INSTANCE_URL: URL for the Redis instance storing rate limit rules.
    * REDIS_RULES_INSTANCE_USERNAME: Username for Redis authentication.
    * REDIS_RULES_INSTANCE_PASSWORD: Password for Redis authentication.
    * REDIS_CLUSTERS_URLS: Comma seperated urls. Ex - 127.0.0.1:6380,127.0.0.1:6381
    * REDIS_CLUSTER_USERNAME: Username for Redis Cluster authentication.
    * REDIS_CLUSTER_PASSWORD: Password for Redis Cluster authentication.
    * **ENABLE_REDIS_FALLBACK** (Optional): Enable in-memory fallback when Redis is unavailable (true/false, default: false)
    * **REDIS_RETRY_INTERVAL** (Optional): Interval for retrying Redis connection when using fallback (default: 30s, format: 30s, 1m, 5m, etc.)


<br>

1. Go to `rate_shield` subfolder.
2. Build docker image using `docker build -t rate-shield-backend .`

3. Once docker image is built you can run it using below command.

**Standard setup (with Redis):**
```
docker run -d \
  -p 8080:8080 \
  -e RATE_SHIELD_PORT=8080 \
  -e REDIS_RULES_INSTANCE_URL=redis://localhost:6379 \
  -e REDIS_RULES_INSTANCE_USERNAME=user \
  -e REDIS_RULES_INSTANCE_PASSWORD=password \
  -e REDIS_CLUSTERS_URLS=localhost:6380,localhost:6381 \
  rate-shield-app
```

**With Redis fallback enabled:**
```
docker run -d \
  -p 8080:8080 \
  -e RATE_SHIELD_PORT=8080 \
  -e REDIS_RULES_INSTANCE_URL=redis://localhost:6379 \
  -e REDIS_RULES_INSTANCE_USERNAME=user \
  -e REDIS_RULES_INSTANCE_PASSWORD=password \
  -e REDIS_CLUSTERS_URLS=localhost:6380,localhost:6381 \
  -e ENABLE_REDIS_FALLBACK=true \
  -e REDIS_RETRY_INTERVAL=30s \
  rate-shield-app
```

<b>Important: - </b> Value for -p and RATE_SHIELD_PORT must be same.

Now you can access rate shield via localhost:8080 (value passed in RATE_SHIELD_PORT).

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

<br>

---

### 🔄 Redis Fallback Feature

Rate Shield now supports **automatic in-memory fallback** when Redis is unavailable. This feature is optional and disabled by default.

#### When to use fallback?

- ✅ **Development/Testing**: Work without Redis dependency
- ✅ **Temporary Resilience**: Handle short Redis outages gracefully
- ✅ **Single Instance Deployments**: When running one Rate Shield instance

#### When NOT to use fallback?

- ❌ **Production Multi-Instance**: Multiple Rate Shield instances (data won't sync)
- ❌ **Strict Rate Limits**: When precise distributed rate limiting is critical
- ❌ **Long-Term Operation**: Extended periods without Redis

#### Quick Start

1. Set environment variables:
   ```bash
   ENABLE_REDIS_FALLBACK=true
   REDIS_RETRY_INTERVAL=30s  # Optional, default is 30s
   ```

2. Application will:
   - Start successfully even if Redis is down
   - Use in-memory storage for rate limiting
   - Automatically restore Redis when available
   - Log warnings during fallback mode

**📖 For complete documentation, limitations, and best practices, see [REDIS_FALLBACK.md](../REDIS_FALLBACK.md)**
