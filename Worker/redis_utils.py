import redis
import logging
import time

def connect_redis(host, port, passwd):
    redisClient = redis.Redis(host=host, port=port, password=passwd)    
    try:
        redisClient.ping()
    except redis.exceptions.ConnectionError as e:
        logging.info(f"Could not connect to Redis: {e}, waiting 5s")
        time.sleep(5)
        connect_redis(host, port, passwd)
    logging.info("REDIS client initialized")
    return redisClient

def add(a, b) -> int:
    return a+b