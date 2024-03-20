# endpoint Auth-Service

## function Login() - Auth-Service
```rust
string username
string email
string password
```
It is used to pass data to log in to the application to the authentication service.

## function Register() - Auth-Service
```rust
string username
string email
string password
string name
string surname
string date
```
It is used to pass data to register to the application to the authentication service.

## function UsernameExists() - Auth-Service
```rust
string username
```
It is used to check if the username exists in the application.

## function EmailExists() - Auth-Service
```rust
string email
```
It is used to check if the email exists in the application.