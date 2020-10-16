#!/usr/bin/python3

import os
import glob
import time
from signal import signal, pause, SIGALRM
from pynats import NATSClient
from json import dumps

base_dir = '/sys/bus/w1/devices/'
device_folder = glob.glob(base_dir + '28*')[0]
device_file = device_folder + '/temperature'


def read_temp(signum, frame):
    with open(device_file, 'r') as f:
        temp_s = f.readline()
        temp = float(temp_s)/1000.0
    cas = time.strftime("%Y-%m-%dT%H:%M:%S")
    data = {'time': cas, 'data': {'temperature': temp}}
    payload = dumps(data)
    try:
        with NATSClient("nats://rpi3:4222") as client:
            client.publish("rpi3", payload=payload.encode())
        pause()
    except ConnectionRefusedError:
        read_temp(14, None)

with open("/tmp/dallas.pid", "w", encoding="UTF-8") as f:
    f.write(str(os.getpid()))
    os.system("chmod 666 /tmp/dallas.pid")
signal(SIGALRM, read_temp)
pause()
