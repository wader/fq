local sometable = {
	somenil= nil,
	sometrue= true,
	somefalse= false,
	someint= -3,
	somenum= 7.89437298e11,
	somestr= "it's a trap",
}

mycplx = 3.2i

mytbl = sometable

-- f1's upvalue
local a = 123
local b = 666

local f1 = function(x)
	local c = a + b
	return x * c * 2973289 + 38793457897
end

myfunc = f1
myfunc_result = f1(42)
