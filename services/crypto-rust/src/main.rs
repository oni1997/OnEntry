mod crypto;
mod password;
mod types;

use axum::{
    routing::post,
    Router, Server,
};
use std::net::SocketAddr;
use tracing_subscriber;

use crate::types::{EncryptRequest, GeneratePasswordRequest, HashPasswordRequest, VerifyPasswordRequest};

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt::init();
    let app = Router::new()
        .route("/health", axum::routing::get(health))
        .route("/hash-password", post(hash_password))
        .route("/verify-password", post(verify_password))
        .route("/derive-master-key", post(derive_master_key))
        .route("/encrypt", post(encrypt))
        .route("/decrypt", post(decrypt))
        .route("/generate-password", post(generate_password));

    let addr = SocketAddr::from(([0, 0, 0, 0], 8083));
    println!("Crypto service listening on {}", addr);
    Server::bind(&addr)
        .serve(app.into_make_service())
        .await
        .unwrap();
}

async fn health() -> axum::Json<serde_json::Value> {
    axum::Json(serde_json::json!({"status": "ok"}))
}

async fn hash_password(axum::Json(req): axum::Json<HashPasswordRequest>) -> axum::Json<serde_json::Value> {
    match crypto::hash_password(&req.password) {
        Ok(result) => axum::Json(serde_json::json!({
            "hash": result.hash,
            "salt": result.salt
        })),
        Err(_) => axum::Json(serde_json::json!({"error": "Failed to hash password"})),
    }
}

async fn verify_password(axum::Json(req): axum::Json<VerifyPasswordRequest>) -> axum::Json<serde_json::Value> {
    match crypto::verify_password(&req.password, &req.hash) {
        Ok(valid) => axum::Json(serde_json::json!({"valid": valid})),
        Err(_) => axum::Json(serde_json::json!({"error": "Failed to verify password"})),
    }
}

async fn derive_master_key(axum::Json(req): axum::Json<types::DeriveMasterKeyRequest>) -> axum::Json<serde_json::Value> {
    match crypto::derive_master_key(&req.password, &req.salt) {
        Ok(key) => axum::Json(serde_json::json!({"master_key": key})),
        Err(_) => axum::Json(serde_json::json!({"error": "Failed to derive master key"})),
    }
}

async fn encrypt(axum::Json(req): axum::Json<EncryptRequest>) -> axum::Json<serde_json::Value> {
    match crypto::encrypt(&req.plaintext, &req.key) {
        Ok(result) => axum::Json(serde_json::json!({
            "ciphertext": result.ciphertext,
            "nonce": result.nonce
        })),
        Err(_) => axum::Json(serde_json::json!({"error": "Failed to encrypt"})),
    }
}

async fn decrypt(axum::Json(req): axum::Json<EncryptRequest>) -> axum::Json<serde_json::Value> {
    match crypto::decrypt(&req.ciphertext, &req.nonce, &req.key) {
        Ok(plaintext) => axum::Json(serde_json::json!({"plaintext": plaintext})),
        Err(_) => axum::Json(serde_json::json!({"error": "Failed to decrypt"})),
    }
}

async fn generate_password(axum::Json(req): axum::Json<GeneratePasswordRequest>) -> axum::Json<serde_json::Value> {
    match password::generate(req) {
        Ok(password) => axum::Json(serde_json::json!({"password": password})),
        Err(_) => axum::Json(serde_json::json!({"error": "Failed to generate password"})),
    }
}