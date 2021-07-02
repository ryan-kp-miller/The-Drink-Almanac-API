import os
from json import loads
from pytest import fixture
from api.app import create_app

TEST_USERNAME = "test"
TEST_PASSWORD = "test"

TEST_CREDENTIALS_PAYLOAD = {
    "username": TEST_USERNAME,
    "password": TEST_PASSWORD
}

# setting scope to module allows us to use the same test database
# across all tests in this file, so users persist across tests
# this lets us reuse a single JWT access token in later tests
@fixture()
def client():
    # create the Flask app and yield the test client
    app = create_app("testing")
    with app.app_context():
        with app.test_client() as client:
            yield client
    
    # remove the test SQLite DB after all tests are run
    os.remove(os.path.join(os.getcwd(), "test_db.db"))


def get_auth_header(client, credentials):
    """
    create a new account and return an authorization header dict
    based on the credentials given

    Args:
        client (object): flask test client
        credentials (dict): contains the username and password 
            for the new account
    
    Returns:
        dict with the format: {"Authorization": "Bearer <access_token>"} 
    """
    response = client.post('user/register', json=credentials)
    response = client.post('user/login', json=credentials)
    return {"Authorization": f"Bearer {loads(response.data)['access_token']}"}

