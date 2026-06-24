# user-authentication Specification

## Purpose
TBD - created by archiving change build-core-gin-template. Update Purpose after archive.
## Requirements
### Requirement: User login
The system SHALL expose `POST /api/v1/user/login` and authenticate a user by required username and password fields.

#### Scenario: Successful login
- **WHEN** the request contains valid credentials for an active user
- **THEN** the system SHALL return HTTP 200 with a JWT token in the response data

#### Scenario: Invalid credentials
- **WHEN** the username does not exist or the password does not match its bcrypt hash
- **THEN** the system SHALL return an authentication failure response and SHALL NOT issue a token

#### Scenario: Invalid login payload
- **WHEN** the login request omits either username or password
- **THEN** the system SHALL return HTTP 400 with the shared failure response format

### Requirement: Bearer token authentication
The system SHALL protect authenticated routes by parsing a JWT from the `Authorization: Bearer <token>` header and storing its `user_id` claim in the Gin request context.

#### Scenario: Valid bearer token
- **WHEN** a request to a protected route includes a valid, unexpired bearer token
- **THEN** the middleware SHALL set `user_id` in the request context and continue the request

#### Scenario: Missing or invalid bearer token
- **WHEN** a request to a protected route has no bearer token, an invalid token, or an expired token
- **THEN** the middleware SHALL return HTTP 401 in the shared failure response format and abort the request

### Requirement: Authenticated user information lookup
The system SHALL expose `GET /api/v1/user/info` as a protected route and return the authenticated user's safe profile fields.

#### Scenario: Successful user information lookup
- **WHEN** an authenticated request has a valid `user_id` claim associated with an existing user
- **THEN** the system SHALL return the user's ID, username, nickname, and status and SHALL NOT return the password hash

#### Scenario: Missing user record
- **WHEN** the authenticated user ID does not map to a user record
- **THEN** the system SHALL return a failure response without exposing persistence details

