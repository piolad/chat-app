fn main() -> Result<(), Box<dyn std::error::Error>> {
    tonic_build::compile_protos("protos/auth.proto")?;
    tonic_build::compile_protos("protos/active_sessions.proto")?;
    Ok(())
}