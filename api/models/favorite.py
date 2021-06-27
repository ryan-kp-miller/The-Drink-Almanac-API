from api.db import db

class FavoriteModel(db.Model):
    # declare the table name, columns, and relationships for SQLAlchemy
    __tablename__ = 'favorites'
    id = db.Column(db.Integer, primary_key=True)
    drink_id = db.Column(db.Integer)
    user_id = db.Column(db.Integer, db.ForeignKey('users.id'))
    user = db.relationship('UserModel')

    def __init__(self, user_id, drink_id):
        self.user_id = user_id
        self.drink_id = drink_id

    @classmethod
    def find_favorites_by_user_id(cls, user_id):
        return cls.query.filter_by(user_id=user_id)

    @classmethod
    def find_by_id(cls, _id):
        return cls.query.filter_by(id=_id).first()

    @classmethod
    def find_by_user_and_drink_ids(cls, user_id, drink_id):
        return cls.query.filter_by(user_id=user_id, drink_id=drink_id).first()

    def json(self):
        return {
            'id': self.id,
            'user_id':self.user_id, 
            'drink_id': self.drink_id,
        }

    def save_to_db(self):
        db.session.add(self)
        db.session.commit()

    def delete_from_db(self):
        db.session.delete(self)
        db.session.commit()
