package impl

const (
	sqlInsertResource = `INSERT INTO resource (
		id,resource_type,vendor,region,zone,create_at,expire_at,category,type,
		name,description,status,update_at,sync_at,sync_accout,public_ip,
		private_ip,pay_type,describe_hash,resource_hash,secret_id,domain,
		namespace,env,usage_mode
	) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);`

	// 定义的用于变更Informtion属性
	sqlUpdateResource = `UPDATE resource SET
		expire_at=?,category=?,type=?,name=?,description=?,
		status=?,update_at=?,sync_at=?,sync_accout=?,
		public_ip=?,private_ip=?,pay_type=?,describe_hash=?,resource_hash=?,
		secret_id=?,namespace=?,env=?,usage_mode=?
	WHERE id = ?`

	sqlDeleteResource = `DELETE FROM resource WHERE id = ?;`

	// SELECT r.* FROM resource r LEFT JOIN resource_tag t ON r.id=t.resource_id WHERE t.t_key='xx', t.t_value='xxx';
	sqlQueryResource = `SELECT r.* FROM resource r %s JOIN resource_tag t ON r.id = t.resource_id`

	// 	-- resourceA   t1=v1  t2=v2
	// -- resourceA  t1=v1
	// -- resourceA  t2=v2
	// -- 使用DISTINCT对字段去重
	// -- 用于分页时使用
	sqlCountResource = `SELECT COUNT(DISTINCT r.id) FROM resource r %s JOIN resource_tag t ON r.id = t.resource_id`
)