package globals

const (
	DefaultBufferSize int = 16384

	GhostServerName    string = "Ghost STOMP server"
	GhostVersionNumber string = "0.1"

	// header related limitations

	MaxHeaderSize      int = 512
	MaxHaederKeyLength int = 64
	MaxHaederValLength int = 256

	// heartbeat related constants

	HeartbeatsSendingInterval int = 3000 // in milliseconds, 3s
	HeartbeatsMinimalInterval int = 100
)
