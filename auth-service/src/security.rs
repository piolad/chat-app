use rsa::{RsaPrivateKey, RsaPublicKey};
use rand::rngs::OsRng;

// Function to generate RSA keypair
pub fn generate_rsa_keypair() -> (RsaPrivateKey, RsaPublicKey) {
    let mut rng = OsRng;
    let bits = 256; // The smaller the value, the faster the key generation
    let private_key = RsaPrivateKey::new(&mut rng, bits).expect("Failed to generate private key");
    let public_key = RsaPublicKey::from(&private_key);
 
    (private_key, public_key)
}
