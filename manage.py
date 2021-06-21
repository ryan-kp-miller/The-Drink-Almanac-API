from api.app import create_app
from api.db import db


app = create_app()

# before any request to the API, this function will be called 
# and will create the data.db file and all the tables within the db (unless they already exist)
@app.before_first_request
def create_tables():
    db.create_all()

db.init_app(app)


if __name__ == "__main__":
    app.run(debug=True)