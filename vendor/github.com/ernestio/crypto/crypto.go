package crypto

// Crypto : Interface for each encrypt tool
type Crypto interface {
	Encrypt([]byte, []byte) ([]byte, error)
	Dencrypt([]byte, []byte) ([]byte, error)
	KeyValidation([]byte) bool
}
