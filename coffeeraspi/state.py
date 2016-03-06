class State():
    class NoMugsException(Exception): pass

    def __init__(self, interface,
            mug_capacity=6, default_mug_fill=True, current_slot=0):
        self.interface = interface
        self.mug_slots = [default_mug_fill] * mug_capacity

        self.current_slot = current_slot

    def enforce(self):
        """Attempts to enforce current known state by performing hardware
        operations."""
        self.interface.mug_seek(self.current_slot)

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
