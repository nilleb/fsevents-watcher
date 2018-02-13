#!/usr/bin/env python
#
# Replacement `MtimeFileWatcher` for App Engine SDK's dev_appserver.py,
# designed for OS X. Improves upon existing file watcher (under OS X) in
# numerous ways:
#
#   - Uses FSEvents API to watch for changes instead of polling. This saves a
#     dramatic amount of CPU, especially in projects with several modules.
#   - Tries to be smarter about which modules reload when files change, only
#     modified module should reload.
#
import os
import time
from os.path import abspath, join
import fsevents_watcher
from ConfigParser import ConfigParser, NoSectionError

# Only watch for changes to .go, .py or .yaml files
WATCHED_EXTENSIONS = set(['.go', '.py', '.yaml'])

def find_upwards(file_name, start_at=os.getcwd()):
    cur_dir = start_at
    while True:
        file_list = os.listdir(cur_dir)
        parent_dir = os.path.dirname(cur_dir)
        if file_name in file_list:
            return cur_dir
        else:
            if cur_dir == parent_dir:
                return None
            else:
                cur_dir = parent_dir

class MtimeFileWatcher(object):
    SUPPORTS_MULTIPLE_DIRECTORIES = True

    def __init__(self, directories, **kwargs):
        self._changes = _changes = []
        # Path to current module
        module_dir = directories[0]

        watched_extensions = WATCHED_EXTENSIONS
        setup_cfg_path = find_upwards("setup.cfg")
        if setup_cfg_path:
            config = ConfigParser()
            try:
                config_value = config.get('appengine:mtime_file_watcher', 'watched_extensions')
            except NoSectionError:
                watched_extensions = WATCHED_EXTENSIONS
            else:
                try:
                    watched_extensions = set(config_value)
                except TypeError:
                    watched_extensions = WATCHED_EXTENSIONS

        # Paths to watch
        paths = [module_dir]

        def callback(path, flags):
            # Get extension
            try:
                ext = os.path.splitext(path)[1]
            except IndexError:
                ext = None

            # Add to changes if we're watching a file with this extension.
            if ext in watched_extensions:
                _changes.append(path)

        fsevents_watcher.schedule(callback, paths)

    def start(self):
        fsevents_watcher.start()

    def changes(self, timeout=None):
        time.sleep(1)
        changed = set(self._changes)
        del self._changes[:]
        return changed

    def quit(self):
        fsevents_watcher.stop()
