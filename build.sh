. ./set_python_home.sh
export PKG_CONFIG_PATH=$PKG_CONFIG_PATH:$PYTHON_HOME/lib/pkgconfig
set -eu
go build -buildmode=c-shared -o fsevents_watcher.so
