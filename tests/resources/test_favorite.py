from json import loads
import pytest
from tests.resources import (
    client, get_auth_header, TEST_CREDENTIALS_PAYLOAD
)


class TestFavorite:
    TEST_DRINK_ID = 11000
    EXPECTED_KEYS = ['id', 'user_id', 'drink_id']

    @pytest.fixture(autouse=True)
    def _set_auth_header(self, client):
        self.auth_header = get_auth_header(client, TEST_CREDENTIALS_PAYLOAD)

    def test_post_correct_args(self, client):
        response = client.post(
            f'/favorite/{self.TEST_DRINK_ID}', 
            headers=self.auth_header
        )

        assert response.status_code == 201
        data = loads(response.data)
        for key in data.keys():
            assert key in self.EXPECTED_KEYS
        assert data['drink_id'] == self.TEST_DRINK_ID
        # first user/favorite created, so their id's should always be 1
        assert data['user_id'] == 1 
        assert data['id'] == 1  

    def test_post_duplicate_favorite(self, client):
        for _ in range(2):
            response = client.post(
                f'/favorite/{self.TEST_DRINK_ID}', 
                headers=self.auth_header
            )

        assert response.status_code == 400
        assert b'User has already favorited this drink' in response.data
    
    def test_post_missing_auth_header(self, client):
        response = client.post(f'/favorite/{self.TEST_DRINK_ID}')
        assert response.status_code == 401
        data = loads(response.data)
        assert 'Missing Authorization Header' in data['msg']
        
    def test_post_invalid_access_token(self, client):
        # access token takend from jwt.io
        bad_auth_header = {
            "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ"
            "zdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM"
            "5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
        }
        response = client.post(
            f"/favorite/{self.TEST_DRINK_ID}",
            headers=bad_auth_header
        )
        assert response.status_code == 422
        
    def test_post_missing_drink_id(self, client):
        response = client.post('/favorite/', headers=self.auth_header)
        assert response.status_code == 404

    def test_post_invalid_drink_id(self, client):
        response = client.post('/favorite/test', headers=self.auth_header)
        assert response.status_code == 404

    def test_get_correct_args(self, client):
        client.post(
            f'/favorite/{self.TEST_DRINK_ID}', 
            headers=self.auth_header
        )
        response = client.get(
            f'/favorite/{self.TEST_DRINK_ID}', 
            headers=self.auth_header
        )

        assert response.status_code == 200
        data = loads(response.data)
        for key in data.keys():
            assert key in self.EXPECTED_KEYS
        assert data['drink_id'] == self.TEST_DRINK_ID
        # first user/favorite created, so their id's should always be 1
        assert data['user_id'] == 1 
        assert data['id'] == 1  
    
    def test_get_missing_auth_header(self, client):
        response = client.get(f'/favorite/{self.TEST_DRINK_ID}')
        assert response.status_code == 401
        data = loads(response.data)
        assert 'Missing Authorization Header' in data['msg']
        
    def test_get_invalid_access_token(self, client):
        # access token takend from jwt.io
        bad_auth_header = {
            "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ"
            "zdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM"
            "5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
        }
        response = client.get(
            f"/favorite/{self.TEST_DRINK_ID}",
            headers=bad_auth_header
        )
        assert response.status_code == 422
        
    def test_get_missing_drink_id(self, client):
        response = client.get('/favorite/', headers=self.auth_header)
        assert response.status_code == 404

    def test_get_invalid_drink_id(self, client):
        response = client.get('/favorite/test', headers=self.auth_header)
        assert response.status_code == 404

    def test_delete_correct_args(self, client):
        client.post(
            f'/favorite/{self.TEST_DRINK_ID}', 
            headers=self.auth_header
        )
        response = client.delete(
            f'/favorite/{self.TEST_DRINK_ID}', 
            headers=self.auth_header
        )

        assert response.status_code == 200
        assert b"Success" in response.data

        response = client.get(
            f'/favorite/{self.TEST_DRINK_ID}', 
            headers=self.auth_header
        )
        assert response.status_code == 404
    
    def test_delete_nonexistent_favorite(self, client):
        response = client.delete(
            '/favorite/1', 
            headers=self.auth_header
        )

        assert response.status_code == 404
        assert b'Favorite not found' in response.data
    
    def test_delete_missing_auth_header(self, client):
        response = client.delete(f'/favorite/{self.TEST_DRINK_ID}')
        assert response.status_code == 401
        data = loads(response.data)
        assert 'Missing Authorization Header' in data['msg']
        
    def test_delete_invalid_access_token(self, client):
        # access token takend from jwt.io
        bad_auth_header = {
            "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ"
            "zdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM"
            "5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
        }
        response = client.delete(
            f"/favorite/{self.TEST_DRINK_ID}",
            headers=bad_auth_header
        )
        assert response.status_code == 422
        
    def test_delete_missing_drink_id(self, client):
        response = client.delete('/favorite/', headers=self.auth_header)
        assert response.status_code == 404

    def test_delete_invalid_drink_id(self, client):
        response = client.delete('/favorite/test', headers=self.auth_header)
        assert response.status_code == 404
