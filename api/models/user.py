from api.db import db
from api.models.favorite import FavoriteModel

class UserModel(db.Model):
    # declare the table name, columns, and relationships for SQLAlchemy
    __tablename__ = 'users'
    id = db.Column(db.Integer, primary_key=True)
    username = db.Column(db.String(80))
    password = db.Column(db.String(80))

    def __init__(self, username, password):
        self.username = username
        self.password = password

    def json(self):
        return {
            'id': self.id,
            'username': self.username,
            'favorites': [fav.json()['drink_id'] for fav in FavoriteModel.find_favorites_by_user_id(self.id)],
        }
    
    def save_to_db(self):
        db.session.add(self)
        db.session.commit()

    def delete_from_db(self):
        # retrieve and remove all the favorites for that user
        FavoriteModel.query.filter_by(user_id=self.id).delete()

        # remove the user from the db
        db.session.delete(self)
        db.session.commit()

    @classmethod
    def find_by_username(cls, username):
        return cls.query.filter_by(username=username).first()

    @classmethod
    def find_by_id(cls, _id):
        return cls.query.filter_by(id=_id).first()
