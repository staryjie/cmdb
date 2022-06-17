package impl

import (
	"context"
	"database/sql"

	"github.com/infraboard/mcube/exception"
	"github.com/infraboard/mcube/sqlbuilder"
	"github.com/staryjie/cmdb/apps/secret"
	"github.com/staryjie/cmdb/conf"
)

func (s *service) CreateSecret(ctx context.Context, req *secret.CreateSecretRequest) (
	*secret.Secret, error) {
	sct, err := secret.NewSecret(req)
	if err != nil {
		return nil, exception.NewBadRequest("validate create secret error, %s", err)
	}

	stmt, err := s.db.PrepareContext(ctx, insertSecretSQL)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// TODO: 入库之前先加密数据
	if err := sct.Data.EncryptAPISecret(conf.C().App.EncryptKey); err != nil {
		s.log.Warnf("Encrypt Api Key error: %s", err)
	}

	_, err = stmt.ExecContext(ctx,
		sct.Id, sct.CreateAt, sct.Data.Description, sct.Data.Vendor, sct.Data.Address,
		sct.Data.AllowRegionString(), sct.Data.CrendentialType, sct.Data.ApiKey,
		sct.Data.ApiSecret, sct.Data.RequestRate,
	)
	if err != nil {
		return nil, err
	}

	return sct, nil
}

func (s *service) QuerySecret(ctx context.Context, req *secret.QuerySecretRequest) (
	*secret.SecretSet, error) {
	query := sqlbuilder.NewQuery(querySecretSQL)

	if req.Keywords != "" {
		query.Where("description LIKE ? OR api_key = ?", "%"+req.Keywords+"%", req.Keywords)
	}

	querySQL, args := query.Order("create_at").Desc().Limit(req.Page.ComputeOffset(), uint(req.Page.PageSize)).BuildQuery()
	s.log.Debugf("Sql: %s, args: %v", querySQL, args)

	queryStmt, err := s.db.PrepareContext(ctx, querySQL)
	if err != nil {
		return nil, exception.NewInternalServerError("Prepare query secret errpr, %s", err.Error())
	}
	defer queryStmt.Close()

	rows, err := queryStmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, exception.NewInternalServerError(err.Error())
	}
	defer rows.Close()

	secretSet := secret.NewSecretSet()
	allowRegions := ""
	for rows.Next() {
		sct := secret.NewDefaultSecret()
		err := rows.Scan(
			&sct.Id, &sct.CreateAt, &sct.Data.Description, &sct.Data.Vendor, &sct.Data.Address,
			&allowRegions, &sct.Data.CrendentialType, &sct.Data.ApiKey, &sct.Data.ApiSecret,
			&sct.Data.RequestRate,
		)
		if err != nil {
			return nil, exception.NewInternalServerError("Query secret error, %s", err.Error())
		}

		sct.Data.LoadAllowRegionFromString(allowRegions)
		sct.Data.Desense()
		secretSet.Add(sct)
	}

	// 获取total SELECT COUNT(*) FROMT t Where ....
	countSQL, args := query.BuildCount()
	countStmt, err := s.db.PrepareContext(ctx, countSQL)
	if err != nil {
		return nil, exception.NewInternalServerError(err.Error())
	}
	defer countStmt.Close()

	err = countStmt.QueryRowContext(ctx, args...).Scan(&secretSet.Total)
	if err != nil {
		return nil, exception.NewInternalServerError(err.Error())
	}

	return secretSet, nil
}

func (s *service) DescribeSecret(ctx context.Context, req *secret.DescribeSecretRequest) (
	*secret.Secret, error) {
	query := sqlbuilder.NewQuery(querySecretSQL)
	querySQl, args := query.Where("id = ?", req.Id).BuildQuery()
	// 打印构建好的sql和参数
	s.log.Debugf("SQL: %s, args: %v", querySQl, args)

	queryStmt, err := s.db.PrepareContext(ctx, querySQl)
	if err != nil {
		return nil, exception.NewInternalServerError("Prepare Query secret error, %s", err)
	}
	defer queryStmt.Close()

	sct := secret.NewDefaultSecret()
	allowRegions := ""
	err = queryStmt.QueryRowContext(ctx, args...).Scan(
		&sct.Id, &sct.CreateAt, &sct.Data.Description, &sct.Data.Vendor, &sct.Data.Address,
		&allowRegions, &sct.Data.CrendentialType, &sct.Data.ApiKey, &sct.Data.ApiSecret,
		&sct.Data.RequestRate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewNotFound("%#v not found!", req)
		}
		return nil, exception.NewInternalServerError("Describe secret error, %s", err.Error())
	}

	sct.Data.LoadAllowRegionFromString(allowRegions)

	return sct, nil
}

func (s *service) DeleteSecret(ctx context.Context, req *secret.DeleteSecretRequest) (
	*secret.Secret, error) {
	sct, err := s.DescribeSecret(ctx, secret.NewDescribeSecretRequest(req.Id))
	if err != nil {
		return nil, err
	}

	if err := s.deleteSecret(ctx, sct); err != nil {
		return nil, err
	}

	return sct, nil
}
