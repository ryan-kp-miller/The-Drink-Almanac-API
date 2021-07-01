import json
import pytest
from tests.resources import client, TEST_USERNAME, TEST_PASSWORD

# initializing global access_token variable 
# to be set in TestUserLogin and used in TestUser
access_token = ""


class TestUserRegister:
    def test_register_correct_args(self, client):
        response = client.post(
            "/user/register",
            json={
                "username": TEST_USERNAME,
                "password": TEST_PASSWORD
            }
        )
        assert response.status_code == 201

        data = json.loads(response.data)
        actual_keys = data.keys()
        expected_keys = ["id", "username", "favorites"]
        assert len(actual_keys) == len(expected_keys)
        for key in expected_keys:
            assert key in actual_keys 
        
        assert data["username"] == TEST_USERNAME
        assert isinstance(data['id'], int)
        assert isinstance(data['favorites'], list)
        assert len(data['favorites']) == 0

    def test_register_existing_username(self, client):
        # user created in previous test
        response = client.post(
            "/user/register",
            json={
                "username": TEST_USERNAME,
                "password": TEST_PASSWORD
            }
        )
        assert response.status_code == 400
        data = json.loads(response.data)
        assert data['message'] == f'A user with the username {TEST_USERNAME} already exists'

    def test_register_missing_args(self, client):
        # missing username
        response = client.post(
            'user/register',
            json={
                'username': TEST_USERNAME
            }
        )
        assert response.status_code == 400
        data = json.loads(response.data)
        assert "password" in data['errors'].keys()
        assert ("Missing required parameter in the JSON body" in 
            data['errors']['password'])
        
        # missing password
        response = client.post(
            'user/register',
            json={
                'password': TEST_PASSWORD
            }
        )
        assert response.status_code == 400
        data = json.loads(response.data)
        assert "username" in data['errors'].keys()
        assert ("Missing required parameter in the JSON body" in 
            data['errors']['username'])


class TestUserLogin:
    def test_login_correct_args(self, client):
        response = client.post(
            "/user/login",
            json={
                "username": TEST_USERNAME,
                "password": TEST_PASSWORD
            }
        )
        assert response.status_code == 200
        data = json.loads(response.data)
        actual_keys = data.keys()
        expected_keys = ["access_token", "refresh_token"]
        assert len(actual_keys) == len(expected_keys)
        for key in actual_keys:
            assert key in expected_keys
            assert len(data[key].split(".")) == 3 # very basic check for JWT format

        global access_token
        access_token = data['access_token']

    def test_login_nonexistent_username(self, client):
        # user created in previous test
        bad_username = "asdf"
        response = client.post(
            "/user/login",
            json={
                "username": bad_username,
                "password": TEST_PASSWORD
            }
        )
        assert response.status_code == 404
        data = json.loads(response.data)
        assert data['message'] == f'User with the username {bad_username} not found'

    def test_login_missing_args(self, client):
        # missing username
        response = client.post(
            'user/login',
            json={
                'username': TEST_USERNAME
            }
        )
        assert response.status_code == 400
        data = json.loads(response.data)
        assert "password" in data['errors'].keys()
        assert ("Missing required parameter in the JSON body" in 
            data['errors']['password'])
        
        # missing password
        response = client.post(
            'user/login',
            json={
                'password': TEST_PASSWORD
            }
        )
        assert response.status_code == 400
        data = json.loads(response.data)
        assert "username" in data['errors'].keys()
        assert ("Missing required parameter in the JSON body" in 
            data['errors']['username'])


class TestUser:
    pass

