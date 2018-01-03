package main

// #cgo pkg-config: python-2.7
// #define Py_LIMITED_API
// #include <Python.h>
// int PyArg_ParseTuple_LL(PyObject *, long long *, long long *);
import "C"

//export schedule
func schedule(self, args *C.PyObject) *C.PyObject {
    var a, b C.longlong
    if C.PyArg_ParseTuple_LL(args, &a, &b) == 0 {
        return nil
    }
    return C.PyLong_FromLongLong(a + b)
}

//export unschedule
func unschedule(self, args *C.PyObject) *C.PyObject {
    var a, b C.longlong
    if C.PyArg_ParseTuple_LL(args, &a, &b) == 0 {
        return nil
    }
    return C.PyLong_FromLongLong(a + b)
}

//export stop
func stop(self, args *C.PyObject) *C.PyObject {
    var a, b C.longlong
    if C.PyArg_ParseTuple_LL(args, &a, &b) == 0 {
        return nil
    }
    return C.PyLong_FromLongLong(a + b)
}

//export loop
func loop(self, args *C.PyObject) *C.PyObject {
    var a, b C.longlong
    if C.PyArg_ParseTuple_LL(args, &a, &b) == 0 {
        return nil
    }
    return C.PyLong_FromLongLong(a + b)
}

func main() {}
