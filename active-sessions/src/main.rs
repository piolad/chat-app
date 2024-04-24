use tonic::{transport::Server, Request, Response, Status};
use async_trait::async_trait;
use redis::AsyncCommands;

use active_sessions::{UserData, UserDataResponse, ActiveSessions};

pub mod active_sessions {
    tonic::include_proto!("active-sessions");
}

#[derive(Default)]
pub struct ActiveSessionsService;

#[async_trait]
impl ActiveSessions for ActiveSessionsService {
    async fn add_user(
        &self,
        request: Request<UserData>,             //Userdata from auth-service
    )-> Result<Response<String>, Status> {      

        let user_data = request.into_inner();

        let session_token = generate_session_token();
        
        let redis_url = "redis://redis:6379/";
        let redis_client = redis::Client::open(redis_url)
            .map_err(|e| {
                eprintln!("Failed to connect to Redis: {:?}", e);
                Status::internal("Failed to connect to Redis")
            })?;

        let mut redis_con = redis_client.get_async_connection().await
            .map_err(|e| {
                eprintln!("Failed to get Redis connection: {:?}", e);
                Status::internal("Failed to get Redis connection")
            })?;

        let _: () = redis_con.hset_multiple(
            &user_data.username,
            &[
                ("username", &user_data.username),
                ("email", &user_data.email),
                ("key", &user_data.key),
                ("location", &user_data.location),
                ("session_token", &session_token),
            ]
        ).await 
            .map_err(|e| {
                eprintln!("Failed to set data in Redis: {:?}", e);
                Status::internal("Failed to set data in Redis")
            })?;
        
        // this is done to correctly format the gRPC response
        let response = UserDataResponse{
            session_token: session_token.clone(),
        }; 

        Ok(Response::new(response))
    } 
}



fn generate_session_token() -> String {
    //random token generation 
    use rand::distributions::Alphanumeric;
    use rand::{thread_rng, Rng};
    thread_rng().sample_iter(&Alphanumeric).take(30).collect()
}


#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr = "[::1]:50051".parse()?;
    let active_sessions_service = ActiveSessionsService::default();

    Server::builder()
        .add_service(ActiveSessionsServer::new(active_sessions_service))
        .serve(addr)
        .await?;
    Ok(())
}
