# Make LevelDB data: both uncompressed and compressed.
# Usage: python3 make_ldb.py

import io
import json
import os

import plyvel  # pip install plyvel
import snappy  # pip install python-snappy


def main():
    make("./lorem.json", "./uncompressed.ldb", compression=None, reopen=True)
    make("./lorem.json", "./snappy.ldb", compression="snappy", reopen=True)
    make(make_sample_json(), "./repeats.ldb", rounds=2, reopen=True)
    make("./lorem.json", "./log_only.ldb", compression=None)


def make(input_filepath, output_filepath, compression="snappy", reopen=False, rounds=1):
    if os.path.exists(output_filepath):
        raise FileExistsError(f"The file {output_filepath} already exists.")
    # make a .ldb file and a .log file within
    db = plyvel.DB(output_filepath, compression=compression, create_if_missing=True)
    data = read_json(input_filepath)
    for i in range(rounds):
        for key, value in data.items():
            db.put(key.encode(), value.encode())
    db.close()
    if reopen:
        # reopen, so a .ldb file is generated within the .ldb directory;
        # otherwise there's a .log file only with the fresh changes.
        db = plyvel.DB(output_filepath, compression=compression)
        db.close()


def make_sample_json():
    # in table data blocks, this splits the shared key
    # inside the sequence_number of the internal key;
    # see `readInternalKey` in leveldb_table.go for details.
    data = {}
    for i in range(0x100):
        data[f"lorem.{chr(i)}"] = "ipsum"
    result = io.StringIO()
    json.dump(data, result)
    result.seek(0)
    return result


# Helpers


def compress(value):
    return snappy.compress(value)


def decompress(value):
    return snappy.decompress(value)


def read_json(filepath_or_buffer):
    if hasattr(filepath_or_buffer, "read"):
        return json.load(filepath_or_buffer)
    with open(filepath_or_buffer, "r") as file:
        return json.load(file)


main()
