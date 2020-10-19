import firebase_admin
from firebase_admin import firestore
from pynats import NATSClient, NATSMessage
from pynats.exceptions import NATSReadSocketError
from json import loads
import time

if __name__ == "__main__":
    fireApp = firebase_admin.initialize_app()
    db = firestore.client()

    def fireData(msg: NATSMessage):
        jsonData = msg.payload.decode("utf-8")
        data = loads(jsonData)
        parseTime = data['time'].split('T')
        day = db.collection('room').document(parseTime[0])
        actual = db.collection('room').document('now')
        day.set({parseTime[1]: data['data']}, merge=True)
        actual.set(data['data'])
        check = day.get()
        try:
            check.get('timestamp')
        except KeyError:
            day.set({'timestamp': int(time.mktime(time.localtime()))})


if __name__ == "__main__":
    with NATSClient("nats://rpi3:4222") as ns:
        sub = ns.subscribe("rpi3", callback=fireData)
        ns.wait()
