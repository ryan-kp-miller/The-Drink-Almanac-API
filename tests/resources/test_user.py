from json import loads
import pytest
from tests.resources import (
    client, get_auth_header, TEST_USERNAME, TEST_PASSWORD, 
    TEST_CREDENTIALS_PAYLOAD
)


class TestUserRegister:
    def test_post_correct_args(self, client):
        response = client.post(
            "/user/register",
            json=TEST_CREDENTIALS_PAYLOAD
        )
        assert response.status_code == 201

        data = loads(response.data)
        actual_keys = data.keys()
        expected_keys = ["id", "username", "favorites"]
        assert len(actual_keys) == len(expected_keys)
        for key in expected_keys:
            assert key in actual_keys 
        
        assert data["username"] == TEST_USERNAME
        assert isinstance(data['id'], int)
        assert isinstance(data['favorites'], list)
        assert len(data['favorites']) == 0

    def test_post_existing_username(self, client):
        # user created in previous test
        for _ in range(2):
            response = client.post(
                "/user/register",
                json=TEST_CREDENTIALS_PAYLOAD
            )
        assert response.status_code == 400
        data = loads(response.data)
        assert data['message'] == f'A user with the username {TEST_USERNAME} already exists'

    def test_post_missing_args(self, client):
        # missing username
        response = client.post(
            'user/register',
            json={
                'username': TEST_USERNAME
            }
        )
        assert response.status_code == 400
        data = loads(response.data)
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
        data = loads(response.data)
        assert "username" in data['errors'].keys()
        assert ("Missing required parameter in the JSON body" in 
            data['errors']['username'])


class TestUserLogin:
    @pytest.fixture(autouse=True)
    def _set_auth_header(self, client):
        self.auth_header = get_auth_header(client, TEST_CREDENTIALS_PAYLOAD)

    def test_post_correct_args(self, client):
        response = client.post(
            "/user/login",
            json=TEST_CREDENTIALS_PAYLOAD
        )
        assert response.status_code == 200
        data = loads(response.data)
        actual_keys = data.keys()
        expected_keys = ["access_token", "refresh_token"]
        assert len(actual_keys) == len(expected_keys)
        for key in actual_keys:
            assert key in expected_keys
            assert len(data[key].split(".")) == 3 # very basic check for JWT format

        global access_token
        access_token = data['access_token']

    def test_post_nonexistent_username(self, client):
        # user created in previous test
        bad_username = "asdf"
        response = client.post(
            "/user/login",
            json={
                'username': bad_username,
                'password': TEST_PASSWORD
            }
        )
        assert response.status_code == 404
        data = loads(response.data)
        assert data['message'] == f'User with the username {bad_username} not found'

    def test_post_missing_args(self, client):
        # missing username
        response = client.post(
            'user/login',
            json={
                'username': TEST_USERNAME
            }
        )
        assert response.status_code == 400
        data = loads(response.data)
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
        data = loads(response.data)
        assert "username" in data['errors'].keys()
        assert ("Missing required parameter in the JSON body" in 
            data['errors']['username'])


class TestUser:
    @pytest.fixture(autouse=True)
    def _set_auth_header(self, client):
        self.auth_header = get_auth_header(client, TEST_CREDENTIALS_PAYLOAD)

    def test_get_correct_args(self, client):
        response = client.get(
            "/user",
            headers=self.auth_header
        )
        assert response.status_code == 200
        data = loads(response.data)
        actual_keys = data.keys()
        expected_keys = ["id", "username", "favorites"]
        assert len(actual_keys) == len(expected_keys)
        for key in expected_keys:
            assert key in actual_keys 

        assert data["username"] == TEST_USERNAME
        assert isinstance(data['id'], int)
        assert isinstance(data['favorites'], list)
        assert len(data['favorites']) == 0

    def test_get_missing_auth_header(self, client):
        response = client.get('user')
        assert response.status_code == 401
        data = loads(response.data)
        assert 'Missing Authorization Header' in data['msg']
        
    def test_get_incorrect_access_token(self, client):
        # access token takend from jwt.io
        bad_auth_header = {
            "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ"
            "zdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM"
            "5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
        }
        response = client.get(
            "/user",
            headers=bad_auth_header
        )
        assert response.status_code == 422
        
    def test_delete_nonexistent_username(self, client):
        # user created in previous test
        bad_username = "asdf"
        response = client.delete(
            "/user",
            json={
                'username': bad_username,
                'password': TEST_PASSWORD
            }
        )
        assert response.status_code == 404
        data = loads(response.data)
        assert data['message'] == f'User with the username {bad_username} not found'

    def test_delete_missing_args(self, client):
        # missing username
        response = client.delete(
            'user',
            json={
                'username': TEST_USERNAME
            }
        )
        assert response.status_code == 400
        data = loads(response.data)
        assert "password" in data['errors'].keys()
        assert ("Missing required parameter in the JSON body" in 
            data['errors']['password'])
        
        # missing password
        response = client.delete(
            'user',
            json={
                'password': TEST_PASSWORD
            }
        )
        assert response.status_code == 400
        data = loads(response.data)
        assert "username" in data['errors'].keys()
        assert ("Missing required parameter in the JSON body" in 
            data['errors']['username'])

    def test_delete_correct_args(self, client):
        response = client.delete(
            "/user",
            json=TEST_CREDENTIALS_PAYLOAD,
            headers=self.auth_header
        )
        assert response.status_code == 200

