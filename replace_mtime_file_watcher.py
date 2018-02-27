from __future__ import print_function
import argparse
import glob
from datetime import datetime
import hashlib
import shutil
import os.path
from distutils.spawn import find_executable

reference_file = "mtime_file_watcher_replacement.py"

def md5(fname):
    hash_md5 = hashlib.md5()
    with open(fname, "rb") as f:
        for chunk in iter(lambda: f.read(4096), b""):
            hash_md5.update(chunk)
    return hash_md5.hexdigest()


def detect_mtime_file_watcher_path():
    gcloud_path = find_executable("gcloud")
    cloud_sdk_path = os.path.abspath(os.path.join(gcloud_path, "../.."))
    relative = "platform/google_appengine/google/appengine/tools/devappserver2/mtime_file_watcher.py"
    mtime_file_watcher = "{}/{}".format(cloud_sdk_path, relative)
    if not os.path.isfile(mtime_file_watcher):
        print("sorry, unable to detect the path of mtime_file_watcher.py")
        return None
    return mtime_file_watcher

def replace(skip_backup=False):
    mtime_file_watcher = detect_mtime_file_watcher_path()
    if not mtime_file_watcher:
        exit()
    native_module = "fsevents_watcher.so"
    if not os.path.exists(native_module):
        print("sorry, I am unable to find the native module ({}). have you built it?".format(native_module))
        exit()

    reference_checksum = md5(reference_file)
    in_place_checksum = md5(mtime_file_watcher)

    native_module_destination = os.path.abspath(
        os.path.join(mtime_file_watcher, '../{}'.format(native_module)))
    native_module_reference_checksum = md5(native_module)
    native_module_in_place_checksum = md5(native_module_destination)

    if reference_checksum != in_place_checksum or native_module_reference_checksum != native_module_in_place_checksum:
        ts = datetime.now().isoformat().replace('-', '').replace(':', '').replace('.', '')
        if not skip_backup:
            shutil.copy(mtime_file_watcher, "mtime_file_watcher_backup_{}.py".format(ts))
        shutil.copy(reference_file, mtime_file_watcher)
        shutil.copy(native_module, native_module_destination)
        print("The replacement mtime_file_watcher.py has been copied.")
    else:
        print("Looks like the replacement is already in place.")

def restore():
    try:
        backup_file = sorted(glob.glob("mtime_file_watcher_backup_*.py"))[0]
    except IndexError:
        print('Unable to find a backup file.')
        exit()

    mtime_file_watcher = detect_mtime_file_watcher_path()
    if not mtime_file_watcher:
        exit()

    shutil.copy(backup_file, mtime_file_watcher)
    print('The file {} has been restored.')

parser = argparse.ArgumentParser(description='Replace and restore the AppEngine mtime_file_watcher.')
parser.add_argument('action', help='what to do: `replace` or `restore`')
parser.add_argument(
    '--skip-backup', nargs='?', help='should we backup the existing file?')

ns = parser.parse_args()
if ns.action == 'replace':
    replace(skip_backup=ns.skip_backup)
elif ns.action == 'restore':
    restore()
else:
    print("sorry, I haven't understood your command! ({})".format(ns.action))
    parser.print_help()
