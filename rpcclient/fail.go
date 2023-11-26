package rpcclient

type FailMode int

const (
	Failover FailMode = iota
	Failfast
	Failretry
)
