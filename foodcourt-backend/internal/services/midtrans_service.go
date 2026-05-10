package services

import (
	"errors"
	"fmt"
	"time"
	"log"

	"github.com/google/uuid"
	"github.com/samuelhtb/food-court-system/foodcourt-backend/internal/config"
	"github.com/samuelhtb/food-court-system/foodcourt-backend/internal/models"
	"github.com/samuelhtb/food-court-system/foodcourt-backend/internal/repositories"
	"github.com/samuelhtb/food-court-system/foodcourt-backend/pkg/midtrans"
)

type MidtransService interface {
	CreateSnapTransaction(order *models.Order) (string, error)
	ProcessNotification(notif midtrans.Notification) error
	VerifyPayment(orderID uuid.UUID) error
}

type midtransService struct {
	orderRepo  repositories.OrderRepository
	snapClient *midtrans.SnapClient
}

func NewMidtransService(orderRepo repositories.OrderRepository, cfg *config.MidtransConfig) MidtransService {
	snapClient := midtrans.NewSnapClient(cfg.ServerKey, cfg.IsProduction)
	return &midtransService{
		orderRepo:  orderRepo,
		snapClient: snapClient,
	}
}

func (s *midtransService) CreateSnapTransaction(order *models.Order) (string, error) {
	// Build unique Midtrans order ID (Max 50 chars)
	// FC- (3) + UUID (36) + - (1) + Unix (10) = 50 chars
	midtransOrderID := fmt.Sprintf("FC-%s-%d", order.ID.String(), time.Now().Unix())

	// Build item details
	var itemDetails []midtrans.ItemDetail
	for _, so := range order.SubOrders {
		for _, item := range so.OrderItems {
			itemDetails = append(itemDetails, midtrans.ItemDetail{
				ID:       item.MenuID.String(),
				Price:    int64(item.PriceAtOrder),
				Quantity: item.Quantity,
				Name:     truncate(item.Menu.Name, 50),
			})
		}
	}

	snapReq := midtrans.SnapRequest{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:     midtransOrderID,
			GrossAmount: int64(order.TotalAmount),
		},
		ItemDetails: itemDetails,
		CustomerDetails: &midtrans.CustomerDetails{
			FirstName: order.CustomerName,
		},
	}

	log.Printf("Membuat transaksi Midtrans untuk OrderID: %s, Amount: %d", order.ID, int64(order.TotalAmount))

	snapResp, err := s.snapClient.CreateTransaction(snapReq)
	if err != nil {
		log.Printf("Error dari Midtrans API: %v", err)
		return "", fmt.Errorf("midtrans error: %w", err)
	}

	log.Printf("Berhasil membuat transaksi Midtrans. Token: %s", snapResp.Token)

	// Update order with midtrans info
	err = s.orderRepo.UpdateMidtransInfo(order.ID, snapResp.Token, midtransOrderID)
	if err != nil {
		log.Printf("Gagal update info Midtrans di DB: %v", err)
		return "", fmt.Errorf("failed to update order with midtrans info: %w", err)
	}

	return snapResp.Token, nil
}

func (s *midtransService) ProcessNotification(notif midtrans.Notification) error {
	if !s.snapClient.VerifySignature(notif) {
		return errors.New("invalid midtrans signature")
	}

	order, err := s.orderRepo.FindByMidtransOrderID(notif.OrderID)
	if err != nil {
		return err
	}

	newStatus := midtrans.ResolveTransactionStatus(notif.TransactionStatus, notif.FraudStatus)
	
	return s.orderRepo.UpdatePaymentStatus(order.ID, newStatus)
}

func (s *midtransService) VerifyPayment(orderID uuid.UUID) error {
	return s.orderRepo.UpdatePaymentStatus(orderID, "paid")
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
