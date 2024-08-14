# Choice-Assignment

## Prerequisites

- Go
- Git
- MySQL

## Dependencies List

This project is used by the following companies:

- github.com/gin-gonic/gin v1.10.0 
- gorm.io/gorm v1.25.11 
- gorm.io/driver/mysql v1.5.7 
- github.com/xuri/excelize/v2 v2.8.1 
- github.com/go-redis/redis/v8 v8.11.5


## Setup Instructions

1. **Clone the Repository**

   ```sh
   git clone https://github.com/vinayak-chavan/touch-test.git
   cd touch-test
   ```

2. **Install Dependencies**

   ```sh
   go mod tidy
   ```

3. **Set Up Environment Variables**

   Create .env file in the root directory:
   ```sh 
   DB_USERNAME=username
   DB_PASSWORD=password
   DB_HOST=127.0.0.1
   DB_PORT=3306
   DB_NAME=database-name
   PORT=8000
   REDIS_HOST=127.0.0.1
   REDIS_PORT=6379
   ```

4. **Run the Project**

   ```sh
   go run main.go
   ```

## Run API

In root folder there is file with name `choice-test.postman_collection.json` exist. Import this file in postman and you can execute the APIs.

## API Reference

#### Upload sheet

```http
  POST /api/v1/users
```

#### Get all users

```http
  GET /api/v1/users
```

#### Get user

```http
  GET /api/v1/users/${id}
```

#### Delete user

```http
  DELETE /api/v1/users/:id
```

#### Get item

```http
  UPDATE /api/v1/users/:id

  {
    "first_name": "John",
    "last_name": "Doe",
    "company_name": "New Company",
    "address": "123 New St",
    "city": "New York",
    "county": "New York",
    "postal": "10001",
    "phone": "123-456-7890",
    "email": "john.doe@example.com",
    "web": "https://example.com"
}
```
