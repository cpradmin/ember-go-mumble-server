package mumble

// Channel represents a Mumble channel in the channel tree.
// Channel ID 0 is the root channel.
type Channel struct {
	ID          uint32
	ParentID    uint32
	Name        string
	Description string
	Position    int32
	MaxUsers    uint32
	IsTemporary bool
	Links       []uint32
	InheritACL  bool // inherit ACLs from parent; false = override
	// RequireTunnel forces all voice in this channel over the TLS tunnel
	// (TLS-grade AES-256-GCM) rather than direct UDP, for any client.
	RequireTunnel bool
}
