# fsevents_watcher
A python extension to gather filesystem events generated by the operating system (macOS) for a specific path.
Have a look at `macos_watcher.py` for an example.

## requirements
- homebrew (used to pull other requirements)
- xcode devutils
- golang (used to compile the go code)
- pkg-config (used to generate the header files)

## bootstrap
```
xcode-select --install  # installs xcode devutils
# make sure you have set your GOPATH
brew install pkg-config golang
mkdir -p $GOPATH/src/github.com/nilleb
cd $GOPATH/src/github.com/nilleb
git clone https://github.com/nilleb/fsevents-watcher
cd fsevents-watcher
# which python do you use? if a homebrew one, you're ready to go.
# otherwise, edit set-python-home.sh
./build.sh
# now you can launch the example (it just notifies you about the events in the current folder) by typing
./launch.sh
# if you want to use the included mtime_file_watcher for your AppEngine dev_appserver.py
sudo python replace_mtime_file_watcher.py replace
# if you want to restore the original mtime_file_watcher.py,
sudo python replace_mtime_file_watcher.py restore
```
now you can start `dev_appserver.py` as usual

## notes
has been tested on
- Darwin 17.4.0 Darwin Kernel Version 17.4.0: Sun Dec 17 09:19:54 PST 2017; root:xnu-4570.41.2~1/RELEASE_X86_64 x86_64 i386 MacBookPro12,1 Darwin (homebrew python 2.7.13)
- Darwin 17.4.0 Darwin Kernel Version 17.4.0: Sun Dec 17 09:19:54 PST 2017; root:xnu-4570.41.2~1/RELEASE_X86_64 x86_64 i386 MacBookPro13,1 Darwin (System Default Python Interpreter - 2.7.10 as of writing)
- Darwin 16.7.0 Darwin Kernel Version 16.7.0: Thu Jan 11 22:59:40 PST 2018; root:xnu-3789.73.8~1/RELEASE_X86_64 x86_64 i386 MacBookPro12,1 Darwin

## troubleshooting
### python version mismatch
```
Fatal Python error: PyThreadState_Get: no current thread
./launch.sh: line 4: 52141 Abort trap: 6           $PYTHON_HOME/bin/python macos_watcher.py
```
Execute a `otool -L fsevents_watcher.so` and verify that the python path is the one of the python executable you are using to launch the code.

### gcloud components update
If you update the gcloud components, you shall re-replace the `mtime_file_watcher.py` again..