use crate::security::hash_password; // verify_password not used

pub async fn seed_demo_user(client: &tokio_postgres::Client, bcrypt_cost: u32) {
    // hash once
    let hashed = hash_password("8rud!", bcrypt_cost);

    // Use parameters instead of interpolating into SQL
    let add_user_query = r#"
        INSERT INTO users (email, username, hashed_password, first_name, last_name, date)
        VALUES ($1, $2, $3, $4, $5, $6)
    "#;

    // Try insert; if it violates unique constraints, just log and continue
    match client
        .execute(
            add_user_query,
            &[
                &"brud@brud.pl",
                &"brud",
                &hashed,
                &"Brudas",
                &"Brudowski",
                &"2004-01-01",
            ],
        )
        .await
    {
        Ok(_) => println!("Demo user seeded."),
        Err(e) => eprintln!("Demo user seed: {:?}", e),
    }
}
