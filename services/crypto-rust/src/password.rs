use rand::distributions::{Distribution, Uniform};
use rand::seq::SliceRandom;
use rand::thread_rng;
use rand_distr::Alphanumeric;

const UPPERCASE: &[u8] = b"ABCDEFGHIJKLMNOPQRSTUVWXYZ";
const LOWERCASE: &[u8] = b"abcdefghijklmnopqrstuvwxyz";
const NUMBERS: &[u8] = b"0123456789";
const SYMBOLS: &[u8] = b"!@#$%^&*()_+-=[]{}|;:,.<>?";
const SIMILAR: &[u8] = b"il1Lo0O";

pub fn generate(req: GeneratePasswordRequest) -> Result<String, String> {
    let length = req.length.clamp(4, 128);
    let mut rng = thread_rng();

    let mut charset = Vec::new();
    if req.uppercase { charset.extend_from_slice(UPPERCASE); }
    if req.lowercase { charset.extend_from_slice(LOWERCASE); }
    if req.numbers { charset.extend_from_slice(NUMBERS); }
    if req.symbols { charset.extend_from_slice(SYMBOLS); }

    if charset.is_empty() {
        charset.extend_from_slice(UPPERCASE);
        charset.extend_from_slice(LOWERCASE);
        charset.extend_from_slice(NUMBERS);
        charset.extend_from_slice(SYMBOLS);
    }

    let mut password = Vec::with_capacity(length);

    if req.pronounceable {
        password.extend(generate_pronounceable(length, &charset));
    } else {
        let dist = Uniform::from(0..charset.len());
        for _ in 0..length {
            password.push(charset[dist.sample(&mut rng)]);
        }
    }

    if req.exclude_similar {
        password.retain(|&c| !SIMILAR.contains(&c));
        while password.len() < length {
            let c = charset[Uniform::from(0..charset.len()).sample(&mut rng)];
            if !SIMILAR.contains(&c) {
                password.push(c);
            }
        }
    }

    password.shuffle(&mut rng);

    Ok(String::from_utf8(password).unwrap_or_default())
}

fn generate_pronounceable(length: usize, charset: &[u8]) -> Vec<u8> {
    let mut rng = thread_rng();
    let consonants = b"bcdfghjklmnpqrstvwxyz";
    let vowels = b"aeiou";
    let mut result = Vec::with_capacity(length);

    let mut i = 0;
    while i < length {
        let c = consonants[Uniform::from(0..consonants.len()).sample(&mut rng)];
        result.push(c);
        i += 1;
        if i >= length { break; }
        let v = vowels[Uniform::from(0..vowels.len()).sample(&mut rng)];
        result.push(v);
        i += 1;
        if i >= length { break; }
        let n = NUMBERS[Uniform::from(0..NUMBERS.len()).sample(&mut rng)];
        result.push(n);
        i += 1;
    }

    result.truncate(length);
    result
}

#[derive(serde::Deserialize)]
pub struct GeneratePasswordRequest {
    pub length: usize,
    pub uppercase: bool,
    pub lowercase: bool,
    pub numbers: bool,
    pub symbols: bool,
    pub exclude_similar: bool,
    pub pronounceable: bool,
}