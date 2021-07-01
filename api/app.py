from flask import Flask
from flask_restx import Api
from flask_jwt_extended import JWTManager

from api.resources.user import api as user_namespace
from api.resources.favorite import api as favorite_namespace
from api.config import app_config

def create_app(config_name):
    app = Flask(__name__)
    app.config.from_object(app_config[config_name])
    jwt = JWTManager(app)
    api = Api(
        app, 
        title="The Drink Almanac REST API",
        description="Manage accounts and add or remove favorited drinks",
    )

    api.add_namespace(user_namespace)
    api.add_namespace(favorite_namespace)
    
    return app
