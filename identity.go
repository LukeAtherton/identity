package identity

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
	"time"
)

var (
	seq  uint32
	node = getNodeUint32()
)

type ID [16]byte

func NewUUID() ID {
	id, _ := NewRandomUUID()

	return id
}

func NewRandomUUID() (id ID, err error) {
	id = make([]byte, 16)

	_, err = rand.Read(id[:])
	if err != nil {
		return
	}

	return
}

// 8 bytes of UNIXNANO
// 4 bytes of hardware address
// 4 bytes of counter
func NewSequentialUUID() (id ID, err error) {
	id = make([]byte, 16)

	nano := time.Now().UnixNano()
	incr := atomic.AddUint32(&seq, 1)

	binary.BigEndian.PutUint64(id[0:], uint64(nano))
	binary.BigEndian.PutUint32(id[8:], node)
	binary.BigEndian.PutUint32(id[12:], incr)

	return
}

func getNodeUint32() uint32 {
	n := uuid.NodeID()
	return binary.BigEndian.Uint32(n)
}

func (id ID) Equals(id2 ID) bool {
	return bytes.Equal(id, id2)
}

func (u ID) Bytes() []byte {
	return []byte(u)
}

// Returns unparsed version of the generated UUID sequence.
func (id ID) String() string {
	idBytes := id.Bytes()
	if len(idBytes) != 16 {
		return ""
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", idBytes[0:4], idBytes[4:6], idBytes[6:8], idBytes[8:10], idBytes[10:])
}

func Parse(idString string) ID {
	g := strings.Replace(string(idString), "-", "", -1)
	g = strings.Replace(g, "\"", "", -1)
	b, err := hex.DecodeString(g)
	if err != nil {
		fmt.Printf("decode: error while decoding uuid: %v", err)
	}
	return b
}

//implements TextMarshaler for text encoding
func (id ID) MarshalText() (text []byte, err error) {
	return []byte(id.String()), nil
}

//implements TextUnmarshaler for text encoding
func (id *ID) UnmarshalText(text []byte) error {
	decoded := DecodeIdString(string(text))
	*id = decoded
	return nil
}

//implements JSONUnmarshaler for json dencoding
func (id *ID) UnmarshalJSON(text []byte) (err error) {
	decoded := DecodeIdString(string(text))
	*id = decoded
	return
}
