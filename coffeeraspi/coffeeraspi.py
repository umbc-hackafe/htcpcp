#!env/bin/python3

import argparse
import asyncio
import json
import socket
import websockets
import datetime

import teensy
import messages

def log(message):
    print('{}: {}'.format(
        datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S'),
        message))

async def contact_server(server, name, coffee_queue, reconnect=True):
    while True:
        try:
            async with websockets.connect('ws://' + server + '/ws') as sock:
                await sock.send(json.dumps(dict(
                    message='Hello',
                    name=name,
                    id=None, # In theory we would provide a unique ID for each machine, but we only have one...
                    )))

                resp = await sock.recv()
                log(json.loads(resp))
                log('Connected with server {}'.format(server))

                while True:
                    # Wait for orders.
                    order_msg = json.loads(await sock.recv())
                    order = messages.DrinkOrder.deserialize(order_msg)
                    coffee_queue.put_nowait(order)
                    log('Enqueued order {}'.format(order))
        except (OSError, websockets.exceptions.ConnectionClosed) as e:
            log('Error connecting to server: {}'.format(e))

        log('Lost connection with server, retrying in 5')
        await asyncio.sleep(5)
        if not reconnect:
            return

async def serial_consumer(serial_device_name, coffee_queue, mock=False):
    with teensy.Interface(serial_device_name, mock=mock) as interface:
        while True:
            order = await coffee_queue.get()
            log('Preparing order {}'.format(order))

def main(args):
    loop = asyncio.get_event_loop()
    coffee_queue = asyncio.Queue(loop=loop)
    loop.run_until_complete(asyncio.gather(
        contact_server(args.server, args.name, coffee_queue),
        serial_consumer(args.serial, coffee_queue, mock=args.mock)))
    loop.close()

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='Client for connecting to AWS')
    parser.add_argument('server', help='The server to connect to')
    parser.add_argument(
        '-n', '--name', default=socket.gethostname(),
        help='The name of this client coffee machine'
    )
    parser.add_argument('-s', '--serial', default=None,
            help='The serial device to use, or the first one detected')
    parser.add_argument('-S', '--mock', action='store_true',
            help='Mock the socket device instead of using a real one')

    main(parser.parse_args())
