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

PyObject * start(PyObject *, PyObject *);
PyObject * stop(PyObject *, PyObject *);
PyObject * schedule(PyObject *, PyObject *);

// Workaround missing variadic function support
// https://github.com/golang/go/issues/975
int PyArg_ParseTuple_ourArgs(PyObject* args, PyObject* callback, PyObject* paths) {
    return PyArg_ParseTuple(args, "00:schedule", callback, paths);
}

PyObject* PyObject_CallFunction_ourArgs(PyObject* _callback, char* path, long long flags) {
    return PyObject_CallFunction(_callback, "OO", path, flags);
}

static PyMethodDef methods[] = {
    {"start", start, METH_VARARGS, "Start watching."},
    {"stop", stop, METH_NOARGS, "Stop the watcher."},
    {"schedule", schedule, METH_VARARGS, "Setup the watcher for the given path(s)."},
    {NULL, NULL, 0, NULL}
};

static char doc[] = "Low-level FSEvent interface.";

MOD_INIT(fsevents_watcher) {
    PyObject* mod;
    MOD_DEF(mod, "fsevents_watcher", doc, methods);
    return MOD_SUCCESS_VAL(mod);
}
