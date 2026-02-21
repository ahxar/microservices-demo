use tonic::{Request, Response, Status};
use uuid::Uuid;

use crate::service::PaymentService;

// Include generated proto code
pub mod proto {
    pub mod common {
        pub mod v1 {
            tonic::include_proto!("common.v1");
        }
    }
    pub mod payment {
        pub mod v1 {
            tonic::include_proto!("payment.v1");
        }
    }
}

use proto::payment::v1::payment_service_server::{PaymentService as PaymentServiceTrait, PaymentServiceServer};
use proto::payment::v1::*;
use proto::common::v1 as common;

pub struct PaymentServer {
    service: PaymentService,
}

impl PaymentServer {
    pub fn new(service: PaymentService) -> Self {
        Self { service }
    }
}

#[tonic::async_trait]
impl PaymentServiceTrait for PaymentServer {
    async fn add_payment_method(
        &self,
        request: Request<AddPaymentMethodRequest>,
    ) -> Result<Response<PaymentMethod>, Status> {
        let req = request.into_inner();

        let user_id = Uuid::parse_str(&req.user_id)
            .map_err(|_| Status::invalid_argument("Invalid user ID"))?;

        let method_type = match PaymentMethodType::try_from(req.r#type) {
            Ok(PaymentMethodType::Card) => "card",
            Ok(PaymentMethodType::Bank) => "bank",
            _ => "card",
        };

        let payment_method = self
            .service
            .add_payment_method(
                user_id,
                method_type,
                &req.card_number,
                req.exp_month,
                req.exp_year,
                &req.cvv,
                req.is_default,
            )
            .await
            .map_err(|e| Status::internal(format!("Failed to add payment method: {}", e)))?;

        Ok(Response::new(PaymentMethod {
            id: payment_method.id.to_string(),
            user_id: payment_method.user_id.to_string(),
            r#type: PaymentMethodType::Card as i32,
            token: payment_method.token,
            last_four: payment_method.last_four,
            brand: payment_method.brand,
            exp_month: payment_method.exp_month,
            exp_year: payment_method.exp_year,
            is_default: payment_method.is_default,
            created_at: payment_method.created_at.to_rfc3339(),
        }))
    }

    async fn list_payment_methods(
        &self,
        request: Request<ListPaymentMethodsRequest>,
    ) -> Result<Response<ListPaymentMethodsResponse>, Status> {
        let req = request.into_inner();

        let user_id = Uuid::parse_str(&req.user_id)
            .map_err(|_| Status::invalid_argument("Invalid user ID"))?;

        let methods = self
            .service
            .list_payment_methods(user_id)
            .await
            .map_err(|e| Status::internal(format!("Failed to list payment methods: {}", e)))?;

        let payment_methods = methods
            .into_iter()
            .map(|m| PaymentMethod {
                id: m.id.to_string(),
                user_id: m.user_id.to_string(),
                r#type: PaymentMethodType::Card as i32,
                token: m.token,
                last_four: m.last_four,
                brand: m.brand,
                exp_month: m.exp_month,
                exp_year: m.exp_year,
                is_default: m.is_default,
                created_at: m.created_at.to_rfc3339(),
            })
            .collect();

        Ok(Response::new(ListPaymentMethodsResponse {
            payment_methods,
        }))
    }

    async fn delete_payment_method(
        &self,
        request: Request<DeletePaymentMethodRequest>,
    ) -> Result<Response<common::Empty>, Status> {
        let req = request.into_inner();

        let id = Uuid::parse_str(&req.id)
            .map_err(|_| Status::invalid_argument("Invalid payment method ID"))?;

        let user_id = Uuid::parse_str(&req.user_id)
            .map_err(|_| Status::invalid_argument("Invalid user ID"))?;

        self.service
            .delete_payment_method(id, user_id)
            .await
            .map_err(|e| Status::internal(format!("Failed to delete payment method: {}", e)))?;

        Ok(Response::new(common::Empty {}))
    }

    async fn charge(
        &self,
        request: Request<ChargeRequest>,
    ) -> Result<Response<ChargeResponse>, Status> {
        let req = request.into_inner();

        let order_id = Uuid::parse_str(&req.order_id)
            .map_err(|_| Status::invalid_argument("Invalid order ID"))?;

        let user_id = Uuid::parse_str(&req.user_id)
            .map_err(|_| Status::invalid_argument("Invalid user ID"))?;

        let payment_method_id = Uuid::parse_str(&req.payment_method_id)
            .map_err(|_| Status::invalid_argument("Invalid payment method ID"))?;

        let amount = req.amount.ok_or_else(|| Status::invalid_argument("Amount is required"))?;

        let (success, transaction, error_message) = self
            .service
            .charge(
                order_id,
                user_id,
                payment_method_id,
                amount.amount_cents,
                &amount.currency,
                &req.idempotency_key,
            )
            .await
            .map_err(|e| Status::internal(format!("Failed to process charge: {}", e)))?;

        Ok(Response::new(ChargeResponse {
            success,
            transaction: Some(Transaction {
                id: transaction.id.to_string(),
                order_id: transaction.order_id.to_string(),
                user_id: transaction.user_id.to_string(),
                payment_method_id: transaction.payment_method_id.map(|id| id.to_string()).unwrap_or_default(),
                amount: Some(common::Money {
                    amount_cents: transaction.amount_cents,
                    currency: transaction.currency,
                }),
                status: match transaction.status.as_str() {
                    "succeeded" => TransactionStatus::Succeeded as i32,
                    "failed" => TransactionStatus::Failed as i32,
                    "refunded" => TransactionStatus::Refunded as i32,
                    _ => TransactionStatus::Pending as i32,
                },
                r#type: TransactionType::Charge as i32,
                provider_ref: transaction.provider_ref.unwrap_or_default(),
                idempotency_key: transaction.idempotency_key.unwrap_or_default(),
                created_at: transaction.created_at.to_rfc3339(),
            }),
            error_message: error_message.unwrap_or_default(),
        }))
    }

    async fn refund(
        &self,
        request: Request<RefundRequest>,
    ) -> Result<Response<RefundResponse>, Status> {
        let req = request.into_inner();

        let transaction_id = Uuid::parse_str(&req.transaction_id)
            .map_err(|_| Status::invalid_argument("Invalid transaction ID"))?;

        let amount = req.amount.ok_or_else(|| Status::invalid_argument("Amount is required"))?;

        let (success, transaction, error_message) = self
            .service
            .refund(transaction_id, amount.amount_cents, &req.reason)
            .await
            .map_err(|e| Status::internal(format!("Failed to process refund: {}", e)))?;

        Ok(Response::new(RefundResponse {
            success,
            transaction: Some(Transaction {
                id: transaction.id.to_string(),
                order_id: transaction.order_id.to_string(),
                user_id: transaction.user_id.to_string(),
                payment_method_id: transaction.payment_method_id.map(|id| id.to_string()).unwrap_or_default(),
                amount: Some(common::Money {
                    amount_cents: transaction.amount_cents,
                    currency: transaction.currency,
                }),
                status: match transaction.status.as_str() {
                    "succeeded" => TransactionStatus::Succeeded as i32,
                    "failed" => TransactionStatus::Failed as i32,
                    "refunded" => TransactionStatus::Refunded as i32,
                    _ => TransactionStatus::Pending as i32,
                },
                r#type: TransactionType::Refund as i32,
                provider_ref: transaction.provider_ref.unwrap_or_default(),
                idempotency_key: transaction.idempotency_key.unwrap_or_default(),
                created_at: transaction.created_at.to_rfc3339(),
            }),
            error_message: error_message.unwrap_or_default(),
        }))
    }

    async fn get_transaction(
        &self,
        request: Request<GetTransactionRequest>,
    ) -> Result<Response<Transaction>, Status> {
        let req = request.into_inner();

        let id = Uuid::parse_str(&req.id)
            .map_err(|_| Status::invalid_argument("Invalid transaction ID"))?;

        let transaction = self
            .service
            .get_transaction(id)
            .await
            .map_err(|e| Status::internal(format!("Failed to get transaction: {}", e)))?
            .ok_or_else(|| Status::not_found("Transaction not found"))?;

        Ok(Response::new(Transaction {
            id: transaction.id.to_string(),
            order_id: transaction.order_id.to_string(),
            user_id: transaction.user_id.to_string(),
            payment_method_id: transaction.payment_method_id.map(|id| id.to_string()).unwrap_or_default(),
            amount: Some(common::Money {
                amount_cents: transaction.amount_cents,
                currency: transaction.currency,
            }),
            status: match transaction.status.as_str() {
                "succeeded" => TransactionStatus::Succeeded as i32,
                "failed" => TransactionStatus::Failed as i32,
                "refunded" => TransactionStatus::Refunded as i32,
                _ => TransactionStatus::Pending as i32,
            },
            r#type: match transaction.transaction_type.as_str() {
                "refund" => TransactionType::Refund as i32,
                _ => TransactionType::Charge as i32,
            },
            provider_ref: transaction.provider_ref.unwrap_or_default(),
            idempotency_key: transaction.idempotency_key.unwrap_or_default(),
            created_at: transaction.created_at.to_rfc3339(),
        }))
    }
}

pub fn create_payment_service(service: PaymentService) -> PaymentServiceServer<PaymentServer> {
    PaymentServiceServer::new(PaymentServer::new(service))
}
