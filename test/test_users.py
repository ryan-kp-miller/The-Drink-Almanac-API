import requests

API_BASE_URL = "http://127.0.0.1:5000/"

# register a new user
payload = {
    'username': "test",
    'password': "test",
}

response = requests.post(API_BASE_URL + "register", json=payload)
data = response.json()

print(f"Response Status Code: {response.status_code}")
if 'message' in data:
    print(f"Server message: {data['message']}")