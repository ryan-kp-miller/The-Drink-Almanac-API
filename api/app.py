from flask import Flask
from flask_restx import Api
from flask_jwt_extended import JWTManager

from api.db import db
from api.config import app_config
from api.resources.user import api as user_namespace
from api.resources.favorite import api as favorite_namespace

def create_app(config_name):
    app = Flask(__name__)
    app.config.from_object(app_config[config_name])
    jwt = JWTManager(app)
    db.init_app(app)

    # before any request to the API, this function will be called 
    # and will create the data.db file and all the tables within the db (unless they already exist)
    @app.before_first_request
    def create_tables():
        db.create_all()

    api = Api(
        app, 
        title="The Drink Almanac REST API",
        description="Manage accounts and add or remove favorited drinks",
    )

    api.add_namespace(user_namespace)
    api.add_namespace(favorite_namespace)
    
    return app
