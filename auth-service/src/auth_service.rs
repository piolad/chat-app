use std::sync::Arc;
use tonic::{Request, Response, Status};
use tonic::transport::Endpoint;

use crate::config::Config;
use crate::proto;

use crate::security::{generate_rsa_keypair, hash_password, verify_password};


#[derive(Debug, Clone)]
pub struct AuthService {
    client: Arc<tokio_postgres::Client>,
    active_sessions_url: String,
    default_location: String,
    bcrypt_cost: u32,
}

impl AuthService {
    pub fn new(client: tokio_postgres::Client, cfg: &Config) -> Self {
        Self {
            client: Arc::new(client), 
            active_sessions_url: cfg.active_sessions_url.clone(),
            default_location: cfg.default_location.clone(),
            bcrypt_cost: cfg.bcrypt_cost,
        }
    }
}

#[tonic::async_trait]
impl crate::proto::auth_server::Auth for AuthService {
    async fn login(
        &self,
        request: Request<proto::LoginRequest>,
    ) -> Result<Response<proto::LoginResponse>, Status> {
        let request = request.into_inner();
        let password = request.password;

        let login_identifier = match request.login_data {
            Some(proto::login_request::LoginData::Username(username)) => username,
            Some(proto::login_request::LoginData::Email(email)) => email,
            None => {
                return Err(Status::invalid_argument("Username or email is required"));
            }
        };

        let user_query = "SELECT username, email, hashed_password FROM users WHERE email = $1 OR username = $1";
        let row = match self.client.query_one(user_query, &[&login_identifier]).await {
            Ok(row) => row,
            Err(_) => {
                return Err(Status::not_found("User not found"));
            }
        };

        let username: String = row.get(0);
        let email: String = row.get(1);
        let hashed_password: String = row.get(2);

        if verify_password(&password, &hashed_password) && (login_identifier == email || login_identifier == username) {
            // Use configured Active Sessions URL
            let channel = Endpoint::from_shared(self.active_sessions_url.clone())
                .map_err(|e| {
                    eprintln!("Invalid active sessions URL: {:?}", e);
                    Status::internal("Invalid active sessions URL")
                })?
                .connect()
                .await
                .map_err(|e| {
                    eprintln!("Failed to connect to active sessions service: {:?}", e);
                    Status::internal("Failed to connect to active sessions service")
                })?;

            let mut client = proto::active_sessions_client::ActiveSessionsClient::new(channel);

            let request = tonic::Request::new(proto::UserData {
                username: username.clone(),
                email: email.to_string(),
                // Use configured default location
                location: self.default_location.clone(),
            });

            let idsession = client.add_user(request).await?.into_inner().session_token;
            println!("idsession: {}", idsession);

            let reply = proto::LoginResponse {
                status: "Success".to_string(),
                token: "delete".to_string(), // legacy field
                idsession: idsession.to_string(),
            };
            Ok(Response::new(reply))
        } else {
            Err(Status::unauthenticated("Invalid credentials"))
        }
    }

    async fn register(
        &self,
        request: Request<proto::RegisterRequest>,
    ) -> Result<Response<proto::RegisterResponse>, Status> {
        let request = request.into_inner();
        let firstname = request.firstname;
        let lastname = request.lastname;

        let email = request.email;
        let username = request.username;
        // Hash with configured bcrypt cost
        let hashed_password = hash_password(&request.password, self.bcrypt_cost);

        generate_rsa_keypair(); // test

        let user_query = r#"INSERT INTO users (email, username, hashed_password, first_name, last_name, date) VALUES ($1, $2, $3, $4, $5, $6)"#;
        match self.client.execute(user_query, &[&email, &username, &hashed_password, &firstname, &lastname, &request.birthdate]).await {
            Ok(_) => {
                let reply = proto::RegisterResponse {
                    status: "Success".to_string(),
                };
                Ok(Response::new(reply))
            }
            Err(e) => {
                if e.code() == Some(&tokio_postgres::error::SqlState::UNIQUE_VIOLATION) {
                    let error_msg = format!("User with email '{}' or username '{}' already exists", email, username);
                    return Err(Status::already_exists(error_msg));
                }
                eprintln!("Error executing SQL query: {:?}", e);
                Err(Status::internal("Error executing SQL query"))
            }
        }
    }
}
