use tonic::{transport::Server, Request, Response, Status};
use rand::distributions::Alphanumeric;
use rand::{thread_rng, Rng};
use async_trait::async_trait;
use redis::AsyncCommands;

mod config;
use config::Config;

pub mod active_sessions {
    tonic::include_proto!("active_sessions");
}

use active_sessions::{UserData, UserDataResponse};
use active_sessions::active_sessions_server::{ActiveSessions, ActiveSessionsServer};

#[derive(Default)]
pub struct ActiveSessionsService {
    redis_url: String,
}

impl ActiveSessionsService {
    pub fn new(cfg: &Config) -> Self {
        Self {
            redis_url: cfg.redis_url.clone(),
        }
    }
}

#[async_trait]
impl ActiveSessions for ActiveSessionsService {
    async fn add_user(
        &self,
        request: Request<UserData>,
    ) -> Result<Response<UserDataResponse>, Status> { 

        let user_data = request.into_inner();
        let session_token = generate_session_token();
        
        let redis_client = redis::Client::open(self.redis_url.as_str()).map_err(|e| {
            eprintln!("Failed to connect to Redis: {:?}", e);
            Status::internal("Failed to connect to Redis")
        })?;
        
        let mut redis_con = redis_client.get_async_connection().await.map_err(|e| {
            eprintln!("Failed to get Redis connection: {:?}", e);
            Status::internal("Failed to get Redis connection")
        })?;

        let delete_str = String::from("delete");

        let _: () = redis_con.hset_multiple(
            &user_data.username,
            &[
                ("username", &user_data.username),
                ("email", &user_data.email),
                ("token", &delete_str),
                ("location", &user_data.location),
                ("session_token", &session_token),
            ]
        ).await.map_err(|e| {
            eprintln!("Failed to set data in Redis: {:?}", e);
            Status::internal("Failed to set data in Redis")
        })?;

        let response = UserDataResponse {
            session_token: session_token.clone(),
        };
        
        Ok(Response::new(response))
    }
}

fn generate_session_token() -> String {
    thread_rng()
        .sample_iter(&Alphanumeric)
        .take(30)
        .map(char::from)
        .collect()
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let cfg = Config::from_env()?;
    let addr = cfg.server_addr;

    let active_sessions_service = ActiveSessionsService::new(&cfg);

    println!("Server starting on {addr} (redis_url set: {})", !cfg.redis_url.is_empty());

    Server::builder()
        .add_service(ActiveSessionsServer::new(active_sessions_service)) 
        .serve(addr)
        .await
        .map_err(|e|{
            eprint!("gRPC server failed: {e:?}");
            e
        })?;

    println!("Server started");

    Ok(())
}
