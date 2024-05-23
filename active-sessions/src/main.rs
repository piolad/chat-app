use tonic::{transport::Server, Request, Response, Status};
use rand::distributions::Alphanumeric;
use rand::{thread_rng, Rng};
use async_trait::async_trait;
use redis::AsyncCommands;

pub mod active_sessions {
    tonic::include_proto!("active_sessions");
}

use active_sessions::{UserData, UserDataResponse, IdSessionRequest, IdSessionResponse};
use active_sessions::active_sessions_server::{ActiveSessions, ActiveSessionsServer};

#[derive(Default)]
pub struct ActiveSessionsService;


#[async_trait]
impl ActiveSessions for ActiveSessionsService {
    async fn get_session_id(
        &self, 
        request: Request<IdSessionRequest>,
    ) -> Result<Response<IdSessionResponse>, Status> {

        tokio::time::sleep(tokio::time::Duration::from_secs(0)).await;

        let response = IdSessionResponse {
            status: "OK".to_string(),
            idsession: generate_session_token().to_string(),
        };

        Ok(Response::new(response))
    }

    async fn add_user(
        &self,
        request: Request<UserData>,
    ) -> Result<Response<UserDataResponse>, Status> { 

        let user_data = request.into_inner();
        let session_token = generate_session_token();
        
        let redis_url = "redis://redis:6379/";
        let redis_client = redis::Client::open(redis_url).map_err(|e| {
            eprintln!("Failed to connect to Redis: {:?}", e);
            Status::internal("Failed to connect to Redis")
        })?;
        
        println!("Connected to Redis");

        let mut redis_con = redis_client.get_async_connection().await.map_err(|e| {
            eprintln!("Failed to get Redis connection: {:?}", e);
            Status::internal("Failed to get Redis connection")
        })?;

        println!("Obtained Redis connection");

        let _: () = redis_con.hset_multiple(
            &user_data.username,
            &[
                ("username", &user_data.username),
                ("email", &user_data.email),
                ("key", &user_data.key),
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
        
        println!("User '{}' added successfully with session token '{}'", user_data.username, session_token);

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
    let addr = "0.0.0.0:50053".parse()?;
    let active_sessions_service = ActiveSessionsService::default();

    println!("Server starting on {}", addr);

    Server::builder()
        .add_service(ActiveSessionsServer::new(active_sessions_service)) 
        .serve(addr)
        .await?;

    println!("Server started");

    Ok(())
}
