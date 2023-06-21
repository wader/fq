local sometable = {
	true,
	false,
	nil,
	437784932,
	4.23748378e-6,

	somenil= nil,
	sometrue= true,
	somefalse= false,
	someint= -3,
	somenum= 7.89437298e11,
	somestr= "uwu",

	[2.74389]= "key is a num",
	[-1337]= "key is an int",
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
