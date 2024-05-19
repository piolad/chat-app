use tonic::{transport::Server, Request, Response, Status};
use tokio_postgres::NoTls;
use dotenv::dotenv;
//use bcrypt::{hash, verify};
use proto::auth_server::{Auth, AuthServer};
use rand::{Rng, thread_rng};
use rand::distributions::Alphanumeric;

mod proto {
    tonic::include_proto!("auth");
}

#[derive(Debug)]
struct AuthService {
    client: tokio_postgres::Client,
}

impl AuthService {
    fn new(client: tokio_postgres::Client) -> Self {
        Self { client }
    }
}

fn generate_token(length: usize) -> String {
    let rng = thread_rng();
    let token: String = rng.sample_iter(&Alphanumeric)
        .take(length)
        .map(|c| c as char) // Convert each u8 to char
        .collect(); // Collect characters into a String
    token
}

#[tonic::async_trait]
impl Auth for AuthService {
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

        let user_query = "SELECT email,username,hashed_password FROM users WHERE email = $1 OR username = $1";
        let row = match self.client.query_one(user_query, &[&login_identifier]).await {
            Ok(row) => row,
            Err(_) => {
                return Err(Status::not_found("User not found"));
            }
        };

        let email: String = row.get(0);
        let username: String = row.get(1);
        let hashed_password: String = row.get(2);

        if verify_password(&password, &hashed_password) && (login_identifier == email || login_identifier == username) {
            let reply = proto::LoginResponse {
                status: "Success".to_string(),
                token: generate_token(32).to_string(),
                idsession: "idsession".to_string(),     //this i get from acctive sessions
            };
            return Ok(Response::new(reply));
        } else {
            return Err(Status::unauthenticated("Invalid credentials"));
        }
    }

    // dodanie tworzenia tokenu odszyfrowującego wiadomości przy rejestracji dodanie pola do bazy danych

    async fn register(
        &self,
        request: Request<proto::RegisterRequest>,
    ) -> Result<Response<proto::RegisterResponse>, Status> {
        let request = request.into_inner();
        let firstname = request.firstname;
        let lastname = request.lastname;

        let email = request.email;
        let username = request.username;
        let hashed_password = hash_password(&request.password);

        let user_query = r#"INSERT INTO users (email, username, hashed_password, first_name, last_name, date) VALUES ($1, $2, $3, $4, $5, $6)"#;
        match self.client.execute(user_query, &[&email, &username, &hashed_password, &firstname, &lastname, &request.birthdate]).await {
            Ok(_) => {
                let reply = proto::RegisterResponse {
                    status: "Success".to_string(),
                };
                return Ok(Response::new(reply));
            }
            Err(e) => {
                if e.code() == Some(&tokio_postgres::error::SqlState::UNIQUE_VIOLATION) {
                    // Assuming there's only one unique constraint for both email and username
                    let error_msg = format!("User with email '{}' or username '{}' already exists", email, username);
                    return Err(Status::already_exists(error_msg));
                }
                eprintln!("Error executing SQL query: {:?}", e);
                return Err(Status::internal("Error executing SQL query"));
            }
        }
    }
}

fn hash_password(password: &str) -> String {
    bcrypt::hash(password, bcrypt::DEFAULT_COST).expect("Failed to hash password")
}

fn verify_password(password: &str, hashed_password: &str) -> bool {
    bcrypt::verify(password, hashed_password).expect("Failed to verify password")
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    dotenv().ok();

    let addr = "0.0.0.0:50051".parse().unwrap(); // gRPC listening address

    let database_url = "postgres://postgres:mysecretpassword@auth-service-db/postgres";

    let (client, connection) = tokio_postgres::connect(&database_url, NoTls).await?;

    tokio::spawn(async move {
        if let Err(e) = connection.await {
            eprintln!("connection error: {}", e);
        }
    });

    let table_creation_query = r#"
        CREATE TABLE IF NOT EXISTS users (
            Id SERIAL PRIMARY KEY,
            email VARCHAR(255) NOT NULL UNIQUE,
            username VARCHAR(50) NOT NULL UNIQUE,
            hashed_password TEXT NOT NULL,
            first_name VARCHAR(100),
            last_name VARCHAR(100),
            date VARCHAR(250) NOT NULL
        )"#;

    client.execute(table_creation_query, &[]).await?;

    let add_user_query = format!(
        r#"
            INSERT INTO users (email, username, hashed_password, first_name, last_name, date)
            VALUES ('brud@brud.pl', 'brud', '{}', 'Brudas', 'Brudowski', '2004-01-01')
        "#,
        hash_password("8rud!")
    );

    client.execute(add_user_query.as_str(), &[]).await?;

    println!("Server listening on {}", addr);
    Server::builder()
        .add_service(AuthServer::new(AuthService::new(client)))
        .serve(addr)
        .await?;

    Ok(())
}
