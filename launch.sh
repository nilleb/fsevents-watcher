set -eu
. ./set_python_home.sh
./build.sh
$PYTHON_HOME/bin/python macos_watcher.py
