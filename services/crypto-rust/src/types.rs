#[derive(serde::Deserialize)]
pub struct HashPasswordRequest {
    pub password: String,
}

#[derive(serde::Deserialize)]
pub struct VerifyPasswordRequest {
    pub password: String,
    pub hash: String,
}

#[derive(serde::Deserialize)]
pub struct DeriveMasterKeyRequest {
    pub password: String,
    pub salt: String,
}

#[derive(serde::Deserialize)]
pub struct EncryptRequest {
    pub plaintext: String,
    pub ciphertext: String,
    pub nonce: String,
    pub key: String,
}