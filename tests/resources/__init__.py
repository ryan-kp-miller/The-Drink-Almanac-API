import os
from pytest import fixture
from api.app import create_app

TEST_USERNAME = "test"
TEST_PASSWORD = "test"

# setting scope to module allows us to use the same test database
# across all tests in this file, so users persist across tests
# this lets us reuse a single JWT access token in later tests
@fixture(scope="module")
def client():
    # create the Flask app and yield the test client
    app = create_app("testing")
    with app.app_context():
        with app.test_client() as client:
            yield client
    
    # remove the test SQLite DB after all tests are run
    os.remove(os.path.join(os.getcwd(), "test_db.db"))
