type MessageType int

const (
    Configuration MessageType = iota // Message for configuration updates (e.g., root bridge ID)
    PortState                       // Message for port state updates
)
