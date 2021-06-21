from flask_restful import Resource, reqparse
from werkzeug.security import safe_str_cmp
from api.models.user import UserModel


_user_parser = reqparse.RequestParser()
_user_parser.add_argument('username', type=str, required=True, help="Username for the new user.")
_user_parser.add_argument('password', type=str, required=True, help="Password for the new user.")


class UserRegister(Resource):
    def post(self):
        data = _user_parser.parse_args()
        
        user = UserModel.find_by_username(data['username'])
        if user:
            return {'message': f'A user with the username "{data["username"]}" already exists.'}, 400
        
        user = UserModel(**data)
        user.save_to_db()

        return {'message': 'User created successfully.'}, 201


class User(Resource):
    @classmethod
    def get(cls, user_id):
        user = UserModel.find_by_id(user_id)
        if user:
            return user.json()
        return {'message': f'User with id {user_id} not found.'}, 404

    @classmethod
    def delete(cls, user_id):
        user = UserModel.find_by_id(user_id)
        if not user:
            return {'message': f'User with id {user_id} not found.'}, 404
            
        user.delete_from_db()    
        return {'message': f'User with id {user_id} was deleted.'}, 200

