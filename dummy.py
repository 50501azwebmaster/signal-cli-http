#!/usr/bin/python

# Dummy script used as a stand-in for singal-cli which replies to any JSON
# request with a JSON that has the same ID.

import sys
import json

print(json.dumps({"method":"receive","params":{"envelope":{"source":"67a13c3e-8d29-2539-ce8e-41129c349d6d"},"data":"stuff"}}))

for line in sys.stdin:
	try:
		data = json.loads(line.strip())
		if 'id' in data: print(json.dumps({'id': data['id']}))
	except:
		print("ERROR")