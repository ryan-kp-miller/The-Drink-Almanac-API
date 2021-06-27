from flask_jwt_extended.view_decorators import jwt_required
from flask_restx import Resource, reqparse, Namespace
from flask_jwt_extended.utils import create_access_token, create_refresh_token, get_jwt_identity
from api.models.user import UserModel


api = Namespace('user', description="Manage accounts and logins")
_user_parser = reqparse.RequestParser()
_user_parser.add_argument(
    'username', 
    type=str, 
    required=True, 
    help="Username for the new user.",
    location="json"
)
_user_parser.add_argument(
    'password', 
    type=str, 
    required=True, 
    help="Password for the new user.",
    location="json"
)


@api.route("/register")
class UserRegister(Resource):
    @api.expect(_user_parser)
    def post(self):
        data = _user_parser.parse_args()
        
        user = UserModel.find_by_username(data['username'])
        if user:
            return {'message': f'A user with the username "{data["username"]}" already exists.'}, 400
        
        user = UserModel(**data)
        user.save_to_db()

        return user.json(), 201


@api.route("")
class User(Resource):
    @classmethod
    @jwt_required()
    @api.doc(security="apiKey")
    def get(cls):
        user_id = get_jwt_identity()
        user = UserModel.find_by_id(user_id)           
        return user.json(), 200

    @classmethod
    @api.expect(_user_parser)
    def delete(cls):
        data = _user_parser.parse_args()
        user = UserModel.find_by_username(data['username'])
        if not user:
            return {'message': f"User with username {data['username']} not found."}, 404
            
        if user.password != data['password']:
            return {'message': f"Password was incorrect"}, 400

        user.delete_from_db()    
        return {'message': f"User with id {data['username']} was deleted."}, 200


@api.route("/login")
class UserLogin(Resource):
    @classmethod
    @api.expect(_user_parser)
    def post(cls):
        data = _user_parser.parse_args()
        user = UserModel.find_by_username(data['username'])
        if not user:
            return {'message': f"User with username {data['username']} not found."}, 404
            
        if user.password != data['password']:
            return {'message': f"Password was incorrect"}, 400
        
        access_token = create_access_token(identity=user.id, fresh=True)
        refresh_token = create_refresh_token(identity=user.id)
        return {
            'access_token': access_token,
            'refresh_token': refresh_token
        }, 200
