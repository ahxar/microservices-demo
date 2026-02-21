use tonic::transport::Server;
use tracing_subscriber;

mod config;
mod repository;
mod server;
mod service;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    tracing_subscriber::fmt::init();

    let config = config::Config::from_env();

    let repo = repository::PaymentRepository::new(&config.database_url).await?;

    let payment_service = service::PaymentService::new(repo);

    let addr = format!("0.0.0.0:{}", config.port).parse()?;

    tracing::info!("Payment Service starting on {}", addr);

    Server::builder()
        .add_service(server::create_payment_service(payment_service))
        .serve(addr)
        .await?;

    Ok(())
}
