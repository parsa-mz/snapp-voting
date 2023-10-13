# Snapp Voting

## Introduction 
Simple high-performance Voting API backend built with Go (gin 🍸) and PostgreSQL.

## Features
- JWT Authentication: Securely authenticate users and protect your endpoints.
- User Management: Register, get, and update user profiles.
- Voting System: Users can vote for their favorite option.
- Middleware Support: Enhance and simplify HTTP request handling.
- Minio Storage Integration: Store and retrieve files efficiently with Minio.
- Error Management with Sentry

## Prerequisites
- Go 
- PostgreSQL 
- Minio Server

## Setup & Installation

1. **Install Go**
   ```bash
   brew install go
   ```

2. **Set Up Environment Variables**
   Create a `.env` file in the root directory and fill in the necessary details:
   ```
   DB_HOST=localhost
   DB_USER=your_db_user
   DB_PASS=your_db_password
   DB_NAME=voting_db
   
   JWT_SECRET=<your_jwt_secret>
   
   MINIO_STORAGE_ENDPOINT=<minio.yourdomain.com>
   MINIO_STORAGE_ACCESS=<your_minio_access_key>
   MINIO_STORAGE_SECRET=<your_minio_secret_key>
   ```

3. **Install Dependencies**
   ```bash
   go mod download
   ```

4. **Run the Server**
   ```bash
   go run main.go
   ```


## Middlewares
- **Authentication**: Ensures that the user is authenticated using JWT.
- Other middlewares can be added as well (like logging, etc.)

## License
This project is licensed under the MIT License.


---

Happy Voting!