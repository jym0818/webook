--你的验证码在 Redis上的key local定义一个局部变量
-- phone_code:login:15904922108
local key = KEYS[1]
--验证次数 我们一个验证码最多重复三次 我们记录验证了几次
-- phone_code:login:15904922108:cnt
local cntKey = key..":cnt"
-- 你的验证码 123456
local val = ARGV[1]
--过期时间
local ttl = tonumber(redis.call("ttl",key))
if ttl == -1 then
    --key存在 但是过期时间没有------系统出错了
    return -2
    --540 = 600-60
elseif ttl ==-2 or ttl<540 then
    redis.call("set",key,val)
    redis.call("expire",key,600)
    redis.call("set",cntKey,3)
    redis.call("expire",cntKey,600)
    return 0
else
    --发送频繁
    return -1
end