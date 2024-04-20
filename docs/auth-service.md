# endpoint Auth-Service

## function Login() - Auth-Service
It is used to pass data to log in to the application to the authentication service.
```rust
string username
string email
string password
```
returns
```rust
enum/string status
string keyofsesion
```

## function Register() - Auth-Service
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

## function UsernameExists() - Auth-Service
It is used to check if the username exists in the application.
```rust
string username
```
returns
```rust
bool exists
```

## function EmailExists() - Auth-Service
It is used to check if the email exists in the application.
```rust
string email
```
returns
```rust
bool exists
```