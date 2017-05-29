import subprocess
from flask import Flask, jsonify, request

app = Flask(__name__)


@app.route('/')
def main_route():
    return 'Forbidden to route.', 403

@app.route('/api')
def main_api_route():
    return 'Forbidden to route.', 403

@app.route('/api/v1/')
def get_api_v1():
    return 'Forbidden to route.', 403

@app.route('/api/v1/proxy/users',  methods=['GET'] )
def proxy_users_quantity():
    result = subprocess.Popen(['./check_proxy_users.sh'], stdout=subprocess.PIPE)
    q_users = result.stdout.readlines(-1)[0]
    return str(q_users), 201


if __name__ == '__main__':
    app.run()
