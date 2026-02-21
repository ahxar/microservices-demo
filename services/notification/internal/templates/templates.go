package templates

import (
	"bytes"
	"fmt"
	"html/template"
)

// OrderConfirmationData holds data for order confirmation emails
type OrderConfirmationData struct {
	FirstName       string
	OrderID         string
	OrderDate       string
	Items           []OrderItem
	Subtotal        string
	Shipping        string
	Tax             string
	Total           string
	ShippingAddress Address
}

type OrderItem struct {
	ProductName string
	Quantity    int32
	UnitPrice   string
	TotalPrice  string
}

type Address struct {
	Street  string
	City    string
	State   string
	ZipCode string
	Country string
}

// ShippingUpdateData holds data for shipping update emails
type ShippingUpdateData struct {
	FirstName      string
	OrderID        string
	TrackingNumber string
	Carrier        string
	Status         string
	EstimatedDate  string
}

// WelcomeData holds data for welcome emails
type WelcomeData struct {
	FirstName string
	Email     string
}

// PasswordResetData holds data for password reset emails
type PasswordResetData struct {
	FirstName  string
	ResetToken string
	ResetURL   string
	ExpiresIn  string
}

const orderConfirmationTemplate = `
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #4CAF50; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background-color: #f9f9f9; }
        .order-details { background-color: white; padding: 15px; margin: 15px 0; border-radius: 5px; }
        .item { border-bottom: 1px solid #eee; padding: 10px 0; }
        .total { font-size: 1.2em; font-weight: bold; margin-top: 15px; padding-top: 15px; border-top: 2px solid #4CAF50; }
        .footer { text-align: center; padding: 20px; color: #666; font-size: 0.9em; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Order Confirmation</h1>
        </div>
        <div class="content">
            <p>Hi {{.FirstName}},</p>
            <p>Thank you for your order! We're getting it ready to ship.</p>

            <div class="order-details">
                <h2>Order #{{.OrderID}}</h2>
                <p><strong>Order Date:</strong> {{.OrderDate}}</p>

                <h3>Items:</h3>
                {{range .Items}}
                <div class="item">
                    <p><strong>{{.ProductName}}</strong></p>
                    <p>Quantity: {{.Quantity}} × {{.UnitPrice}} = {{.TotalPrice}}</p>
                </div>
                {{end}}

                <div class="total">
                    <p>Subtotal: {{.Subtotal}}</p>
                    <p>Shipping: {{.Shipping}}</p>
                    <p>Tax: {{.Tax}}</p>
                    <p style="font-size: 1.3em;">Total: {{.Total}}</p>
                </div>

                <h3>Shipping Address:</h3>
                <p>
                    {{.ShippingAddress.Street}}<br>
                    {{.ShippingAddress.City}}, {{.ShippingAddress.State}} {{.ShippingAddress.ZipCode}}<br>
                    {{.ShippingAddress.Country}}
                </p>
            </div>
        </div>
        <div class="footer">
            <p>This is an automated email. Please do not reply.</p>
        </div>
    </div>
</body>
</html>
`

const shippingUpdateTemplate = `
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #2196F3; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background-color: #f9f9f9; }
        .tracking { background-color: white; padding: 15px; margin: 15px 0; border-radius: 5px; }
        .footer { text-align: center; padding: 20px; color: #666; font-size: 0.9em; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Shipping Update</h1>
        </div>
        <div class="content">
            <p>Hi {{.FirstName}},</p>
            <p>Your order has been shipped!</p>

            <div class="tracking">
                <h2>Order #{{.OrderID}}</h2>
                <p><strong>Status:</strong> {{.Status}}</p>
                <p><strong>Carrier:</strong> {{.Carrier}}</p>
                <p><strong>Tracking Number:</strong> {{.TrackingNumber}}</p>
                {{if .EstimatedDate}}
                <p><strong>Estimated Delivery:</strong> {{.EstimatedDate}}</p>
                {{end}}
            </div>
        </div>
        <div class="footer">
            <p>This is an automated email. Please do not reply.</p>
        </div>
    </div>
</body>
</html>
`

const welcomeTemplate = `
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #9C27B0; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background-color: #f9f9f9; }
        .footer { text-align: center; padding: 20px; color: #666; font-size: 0.9em; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Welcome!</h1>
        </div>
        <div class="content">
            <p>Hi {{.FirstName}},</p>
            <p>Welcome to our platform! We're excited to have you with us.</p>
            <p>Your account ({{.Email}}) has been successfully created.</p>
            <p>Start exploring our products and enjoy shopping!</p>
        </div>
        <div class="footer">
            <p>This is an automated email. Please do not reply.</p>
        </div>
    </div>
</body>
</html>
`

const passwordResetTemplate = `
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #FF5722; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background-color: #f9f9f9; }
        .button { display: inline-block; padding: 12px 24px; background-color: #FF5722; color: white; text-decoration: none; border-radius: 5px; margin: 20px 0; }
        .footer { text-align: center; padding: 20px; color: #666; font-size: 0.9em; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Password Reset Request</h1>
        </div>
        <div class="content">
            <p>Hi {{.FirstName}},</p>
            <p>We received a request to reset your password.</p>
            <p>Click the button below to reset your password:</p>
            <p style="text-align: center;">
                <a href="{{.ResetURL}}" class="button">Reset Password</a>
            </p>
            <p>This link will expire in {{.ExpiresIn}}.</p>
            <p>If you didn't request this, please ignore this email.</p>
        </div>
        <div class="footer">
            <p>This is an automated email. Please do not reply.</p>
        </div>
    </div>
</body>
</html>
`

// RenderOrderConfirmation renders the order confirmation email
func RenderOrderConfirmation(data OrderConfirmationData) (string, error) {
	tmpl, err := template.New("order_confirmation").Parse(orderConfirmationTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// RenderShippingUpdate renders the shipping update email
func RenderShippingUpdate(data ShippingUpdateData) (string, error) {
	tmpl, err := template.New("shipping_update").Parse(shippingUpdateTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// RenderWelcome renders the welcome email
func RenderWelcome(data WelcomeData) (string, error) {
	tmpl, err := template.New("welcome").Parse(welcomeTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// RenderPasswordReset renders the password reset email
func RenderPasswordReset(data PasswordResetData) (string, error) {
	tmpl, err := template.New("password_reset").Parse(passwordResetTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
