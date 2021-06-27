from flask_restx import Resource, reqparse
from flask_jwt_extended import jwt_required
from flask_jwt_extended.utils import get_jwt_identity

from api.models.favorite import FavoriteModel
from api.models.user import UserModel

parser = reqparse.RequestParser()
parser.add_argument('drink_id', int, required=True, help="The id of the drink that is being favorited")

class Favorite(Resource):
    @jwt_required()
    def get(self):
        user_id = get_jwt_identity()
        data = parser.parse_args()
        favorite = FavoriteModel.find_by_user_and_drink_ids(user_id, data['drink_id'])
        if favorite:
            return favorite.json(), 200
        return {'message': f'Favorite not found'}, 404

    @jwt_required()
    def post(self):
        user_id = get_jwt_identity()
        data = parser.parse_args()
        user = UserModel.find_by_id(user_id)
        if not user:
            return {'message': f"No user with the id {user_id}"}

        favorite = FavoriteModel.find_by_user_and_drink_ids(user_id, data['drink_id'])
        if favorite:
            return {'message': 'The user has already favorited this drink'}, 400
        favorite = FavoriteModel(user_id, data['drink_id'])
        favorite.save_to_db()
        return favorite.json(), 201

    @jwt_required()
    def delete(self):
        user_id = get_jwt_identity()
        data = parser.parse_args()
        favorite = FavoriteModel.find_by_user_and_drink_ids(user_id, data['drink_id'])
        if favorite:
            favorite.delete_from_db()
            return {'message': f'The favorite was deleted from the DB'}, 200
        return {'message': f'Favorite not found'}, 404
