use chrono::{DateTime, Utc};
use sqlx::{PgPool, Row};
use uuid::Uuid;

#[derive(Debug, Clone)]
pub struct PaymentMethod {
    pub id: Uuid,
    pub user_id: Uuid,
    pub method_type: String,
    pub token: String,
    pub last_four: String,
    pub brand: String,
    pub exp_month: i32,
    pub exp_year: i32,
    pub is_default: bool,
    pub created_at: DateTime<Utc>,
}

#[derive(Debug, Clone)]
pub struct Transaction {
    pub id: Uuid,
    pub order_id: Uuid,
    pub user_id: Uuid,
    pub payment_method_id: Option<Uuid>,
    pub amount_cents: i64,
    pub currency: String,
    pub status: String,
    pub transaction_type: String,
    pub provider_ref: Option<String>,
    pub idempotency_key: Option<String>,
    pub created_at: DateTime<Utc>,
}

pub struct PaymentRepository {
    pool: PgPool,
}

impl PaymentRepository {
    pub async fn new(database_url: &str) -> Result<Self, sqlx::Error> {
        let pool = PgPool::connect(database_url).await?;
        Ok(Self { pool })
    }

    pub async fn add_payment_method(
        &self,
        user_id: Uuid,
        method_type: &str,
        token: &str,
        last_four: &str,
        brand: &str,
        exp_month: i32,
        exp_year: i32,
        is_default: bool,
    ) -> Result<PaymentMethod, sqlx::Error> {
        let mut tx = self.pool.begin().await?;

        // If this is default, unset other defaults
        if is_default {
            sqlx::query("UPDATE payment_methods SET is_default = false WHERE user_id = $1")
                .bind(user_id)
                .execute(&mut *tx)
                .await?;
        }

        let row = sqlx::query(
            r#"
            INSERT INTO payment_methods (user_id, type, token, last_four, brand, exp_month, exp_year, is_default)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
            RETURNING id, user_id, type, token, last_four, brand, exp_month, exp_year, is_default, created_at
            "#,
        )
        .bind(user_id)
        .bind(method_type)
        .bind(token)
        .bind(last_four)
        .bind(brand)
        .bind(exp_month)
        .bind(exp_year)
        .bind(is_default)
        .fetch_one(&mut *tx)
        .await?;

        tx.commit().await?;

        Ok(PaymentMethod {
            id: row.get("id"),
            user_id: row.get("user_id"),
            method_type: row.get("type"),
            token: row.get("token"),
            last_four: row.get("last_four"),
            brand: row.get("brand"),
            exp_month: row.get("exp_month"),
            exp_year: row.get("exp_year"),
            is_default: row.get("is_default"),
            created_at: row.get("created_at"),
        })
    }

    pub async fn list_payment_methods(&self, user_id: Uuid) -> Result<Vec<PaymentMethod>, sqlx::Error> {
        let rows = sqlx::query(
            r#"
            SELECT id, user_id, type, token, last_four, brand, exp_month, exp_year, is_default, created_at
            FROM payment_methods
            WHERE user_id = $1
            ORDER BY is_default DESC, created_at DESC
            "#,
        )
        .bind(user_id)
        .fetch_all(&self.pool)
        .await?;

        let methods = rows
            .into_iter()
            .map(|row| PaymentMethod {
                id: row.get("id"),
                user_id: row.get("user_id"),
                method_type: row.get("type"),
                token: row.get("token"),
                last_four: row.get("last_four"),
                brand: row.get("brand"),
                exp_month: row.get("exp_month"),
                exp_year: row.get("exp_year"),
                is_default: row.get("is_default"),
                created_at: row.get("created_at"),
            })
            .collect();

        Ok(methods)
    }

    pub async fn delete_payment_method(&self, id: Uuid, user_id: Uuid) -> Result<(), sqlx::Error> {
        sqlx::query("DELETE FROM payment_methods WHERE id = $1 AND user_id = $2")
            .bind(id)
            .bind(user_id)
            .execute(&self.pool)
            .await?;

        Ok(())
    }

    pub async fn create_transaction(
        &self,
        order_id: Uuid,
        user_id: Uuid,
        payment_method_id: Option<Uuid>,
        amount_cents: i64,
        currency: &str,
        status: &str,
        transaction_type: &str,
        provider_ref: Option<String>,
        idempotency_key: Option<String>,
    ) -> Result<Transaction, sqlx::Error> {
        let row = sqlx::query(
            r#"
            INSERT INTO transactions (order_id, user_id, payment_method_id, amount_cents, currency, status, type, provider_ref, idempotency_key)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
            RETURNING id, order_id, user_id, payment_method_id, amount_cents, currency, status, type, provider_ref, idempotency_key, created_at
            "#,
        )
        .bind(order_id)
        .bind(user_id)
        .bind(payment_method_id)
        .bind(amount_cents)
        .bind(currency)
        .bind(status)
        .bind(transaction_type)
        .bind(provider_ref)
        .bind(idempotency_key)
        .fetch_one(&self.pool)
        .await?;

        Ok(Transaction {
            id: row.get("id"),
            order_id: row.get("order_id"),
            user_id: row.get("user_id"),
            payment_method_id: row.get("payment_method_id"),
            amount_cents: row.get("amount_cents"),
            currency: row.get("currency"),
            status: row.get("status"),
            transaction_type: row.get("type"),
            provider_ref: row.get("provider_ref"),
            idempotency_key: row.get("idempotency_key"),
            created_at: row.get("created_at"),
        })
    }

    pub async fn get_transaction(&self, id: Uuid) -> Result<Option<Transaction>, sqlx::Error> {
        let row = sqlx::query(
            r#"
            SELECT id, order_id, user_id, payment_method_id, amount_cents, currency, status, type, provider_ref, idempotency_key, created_at
            FROM transactions
            WHERE id = $1
            "#,
        )
        .bind(id)
        .fetch_optional(&self.pool)
        .await?;

        Ok(row.map(|r| Transaction {
            id: r.get("id"),
            order_id: r.get("order_id"),
            user_id: r.get("user_id"),
            payment_method_id: r.get("payment_method_id"),
            amount_cents: r.get("amount_cents"),
            currency: r.get("currency"),
            status: r.get("status"),
            transaction_type: r.get("type"),
            provider_ref: r.get("provider_ref"),
            idempotency_key: r.get("idempotency_key"),
            created_at: r.get("created_at"),
        }))
    }

    pub async fn check_idempotency(&self, key: &str) -> Result<Option<Transaction>, sqlx::Error> {
        let row = sqlx::query(
            r#"
            SELECT id, order_id, user_id, payment_method_id, amount_cents, currency, status, type, provider_ref, idempotency_key, created_at
            FROM transactions
            WHERE idempotency_key = $1
            "#,
        )
        .bind(key)
        .fetch_optional(&self.pool)
        .await?;

        Ok(row.map(|r| Transaction {
            id: r.get("id"),
            order_id: r.get("order_id"),
            user_id: r.get("user_id"),
            payment_method_id: r.get("payment_method_id"),
            amount_cents: r.get("amount_cents"),
            currency: r.get("currency"),
            status: r.get("status"),
            transaction_type: r.get("type"),
            provider_ref: r.get("provider_ref"),
            idempotency_key: r.get("idempotency_key"),
            created_at: r.get("created_at"),
        }))
    }
}
