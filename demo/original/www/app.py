from flask import Flask, request, render_template
import json
import sys

app = Flask(__name__)


@app.route('/', methods=['GET'])
def callback():
    user_id = request.args.get("userId")
    if user_id:
        user = user_info.get(user_id)
        if user:
            return render_template('index.tmpl', user=user)
        else:
            return render_template('index.tmpl', error="no user with id: "+user_id)
    else:
        return render_template('index.tmpl')


def load_user_info(data_file):
    with open(data_file, 'r') as user_data:
        global user_info
        user_info = json.load(user_data)


if __name__ == '__main__':
    data_file = (len(sys.argv) == 2 and sys.argv[1] or "/data/users.json")
    load_user_info(data_file)
    app.run(debug=True, host="0.0.0.0")
