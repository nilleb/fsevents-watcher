export PKG_CONFIG_PATH=$PKG_CONFIG_PATH:/System/Library/Frameworks/Python.framework/Versions/2.7/lib/pkgconfig
go get github.com/nilleb/fsevents
go build -buildmode=c-shared -o fsevents_watcher.so
