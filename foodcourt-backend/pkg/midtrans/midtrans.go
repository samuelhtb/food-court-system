package midtrans

import (
	"bytes"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	SandboxSnapURL    = "https://app.sandbox.midtrans.com/snap/v1/transactions"
	ProductionSnapURL = "https://app.midtrans.com/snap/v1/transactions"
)

// SnapClient handles communication with the Midtrans Snap API
type SnapClient struct {
	ServerKey    string
	IsProduction bool
}

// NewSnapClient creates a new Midtrans Snap client
func NewSnapClient(serverKey string, isProduction bool) *SnapClient {
	return &SnapClient{
		ServerKey:    serverKey,
		IsProduction: isProduction,
	}
}

// TransactionDetails represents the transaction_details field
type TransactionDetails struct {
	OrderID     string `json:"order_id"`
	GrossAmount int64  `json:"gross_amount"`
}

// CustomerDetails represents customer info
type CustomerDetails struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
}

// ItemDetail represents an item in the transaction
type ItemDetail struct {
	ID       string `json:"id,omitempty"`
	Price    int64  `json:"price"`
	Quantity int    `json:"quantity"`
	Name     string `json:"name"`
}

// Callbacks represents redirect URLs after payment
type Callbacks struct {
	Finish  string `json:"finish,omitempty"`
	Error   string `json:"error,omitempty"`
	Pending string `json:"pending,omitempty"`
}

// SnapRequest is the request body for creating a Snap transaction
type SnapRequest struct {
	TransactionDetails TransactionDetails `json:"transaction_details"`
	CustomerDetails    *CustomerDetails   `json:"customer_details,omitempty"`
	ItemDetails        []ItemDetail       `json:"item_details,omitempty"`
	Callbacks          *Callbacks         `json:"callbacks,omitempty"`
}

// SnapResponse is the response from Midtrans Snap API
type SnapResponse struct {
	Token       string `json:"token"`
	RedirectURL string `json:"redirect_url"`
}

// Notification represents a payment notification from Midtrans
type Notification struct {
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	TransactionID     string `json:"transaction_id"`
	StatusMessage     string `json:"status_message"`
	StatusCode        string `json:"status_code"`
	SignatureKey      string `json:"signature_key"`
	OrderID           string `json:"order_id"`
	MerchantID        string `json:"merchant_id"`
	GrossAmount       string `json:"gross_amount"`
	FraudStatus       string `json:"fraud_status"`
	PaymentType       string `json:"payment_type"`
}

// getBaseURL returns the appropriate API URL based on environment
func (c *SnapClient) getBaseURL() string {
	if c.IsProduction {
		return ProductionSnapURL
	}
	return SandboxSnapURL
}

// CreateTransaction creates a new Snap transaction and returns the token
func (c *SnapClient) CreateTransaction(req SnapRequest) (*SnapResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.getBaseURL(), bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.SetBasicAuth(c.ServerKey, "")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to Midtrans: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("midtrans API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var snapResp SnapResponse
	if err := json.Unmarshal(respBody, &snapResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &snapResp, nil
}

// VerifySignature verifies the signature key from a Midtrans notification
// Signature = SHA512(order_id + status_code + gross_amount + server_key)
func (c *SnapClient) VerifySignature(notif Notification) bool {
	raw := notif.OrderID + notif.StatusCode + notif.GrossAmount + c.ServerKey
	hash := sha512.Sum512([]byte(raw))
	expected := fmt.Sprintf("%x", hash)
	return expected == notif.SignatureKey
}

// ResolveTransactionStatus maps Midtrans transaction_status to our order status
func ResolveTransactionStatus(transactionStatus, fraudStatus string) string {
	switch transactionStatus {
	case "capture":
		if fraudStatus == "accept" {
			return "paid"
		}
		return "pending"
	case "settlement":
		return "paid"
	case "pending":
		return "pending"
	case "deny", "cancel", "expire":
		return "cancelled"
	case "refund", "partial_refund":
		return "cancelled"
	default:
		return "pending"
	}
}
