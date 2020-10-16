import firebase_admin
from firebase_admin import firestore
from pynats import NATSClient, NATSMessage
from json import loads

if __name__ == "__main__":
    fireApp = firebase_admin.initialize_app()
    db = firestore.client()

    def fireData(msg: NATSMessage):
        jsonData = msg.payload.decode("utf-8")
        data = loads(jsonData)
        parseTime = data['time'].split()
        day = db.collection('room1').document(parseTime[0])
        day.set({parseTime[1]: {'temperature': data['temperature']}}, merge=True)
        print(data)

    with NATSClient("nats://rpi3:4222") as ns:
        sub = ns.subscribe("rpi3", callback=fireData)
        ns.wait()
