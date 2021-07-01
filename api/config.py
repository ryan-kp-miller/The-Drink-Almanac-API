import os
from datetime import timedelta


class Config(object):
    SECRET = os.getenv('SECRET_KEY')

    # turn off flask_sqlalchemy modification tracker so we can use SQLAlchemy's mod tracker, which is better
    SQLALCHEMY_TRACK_MODIFICATIONS = False 
    
    # allows Flask extensions like flask_jwt to raise their own errors,
    # as opposed to Flask just returning a 500 for all errors
    PROPAGATE_EXCEPTIONS = True

    # set access tokens to expire in 1 hour
    JWT_ACCESS_TOKEN_EXPIRES = timedelta(hours=1)

    
    # heroku uses an old postgres dialect that is no longer supported by flask-sqlalchemy
    # so manually switching to the new one
    SQLALCHEMY_DATABASE_URI = os.getenv('DATABASE_URL').replace("postgres://", "postgresql://")


class DevelopmentConfig(Config):
    DEBUG = True


class TestingConfig(Config):
    TESTING = True
    SQLALCHEMY_DATABASE_URI = 'sqlite:///../test_db.db'
    DEBUG = True


class ProductionConfig(Config):
    DEBUG = False
    TESTING = False


app_config = {
    'development': DevelopmentConfig,
    'testing': TestingConfig,
    'production': ProductionConfig,
}
