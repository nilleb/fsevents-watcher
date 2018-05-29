. ./set_python_home.sh
export PKG_CONFIG_PATH=$PKG_CONFIG_PATH:$PYTHON_HOME/lib/pkgconfig
set -eu
go get github.com/fsnotify/fsevents
go build -buildmode=c-shared -o fsevents_watcher.so  # in case of problems, add here -x
