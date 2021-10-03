import redis

def test_add():
    assert(redis.add(4, 5)==9)