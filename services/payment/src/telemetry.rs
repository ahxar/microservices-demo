use std::sync::Arc;

use axum::{
    extract::State,
    http::{header, HeaderValue, StatusCode},
    response::IntoResponse,
    routing::get,
    Router,
};
use prometheus::{Encoder, IntCounterVec, Opts, Registry, TextEncoder};

#[derive(Clone)]
pub struct Telemetry {
    registry: Registry,
    grpc_requests_total: IntCounterVec,
}

impl Telemetry {
    pub fn new() -> Result<Self, prometheus::Error> {
        let registry = Registry::new();

        let grpc_requests_total = IntCounterVec::new(
            Opts::new(
                "payment_grpc_requests_total",
                "Total number of payment gRPC requests",
            ),
            &["method"],
        )?;

        registry.register(Box::new(grpc_requests_total.clone()))?;

        Ok(Self {
            registry,
            grpc_requests_total,
        })
    }

    pub fn record_grpc_request(&self, method: &str) {
        self.grpc_requests_total.with_label_values(&[method]).inc();
    }

    fn gather_metrics(&self) -> Result<Vec<u8>, prometheus::Error> {
        let metric_families = self.registry.gather();
        let mut buffer = Vec::new();
        TextEncoder::new().encode(&metric_families, &mut buffer)?;
        Ok(buffer)
    }
}

pub async fn run_http_server(
    telemetry: Arc<Telemetry>,
    metrics_port: &str,
) -> Result<(), Box<dyn std::error::Error + Send + Sync>> {
    let app = Router::new()
        .route("/metrics", get(metrics_handler))
        .route("/healthz", get(health_handler))
        .route("/readyz", get(ready_handler))
        .with_state(telemetry);

    let listener = tokio::net::TcpListener::bind(format!("0.0.0.0:{metrics_port}")).await?;
    axum::serve(listener, app).await?;

    Ok(())
}

async fn metrics_handler(State(telemetry): State<Arc<Telemetry>>) -> impl IntoResponse {
    match telemetry.gather_metrics() {
        Ok(bytes) => {
            let mut headers = axum::http::HeaderMap::new();
            headers.insert(
                header::CONTENT_TYPE,
                HeaderValue::from_static("text/plain; version=0.0.4; charset=utf-8"),
            );
            (StatusCode::OK, headers, bytes).into_response()
        }
        Err(err) => (
            StatusCode::INTERNAL_SERVER_ERROR,
            format!("failed to encode metrics: {err}"),
        )
            .into_response(),
    }
}

async fn health_handler() -> impl IntoResponse {
    (StatusCode::OK, "ok")
}

async fn ready_handler() -> impl IntoResponse {
    (StatusCode::OK, "ready")
}
