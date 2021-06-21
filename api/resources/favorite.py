from flask_restful import Resource, reqparse
from api.models.favorite import FavoriteModel

parser = reqparse.RequestParser()
parser.add_argument('user_id',  int, required=True, help="The id of the user that is favoriting the drink")
parser.add_argument('drink_id', int, required=True, help="The id of the drink that is being favorited")

class Favorite(Resource):
    def get(self, user_id):
        favorites = FavoriteModel.find_favorites_by_user_id(user_id)
        if favorites:
            return {'favorites': [fav.json() for fav in favorites]}, 200
        return {'message': f'user_id {user_id} has no favorites'}, 404

    def post(self):
        data = parser.parse_args()
        favorite = FavoriteModel.find_by_user_and_drink_ids(**data)
        if favorite:
            return {'message', 'User has already favorited this drink'}
        favorite.save_to_db()
        return favorite.json()

    def delete(self, _id):
        favorite = FavoriteModel.find_by_id(_id)
        if favorite:
            favorite.delete_from_db()
            return {'message': f'Favorite with id {_id} was deleted from the DB'}, 200
        return {'message': f'Favorite with id {_id} not found'}, 404