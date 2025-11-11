use tonic::{transport::Server, Request, Response, Status};
use tonic::transport::Endpoint;

use tokio_postgres::NoTls;
use tokio::time::{sleep, Duration};

use std::net::SocketAddr;
use std::sync::Arc;

mod config;
use config::Config;

mod seed_demo;

mod security;
use security::{generate_rsa_keypair, hash_password, verify_password};

mod auth_service;
use crate::auth_service::AuthService;

mod proto {
    tonic::include_proto!("auth");
    tonic::include_proto!("active_sessions");
}
use crate::proto::auth_server::AuthServer;


async fn connect_with_retry(db_url: &str, max_retries: u32) -> Result<tokio_postgres::Client, tokio_postgres::Error> {
    let mut attempt = 0u32;
    let mut backoff = Duration::from_millis(200);

    loop {
        match tokio_postgres::connect(db_url, NoTls).await {
            Ok((client, connection)) => {
                // keep the connection running
                tokio::spawn(async move {
                    if let Err(e) = connection.await {
                        eprintln!("connection error: {}", e);
                    }
                });

                // ping to ensure the server is fully ready
                match client.simple_query("SELECT 1").await {
                    Ok(_) => return Ok(client),
                    Err(e) => {
                        attempt += 1;
                        if attempt >= max_retries {
                            return Err(e);
                        }
                    }
                }
            }
            Err(e) => {
                attempt += 1;
                if attempt >= max_retries {
                    return Err(e);
                }
            }
        }

        // exponential backoff up to ~5s
        sleep(backoff).await;
        let next_ms = (backoff.as_millis() * 2).min(5000) as u64;
        backoff = Duration::from_millis(next_ms);
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let cfg = Config::from_env()?;

    let addr: SocketAddr = cfg.server_addr;
    let client = connect_with_retry(&cfg.database_url, 30).await?;

    // Create tables
    let table_creation_users = r#"
        CREATE TABLE IF NOT EXISTS users (
            Id SERIAL PRIMARY KEY,
            email VARCHAR(255) NOT NULL UNIQUE,
            username VARCHAR(50) NOT NULL UNIQUE,
            hashed_password TEXT NOT NULL,
            first_name VARCHAR(100),
            last_name VARCHAR(100),
            date VARCHAR(250) NOT NULL
        )"#;
    client.execute(table_creation_users, &[]).await?;

    let table_creation_keys = r#"
        CREATE TABLE IF NOT EXISTS keys (
            Id_user SERIAL PRIMARY KEY,
            public_key TEXT NOT NULL,
            private_key TEXT NOT NULL
        )"#;
    client.execute(table_creation_keys, &[]).await?;

    // Seed optional demo user
    if cfg.seed_demo_user {
        seed_demo::seed_demo_user(&client, cfg.bcrypt_cost).await;
    }

    println!("Server listening on {}", addr);
    Server::builder()
        .add_service(AuthServer::new(AuthService::new(client, &cfg)))
        .serve(addr)
        .await?;
    Ok(())
}


// wyrzucić token i expiration date dla niego
// dodać tabele do kluczy
