#!/usr/bin/python

# Dummy script used as a stand-in for singal-cli which replies to any JSON
# request with a JSON that has the same ID.

import sys
import json

for line in sys.stdin:
	try:
		data = json.loads(line.strip())
		if 'id' in data: print(json.dumps({'id': data['id']}))
	except:
		print("ERROR")