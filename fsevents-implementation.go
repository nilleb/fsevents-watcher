package main

// #cgo pkg-config: python-2.7
// #define Py_LIMITED_API
// #include <Python.h>
// int PyArg_ParseTuple_ourArgs(PyObject *, PyObject **, PyObject **);
// PyObject* PyArg_BuildNone();
// PyObject* PyArg_BuildCallbackArguments(char* path, char* flags);
// PyObject* PyObject_CallFunction_ourArgs(PyObject*, char*, char*);
import "C"
import (
	"log"
	"time"

	"github.com/nilleb/fsevents"
)

var es *fsevents.EventStream
var _callback *C.PyObject

// PyListOfStrings convert a list of objects to a list of strings
func PyListOfStrings(listObj *C.PyObject) []string {
	numLines := int(C.PyList_Size(listObj))

	log.Printf("number of lines: %d", numLines)
	if numLines < 0 {
		return nil
	}

	aList := []string{}
	for i := 0; i < numLines; i++ {
		strObj := C.PyList_GetItem(listObj, C.Py_ssize_t(i))
		aList = append(aList, PyString_AsString(strObj))
	}
	return aList
}

// PyString_AsString convert a python string to a Go string
func PyString_AsString(s *C.PyObject) string {
	return C.GoString(C.PyString_AsString(s))
}

//export schedule
func schedule(self, args *C.PyObject) *C.PyObject {
	var argPaths *C.PyObject
	success := C.PyArg_ParseTuple_ourArgs(args, &_callback, &argPaths)

	log.Printf("args, %s", C.PyObject_Repr(args))
	log.Printf("converted: %v, %s", _callback, C.PyObject_Repr(argPaths))

	if success == 0 {
		return nil
	}
	paths := PyListOfStrings(argPaths)
	if paths == nil {
		log.Fatal("Sorry, you should pass a slist of paths as second argument.")
		return nil
	}

	path := paths[0]
	log.Printf("Setting up an eventstream for %s", path)
	dev, err := fsevents.DeviceForPath(path)
	if err != nil {
		log.Fatalf("Failed to retrieve device for path: %v", err)
	}
	es = &fsevents.EventStream{
		Paths:   paths,
		Latency: 500 * time.Millisecond,
		Device:  dev,
		Flags:   fsevents.FileEvents | fsevents.WatchRoot}

	return C.PyArg_BuildNone()
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
	note := createNote(event)
	// args = C.PyArg_BuildCallbackArguments(C.CString(event.Path), C.CString(note))
	C.PyObject_CallFunction_ourArgs(_callback, C.CString(event.Path), C.CString(note))
}

func createNote(event fsevents.Event) string {
	note := ""
	for bit, description := range noteDescription {
		if event.Flags&bit == bit {
			note += description + " "
		}
	}
	return note
}

func logEvent(event fsevents.Event) {
	note := createNote(event)
	log.Printf("EventID: %d Path: %s Flags: %s", event.ID, event.Path, note)
}

//export stop
func stop(self *C.PyObject) *C.PyObject {
	es.Stop()
	return C.PyArg_BuildNone()
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
	return C.PyArg_BuildNone()
}

func main() {}
