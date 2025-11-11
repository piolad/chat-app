use std::env;
use std::net::SocketAddr;

#[derive(Debug, Clone)]
pub struct Config {
    pub server_addr: SocketAddr,
    pub redis_url: String,
}

impl Config {
    pub fn from_env() -> Result<Self, Box<dyn std::error::Error>> {
        // .env is for local only; in containers/prod, just set env vars.
        dotenv::dotenv().ok();

        let server_addr: SocketAddr = env::var("SERVER_ADDR")
            .unwrap_or_else(|_| "0.0.0.0:50053".into())
            .parse()?;

        let redis_url = env::var("REDIS_URL")
            .unwrap_or_else(|_| "redis://active-sessions-db:6379/".into());

        Ok(Self {
            server_addr,
            redis_url,
        })
    }
}
