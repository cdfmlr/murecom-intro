import json
import re
import sys

PATTERN = r'word: (.*?)\n'

if not len(sys.argv):
    print("usage: seen_words_from_last_run.py input_file")
    exit(1)
file = sys.argv[1]

with open(file) as f:
    finds = re.findall(PATTERN, f.read())

print(json.dumps(finds))
