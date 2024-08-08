# Todo List in Golang

## How to Use:
1. **Setup:** Ensure you have Go installed on your system. Clone the repository from GitHub.
```bash
git clone https://github.com/SamuelMolling/todo-list-golang.git
cd todo-list-golang
```
2. **MongoDB Configuration:** Set up your cluster and configure your connection string on line 32.
```
mongoURI := "mongodb+srv://<USERNAME>:<PASSWORD>@<CLUTER_ENDPOINT>/?retryWrites=true&w=majority&appName=demo1"
```
3. Run the Application: Execute the following command to start the server.
```bash
go run main.go
```
4. Access the Application: Open your browser and navigate to `http://localhost:8080` to interact with the Todo List application.
