use tonic::{transport::Server, Request, Response, Status};
use std::env;
use tokio_postgres::NoTls;
use dotenv::dotenv;
use bcrypt::{hash, verify};
use proto::auth_server::{Auth, AuthServer};

mod proto{
    tonic::include_proto!("auth");
}

fn hash_password(password: &str) -> String {
    bcrypt::hash(password, bcrypt::DEFAULT_COST).expect("Failed to hash password")
}

fn verify_password(password: &str, hashed_password: &str) -> bool {
    bcrypt::verify(password, hashed_password).expect("Failed to verify password")
}

#[derive(Debug, Default)]
struct AuthService {}

#[tonic::async_trait]
impl Auth for AuthService{
    async fn login(&self, request: Request<proto::LoginRequest>) -> Result<Response<proto::LoginResponse>, Status> {
        let request = request.into_inner();

        let login_identifier = match request.login_data {
            Some(proto::login_request::LoginData::Username(username)) => username,
            Some(proto::login_request::LoginData::Email(email)) => email,
            None => {
                return Err(Status::invalid_argument("Username or email is required"));
            }
        };

        let password = request.password;
        let mix = format!("{}{}", login_identifier, password);
        
        let reply = proto::LoginResponse {
            status: mix.to_string(),
        };
        Ok(Response::new(reply))
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr = "0.0.0.0:50051".parse().unwrap();//grpc listening address

    let database_url = "postgres://postgres:mysecretpassword@auth-service-db/postgres";

    let (client, connection) =
        tokio_postgres::connect(&database_url, NoTls).await?;

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
            date DATE NOT NULL
    )"#;

    client
        .execute(table_creation_query, &[])
        .await?;

    let addUserQuery = format!(
        r#"
            INSERT INTO users (email, username, hashed_password, first_name, last_name, date)
            VALUES ('brud@brud.pl', 'brud', '{}', 'Brudas', 'Brudowski', '2004-01-01')
        "#, hash_password("8rud!")
    );

    client
        .execute(addUserQuery.as_str(), &[]).await?;

    println!("Server listening on {}", addr);
    Server::builder()
        .add_service(AuthServer::new(AuthService::default()))
        .serve(addr)
        .await?;
    
    Ok(())
}
