local pattern = KEYS[1] -- 输入的 pattern
local result = {}
local cursor = "0" -- 声明并初始化 cursor 变量
repeat
    local data = redis.call("SCAN", cursor, "MATCH", pattern)
    cursor = data[1]
    local keys = data[2]

    for _, key in ipairs(keys) do
        local value = redis.call("GET", key)

        ---- 解析 JSON 字符串为 Lua 表
        --local decoded_value = cjson.decode(value)
        --
        ---- 检查解析后的值是否为表，如果不是，将其转换为表
        --if type(decoded_value) ~= "table" then
        --    decoded_value = {name = decoded_value}
        --end
        --
        ---- 将解析后的值添加到 result 数组中
        --table.insert(result, decoded_value)
        table.insert(result, value)
    end
until cursor == "0"

return result