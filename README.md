## 🚀 **RateShield**

A fully customizable rate limiter designed to apply rate limiting on individual APIs with specific rules.


#### 📊 **Dashboard Overview**

![RateShield Dashboard](https://raw.githubusercontent.com/x-sushant-x/Rate-Shield/main/assets/main.png)


![RateShield Edit Rule](https://raw.githubusercontent.com/x-sushant-x/Rate-Shield/main/assets/Edit%20Rule.png)

___

#### 🎯 **Why RateShield?**

Why not? With some free time on hand, RateShield was created to explore the potential of building a versatile rate-limiting solution. What started as a side project is evolving into a powerful tool for developers.

---

#### 🌟 **Key Features**

- **Customizable Limiting:** <br>
   Tailor rate limiting rules to each API endpoint according to your needs.
   
- **Intuitive Dashboard:** <br>
   A user-friendly interface to monitor and manage all your rate limits effectively.
   
- **Easy Integration:** <br>
   Plug-and-play middleware that seamlessly integrates into your existing infrastructure.

---

#### ⚙️ **Use Cases**

- **Preventing Abuse:**  
  Control the number of requests your APIs can handle to prevent misuse and malicious activities.
  
- **Cost Management:**  
  Manage third-party API calls efficiently to avoid unexpected overages.

---

#### 🚀 **Supported Rate Limiting Algorithms**

- **Token Bucket**
- **Fixed Window Counter**
- **Sliding Window**

---

#### 🪧 Usage Guide
  Check out this [document](https://github.com/x-sushant-x/Rate-Shield/tree/main/rate_shield/documentation).

---

### A detailed blog post about its working.
[Read Here](https://beyondthesyntax.substack.com/p/i-made-a-configurable-rate-limiter)

---

#### How it works?
<img src="https://raw.githubusercontent.com/x-sushant-x/Rate-Shield/main/assets/architecture.png"></img>

---

#### 📋 **Prerequisites**

Before setting up RateShield, ensure you have the following:

- **Docker & Docker Compose** - Required for running the application and Redis infrastructure
- **Redis Stack** (included in docker-compose):
  - Single Redis instance for storing rate limit rules (port 6379)
  - Redis Cluster with 6 nodes (3 masters + 3 replicas) for distributed rate limiting (ports 7000-7005)
  - ReJSON module enabled (automatically included with redis/redis-stack image)
- **Slack Integration** (for error notifications):
  - Slack Bot Token
  - Slack Channel ID

**Quick Start with Docker Compose:**

The provided `docker-compose.yml` includes everything you need - no separate Redis setup required!

```bash
cd rate_shield
cp .env.example .env
# Edit .env with your Slack credentials
docker-compose up -d
```

For detailed setup instructions, see [SETUP.md](SETUP.md).

 ---

#### 🤝 **Contributing**

Interested in contributing? We'd love your help! Check out our [Contribution Guidelines](https://github.com/x-sushant-x/Rate-Shield/blob/main/CONTRIBUTION.md) to get started.

---
