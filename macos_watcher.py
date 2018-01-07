"""
Sample caller for the native fsevents_watcher module
"""
import os
import fsevents_watcher

def callback(path, flags):
    """
    A sample callback. It will be called whenever an fsevent for the given path
    will be generated.
    """
    logging.info(path)
    logging.info(flags)

fsevents_watcher.schedule(callback, os.path.abspath("."))
fsevents_watcher.start()
raw_input('Press [Enter] to stop listening.')
fsevents_watcher.stop()
