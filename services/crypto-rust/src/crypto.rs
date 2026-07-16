use argon2::{
    password_hash::{
        rand_core::OsRng,
        PasswordHash, PasswordHasher, PasswordVerifier, SaltString
    },
    Argon2, Algorithm, Params, Version,
};
use base64::{engine::general_purpose, Engine as _};
use aes_gcm::{
    aead::{Aead, KeyInit, OsRng as AesOsRng},
    Aes256Gcm, Nonce,
};
use hex;

pub struct HashResult {
    pub hash: String,
    pub salt: String,
}

pub struct EncryptResult {
    pub ciphertext: String,
    pub nonce: String,
}

pub fn hash_password(password: &str) -> Result<HashResult, String> {
    let salt = SaltString::generate(&mut OsRng);
    let argon2 = Argon2::default();
    let password_hash = PasswordHasher::hash_password(&argon2, password.as_bytes(), &salt)
        .map_err(|e| e.to_string())?
        .to_string();

    Ok(HashResult {
        hash: password_hash,
        salt: salt.to_string(),
    })
}

pub fn verify_password(password: &str, hash: &str) -> Result<bool, String> {
    let parsed_hash = PasswordHash::new(hash).map_err(|e| e.to_string())?;
    let argon2 = Box::leak(Box::new(Argon2::default()));
    Ok(parsed_hash
        .verify_password(&[argon2], password.as_bytes())
        .is_ok())
}

pub fn derive_master_key(password: &str, salt: &str) -> Result<String, String> {
    let salt_bytes = salt.as_bytes();
    let mut key = [0u8; 32];

    Argon2::default()
        .hash_password_into(password.as_bytes(), salt_bytes, &mut key)
        .map_err(|e| e.to_string())?;

    Ok(hex::encode(key))
}

pub fn encrypt(plaintext: &str, key_hex: &str) -> Result<EncryptResult, String> {
    let key_bytes = hex::decode(key_hex).map_err(|e| e.to_string())?;
    if key_bytes.len() != 32 {
        return Err("Invalid key length".to_string());
    }

    let cipher = Aes256Gcm::new_from_slice(&key_bytes).map_err(|e| e.to_string())?;
    let nonce_bytes = AesOsRng.next_bytes(12);
    let nonce = Nonce::from_slice(&nonce_bytes);
    let ciphertext = cipher.encrypt(nonce, plaintext.as_bytes())
        .map_err(|e| e.to_string())?;

    Ok(EncryptResult {
        ciphertext: general_purpose::STANDARD.encode(&ciphertext),
        nonce: general_purpose::STANDARD.encode(&nonce_bytes),
    })
}

pub fn decrypt(ciphertext_b64: &str, nonce_b64: &str, key_hex: &str) -> Result<String, String> {
    let key_bytes = hex::decode(key_hex).map_err(|e| e.to_string())?;
    if key_bytes.len() != 32 {
        return Err("Invalid key length".to_string());
    }

    let cipher = Aes256Gcm::new_from_slice(&key_bytes).map_err(|e| e.to_string())?;
    let ciphertext = general_purpose::STANDARD.decode(ciphertext_b64).map_err(|e| e.to_string())?;
    let nonce_bytes = general_purpose::STANDARD.decode(nonce_b64).map_err(|e| e.to_string())?;
    let nonce = Nonce::from_slice(&nonce_bytes);

    let plaintext = cipher.decrypt(nonce, ciphertext.as_ref())
        .map_err(|e| e.to_string())?;

    String::from_utf8(plaintext).map_err(|e| e.to_string())
}