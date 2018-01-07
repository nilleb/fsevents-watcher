package main

// #cgo pkg-config: python-2.7
// #define Py_LIMITED_API
// #include <Python.h>
// int PyArg_ParseTuple_ourArgs(PyObject *, PyObject *, PyObject *);
// PyObject* PyObject_CallFunction_ourArgs(PyObject*, char*, long long);
import "C"
import (
    "log"
    "time"

    "github.com/fsnotify/fsevents"
)

var es *fsevents.EventStream
var _callback *C.PyObject

func PyListOfStrings(listObj *C.PyObject) []string {
    numLines := int(C.PyList_Size(listObj))

    if (numLines < 0) { return nil; }

    aList := []string{}
    for i := 0; i < numLines; i++ {
        strObj := C.PyList_GetItem(listObj, C.Py_ssize_t(i));
        aList = append(aList, PyString_AsString(strObj))
    }
    return aList
}

func PyString_AsString(s *C.PyObject) string {
    return C.GoString(C.PyString_AsString(s))
}

//export schedule
func schedule(self, args *C.PyObject) *C.PyObject {
    var argPaths *C.PyObject
    if C.PyArg_ParseTuple_ourArgs(args, &_callback, &argPaths) != 0 {
        return nil
    }
    paths := PyListOfStrings(argPaths)

    path := paths[0]
    dev, err := fsevents.DeviceForPath(path)
    if err != nil {
        log.Fatalf("Failed to retrieve device for path: %v", err)
    }
    es = &fsevents.EventStream{
        Paths:   paths,
        Latency: 500 * time.Millisecond,
        Device:  dev,
        Flags:   fsevents.FileEvents | fsevents.WatchRoot}

    return nil
}

var noteDescription = map[fsevents.EventFlags]string{
    fsevents.MustScanSubDirs: "MustScanSubdirs",
    fsevents.UserDropped:     "UserDropped",
    fsevents.KernelDropped:   "KernelDropped",
    fsevents.EventIDsWrapped: "EventIDsWrapped",
    fsevents.HistoryDone:     "HistoryDone",
    fsevents.RootChanged:     "RootChanged",
    fsevents.Mount:           "Mount",
    fsevents.Unmount:         "Unmount",

    fsevents.ItemCreated:       "Created",
    fsevents.ItemRemoved:       "Removed",
    fsevents.ItemInodeMetaMod:  "InodeMetaMod",
    fsevents.ItemRenamed:       "Renamed",
    fsevents.ItemModified:      "Modified",
    fsevents.ItemFinderInfoMod: "FinderInfoMod",
    fsevents.ItemChangeOwner:   "ChangeOwner",
    fsevents.ItemXattrMod:      "XAttrMod",
    fsevents.ItemIsFile:        "IsFile",
    fsevents.ItemIsDir:         "IsDir",
    fsevents.ItemIsSymlink:     "IsSymLink",
}

func callTheCallback(event fsevents.Event) {
    if (C.PyObject_CallFunction_ourArgs(_callback, "OO", C.CString(event.Path),
                              event.Flags) == nil) {
        log.Printf("the callback has returned nil")
    }
}

func logEvent(event fsevents.Event) {
    note := ""
    for bit, description := range noteDescription {
        if event.Flags&bit == bit {
            note += description + " "
        }
    }
    log.Printf("EventID: %d Path: %s Flags: %s", event.ID, event.Path, note)
}


//export stop
func stop(self *C.PyObject) *C.PyObject {
    es.Stop()
    return nil
}

//export start
func start(self *C.PyObject) *C.PyObject {
    es.Start()
    ec := es.Events
    go func() {
        for msg := range ec {
            for _, event := range msg {
                callTheCallback(event)
                logEvent(event)
            }
        }
    }()
    return nil
}

func main() {}
