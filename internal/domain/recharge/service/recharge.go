package service

import (
	"context"

	"finance/internal/common/errors"
	"finance/internal/domain/recharge/entity"
	rechargeRepo "finance/internal/domain/recharge/repository"
	userRepo "finance/internal/domain/user/repository"
	supplierRepo "finance/internal/domain/supplier/repository"
)

// RechargeService handles recharge business logic
type RechargeService struct {
	rechargeRepo rechargeRepo.RechargeRepository
	userRepo     userRepo.UserRepository
	supplierRepo supplierRepo.SupplierRepository
}

// NewRechargeService creates a new recharge service
func NewRechargeService(rechargeRepo rechargeRepo.RechargeRepository, userRepo userRepo.UserRepository, supplierRepo supplierRepo.SupplierRepository) *RechargeService {
	return &RechargeService{
		rechargeRepo: rechargeRepo,
		userRepo:     userRepo,
		supplierRepo: supplierRepo,
	}
}

// ProcessUserRecharge processes a user recharge
func (s *RechargeService) ProcessUserRecharge(ctx context.Context, userID int64, amount float64, method string, remark string) error {
	record := entity.NewRechargeRecord(userID, amount, entity.RechargeTypeUser, entity.RechargeMethod(method), remark)
	
	if err := s.rechargeRepo.Create(ctx, record); err != nil {
		return err
	}
	
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return errors.ErrUserNotFound
	}
	
	user.Recharge(amount)
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}
	
	record.Complete()
	return s.rechargeRepo.Update(ctx, record)
}

// ProcessSupplierRecharge processes a supplier recharge
func (s *RechargeService) ProcessSupplierRecharge(ctx context.Context, supplierID int64, amount float64, method string, remark string) error {
	record := entity.NewRechargeRecord(supplierID, amount, entity.RechargeTypeSupplier, entity.RechargeMethod(method), remark)
	
	if err := s.rechargeRepo.Create(ctx, record); err != nil {
		return err
	}
	
	supplier, err := s.supplierRepo.GetByID(ctx, supplierID)
	if err != nil {
		return errors.ErrSupplierNotFound
	}
	
	supplier.AddBalance(amount)
	if err := s.supplierRepo.Update(ctx, supplier); err != nil {
		return err
	}
	
	record.Complete()
	return s.rechargeRepo.Update(ctx, record)
}

// ListUserRecharges lists user recharges
func (s *RechargeService) ListUserRecharges(ctx context.Context, userID int64, limit, offset int) ([]*entity.RechargeRecord, error) {
	return s.rechargeRepo.ListByUserID(ctx, userID, limit, offset)
}
