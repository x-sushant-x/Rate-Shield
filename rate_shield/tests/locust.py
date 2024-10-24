from locust import FastHttpUser, TaskSet, task, between
import random

class RateLimiterLoadTest(TaskSet):
    @task
    def checkLoad(self):
        ip = f"192.168.{random.randint(0, 255)}.{random.randint(0, 255)}"
        endpoint = "/api/v1/test"
        headers = {
            "ip" : ip,
            "endpoint" : endpoint,
            "Accept": "*/*",
            "User-Agent": "LocustLoadTest/1.0"
        }
        self.client.get("/check-limit", headers = headers)

class User(FastHttpUser):
    tasks = [RateLimiterLoadTest]
    wait_time = between(0.1, 0.5) 