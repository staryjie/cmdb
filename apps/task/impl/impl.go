package impl

import (
	"database/sql"

	"github.com/infraboard/mcube/app"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"google.golang.org/grpc"

	"github.com/staryjie/cmdb/apps/host"
	"github.com/staryjie/cmdb/apps/secret"
	"github.com/staryjie/cmdb/apps/task"
	"github.com/staryjie/cmdb/conf"
)

var (
	// Service服务实例
	svr = &service{}
)

type service struct {
	db  *sql.DB
	log logger.Logger
	task.UnimplementedServiceServer

	secret secret.ServiceServer
	host   host.ServiceServer
}

func (s *service) Name() string {
	return task.AppName
}

func (s *service) Config() error {
	db, err := conf.C().MySQL.GetDB()
	if err != nil {
		return err
	}

	s.log = zap.L().Named(s.Name())
	s.db = db
	s.secret = &secretMock{} // 通过Mock对象来解耦依赖
	// s.secret = app.GetGrpcApp(secret.AppName).(secret.ServiceServer)
	s.host = app.GetGrpcApp(host.AppName).(host.ServiceServer)

	return nil
}

func (s *service) Registry(server *grpc.Server) {
	task.RegisterServiceServer(server, svr)
}

func init() {
	app.RegistryGrpcApp(svr)
}
