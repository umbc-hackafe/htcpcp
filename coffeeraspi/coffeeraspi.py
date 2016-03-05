#!env/bin/python3

import argparse
import asyncio
import json
import socket
import websockets

import teensy

async def contact_server():
    async with websockets.connect(server) as sock:
        await sock.send(json.dumps(dict(
            message='Hello',
            name=name,
            id=None, # In theory we would provide a unique ID for each machine, but we only have one...
        )))

        resp = await sock.recv()
        # Handle new response
        print(json.loads(resp))


def main():
    asyncio.get_event_loop().run_until_complete(contact_server())

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='Client for connecting to AWS')
    parser.add_argument('server', help='The server to connect to')
    parser.add_argument(
        '-n', '--name', default=socket.gethostname(),
        help='The name of this client coffee machine'
    )

    args = parser.parse_args()
    server = args.server
    name = args.name

    main()
