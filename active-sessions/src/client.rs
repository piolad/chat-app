use tonic::{transport::Channel, Request, Response, Status};

pub mod active_sessions {
    tonic::include_proto!("active_sessions");
}

use active_sessions::{UserData, UserDataResponse};
use active_sessions::active_sessions_client::ActiveSessionsClient;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    println!("TEST BEFORE CONECTION TO CHANNEL");
    let channel = Channel::from_static("http://localhost:50052")
        .connect()
        .await?;

    let mut client = ActiveSessionsClient::new(channel);

    println!("TEST AFTER CONECTION TO CHANNEL");

    let request = tonic::Request::new(
        UserData {
            username: "Wiktor".to_owned(),
            email: "abc@gmail.com".to_owned(),
            key: "1234".to_owned(),
            location: "Warsaw".to_owned(),
        }
    );

    println!("Client TEST");
    let response = client.add_user(request).await?;

    println!("RESPONSE = {:?}", response);

    Ok(())
}
