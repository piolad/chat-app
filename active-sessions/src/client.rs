use tonic::{transport::Channel, Request, Response, Status};
use active_sessions::{UserData, UserDataResponse, ActiveSessions};

pub mod active_sessions {
    tonic::include_proto!("active-sessions");
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let channel = Channel::from_static("http://[::1]:50052")
        .connect()
        .await?;

    let mut client = ActiveSessionsClient::new(channel);

    let request = tonic::Request::new(
        UserData{
            username: "Wiktor".to_owned(),
            email: "abc@gmail.com".to_owned(),
            key: "1234".to_owned(),
            location: "Warsaw".to_owned(), 
        }
    );

    let response = client.add_user(request).await?; 

    println!("RESPONSE = {:?}", response); 

    Ok(())
}