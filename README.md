# uacl
## Introduction
uacl or user account, creation and log in is the service that allows users to create accounts, log in and authenticate requests from other services.

At the moment is only supports one refresh token, this means if a user logs into an account on a PC and then on a tablet the one on the PC will be logged out. However, it does support multiple autologin tokens.

Cannot deactivate users at the moment. Could alter the password from the database though to effectively lock them out and delete all autologin tokens and refresh tokens.

Note: all usernames are lowercased e.g. If user test exists, attempting to register TeSt will not work.

## Production Environment Variables
```
DATABASE_URL - URL for the database. Should also be connecting to the uacl database
HOST - In case the service needs to run on anything other than 0.0.0.0
PORT - In case the service needs to run on anything other than 80
PRIVATE_KEY - Location of the JWT private key
PUBLIC_KEY - Location of the JWT public key
SECRET - This is the create account secret, to restrict who can create an account. This can be changed and deployed whenever the secret is compromised, existing users are unaffected.
CHATTER_URL - The url of the chatter service
CHATTER_AUTH - The authorization string to authenticate that uacl is talking to it.
AUTOLOGIN_URL - The autologin url the front end uses to sign in accounts
AUTOLOGIN_CREATE_USERS - list of comma seperated usernames that can create autologin tokens. This way only a select users can create and delete autologin tokens
EMAIL_FROM - Email configuration.
EMAIL_PASSWORD - Email configuration.
EMAIL_LEVEL - What level of logs gets sent to the email address.
ALLOWED_ORIGINS - Cors setup.
```
## Endpoints
```
base URL is uacl.emotives.net

GET - /healthz - Standard endpoint that just returns ok and a 200 status code. Can be used to test if the service is up
GET - /authorize - Used to authenticate the Authorization header. The header should be in the form of Bearer {JWT}. Will return 401 if unsuccessful or 200 if successful.
POST - /refresh - Used to refresh the JWT token. The refresh token should be in the request body as a JSON with just a refresh_token field and value.
POST - /user - Creates a user based on the request body. See model/user for JSON fields of the user.
POST - /login - attempts to sign into an account. Only need to supply username and password in the request body

POST - /autologin - creates an autologin token for a user. Usage of the token is unlimited and has no expire date.
POST - /autologin/{token} - attempts to find the {token} in the database and then creates a JWT for the user, signing them in
GET - /autologin - Fetches all autologin tokens
GET - /autologin/latest - Fetches latest autologin token for the same user
GET - /autologin/{token_id} - Fetches specific autologin token for the same user or AUTOLOGIN_CREATE_USERS users
DELETE - /autologin/{token} - deletes the token if it can be found the database. Useful if any login tokens are compromised/lost/want to restrict the account.
```
## Database design
Uses a postgres database.
[See here for latest schema, uses the uacl_db](https://github.com/EmotivesProject/databases)

