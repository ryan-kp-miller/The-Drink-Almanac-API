from flask import Flask
from flask_restful import Api

from api.resources.user import (
    User, UserRegister 
)

def create_app():
    app = Flask(__name__)
    
    # tell SQLAlchemy where to find the database
    app.config['SQLALCHEMY_DATABASE_URI'] = 'sqlite:///data.db'

    # turn off flask_sqlalchemy modification tracker so we can use SQLAlchemy's mod tracker, which is better
    app.config['SQLALCHEMY_TRACK_MODIFICATIONS'] = False 
    
    api = Api(app)
    api.add_resource(User,         '/user')
    api.add_resource(UserRegister, '/register')
    
    return app
