use std::env;

pub struct Config {
    pub port: String,
    pub database_url: String,
}

impl Config {
    pub fn from_env() -> Self {
        Self {
            port: env::var("PORT").unwrap_or_else(|_| "50056".to_string()),
            database_url: env::var("DATABASE_URL").unwrap_or_else(|_| {
                "postgres://postgres:postgres@localhost:5432/payment_db".to_string()
            }),
        }
    }
}
