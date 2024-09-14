package entity

// LoginStatus 表示登录状态
type LoginStatus int

const (
	LoginStatusWaiting LoginStatus = iota
	LoginStatusScanning
	LoginStatusExpired
	LoginStatusSuccess
)
