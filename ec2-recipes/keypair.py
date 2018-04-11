from cffi import FFI

ffi = FFI()



lib = ffi.dlopen("./libmath.dylib")


