#!/usr/bin/python3

import os
import glob
import time
from signal import signal,pause,SIGALRM
from pynats import NATSClient

base_dir = '/sys/bus/w1/devices/'
device_folder = glob.glob(base_dir + '28*')[0]
device_file = device_folder + '/w1_slave'

def read_temp_raw():
    with open(device_file, 'r') as f:
        f = open(device_file, 'r')
        lines = f.readlines()
    return lines

def read_temp(signum, frame):
    lines = read_temp_raw()
    while lines[0].strip()[-3:] != 'YES':
        time.sleep(0.2)
        lines = read_temp_raw()
    equals_pos = lines[1].find('t=')
    if equals_pos != -1:
        temp_string = lines[1][equals_pos+2:]
        temp_c = float(temp_string) / 1000.0
        temp = str(temp_c)

    cas = time.strftime("%Y-%m-%d %H:%M:%S")
    try:
        with NATSClient("nats://rpi3:4222") as client:
            payload = '{"time":"' + cas + '","temperature":' + temp + '}'
            client.publish("rpi3", payload=payload.encode())
        pause()
    except ConnectionRefusedError:
        read_temp(14, None)

with open("/tmp/dallas.pid", "w", encoding="UTF-8") as f:
    f.write(str(os.getpid()))
    os.system("chmod 666 /tmp/dallas.pid")
signal(SIGALRM, read_temp)
pause()
