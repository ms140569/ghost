package constants

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

	// Frame queue size for inbound and outbound

	QueueSizeInbound  int = 50
	QueueSizeOutbound int = 20

	// Default Port Numbers if nothing is given

	DefaultStompPortNumber  int = 7777
	DefaultRestPortNumber   int = 7778
	DefaultTelnetPortNumber int = 7779
)
