"""
Sample caller for the native fsevents_watcher module
"""
import os
import logging
import fsevents_watcher

def callback(path, flags):
    """
    A sample callback. It will be called whenever an fsevent for the given path
    will be generated.
    """
    logging.info("python: callback called: %s: %s", path, flags)
    print("python: callback called: {}: {}".format(path, flags))

fsevents_watcher.schedule(callback, [os.path.abspath(".")])
fsevents_watcher.start()
raw_input('Press [Enter] to stop listening.\n')
fsevents_watcher.stop()
