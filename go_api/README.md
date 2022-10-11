# The Drink Almanac API (Golang version) <!-- omit in toc -->

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
- `/favorites/{user_id}` (will be changed to use access tokens)
- `/favorite/{drink_id}` (will be changed to use access tokens)
  - HTTP Commands Allowed:
    - `POST`: create a new favorite for a given user and drink
    - `DELETE`: delete a favorite for a given user and drink


## To Do

- Troubleshoot why drinkId and userId are no longer coming through since they were changed to strings
- Create endpoints for:
  - User authentication
  - Create new users
  - Create new favorites
    - Add validation logic inside the favorite service to verify that a favorite doesn't exist for 
- Set up Data Transfer Object types for request bodies and responses
- Update existing endpoints to use user authorization (access tokens)
- Write unit tests (80% or higher coverage)
- Add better logging using Gin's built-in logger
- Deploy API
- Set up CI/CD for automatically deploying changes
- Add code coverage badge to repo's README