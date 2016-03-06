class DrinkOrder():
    def __init__(self, mug_size, add_ins, name=None):
        self.mug_size = mug_size
        self.add_ins = add_ins
        self.name = name

    @classmethod
    def deserialize(cls, data):
        return DrinkOrder(data['mug_size'],
                data['add_ins'],
                data.get('name', None))
