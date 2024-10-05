from flask import Flask, request, jsonify
import requests

app = Flask(__name__)

def rate_limit_check():
    endpoint = request.path
    ip = request.remote_addr

    headers = {
        'endpoint' : endpoint,
        'ip' : ip
    }

    try:
        response = requests.get('http://127.0.0.1:8080/check-limit', headers=headers)

        if response.status_code == 429:
            return jsonify({
                'error' : 'TOO MANY REQUESTS'
            }), 429
        
        if response.status_code == 500:
            return jsonify({
                'error' : 'INTERNAL SERVER ERROR'
            }), 500
    except requests.exceptions.RequestException as e:
        return jsonify({
                'error' : 'Rate limit service unavailable'
            }), 500
    
@app.before_request
def before_request():
    rate_limit_response = rate_limit_check()
    if rate_limit_response:
        return rate_limit_response
    


@app.route("/api/v1/process", methods=['GET'])
def process():
    return jsonify({"success": True})

if __name__ == '__main__':
    app.run(port=3002)