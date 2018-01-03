#define Py_LIMITED_API
#include <Python.h>

#if PY_MAJOR_VERSION >= 3
  #define MOD_ERROR_VAL NULL
  #define MOD_SUCCESS_VAL(val) val
  #define MOD_INIT(name) PyMODINIT_FUNC PyInit_##name(void)
  #define MOD_DEF(ob, name, doc, methods) \
          static struct PyModuleDef moduledef = { \
            PyModuleDef_HEAD_INIT, name, doc, -1, methods, }; \
          ob = PyModule_Create(&moduledef);
#else
  #define MOD_ERROR_VAL
  #define MOD_SUCCESS_VAL(val)
  #define MOD_INIT(name) void init##name(void)
  #define MOD_DEF(ob, name, doc, methods) \
          ob = Py_InitModule3(name, methods, doc);
#endif

PyObject * loop(PyObject *, PyObject *);
PyObject * stop(PyObject *, PyObject *);
PyObject * unschedule(PyObject *, PyObject *);
PyObject * schedule(PyObject *, PyObject *);

// Workaround missing variadic function support
// https://github.com/golang/go/issues/975
int PyArg_ParseTuple_LL(PyObject * args, long long * a, long long * b) {
    return PyArg_ParseTuple(args, "LL", a, b);
}

static PyMethodDef methods[] = {
    {"loop", loop, METH_VARARGS, "Start looping."},
    {"stop", stop, METH_O, "Stop the watcher."},
    {"unschedule", unschedule, METH_O, "Reset the watcher for the given path."},
    {"schedule", schedule, METH_VARARGS, "Setup the watcher for the given path."},
    {NULL, NULL, 0, NULL}
};

static char doc[] = "Low-level FSEvent interface.";

PyMODINIT_FUNC
PyInit_foo(void)
{
    PyObject* mod;
    MOD_DEF(mod, "_fsevents", doc, methods);
    return MOD_SUCCESS_VAL(mod);
}
