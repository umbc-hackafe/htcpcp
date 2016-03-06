class State():
    class NoMugsException(Exception): pass

    def __init__(self, interface,
            mug_capacity=6, default_mug_fill=True, current_slot=0):
        self.interface = interface
        self.mug_slots = [default_mug_fill] * mug_capacity

        self.current_slot = current_slot

        self.kcup_loaded = False
        self.kcup_tray_open = False

    def enforce(self):
        """Attempts to enforce current known state by performing hardware
        operations."""
        self.interface.mug_seek(self.current_slot)

        # Ensure that there is no kcup loaded
        if not self.kcup_loaded:
            self.interface.kcup_eject()

        if self.kcup_tray_open:
            self.interface.kcup_tray_open()
        else:
            self.interface.kcup_tray_close()

    def open_mugs(self, except_none=False):
        """Return a list of indices of mugs which are available. If except_none,
        then a NoMugsException will be raised if there are no available mugs."""
        indices = [index for index, available in enumerate(self.mug_slots) if available]
        if except_none and len(indices) == 0:
            raise NoMugsException("no mugs available")

        return indices

    def ready_mug(self):
        # Get the index of the first open mug. Will raise NoMugsException if
        # none are available.
        mug_index = self.open_mugs(except_none=True)[0]

        # Seek to that slot and record that we are there.
        self.interface.mug_seek(mug_index)
        self.current_slot = mug_index

    def ready_kcup(self):
        # Load the cup (Teensy ensures lid is open and ejected)
        self.interface.kcup_load()

        self.kcup_loaded = True

        # Close the tray.
        self.interface.kcup_tray_close()

    def brew(self):
        self.interface.brew()
