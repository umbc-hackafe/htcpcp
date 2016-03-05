import serial


class Interface():
    _func_map = {
            'mug_seek':         b'a',
            'mug_get':          b'b',
            'kcup_tray_open':   b'c',
            'kcup_tray_close':  b'd',
            'kcup_load':        b'e',
            'kcup_eject':       b'f',
            'kcup_count':       b'g',
            'brew':             b'h',
            'addin_insert':     b'i'
            }
    def __init__(self, device):
        self.serial = serial.Serial()
        self.serial.port = device

    def __enter__(self):
        self.open()
        return self
    def __exit__(self):
        self.close()
        return False # Don't suppress exceptions

    def open(self):
        self.serial.open()
    def close(self):
        self.serial.close()
