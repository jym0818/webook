local key = KEYS[1]

local cntKey = key..":cnt"
local inputCode = ARGV[1]
local code = redis.call("get",key)
local cnt = tonumber(redis.call("get",cntKey))

if cnt == nil or cnt <= 0 then
    return -1
end

if code == inputCode then
    return 0
else
    return -2
end