# Human Resources Information System

## Description
This is a simple Human Resources Information System (HRIS).

[Postman Collection](https://documenter.getpostman.com/view/23207346/2sA3duEsLN)

## How to run
```plaintext
1. Clone this repository
2. Run `make go-migrate-up` to migrate the database
3. Run `make server` to start the server on port 8080 (default)
```

swagger will be available at `http://localhost:8080/swagger/index.html`

Login with default user:
```plaintext
username: admin
password: password
```

## Features
- [x] Create user
- [x] Login
- [x] Role
- [x] Permission
- [x] Approval Management
