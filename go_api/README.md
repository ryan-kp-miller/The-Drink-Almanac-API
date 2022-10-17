# The Drink Almanac API (Golang version) <!-- omit in toc -->

[![codecov](https://codecov.io/gh/ryan-kp-miller/The-Drink-Almanac-API/branch/feature%2Fgo-api/graph/badge.svg?token=D5YMAWKNM4)](https://codecov.io/gh/ryan-kp-miller/The-Drink-Almanac-API)

The original API was built using Python and Heroku. This version will use Golang and AWS. The functionality should not change (for now).


## Table of Contents <!-- omit in toc -->

- [How to Run Locally](#how-to-run-locally)
- [Endpoints](#endpoints)
- [To Do](#to-do)


## How to Run Locally

First, create a `.env` file in the top-level directory with the following fields:
```
API_PORT="8000"
AWS_ENDPOINT=http://localstack:4566
AWS_DEFAULT_REGION="us-east-1"
AWS_ACCESS_KEY_ID="local"
AWS_SECRET_ACCESS_KEY="local"
AWS_SESSION_TOKEN="local"
```

Then, run the `make up` command. This will start up docker containers for localstack, dynamodb-admin, and the api. The api will be available at `localhost:8000` and dynamodb-admin at `localhost:8001`.

To stop the api, run the `make down` command.


## Endpoints

- `/user` (will be changed to use access tokens)
- `/user/login` (not implemented yet)
- `/user/register` (not implemented yet)
- `/favorite`
  - HTTP Commands Allowed:
    - `GET`: get all favorites for a user
      - User id provided in the url (will be changed to use access tokens)
    - `POST`: create a new favorite for a given user and drink
      - user id and drink id provided in the request body
    - `DELETE`: delete a favorite
      - favorite id provided in the url


## To Do

- [x] Troubleshoot why drinkId and userId are no longer coming through since they were changed to strings
- Create endpoints for:
  - [x] Create new users
  - [x] Create new favorites
    - [x] Add validation logic inside the favorite service to verify that a favorite doesn't exist before creating a new one
  - [x] Delete a favorite using the favorite's id
  - [x] Fix create favorite post method to actually check if drink or user ids are empty
  - [ ] Add new method to find a favorite by user and drink ids and update the favorite service to use that instead of getting all favorites and then filtering
  - [x] Add method to delete user 
  - [ ] Update delete user method to also delete any favorites associated with that user
    - [ ] Add DeleteFavorites method that takes a slice of id strings and deletes those favorites
    - [ ] Add favorite store field to UserService
  - [ ] User authentication
- [ ] Set up Data Transfer Object types for request bodies and responses
- [ ] Update existing endpoints to use user authorization (access tokens)
- [ ] Move hardcoded table names and other env variables to an app config struct
- [ ] Write unit tests (80% or higher coverage)
  - [ ] Favorites
    - [x] dto
    - [x] handler
    - [x] service
    - [ ] store
  - [ ] Users
    - [x] dto
    - [x] handler
    - [x] service
    - [ ] store
  - [x] Add tests to CICD pipeline
- [ ] Add better logging
- [ ] Swagger docs
- [ ] Deploy API
  - Terraform?
- [ ] Set up CI/CD for automatically deploying changes
- [ ] Add code coverage badge to repo's README