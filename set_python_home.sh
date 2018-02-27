PYTHON_EXECUTABLE=`which python`
if [ "$PYTHON_EXECUTABLE" == "/usr/bin/python" ]; then
    # macOS System python
    export PYTHON_HOME=/System/Library/Frameworks/Python.framework/Versions/2.7
else
    # homebrew python
    export PYTHON_HOME=/usr/local
fi
