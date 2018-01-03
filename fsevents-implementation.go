package main

// #cgo pkg-config: python-2.7
// #define Py_LIMITED_API
// #include <Python.h>
// int PyArg_ParseTuple_str(PyObject *, char *);
import "C"
import (
    //"bufio"
    //"io/ioutil"
    "log"
    //"os"
    //"runtime"
    "time"

    "github.com/fsnotify/fsevents"
)

//export schedule
func schedule(self, args *C.PyObject) *C.PyObject {
    var path *C.char
    if C.PyArg_ParseTuple_str(args, path) == 0 {
        return nil
    }
    dev, err := fsevents.DeviceForPath(C.GoString(path))
    if err != nil {
        log.Fatalf("Failed to retrieve device for path: %v", err)
    }
    es := &fsevents.EventStream{
        Paths:   []string{C.GoString(path)},
        Latency: 500 * time.Millisecond,
        Device:  dev,
        Flags:   fsevents.FileEvents | fsevents.WatchRoot}

    es.Start()
    ec := es.Events
    go func() {
        for msg := range ec {
            for _, event := range msg {
                logEvent(event)
            }
        }
    }()
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
func stop(self, args *C.PyObject) *C.PyObject {
    return nil
}

//export loop
func loop(self, args *C.PyObject) *C.PyObject {
    var path *C.char
    if C.PyArg_ParseTuple_str(args, path) == 0 {
        return nil
    }
    return C.PyLong_FromLongLong(-1)
}

func main() {}
