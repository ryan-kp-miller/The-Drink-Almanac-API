from flask import Flask
from flask_restful import Api

from api.resources.user import (
    User, UserRegister, UserLogin, UserLogout, TokenRefresh 
)

app = Flask(__name__)
api = Api(app)

api.add_resource(User,         '/user')
api.add_resource(UserRegister, '/register')
api.add_resource(UserLogin,    '/login')
api.add_resource(UserLogout,   '/logout')
api.add_resource(TokenRefresh, '/refresh')


if __name__ == "__main__":
    from api.db import db
    db.init_app(app)
    app.run(debug=True)

