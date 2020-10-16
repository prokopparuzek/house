import firebase_admin
from firebase_admin import firestore
from pynats import NATSClient, NATSMessage
from pynats.exceptions import NATSReadSocketError
from json import loads
from time import sleep

if __name__ == "__main__":
    fireApp = firebase_admin.initialize_app()
    db = firestore.client()

    def fireData(msg: NATSMessage):
        jsonData = msg.payload.decode("utf-8")
        data = loads(jsonData)
        parseTime = data['time'].split('T')
        day = db.collection(parseTime[0]).document(parseTime[1])
        actual = db.collection('now').document('current')
        day.set(data['data'],)
        actual.set(data['data'])
        print(data)


def main():
    try:
        with NATSClient("nats://rpi3:4222") as ns:
            sub = ns.subscribe("rpi3", callback=fireData)
            ns.wait()
    except:
        main()

if __name__ == "__main__":
    main()
