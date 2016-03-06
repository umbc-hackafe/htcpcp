import threading
import requests
import pyalexa
import flask
import json

BASE_RESPONSE = {
    "version": "1",
    "response": {}
}

with open(".keys") as f:
    keys = json.load(f)

APP_ID = keys["app_id"]

api = flask.Flask(__name__)

skill = pyalexa.Skill(app_id=APP_ID)

my_drink_id = -1

def do_brew(drink, creams=0, sugar=0):
    global my_drink_id
    data = {
        "size": 6,
        "k_cup": "yes please" if drink == "coffee" else "",
        "tea_bag": "yes please" if drink == "tea" else "",
        "sugar": int(sugar) if sugar else 0,
        "creamer": int(creams) if creams else 0,
    }

    if my_drink_id != -1:
        data["id"] = my_drink_id

    res = requests.post("http://localhost/api/update/drink", json=data)
    res.raise_for_status()

    if "id" in res.json():
        my_drink_id = res.json()["id"]

    requests.post("http://localhost/api/brew/{}/1".format(my_drink_id))

@skill.launch
def launch(request):
    return request.response(end=True, speech="Welcome to Breakfast Time! You can ask me to make coffee, with or without sugar and cream.")

@skill.end
def end(request):
    return request.response(end=True, speech="Thanks for using Breakfast Time!")

@skill.intent("Brew")
def coffee(request):
    creams = request.data().get("Cream", 0)
    sugars = request.data().get("Sugar", 0)
    drink = request.data().get("Drink", "coffee")

    if creams == "?":
        return request.response(end=False, speech="You want how much cream!?")

    if sugars == "?":
        return request.response(end=False, speech="You want how much sugar!?")

    threading.Thread(target=do_brew, args=(drink, creams, sugars)).start()

    return request.response(end=True, speech="Okay! One {} coming right up!".format(drink))

api.add_url_rule('/', 'pyalexa', skill.flask_target, methods=['POST'])
api.run('0.0.0.0', port=8081, debug=True)
