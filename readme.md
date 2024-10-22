# Split Expenses

## How to Run the App

1. **Prerequisites:**
   - Make sure you have Go installed on your machine. You can download it from [golang.org](https://golang.org/dl/).
   - Ensure that you have MongoDB installed and running. You can find installation instructions on the [MongoDB website](https://www.mongodb.com/try/download/community). You can use local instance or create an account in MongoDB Atlas, a cloud based MongoDB service.

   - Setup your .env file
    ```
    DB_USER=dbUser
    DB_PASS=dbPass
    DB_CLUSTER=dbCluster
    DB_NAME=dbName
    PORT=8080
    ```


2. **Clone and run the application:**
   ```bash
   git clone https://github.com/RishanKP/split-expenses.git
   cd split-expenses
   go mod tidy
   go run main.go
    ```
## API DOCUMENTATION

| Method | Endpoint          | Description                                                     |  
|--------|-------------------|-----------------------------------------------------------------|
| POST   | /user/signup          | Create a new user   | 
| POST    | /user/login          | Login with credentials.             | 
| GET    | /user/:id  | Get user details by Id  | 
| POST    | /expense | Create a new expense. (Requires authentication) | 
| GET    | /expense | Get list of expenses of logged in user. (Requires authentication) | 
| GET    | /balance-sheet | Downloads balance sheet of logged in user (Requires authentication) | 

## SAMPLE REQUESTS
You can download the postman collection for sample requests using the link below:

[Download postman collection](./requests.postman_collection.json)

