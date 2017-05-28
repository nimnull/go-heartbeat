import os
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
    test = os.system("netstat -npt | grep ':3000' | grep ESTABLISHED | awk '{ print $5 }' | awk -F ':' '{ print $1 }' | sort -u | wc -l")
    return str(test), 201


if __name__ == '__main__':
    app.run()
