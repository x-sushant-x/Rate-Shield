# Configuration-Based Rate Limiting Rules

This document explains how to use file-based configuration for rate limiting rules in RateShield.

## Overview

RateShield now supports loading rate limiting rules from configuration files instead of Redis. This provides several benefits:

- **Version Control**: Rules can be stored in version control systems
- **Environment Management**: Different configurations for different environments
- **Easier Deployment**: No need to manually configure rules via Dashboard UI
- **Backup & Recovery**: Configuration files serve as backup

## Configuration File Detection

On application startup, RateShield checks for configuration files in the following order:

1. `rules_config.yaml`
2. `rules_config.yml`
3. `rules_config.json`

If any of these files are found, RateShield will use file-based configuration. Otherwise, it falls back to Redis-based rules.

## Configuration File Format

### YAML Format (Recommended)

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

  - strategy: "FIXED WINDOW COUNTER"
    endpoint: "/api/v1/login"
    http_method: "POST"
    allow_on_error: false
    fixed_window_counter_rule:
      max_requests: 5
      window: 60

  - strategy: "SLIDING WINDOW COUNTER"
    endpoint: "/api/v1/upload"
    http_method: "POST"
    allow_on_error: true
    sliding_window_counter_rule:
      max_requests: 20
      window: 300
```

### JSON Format

```json
{
  "rules": [
    {
      "strategy": "TOKEN BUCKET",
      "endpoint": "/api/v1/users",
      "http_method": "GET",
      "allow_on_error": true,
      "token_bucket_rule": {
        "bucket_capacity": 100,
        "token_add_rate": 10,
        "retention_time": 3600
      }
    }
  ]
}
```

## Rule Configuration Fields

### Common Fields

- **strategy**: Rate limiting strategy (`"TOKEN BUCKET"`, `"FIXED WINDOW COUNTER"`, `"SLIDING WINDOW COUNTER"`)
- **endpoint**: API endpoint to apply the rule to
- **http_method**: HTTP method (GET, POST, PUT, DELETE, etc.)
- **allow_on_error**: Whether to allow requests when rate limiter encounters an error

### Token Bucket Rule

```yaml
token_bucket_rule:
  bucket_capacity: 100      # Maximum tokens in bucket
  token_add_rate: 10        # Tokens added per interval
  retention_time: 3600      # Token retention time in seconds
```

### Fixed Window Counter Rule

```yaml
fixed_window_counter_rule:
  max_requests: 5           # Maximum requests per window
  window: 60                # Window size in seconds
```

### Sliding Window Counter Rule

```yaml
sliding_window_counter_rule:
  max_requests: 20          # Maximum requests per window
  window: 300               # Window size in seconds
```

## Usage Examples

### 1. API Rate Limiting

```yaml
rules:
  # Limit API calls to 1000 requests per hour
  - strategy: "SLIDING WINDOW COUNTER"
    endpoint: "/api/v1/*"
    http_method: "GET"
    allow_on_error: true
    sliding_window_counter_rule:
      max_requests: 1000
      window: 3600
```

### 2. Login Protection

```yaml
rules:
  # Limit login attempts to 5 per minute
  - strategy: "FIXED WINDOW COUNTER"
    endpoint: "/api/v1/login"
    http_method: "POST"
    allow_on_error: false
    fixed_window_counter_rule:
      max_requests: 5
      window: 60
```

### 3. File Upload Limiting

```yaml
rules:
  # Use token bucket for file uploads (handles bursts)
  - strategy: "TOKEN BUCKET"
    endpoint: "/api/v1/upload"
    http_method: "POST"
    allow_on_error: true
    token_bucket_rule:
      bucket_capacity: 50
      token_add_rate: 5
      retention_time: 1800
```

## Migration from Redis

To migrate from Redis-based rules to configuration files:

1. Export existing rules from Redis (via Dashboard or API)
2. Convert them to YAML/JSON format
3. Save as `rules_config.yaml` in the application directory
4. Restart the application

## Limitations

When using configuration-based rules:

- **Read-Only**: Rules cannot be modified via Dashboard UI or API
- **No Real-Time Updates**: Application restart required to reload rules
- **No Rule Management**: CRUD operations are not supported

To modify rules, edit the configuration file and restart the application.

## Best Practices

1. **Use YAML**: More readable and supports comments
2. **Version Control**: Store configuration files in your repository
3. **Environment-Specific**: Use different config files for different environments
4. **Validation**: Test configuration files before deployment
5. **Documentation**: Add comments to explain complex rules

## Troubleshooting

### Configuration File Not Found
- Ensure the file is in the application directory
- Check file naming (must be exactly `rules_config.yaml`, `rules_config.yml`, or `rules_config.json`)

### Invalid Configuration
- Check YAML/JSON syntax
- Validate required fields are present
- Ensure numeric values are positive

### Application Startup Errors
- Check application logs for detailed error messages
- Validate rule configurations against the schema
- Ensure all required fields are provided

## Example Files

See the following example files in the project:
- `rules_config.yaml` - YAML format example
- `example_rules_config.json` - JSON format example
