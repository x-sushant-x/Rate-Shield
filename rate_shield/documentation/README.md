<!-- ### How to use?

**Rate Shield** has an endpoint `/check-limit` exposed which accept following headers: 

```
ip: <IP_ADDRESS> 
endpoint: <API_TARGET_API_ENDPOINT>
```

When you send following data to `/check-limit` it will fetch rule defined for that endpoint and apply rate limiting. Once all this process is done it will return with appropriate status code. It may return `200, 429 or 500` codes. 


Based on this response you can decide if you have to hit target API or not.

<br>

**Automate using middleware** 

You can define your custom middleware or intercepters to automate this process. Your middleware will send ip and endpoint to **Rate Shield** and return appropriate response. Below is the cURL of /check-limit.

```
curl -X GET \
  'localhost:8080/check-limit' \
  --header 'Accept: */*' \
  --header 'ip: 127.0.0.1' \
  --header 'endpoint: /api/v1/resource'
```

**Response**
```
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Content-Type, Authorization
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Access-Control-Allow-Origin: *
Rate-Limit: 100
Rate-Limit-Remaining: 97
Date: Sun, 08 Sep 2024 18:39:18 GMT
Content-Length: 0
```

I'm not writing exact code for creating middleware in different languages and framework as you can create it in your favourite language more efficently than me. You can also contribute to this project by creating custom middleware for **Rate Shield** in your favourite langauge/framework and than sharing it here. -->

## Rate Shield Documentation
**Rate Shield** provides a /check-limit endpoint that enables you to implement rate limiting for your API endpoints. By sending specific headers to this endpoint, Rate Shield applies the defined rate limiting rules and returns an appropriate HTTP status code based on the result.

### How to Use
#### Endpoint Details
The /check-limit endpoint accepts the following headers:

* `ip:` <IP_ADDRESS>
* `endpoint:` <API_TARGET_API_ENDPOINT>
<br>

When you send a request with these headers to /check-limit, Rate Shield retrieves the rate limiting rules defined for the specified endpoint and applies them based on the provided IP address. After processing, it returns one of the following HTTP status codes:

* `200 OK:` The request is within the rate limit or **no rules are defined for the endpoint.**
* `429 Too Many Requests:` The rate limit has been exceeded.
* `500 Internal Server Error:` An error occurred during processing.

Based on the response status code, you can decide whether to proceed with the request to your target API.

### Example Request
Below is an example of using cURL to send a request to the /check-limit endpoint:

```
curl -X GET \
  'http://localhost:8080/check-limit' \
  --header 'Accept: */*' \
  --header 'ip: 127.0.0.1' \
  --header 'endpoint: /api/v1/resource'
```

### Example Response
```
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Content-Type, Authorization
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Access-Control-Allow-Origin: *
Rate-Limit: 100
Rate-Limit-Remaining: 97
Date: Sun, 08 Sep 2024 18:39:18 GMT
Content-Length: 0
```

The response headers include rate limit information:

* `Rate-Limit:` The total number of allowed requests.
* `Rate-Limit-Remaining:` The number of remaining requests in the current time window.

### Automating with Middleware
To streamline the rate limiting process, you can create custom middleware or interceptors in your preferred programming language and framework. The middleware should:

1. Send the client's IP address and the requested endpoint to the Rate Shield /check-limit endpoint.
2. Receive the response and handle it accordingly, such as allowing the request to proceed or returning an error message to the client.

While I'm not providing specific code examples for creating middleware in different languages and frameworks, we encourage you to implement it in your environment of choice. Your contributions are valuable; feel free to share your custom middleware implementations with the community to enhance the Rate Shield project.

---

## Redis Fallback Feature

Rate Shield includes an **optional in-memory fallback** mechanism that activates when Redis becomes unavailable.

### Configuration

Enable fallback mode by setting these environment variables:

```bash
ENABLE_REDIS_FALLBACK=true       # Enable fallback (default: false)
REDIS_RETRY_INTERVAL=30s         # Health check interval (default: 30s)
```

### How It Works

1. **Redis Available**: Normal operation using Redis for rate limiting
2. **Redis Fails**: Automatically switches to in-memory storage
3. **Periodic Health Checks**: Monitors Redis availability based on retry interval
4. **Auto-Recovery**: Switches back to Redis when connection is restored

### Important Notes

⚠️ **Limitations:**
- In-memory data is **per-instance** (not distributed across multiple Rate Shield instances)
- Data is **lost on restart** (non-persistent)
- Best for **development**, **testing**, or **temporary resilience**
- **Not recommended** for production multi-instance deployments

✅ **Best For:**
- Development without Redis dependency
- Handling temporary Redis outages
- Single-instance deployments

📖 **Full Documentation**: See [REDIS_FALLBACK.md](../../REDIS_FALLBACK.md) for complete details, logging behavior, and best practices.
