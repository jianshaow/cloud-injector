from flask import Flask, request, render_template

app = Flask(__name__)


@app.route('/index', methods=['GET'])
def callback():
    user_id = request.args['userId']
    return render_template('index.tmpl', user_id=user_id, name="Bob", email="bob@abc.com")


if __name__ == '__main__':
    app.run(debug=True, host="0.0.0.0")
