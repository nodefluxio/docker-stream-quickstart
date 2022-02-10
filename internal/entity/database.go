package entity

// PsqlDBConnOption is struct for option psql database
type PsqlDBConnOption struct {
	URL                 string
	MaxIdleConn         string
	MaxOpenConn         string
	MaxLifetimeInMinute string
}
