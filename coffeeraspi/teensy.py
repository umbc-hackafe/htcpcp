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
    def __init__(self, device, mock=False):
        self.serial = serial.Serial() if not mock else SerialMock()
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

class SerialMock():
    def __init__(self, port = None, *args, **kwargs):
        self._opened = False
        self.is_open = False

        self.port = port
    def __enter__(self):
        self.open()
        return self
    def __exit__(self, exception_type, exception_val, exception_tb):
        self.close()
        return False # Don't suppress exceptions

    def open(self):
        if self._opened:
            raise IOError('Tried to re-open a mocked socket');
        else:
            self._opened = True
            self.is_open = True

    def close(self):
        self.is_open = False

    def read(self, n = 1):
        if self.is_open:
            return bytes([0] * n)
        else:
            raise IOError('Tried to read from a closed mocked socket');

    def write(self, b):
        print('SerialMock \'{}\': {}'.format(self.port, b))
        return len(b)
