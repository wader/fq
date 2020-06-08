-- local function match(path)
--     print("match:", path)

--     return function(params)
--         print("params:", params)
--         -- both path and params are now availble for use here
--     end
-- end

-- match "" {
--     BLA = function()
--         print "ASdsad"
--     end,
--     TJO = 123,
--     ASD = 123
-- }


function section(name, fn)
    print(name)
    fn()
end

local a = section("asd", function()
    return 123
end)

print(a)

int(1, "asd")
uint 2
uint 8 "asd"


