package impl

import (
	"database/sql"

	"github.com/infraboard/mcube/app"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"google.golang.org/grpc"

	"github.com/staryjie/cmdb/apps/host"
	"github.com/staryjie/cmdb/apps/secret"
	"github.com/staryjie/cmdb/conf"
)

var (
	// Service服务实例
	svr = &service{}
)

type service struct {
	db  *sql.DB
	log logger.Logger

	host host.ServiceServer
	secret.UnimplementedServiceServer
}

func (s *service) Name() string {
	return secret.AppName
}

func (s *service) Config() error {
	db, err := conf.C().MySQL.GetDB()
	if err != nil {
		return err
	}

	s.log = zap.L().Named(s.Name())
	s.db = db
	s.host = app.GetGrpcApp(host.AppName).(host.ServiceServer)

	return nil
}

func (s *service) Registry(server *grpc.Server) {
	secret.RegisterServiceServer(server, svr)
}

func init() {
	app.RegistryGrpcApp(svr)
}
