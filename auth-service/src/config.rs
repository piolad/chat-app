use std::env;
use std::net::SocketAddr;

#[derive(Debug, Clone)]
pub struct Config {
    pub server_addr: SocketAddr,
    pub database_url: String,
    pub active_sessions_url: String,
    pub bcrypt_cost: u32,
    pub default_location: String,
    pub seed_demo_user: bool,
}

impl Config {
    pub fn from_env() -> Result<Self, Box<dyn std::error::Error>> {
        // .env is for local only; in containers/prod, just set env vars.
        dotenv::dotenv().ok();

        let server_addr: SocketAddr = env::var("SERVER_ADDR")
            .unwrap_or_else(|_| "0.0.0.0:50051".into())
            .parse()?;

        let database_url = env::var("DATABASE_URL")
            .unwrap_or_else(|_| "postgres://postgres:mysecretpassword@auth-service-db/postgres".into());

        let active_sessions_url = env::var("ACTIVE_SESSIONS_URL")
            .unwrap_or_else(|_| "http://active-sessions:50053".into());

        let bcrypt_cost = env::var("BCRYPT_COST")
            .ok()
            .and_then(|s| s.parse::<u32>().ok())
            .unwrap_or(4);

        let default_location = env::var("DEFAULT_LOCATION")
            .unwrap_or_else(|_| "Warsaw".into());

        // For flags, accept "1/0/true/false"
        let seed_demo_user = env::var("SEED_DEMO_USER")
            .map(|v| matches!(&*v.to_lowercase(), "1" | "true" | "yes"))
            .unwrap_or(false);

        Ok(Self {
            server_addr,
            database_url,
            active_sessions_url,
            bcrypt_cost,
            default_location,
            seed_demo_user,
        })
    }
}
