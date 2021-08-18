package enrollment_test

// import (
// 	"context"
// 	"errors"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// 	"gitlab.com/nodefluxio/cuan/internal/entity"
// 	"gitlab.com/nodefluxio/cuan/internal/presenter"
// 	"gitlab.com/nodefluxio/cuan/internal/repository/mocks"
// 	"gitlab.com/nodefluxio/cuan/internal/usecase/tier"
// 	logutil "gitlab.com/nodefluxio/goutils/pkg/log"
// )

// func TestCreate(t *testing.T) {
// 	logutil.Init("info")
// 	tierRepo := &mocks.Tier{}
// 	tierService := tier.ServiceImpl{
// 		TierRepo: tierRepo,
// 	}
// 	t.Run("Success Create Tier", func(t *testing.T) {
// 		tierRepo.On("Create", mock.Anything, mock.Anything).Once().Return(nil)
// 		err := tierService.Create(context.Background(), &presenter.Tier{
// 			Name:        "Tier 1",
// 			IsUnlimited: true,
// 			MaxHit:      710000,
// 		})
// 		assert.Nil(t, err)
// 	})
// 	t.Run("Failed Create Tier", func(t *testing.T) {
// 		tierRepo.On("Create", mock.Anything, mock.Anything).Once().Return(errors.New("error on create tier"))
// 		err := tierService.Create(context.Background(), &presenter.Tier{})
// 		assert.NotNil(t, err)
// 	})
// }

// func TestUpdate(t *testing.T) {
// 	logutil.Init("info")
// 	tierRepo := &mocks.Tier{}
// 	priceCacheRepo := &mocks.PriceCache{}
// 	priceRepo := &mocks.Price{}
// 	tierService := tier.ServiceImpl{
// 		TierRepo:       tierRepo,
// 		PriceCacheRepo: priceCacheRepo,
// 		PriceRepo:      priceRepo,
// 	}
// 	tierID := uint64(1)
// 	t.Run("Success Update Tier", func(t *testing.T) {
// 		tierRepo.On("Update", mock.Anything, mock.Anything).Once().Return(nil)
// 		price := make([]entity.Price, 0)
// 		priceRepo.On("GetPriceByTierID", mock.Anything, mock.Anything).Once().Return(&price, errors.New("record not found"))
// 		priceCacheRepo.On("ScanByAnalyticID", mock.Anything, mock.Anything).Once().Return([]string{}, nil)
// 		priceCacheRepo.On("DeleteBulkKey", mock.Anything, mock.Anything).Once().Return(nil)
// 		err := tierService.Update(context.Background(), tierID, &presenter.Tier{
// 			Name:        "Tier 1",
// 			IsUnlimited: true,
// 			MaxHit:      710000,
// 		})
// 		assert.Nil(t, err)
// 	})
// 	t.Run("Failed Update Tier", func(t *testing.T) {
// 		tierRepo.On("Update", mock.Anything, mock.Anything).Once().Return(errors.New("error on update tier"))
// 		err := tierService.Update(context.Background(), tierID, &presenter.Tier{})
// 		assert.NotNil(t, err)
// 	})
// }

// func TestDelete(t *testing.T) {
// 	logutil.Init("info")
// 	tierRepo := &mocks.Tier{}
// 	priceRepo := &mocks.Price{}
// 	tierService := tier.ServiceImpl{
// 		TierRepo:  tierRepo,
// 		PriceRepo: priceRepo,
// 	}
// 	t.Run("Success Delete Tier", func(t *testing.T) {
// 		price := make([]entity.Price, 0)
// 		priceRepo.On("GetPriceByTierID", mock.Anything, mock.Anything).Once().Return(&price, errors.New("record not found"))
// 		tierRepo.On("Delete", mock.Anything, mock.Anything).Once().Return(nil)
// 		err := tierService.Delete(context.Background(), 8)
// 		assert.Nil(t, err)
// 	})

// 	t.Run("Failed Delete Tier", func(t *testing.T) {
// 		price := make([]entity.Price, 0)
// 		priceRepo.On("GetPriceByTierID", mock.Anything, mock.Anything).Once().Return(&price, errors.New("record not found"))
// 		tierRepo.On("Delete", mock.Anything, mock.Anything).Once().Return(errors.New("error on delete tier"))
// 		err := tierService.Delete(context.Background(), 8)
// 		assert.NotNil(t, err)
// 	})
// }

// func TestGetAllTierByAnalyticID(t *testing.T) {
// 	logutil.Init("info")
// 	tierRepo := &mocks.Tier{}
// 	tierService := tier.ServiceImpl{
// 		TierRepo: tierRepo,
// 	}
// 	t.Run("Success Get Tier By Analytic ID", func(t *testing.T) {
// 		tierRepo.On("GetAllTierByAnalyticID", mock.Anything, mock.Anything).Once().Return([]*presenter.TierWithPrice{}, nil)
// 		res, err := tierService.GetAllTierByAnalyticID(context.Background(), 1)
// 		assert.NotNil(t, res)
// 		assert.Nil(t, err)
// 	})
// 	t.Run("Failed Get Tier By Analytic ID", func(t *testing.T) {
// 		tierRepo.On("GetAllTierByAnalyticID", mock.Anything, mock.Anything).Once().Return(nil, errors.New("error on get tier by analytic id"))
// 		res, err := tierService.GetAllTierByAnalyticID(context.Background(), 1)
// 		assert.Nil(t, res)
// 		assert.NotNil(t, err)
// 	})
// }

// func TestGetTierByID(t *testing.T) {
// 	logutil.Init("info")
// 	tierRepo := &mocks.Tier{}
// 	tierService := tier.ServiceImpl{
// 		TierRepo: tierRepo,
// 	}
// 	t.Run("Success Get Tier By ID", func(t *testing.T) {
// 		tierRepo.On("GetTierByID", mock.Anything, mock.Anything).Once().Return(&entity.Tier{}, nil)
// 		res, err := tierService.GetTierByID(context.Background(), 1)
// 		assert.NotNil(t, res)
// 		assert.Nil(t, err)
// 	})
// 	t.Run("Failed Get Tier By ID", func(t *testing.T) {
// 		tierRepo.On("GetTierByID", mock.Anything, mock.Anything).Once().Return(nil, errors.New("error on get tier by id"))
// 		res, err := tierService.GetTierByID(context.Background(), 1)
// 		assert.Nil(t, res)
// 		assert.NotNil(t, err)
// 	})
// }

// func TestGetAllTiers(t *testing.T) {
// 	tierRepo := &mocks.Tier{}
// 	tierService := tier.ServiceImpl{
// 		TierRepo: tierRepo,
// 	}
// 	t.Run("Success Get All Tier", func(t *testing.T) {
// 		tierRepo.On("GetAllTiers", mock.Anything).Once().Return([]*entity.Tier{}, nil)
// 		res, err := tierService.GetAllTiers(context.Background())
// 		assert.NotNil(t, res)
// 		assert.Nil(t, err)
// 	})
// }
