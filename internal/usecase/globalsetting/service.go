package globalsetting

import (
	"context"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/presenter"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"

)

// Service interface abstracts the controller layer and should be implemented in controller directory. Controller contains business logics and is independent of any database connection.
type Service interface {
	GetDetail(ctx context.Context) (*entity.GlobalSetting, error)
	CreateOrUpdate(ctx context.Context, postData *presenter.GlobalSettingRequest) (*entity.GlobalSetting, error)
	
}
