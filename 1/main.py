"""
Module main.py Find duplicate files in a given directory.
"""

import argparse
import os
import hashlib
from collections import defaultdict


def get_file_hash(path, chunk_size=8192):
    """Calculate the hash of a file's contents."""

    hasher = hashlib.sha256()
    with open(path, "rb") as file:
        for chunk in iter(lambda: file.read(chunk_size), b""):
            hasher.update(chunk)
    return hasher.hexdigest()


def find_files_with_different_size(directory):
    """
    Walk through all files in all sub-folders in a directory and
    return dict containing non-unique file sizes with their hash and path.
    """

    # Dictionary to store files needed to be hashed (different size)
    file_sizes = defaultdict(list)

    # Walk through the directory and its sub-folders
    for root, _, files in os.walk(directory):
        for file_name in files:
            path = os.path.join(root, file_name)
            this_size = os.stat(path).st_size
            file_sizes[this_size].append(path)

    # Filter files with different size
    filtered_files = {
        non_unique_size: paths
        for non_unique_size, paths in file_sizes.items()
        if len(paths) > 1
    }

    return filtered_files


def find_duplicate_files(dictionary):
    """Find duplicate files in a directory and its sub-folders."""

    file_hashes = defaultdict(list)
    for _, values in dictionary.items():
        for value in values:
            this_hash = get_file_hash(value)
            file_hashes[this_hash].append(value)

    return file_hashes


# Init CLI argument parser:
argParser = argparse.ArgumentParser()
argParser.add_argument(
    "-d",
    "--directory",
    type=str,
    help="Path to root directory to search for duplicate files.",
)
args = argParser.parse_args()

# Scan for duplicates:
directory_to_scan = args.directory
different_size_files = find_files_with_different_size(directory_to_scan)
duplicate_files = find_duplicate_files(different_size_files)

# Print duplicate files
for file_hash, file_paths in duplicate_files.items():
    print(f"Duplicate files with hash {file_hash}:")
    for file_path in file_paths:
        print(file_path)
