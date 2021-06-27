from flask_restx import Resource, reqparse, Namespace
from flask_jwt_extended import jwt_required
from flask_jwt_extended.utils import get_jwt_identity
from api.models.favorite import FavoriteModel
from api.models.user import UserModel

api = Namespace("favorite", description="Add or remove favorited drinks from the user")
parser = reqparse.RequestParser()
parser.add_argument(
    'drink_id', 
    int, 
    required=True, 
    help="The id of the drink that is being favorited",
    location="json"
    
)

@api.route("")
class Favorite(Resource):
    @jwt_required()
    @api.doc(security="apiKey", responses={
        404: 'Favorite not found',
        200: 'Success',
    })
    def get(self):
        user_id = get_jwt_identity()
        data = parser.parse_args()
        favorite = FavoriteModel.find_by_user_and_drink_ids(user_id, data['drink_id'])
        if favorite:
            return favorite.json(), 200
        return {'message': f'Favorite not found'}, 404

    @jwt_required()
    @api.doc(security="apiKey", responses={
        404: 'User for that JWT not found. Please remove the stale JWT',
        400: 'The user has already favorited this drink',
        201: 'Success',
    })
    @api.expect(parser)
    def post(self):
        user_id = get_jwt_identity()
        data = parser.parse_args()
        user = UserModel.find_by_id(user_id)
        if not user:
            return {'message': "User for that JWT not found. Please remove the stale JWT"}, 404

        favorite = FavoriteModel.find_by_user_and_drink_ids(user_id, data['drink_id'])
        if favorite:
            return {'message': 'The user has already favorited this drink'}, 400
        favorite = FavoriteModel(user_id, data['drink_id'])
        favorite.save_to_db()
        return favorite.json(), 201

    @jwt_required()
    @api.doc(security="apiKey", responses={
        404: 'Favorite not found',
        201: 'Success',
    })
    def delete(self):
        user_id = get_jwt_identity()
        data = parser.parse_args()
        favorite = FavoriteModel.find_by_user_and_drink_ids(user_id, data['drink_id'])
        if favorite:
            favorite.delete_from_db()
            return {'message': 'Success'}, 200
        return {'message': 'Favorite not found'}, 404
