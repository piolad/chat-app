use tonic::{transport::Server, Request, Response, Status};
use std::env;
use tokio_postgres::NoTls;
use dotenv::dotenv;

mod proto{      //podobno tak trzeba
    tonic::include_proto!("auth");
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
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

    let addUserQuery = r#"
        INSERT INTO users (email, username, hashed_password, first_name, last_name, date)
        VALUES ('brud@brud.pl', 'brud', '8rud!', 'Brudas', 'Brudowski', '2004-01-01')
    "#;

    client
        .execute(addUserQuery, &[]).await?;
    
    Ok(())
}
