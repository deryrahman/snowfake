package snowfake

const (
	encodeBase58Map = "123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"
)

var decodeBase58Map [256]byte

func init() {
	for i := 0; i < 256; i++ {
		decodeBase58Map[i] = 0xFF
	}
	for i := 0; i < len(encodeBase58Map); i++ {
		decodeBase58Map[encodeBase58Map[i]] = byte(i)
	}
}

// EncodeBase58 encodes uint64 to string based on flickr Base58
func EncodeBase58(id uint64) string {
	// implementation is adopted from https://github.com/bwmarrin/snowflake/blob/c09e69ae59935edf6d85799e858c26de86b04cb3/snowflake.go#L250
	if id < 58 {
		return string(encodeBase58Map[id])
	}

	b := make([]byte, 0, 11)
	for id >= 58 {
		b = append(b, encodeBase58Map[id%58])
		id /= 58
	}
	b = append(b, encodeBase58Map[id])

	for x, y := 0, len(b)-1; x < y; x, y = x+1, y-1 {
		b[x], b[y] = b[y], b[x]
	}

	return string(b)
}

// DecodeBase58 decodes string to uint64 based on flickr Base58
func DecodeBase58(str string) uint64 {
	// implementation is adopted from https://github.com/bwmarrin/snowflake/blob/c09e69ae59935edf6d85799e858c26de86b04cb3/snowflake.go#L271
	var id uint64

	b := []byte(str)
	for i := range b {
		if decodeBase58Map[b[i]] == 0xFF {
			return 0
		}
		id = id*58 + uint64(decodeBase58Map[b[i]])
	}

	return id
}
