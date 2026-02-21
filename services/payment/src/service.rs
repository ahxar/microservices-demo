use crate::repository::{PaymentMethod, PaymentRepository, Transaction};
use sha2::{Digest, Sha256};
use uuid::Uuid;

pub struct PaymentService {
    repo: PaymentRepository,
}

impl PaymentService {
    pub fn new(repo: PaymentRepository) -> Self {
        Self { repo }
    }

    pub async fn add_payment_method(
        &self,
        user_id: Uuid,
        method_type: &str,
        card_number: &str,
        exp_month: i32,
        exp_year: i32,
        cvv: &str,
        is_default: bool,
    ) -> Result<PaymentMethod, Box<dyn std::error::Error>> {
        // Tokenize card number (never store raw)
        let token = self.tokenize_card(card_number, cvv);
        let last_four = &card_number[card_number.len() - 4..];
        let brand = self.detect_card_brand(card_number);

        let payment_method = self
            .repo
            .add_payment_method(
                user_id,
                method_type,
                &token,
                last_four,
                &brand,
                exp_month,
                exp_year,
                is_default,
            )
            .await?;

        Ok(payment_method)
    }

    pub async fn list_payment_methods(
        &self,
        user_id: Uuid,
    ) -> Result<Vec<PaymentMethod>, Box<dyn std::error::Error>> {
        let methods = self.repo.list_payment_methods(user_id).await?;
        Ok(methods)
    }

    pub async fn delete_payment_method(
        &self,
        id: Uuid,
        user_id: Uuid,
    ) -> Result<(), Box<dyn std::error::Error>> {
        self.repo.delete_payment_method(id, user_id).await?;
        Ok(())
    }

    pub async fn charge(
        &self,
        order_id: Uuid,
        user_id: Uuid,
        payment_method_id: Uuid,
        amount_cents: i64,
        currency: &str,
        idempotency_key: &str,
    ) -> Result<(bool, Transaction, Option<String>), Box<dyn std::error::Error>> {
        // Check idempotency
        if let Some(existing_tx) = self.repo.check_idempotency(idempotency_key).await? {
            return Ok((true, existing_tx, None));
        }

        // Mock payment processing
        let (success, provider_ref) = self.process_mock_payment(amount_cents);

        let status = if success { "succeeded" } else { "failed" };
        let error_message = if success {
            None
        } else {
            Some("Mock payment failure".to_string())
        };

        let transaction = self
            .repo
            .create_transaction(
                order_id,
                user_id,
                Some(payment_method_id),
                amount_cents,
                currency,
                status,
                "charge",
                Some(provider_ref),
                Some(idempotency_key.to_string()),
            )
            .await?;

        Ok((success, transaction, error_message))
    }

    pub async fn refund(
        &self,
        transaction_id: Uuid,
        amount_cents: i64,
        _reason: &str,
    ) -> Result<(bool, Transaction, Option<String>), Box<dyn std::error::Error>> {
        // Get original transaction
        let original_tx = self
            .repo
            .get_transaction(transaction_id)
            .await?
            .ok_or("Transaction not found")?;

        // Mock refund processing
        let (success, provider_ref) = self.process_mock_refund(amount_cents);

        let status = if success { "refunded" } else { "failed" };
        let error_message = if success {
            None
        } else {
            Some("Mock refund failure".to_string())
        };

        let refund_transaction = self
            .repo
            .create_transaction(
                original_tx.order_id,
                original_tx.user_id,
                original_tx.payment_method_id,
                amount_cents,
                &original_tx.currency,
                status,
                "refund",
                Some(provider_ref),
                None,
            )
            .await?;

        Ok((success, refund_transaction, error_message))
    }

    pub async fn get_transaction(
        &self,
        id: Uuid,
    ) -> Result<Option<Transaction>, Box<dyn std::error::Error>> {
        let transaction = self.repo.get_transaction(id).await?;
        Ok(transaction)
    }

    fn tokenize_card(&self, card_number: &str, cvv: &str) -> String {
        let mut hasher = Sha256::new();
        hasher.update(card_number.as_bytes());
        hasher.update(cvv.as_bytes());
        let result = hasher.finalize();
        hex::encode(result)
    }

    fn detect_card_brand(&self, card_number: &str) -> String {
        let first_digit = card_number.chars().next().unwrap_or('0');
        match first_digit {
            '4' => "Visa".to_string(),
            '5' => "Mastercard".to_string(),
            '3' => "American Express".to_string(),
            _ => "Unknown".to_string(),
        }
    }

    fn process_mock_payment(&self, amount_cents: i64) -> (bool, String) {
        let success = amount_cents < 100_000;
        let provider_ref = format!("mock_charge_{}", Uuid::new_v4());
        (success, provider_ref)
    }

    fn process_mock_refund(&self, _amount_cents: i64) -> (bool, String) {
        let provider_ref = format!("mock_refund_{}", Uuid::new_v4());
        (true, provider_ref)
    }
}
