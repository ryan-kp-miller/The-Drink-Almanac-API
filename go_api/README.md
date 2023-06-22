# The Drink Almanac API (Golang version) <!-- omit in toc -->

[![codecov](https://codecov.io/gh/ryan-kp-miller/The-Drink-Almanac-API/branch/feature%2Fgo-api/graph/badge.svg?token=D5YMAWKNM4)](https://codecov.io/gh/ryan-kp-miller/The-Drink-Almanac-API)

The original API was built using Python and Heroku. This version will use Golang and AWS. The functionality should not change (for now).


## Table of Contents <!-- omit in toc -->

- [How to Run Locally](#how-to-run-locally)
- [Endpoints](#endpoints)
- [To Do](#to-do)
  - [Unfinished](#unfinished)
  - [Finished](#finished)


## How to Run Locally

First, create a `.env` file in the top-level directory with the following fields:
```
ENV="local"
AWS_DEFAULT_REGION="us-east-1"
AWS_ACCESS_KEY_ID="local"
AWS_SECRET_ACCESS_KEY="local"
AWS_SESSION_TOKEN="local"
JWT_SECRET_KEY="some_secret_key_value" # replace this value with something more secure
```

Then, run the `make up` command. This will start up docker containers for localstack, dynamodb-admin, and the api. The api will be available at `localhost:8000` and dynamodb-admin at `localhost:8001`.

To stop the api, run the `make down` command.


## Endpoints

- `/user`
  - HTTP Commands Allowed:
    - `GET`: get user info using JWT
      - JWT must be stored in `Token` header
    - `POST`: create a new user
    - `DELETE`: delete user account
      - JWT must be stored in `Token` header
- `/user/login`
  - HTTP Commands Allowed:
    - `POST`: log in to user's account
      - JWT is returned in `Token` header
- `/favorite`
  - HTTP Commands Allowed:
    - `GET`: get all favorites for a user
      - User id is retrieved from JWT in the `Token` header
    - `POST`: create a new favorite for a given user and drink
      - Drink id should be provided in the request body
      - User id is retrieved from JWT in the `Token` header
    - `DELETE`: delete a favorite
      - Favorite id provided in the url
      - JWT must be stored in `Token` header


## To Do

### Unfinished

- [ ] Set up lambda handler to be used by each endpoint
- [ ] Switch endpoints to lambdas
- Create endpoints for:
  - [ ] Add new method to find a favorite by user and drink ids and update the favorite service to use that instead of getting all favorites and then filtering
  - [ ] Update delete user method to also delete any favorites associated with that user
    - [ ] Add DeleteFavorites method that takes a slice of id strings and deletes those favorites
    - [ ] Add favorite store field to UserService
  - [ ] Add endpoint for retrieving drink data
- [ ] Add better logging
- [ ] Use API Gateway (separate repo?) for the endpoints
- [ ] Swagger docs
- [ ] Deploy API using Terraform
  - [ ] Set up CI/CD for automatically deploying changes
- [ ] Add code coverage badge to repo's README
- [ ] Add `create_ts` for users and favorites


### Finished

- [x] Troubleshoot why drinkId and userId are no longer coming through since they were changed to strings
- Create endpoints for:
  - [x] Create new users
  - [x] Create new favorites
    - [x] Add validation logic inside the favorite service to verify that a favorite doesn't exist before creating a new one
  - [x] Delete a favorite using the favorite's id
  - [x] Fix create favorite post method to actually check if drink or user ids are empty
  - [x] Add method to delete user 
  - [x] User authentication
  - [x] Update existing endpoints to use user authorization (access tokens)
    - [x] DELETE /user
    - [x] GET /user
    - [x] GET /favorite (remove endpoint for retrieving all favorites for all users)
    - [x] POST /favorite
    - [x] DELETE /favorite
- [x] Write unit tests (80% or higher coverage)
  - [x] Favorites
    - [x] dto
    - [x] handler
    - [x] service
    - [x] store
  - [x] Users
    - [x] dto
    - [x] handler
    - [x] service
    - [x] store
  - [x] Add tests to CICD pipeline
- [x] Set up Data Transfer Object types for request bodies and responses
- [x] Move hardcoded table names and other env variables to an app config struct
- [x] Move auth validation from handlers to a middleware
- [x] Update stores to use query instead of scan where applicable
- [x] More tests for new code
  - [x] `UserService.Login`
  - [x] `UserHandlers.Login`
  - [x] `AuthMiddleware.AuthUser`