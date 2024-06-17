use rsa::{RsaPrivateKey, RsaPublicKey};
use rand::rngs::OsRng;

// Function to generate RSA keypair
pub fn generate_rsa_keypair() -> (RsaPrivateKey, RsaPublicKey) {
    let mut rng = OsRng;
    let bits = 2048; // Choose an appropriate key size
    let private_key = RsaPrivateKey::new(&mut rng, bits).expect("Failed to generate private key");
    let public_key = RsaPublicKey::from(&private_key);

    // println!("Private Key Modulus: {:?}", private_key);
    // println!("Public Key Modulus: {:?}", public_key);
    
    (private_key, public_key)
}
