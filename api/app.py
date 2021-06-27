import os
from datetime import timedelta
from flask import Flask
from flask_restx import Api
from flask_cors import CORS
from flask_jwt_extended import JWTManager

from api.resources.user import User, UserRegister, UserLogin
from api.resources.favorite import Favorite
# from config import DB_URI, SECRET_KEY
DB_URI = os.environ.get('DATABASE_URL')
SECRET_KEY = os.environ.get('SECRET_KEY')

# heroku uses the old postgres dialect that is no longer supported by flask-sqlalchemy
# so manually switching to the new one
DB_URI = DB_URI.replace("postgres://", "postgresql://")

def create_app():
    app = Flask(__name__)
    CORS(app)

    # tell SQLAlchemy where to find the database
    app.config['SQLALCHEMY_DATABASE_URI'] = DB_URI

    # turn off flask_sqlalchemy modification tracker so we can use SQLAlchemy's mod tracker, which is better
    app.config['SQLALCHEMY_TRACK_MODIFICATIONS'] = False 
    
    # allows us to blacklist both access and refresh tokens
    app.config['JWT_BLACKLIST_ENABLED'] = True
    app.config['JWT_BLACKLIST_TOKEN_CHECKS'] = ['access', 'refresh']

    # set access tokens to expire in 1 hour
    app.config["JWT_ACCESS_TOKEN_EXPIRES"] = timedelta(hours=1)

    app.secret_key = SECRET_KEY
    jwt = JWTManager(app)  #create /auth endpoint


    api = Api(
        app, 
        prefix="/api", 
        doc="/api", 
        title="The Drink Almanac REST API",
        description="Manage accounts and add or remove favorited drinks"
    )

    api.add_resource(User,         '/user')
    api.add_resource(UserRegister, '/register')
    api.add_resource(Favorite,     '/favorite')
    api.add_resource(UserLogin,     '/login')
    
    return app
