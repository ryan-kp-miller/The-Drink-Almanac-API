import requests

API_BASE_URL = "http://127.0.0.1:5000/api/"

def print_response(response):
    print(f"Response Status Code: {response.status_code}")
    if response.status_code != 500:
        data = response.json()
        if 'message' in data:
            print(f"Server message: {data['message']}")
        else:
            print("Server data keys: ", data.keys())


user_payload = {
    'username': "test",
    'password': "test",
}

# register a new user if the test user doesn't already exist
response = requests.post(API_BASE_URL + "register", json=user_payload)
print_response(response)
if response.status_code == 201:
    user_id = response.json()['id']

# log in as the user
response = requests.post(API_BASE_URL + "login", json=user_payload)
print_response(response)
access_token = response.json()['access_token']
auth_headers = {'Authorization': f'Bearer {access_token}'}

# get the user's favorite drinks
response = requests.get(API_BASE_URL + "user", headers=auth_headers)
print_response(response)

payload_list = [
    {'drink_id': 11007},
    {'drink_id': 11001},
]

# check if the user already favorited these drinks
# if not, favorite the drink
# then delete that favorite
for payload in payload_list:
    response = requests.get(
        API_BASE_URL + "favorite", 
        json=payload,
        headers=auth_headers
    )
    print_response(response)

    if response.status_code == 404:
        response = requests.post(
            API_BASE_URL + "favorite", 
            json=payload,
            headers=auth_headers
        )
        print_response(response)

    response = requests.delete(
        API_BASE_URL + "favorite", 
        json=payload,
        headers=auth_headers
    )
    print_response(response)
    


# delete the test user
# response = requests.delete(API_BASE_URL + "user", json=user_payload)
# print_response(response)

