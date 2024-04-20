# endpoint Main-Service

## function Login() - Main-Service
It is used to pass data to log in to the application to the authentication service.
```rust
string username
string email
string password
string location
```
returns
```rust
enum/string status
string idsession        //I get this from active-session
```

## function Register() - Main-Service
It is used to pass data to register to the application to the authentication service.
```rust
string username
string email
string password
string name
string surname
string date
```
returns
```rust
enum/string status
```

## function UsernameExists() - Main-Service
It is used to check if the username exists in the application.
```rust
string username
```
returns
```rust
bool exists
```

## function EmailExists() - Main-Service
It is used to check if the email exists in the application.
```rust
string email
```
returns
```rust
bool exists
```

# endpoint Active-Sessions

## function GetActiveSessions() - Active-Sessions
```rust
string username        //chnage the name of the function
string email
string name
string surname
string token          //I just create it, dont save it
string localization   //I get this
```