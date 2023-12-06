# Make LevelDB data: both uncompressed and compressed.
# Usage: python3 make_ldb.py

import json
import os

import plyvel  # pip install plyvel
import snappy  # pip install python-snappy


def main():
    make("./lorem.json", "./uncompressed.ldb", reopen=True)
    make("./lorem.json", "./snappy.ldb", compression="snappy", reopen=True)
    make("./lorem.json", "./log_only.ldb", compression=None)


def make(input_filepath, output_filepath, compression=None, reopen=False):
    if os.path.exists(output_filepath):
        raise FileExistsError(f"The file {output_filepath} already exists.")
    # make a .ldb file and a .log file within
    db = plyvel.DB(output_filepath, compression=compression, create_if_missing=True)
    for key, value in read_json(input_filepath).items():
        db.put(key.encode(), value.encode())
    db.close()
    if reopen:
        # reopen, so a .ldb file is generated within the .ldb directory;
        # otherwise there's a .log file only with the fresh changes.
        db = plyvel.DB(output_filepath, compression=compression)
        db.close()


# Helpers


def compress(value):
    return snappy.compress(value)


def decompress(value):
    return snappy.decompress(value)


def read_json(filepath):
    with open(filepath, "r") as file:
        return json.load(file)


main()
