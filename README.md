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
   
- **Configuration-Based Rules:** <br>
   Define rate limiting rules via YAML/JSON configuration files for version control and easy deployment.
   
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

#### 📝 **Configuration-Based Rules**
  
RateShield now supports loading rate limiting rules from configuration files! This enables:

- **Version Control**: Store rules in your repository
- **Environment Management**: Different configs for different environments  
- **Easy Deployment**: No manual rule configuration needed

**Quick Start:**
1. Create a `rules_config.yaml` file in your application directory
2. Define your rate limiting rules (see [CONFIG_RULES.md](rate_shield/CONFIG_RULES.md) for examples)
3. Start the application - it will automatically detect and use the config file

**Example Configuration:**
```yaml
rules:
  - strategy: "TOKEN BUCKET"
    endpoint: "/api/v1/users"
    http_method: "GET"
    allow_on_error: true
    token_bucket_rule:
      bucket_capacity: 100
      token_add_rate: 10
      retention_time: 3600
```

For detailed documentation, see [CONFIG_RULES.md](rate_shield/CONFIG_RULES.md).

---

### A detailed blog post about its working.
[Read Here](https://beyondthesyntax.substack.com/p/i-made-a-configurable-rate-limiter)

---

#### How it works?
<img src="https://raw.githubusercontent.com/x-sushant-x/Rate-Shield/main/assets/architecture.png"></img>

 ---

#### 🤝 **Contributing**

Interested in contributing? We'd love your help! Check out our [Contribution Guidelines](https://github.com/x-sushant-x/Rate-Shield/blob/main/CONTRIBUTION.md) to get started.

---
