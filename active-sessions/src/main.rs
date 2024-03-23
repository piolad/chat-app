use tonic::{transport::Server, Request, Response, Status};
use async_trait::async_trait;
use redis::AsyncCommands;
use active_sessions::{UserData, ActiveSessions};

pub mod active_sessions {
    tonic::include_proto!("active_sessions");
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
        
        let redis_url = "redis://localhost:6379/"; 
        let redis_client = redis::Client::open("redis_url").unwrap();
        let mut redis_con = redis_client.get_async_connection().await.unwrap();

        let _: () = redis_con.hset_multiple(
            &user_data.username,
            &[
                ("email", &user_data.email),
                ("key", &user_data.key),
                ("locaion", &user_data.location),
                ("session_token", &session_token),
            ]
        ).await.unwrap();
        
        //Returns token generated for the user 
        Ok(Response::new(session_token))
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
