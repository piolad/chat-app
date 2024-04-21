# endpoint Auth-Service
# endpoint Main-Service

## function Login() - Auth-Service
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
string keyofsesion
string idsession        //I get this from active-session
```

## function Register() - Auth-Service
## function Register() - Main-Service
It is used to pass data to register to the application to the authentication service.
```rust
string username
@@ -28,7 +29,7 @@ returns
enum/string status
```

## function UsernameExists() - Auth-Service
## function UsernameExists() - Main-Service
It is used to check if the username exists in the application.
```rust
string username
@@ -38,12 +39,24 @@ returns
bool exists
```

## function EmailExists() - Auth-Service
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