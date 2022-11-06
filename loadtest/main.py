import json

bullet = {
    'tag': 'http_AckNotif',
    'uri': '/api/v1/notification/ack/1',
    'method': 'PUT',
    'headers': {"Accept": "application/json", "Content-Type": "application/json", "Content-Length": "36"},
    'host': 'act-device-api'
}

with open('loadtest/ammo.txt', 'w') as outfile:
    json.dump(bullet, outfile)