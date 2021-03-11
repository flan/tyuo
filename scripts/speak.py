#!/usr/bin/env python3
import requests
import sys

r = requests.post('http://localhost:48100/speak',
    json={
        "ContextId": sys.argv[1],
        "Input": ' '.join(sys.argv[2:]),
    },
    timeout=10.0,
)
print(r.status_code)
print(r.text)
