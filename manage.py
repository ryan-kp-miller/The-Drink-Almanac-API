import os
from api.app import create_app
from api.db import db

config_name = os.getenv('APP_SETTINGS')
app = create_app(config_name)


if __name__ == "__main__":
    app.run(
        port=os.environ.get('PORT', 5000)
    )
