# gin-template-runtime Specification

## Purpose
TBD - created by archiving change build-core-gin-template. Update Purpose after archive.
## Requirements
### Requirement: Configuration-backed service startup
The system SHALL load server, MySQL, Redis, JWT, and log configuration from `etc/config.yaml` before constructing application dependencies.

#### Scenario: Valid startup configuration
- **WHEN** a valid configuration file is supplied at application startup
- **THEN** the application SHALL initialize its dependencies using the configured values and bind the HTTP server to the configured address

#### Scenario: Invalid startup configuration
- **WHEN** required configuration cannot be loaded or decoded
- **THEN** the application SHALL terminate before accepting HTTP requests and report the initialization failure

### Requirement: Layered dependency boundaries
The system SHALL organize request processing as `router -> handler -> service -> repo -> model/db`; each layer MUST only depend on the next lower layer or approved shared packages.

#### Scenario: HTTP request dispatch
- **WHEN** an API request reaches a registered route
- **THEN** its handler SHALL bind and validate HTTP input, invoke a service method, and return a shared response without directly querying GORM or Redis

#### Scenario: Service execution
- **WHEN** a service performs business work
- **THEN** it SHALL receive a standard `context.Context` and SHALL NOT depend on `gin.Context` or write an HTTP response

### Requirement: Infrastructure initialization
The system SHALL initialize MySQL and Redis during startup, configure the MySQL connection pool, and expose dependencies through the application containers.

#### Scenario: Reachable infrastructure
- **WHEN** configured MySQL and Redis services are reachable
- **THEN** startup SHALL initialize both clients before registering application routes

#### Scenario: Unreachable infrastructure
- **WHEN** MySQL or Redis initialization fails
- **THEN** startup SHALL fail before the HTTP server begins serving requests

### Requirement: Database model migration
The system SHALL run the registered GORM migrations after MySQL initialization and before application routes are served.

#### Scenario: User model migration
- **WHEN** the service starts with a reachable MySQL database
- **THEN** the registered user model schema SHALL be created or migrated before user endpoints are available

### Requirement: Structured and Gin request logging
The system SHALL provide structured Zap business logs and write Gin access and error logs to separately named, hourly-rotated files according to log retention configuration.

#### Scenario: Request logging
- **WHEN** Gin handles a request or emits an error
- **THEN** the corresponding access or error log entry SHALL be written through the configured Gin writer

#### Scenario: Log retention
- **WHEN** an hourly log writer rotates to a new hour
- **THEN** log files older than the configured retention period SHALL be eligible for cleanup

### Requirement: Shared response and utility contracts
The system SHALL provide shared JSON responses, bcrypt password hashing, JWT token utilities, and Snowflake ID generation for use by internal modules.

#### Scenario: Successful JSON response
- **WHEN** a handler completes successfully with data
- **THEN** it SHALL return HTTP 200 with a JSON body containing `code` equal to `0`, `message` equal to `success`, and the result in `data`

#### Scenario: Password verification
- **WHEN** a service compares a plaintext password to a stored bcrypt hash
- **THEN** the shared hash utility SHALL return whether the password matches without exposing the stored plaintext

