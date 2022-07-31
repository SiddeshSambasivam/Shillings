# Shillings

_A backend system for a payment service that allows you to send and receive money._

<p align="center">
    <img src="assets/overview.png" />
    <p align="center">Fig 1. <i>An overview of the backend system</i></p>
</p>

### Project Scope

A set of web APIs that provide **payment services** with authentication. A _custom communication protocol (shillings)_ is used by the application and platform layer. The `platform layer` handles all the business logic and the `application layer` handles the client calls. _Shillings_ is a custom protocol on top of TCP written in `Go` to handle platform level services such as authentication, payment, database access, and so on.

## Technical Design Decisions

### 1. Database

**Tables**: `users`, `transactions`, `credentials`

| Table        | Columns                                                                               |
| ------------ | ------------------------------------------------------------------------------------- |
| users        | id, first_name, middle_name, last_name, email, phone, balance, created_at, updated_at |
| transactions | id, sender_id, receiver_id, amount, created_at                                        |
| credentials  | id, user_id, password, salt, updated_at, last_login                                   |

**Stack:** SQL

In addition, `redis` is used to cache the user data and authentication tokens.

#### **Tasks**

-   [ ] Setup SQL database locally (docker)
-   [ ] Setup redis locally (docker)
-   [ ] Populate the database with some data

### 2. Application Layer

| API                        | Method | Description                                |
| -------------------------- | ------ | ------------------------------------------ |
| `/v1/login`                | `POST` | Authenticates the user and returns a token |
| `/v1/signup`               | `PUT`  | Register a new user                        |
| `/v1/pay`                  | `POST` | Makes a payment to another user            |
| `/v1/auth/check_user`      | `POST` | Checks if a user exists                    |
| `/v1/account`              | `GET`  | Gets a user profile                        |
| `/v1/validate_transaction` | `POST` | Checks if a transaction is valid           |

#### **Tasks**

-   [ ] write the API handlers
-   [ ] Write the utility functions to handle protobuf, read and write requests with platform layer

### 3. Platform Layer

| Command  | Value | Function                                   |
| -------- | ----- | ------------------------------------------ |
| `LGN`    | 0     | Authenticates the user and returns a token |
| `SGN`    | 1     | Register a new user                        |
| `PY`     | 2     | Makes a payment to another user            |
| `CHKUSR` | 3     | Checks if a user exists                    |
| `ACCT`   | 4     | Gets a user profile                        |
| `VLDT`   | 5     | Checks if a transaction is valid           |
| `PING`   | 6     | Pings the service                          |

#### **Tasks**

-   [ ] Write the required protobuf messages for the communication protocol
    -   [ ] Compile the protobuf messages with `protoc`
    -   [ ] Copy the generated files to the `/proto` directory
-   [ ] Write the handlers for each command
-   [ ] Setup database handlers
    -   [ ] Write connection pool for the database
-   [ ] Setup redis handlers
    -   [ ] Write connection pool for the redis database
-   [ ] Add health check for platform server.

## Performance requirements:

-   [ ] Backend service should be able to handle a minimum of 10K qps
    -   Load test the backend service with `wrk`
