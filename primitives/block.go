package primitives

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
	"io"

	"github.com/kookehs/watchmen/crypto"
)

// Block represents the common elements shared between various types.
type Block interface {
	// Block
	Balance() Amount
	Delegates() []IBAN
	Hash() (BlockHash, error)
	Previous() BlockHash
	Root() BlockHash
	Sign(*ecdsa.PrivateKey) error
	Source() BlockHash
	Timestamp() int64
	Type() BlockType
	Verify(*ecdsa.PublicKey) (bool, error)

	// Deserialization
	Deserialize(io.Reader) error
	DeserializeJSON(io.Reader) error

	// Serialization
	Serialize(io.Writer) error
	SerializeJSON(io.Writer) error

	// Conversion
	String() (string, error)
	ToJSON() (string, error)
}

// ChangeBlock represents a change in delegates.
type ChangeBlock struct {
	Hashables ChangeHashables `json:"hashables"`
	Signature Signature       `json:"signature"`
}

// NewChangeBlock creates and initializes a ChangeBlock from the given arguments.
func NewChangeBlock(amt Amount, delegates []IBAN, prev BlockHash) *ChangeBlock {
	return &ChangeBlock{
		Hashables: MakeChangeHashables(amt, delegates, prev),
	}
}

// Balance returns the balance associated with this block.
func (cb *ChangeBlock) Balance() Amount {
	return cb.Hashables.Balance
}

// Delegates returns the delegates associated with this block.
func (cb *ChangeBlock) Delegates() []IBAN {
	return cb.Hashables.Delegates
}

// Hash returns the SHA256 hash of the serialized bytes of Hashables.
func (cb *ChangeBlock) Hash() (BlockHash, error) {
	var buffer bytes.Buffer

	if err := cb.Hashables.Serialize(&buffer); err != nil {
		return BlockHashZero, err
	}

	return sha256.Sum256(buffer.Bytes()), nil
}

// Previous returns the previous hash associated with this block.
func (cb *ChangeBlock) Previous() BlockHash {
	return cb.Hashables.Previous
}

// Root returns the previous hash associated with this block.
func (cb *ChangeBlock) Root() BlockHash {
	return cb.Hashables.Previous
}

// Sign signs the block with the given private key.
func (cb *ChangeBlock) Sign(priv *ecdsa.PrivateKey) error {
	hash, err := cb.Hash()

	if err != nil {
		return err
	}

	r, s, err := crypto.Sign(hash[:], priv)

	if err != nil {
		return err
	}

	cb.Signature = MakeSignature(r, s)
	return nil
}

// Source returns the source hash associated with this block.
func (cb *ChangeBlock) Source() BlockHash {
	return BlockHashZero
}

// Timestamp returns the timestamp of when the block was created.
func (cb *ChangeBlock) Timestamp() int64 {
	return cb.Hashables.Timestamp
}

// Type returns the type of this block.
func (cb *ChangeBlock) Type() BlockType {
	return Change
}

// Verify verifies whether this block was signed by the given public key owner.
func (cb *ChangeBlock) Verify(pub *ecdsa.PublicKey) bool {
	hash, err := cb.Hash()

	if err != nil {
		return false
	}

	return crypto.Verify(hash[:], pub, cb.Signature.R, cb.Signature.S)
}

// Deserialize decodes byte data encoded by gob.
func (cb *ChangeBlock) Deserialize(r io.Reader) error {
	decoder := gob.NewDecoder(r)
	return decoder.Decode(cb)
}

// DeserializeJSON decodes JSON data.
func (cb *ChangeBlock) DeserializeJSON(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(cb)
}

// Serialize encodes to byte data using gob.
func (cb *ChangeBlock) Serialize(w io.Writer) error {
	encoder := gob.NewEncoder(w)
	return encoder.Encode(cb)
}

// SerializeJSON encodes to JSON data.
func (cb *ChangeBlock) SerializeJSON(w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(cb)
}

// String returns a JSON encoded string.
func (cb *ChangeBlock) String() (string, error) {
	return cb.ToJSON()
}

// ToJSON returns a JSON encoded string.
func (cb *ChangeBlock) ToJSON() (string, error) {
	bytes, err := json.Marshal(cb)

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// DelegateBlock represents a change in delegates.
type DelegateBlock struct {
	Hashables DelegateHashables `json:"hashables"`
	Signature Signature         `json:"signature"`
}

// NewDelegateBlock creates and initializes a DelegateBlock from the given arguments.
func NewDelegateBlock(amt Amount, delegate bool, prev BlockHash) *DelegateBlock {
	return &DelegateBlock{
		Hashables: MakeDelegateHashables(amt, delegate, prev),
	}
}

// Balance returns the balance associated with this block.
func (db *DelegateBlock) Balance() Amount {
	return db.Hashables.Balance
}

// Delegates returns the delegates associated with this block.
func (db *DelegateBlock) Delegates() []IBAN {
	return nil
}

// Hash returns the SHA256 hash of the serialized bytes of Hashables.
func (db *DelegateBlock) Hash() (BlockHash, error) {
	var buffer bytes.Buffer

	if err := db.Hashables.Serialize(&buffer); err != nil {
		return BlockHashZero, err
	}

	return sha256.Sum256(buffer.Bytes()), nil
}

// Previous returns the previous hash associated with this block.
func (db *DelegateBlock) Previous() BlockHash {
	return db.Hashables.Previous
}

// Root returns the previous hash associated with this block.
func (db *DelegateBlock) Root() BlockHash {
	return db.Hashables.Previous
}

// Sign signs the block with the given private key.
func (db *DelegateBlock) Sign(priv *ecdsa.PrivateKey) error {
	hash, err := db.Hash()

	if err != nil {
		return err
	}

	r, s, err := crypto.Sign(hash[:], priv)

	if err != nil {
		return err
	}

	db.Signature = MakeSignature(r, s)
	return nil
}

// Source returns the source hash associated with this block.
func (db *DelegateBlock) Source() BlockHash {
	return BlockHashZero
}

// Timestamp returns the timestamp of when the block was created.
func (db *DelegateBlock) Timestamp() int64 {
	return db.Hashables.Timestamp
}

// Type returns the type of this block.
func (db *DelegateBlock) Type() BlockType {
	return Delegate
}

// Verify verifies whether this block was signed by the given public key owner.
func (db *DelegateBlock) Verify(pub *ecdsa.PublicKey) bool {
	hash, err := db.Hash()

	if err != nil {
		return false
	}

	return crypto.Verify(hash[:], pub, db.Signature.R, db.Signature.S)
}

// Deserialize decodes byte data encoded by gob.
func (db *DelegateBlock) Deserialize(r io.Reader) error {
	decoder := gob.NewDecoder(r)
	return decoder.Decode(db)
}

// DeserializeJSON decodes JSON data.
func (db *DelegateBlock) DeserializeJSON(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(db)
}

// Serialize encodes to byte data using gob.
func (db *DelegateBlock) Serialize(w io.Writer) error {
	encoder := gob.NewEncoder(w)
	return encoder.Encode(db)
}

// SerializeJSON encodes to JSON data.
func (db *DelegateBlock) SerializeJSON(w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(db)
}

// String returns a JSON encoded string.
func (db *DelegateBlock) String() (string, error) {
	return db.ToJSON()
}

// ToJSON returns a JSON encoded string.
func (db *DelegateBlock) ToJSON() (string, error) {
	bytes, err := json.Marshal(db)

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// OpenBlock represents a openining of an account.
type OpenBlock struct {
	Hashables OpenHashables `json:"hashables"`
	Signature Signature     `json:"signature"`
}

// NewOpenBlock creates and initializes an OpenBlock from the given arguments.
func NewOpenBlock(amt Amount, iban IBAN) *OpenBlock {
	return &OpenBlock{
		Hashables: MakeOpenHashables(amt, iban),
	}
}

// Balance returns the balance associated with this block.
func (ob *OpenBlock) Balance() Amount {
	return ob.Hashables.Balance
}

// Delegates returns the delegates associated with this block.
func (ob *OpenBlock) Delegates() []IBAN {
	return nil
}

// Hash returns the SHA256 hash of the serialized bytes of Hashables.
func (ob *OpenBlock) Hash() (BlockHash, error) {
	var buffer bytes.Buffer
	err := ob.Hashables.Serialize(&buffer)

	if err != nil {
		return BlockHashZero, err
	}

	return sha256.Sum256(buffer.Bytes()), nil
}

// Previous returns the previous hash associated with this block.
func (ob *OpenBlock) Previous() BlockHash {
	return BlockHashZero
}

// Root returns the previous hash associated with this block.
func (ob *OpenBlock) Root() BlockHash {
	return BlockHashZero
}

// Sign signs the block with the given private key.
func (ob *OpenBlock) Sign(priv *ecdsa.PrivateKey) error {
	hash, err := ob.Hash()

	if err != nil {
		return err
	}

	r, s, err := crypto.Sign(hash[:], priv)

	if err != nil {
		return err
	}

	ob.Signature = MakeSignature(r, s)
	return nil
}

// Source returns the source hash associated with this block.
func (ob *OpenBlock) Source() BlockHash {
	return BlockHashZero
}

// Timestamp returns the timestamp of when the block was created.
func (ob *OpenBlock) Timestamp() int64 {
	return ob.Hashables.Timestamp
}

// Type returns the type of this block.
func (ob *OpenBlock) Type() BlockType {
	return Open
}

// Verify verifies whether this block was signed by the given public key owner.
func (ob *OpenBlock) Verify(pub *ecdsa.PublicKey) (bool, error) {
	hash, err := ob.Hash()

	if err != nil {
		return false, err
	}

	return crypto.Verify(hash[:], pub, ob.Signature.R, ob.Signature.S), nil
}

// Deserialize decodes byte data encoded by gob.
func (ob *OpenBlock) Deserialize(r io.Reader) error {
	decoder := gob.NewDecoder(r)
	return decoder.Decode(ob)
}

// DeserializeJSON decodes JSON data.
func (ob *OpenBlock) DeserializeJSON(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(ob)
}

// Serialize encodes to byte data using gob.
func (ob *OpenBlock) Serialize(w io.Writer) error {
	encoder := gob.NewEncoder(w)
	return encoder.Encode(ob)
}

// SerializeJSON encodes to JSON data.
func (ob *OpenBlock) SerializeJSON(w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(ob)
}

// String returns a json encoded string.
func (ob *OpenBlock) String() (string, error) {
	return ob.ToJSON()
}

// ToJSON returns a JSON encoded string.
func (ob *OpenBlock) ToJSON() (string, error) {
	bytes, err := json.Marshal(ob)

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// ReceiveBlock represents the receiving end of a send transaction.
type ReceiveBlock struct {
	Hashables ReceiveHashables `json:"hashables"`
	Signature Signature        `json:"signature"`
}

// NewReceiveBlock creates and initializes a ReceiveBlock from the given arguments.
func NewReceiveBlock(amt Amount, prev, src BlockHash) *ReceiveBlock {
	return &ReceiveBlock{
		Hashables: MakeReceiveHashables(amt, prev, src),
	}
}

// Balance returns the balance associated with this block.
func (rb *ReceiveBlock) Balance() Amount {
	return rb.Hashables.Balance
}

// Delegates returns the delegates associated with this block.
func (rb *ReceiveBlock) Delegates() []IBAN {
	return nil
}

// Hash returns the SHA256 hash of the serialized bytes of Hashables.
func (rb *ReceiveBlock) Hash() (BlockHash, error) {
	var buffer bytes.Buffer

	if err := rb.Hashables.Serialize(&buffer); err != nil {
		return BlockHashZero, err
	}

	return sha256.Sum256(buffer.Bytes()), nil
}

// Previous returns the previous hash associated with this block.
func (rb *ReceiveBlock) Previous() BlockHash {
	return rb.Hashables.Previous
}

// Root returns the previous hash associated with this block.
func (rb *ReceiveBlock) Root() BlockHash {
	return rb.Hashables.Previous
}

// Sign signs the block with the given private key.
func (rb *ReceiveBlock) Sign(priv *ecdsa.PrivateKey) error {
	hash, err := rb.Hash()

	if err != nil {
		return err
	}

	r, s, err := crypto.Sign(hash[:], priv)

	if err != nil {
		return err
	}

	rb.Signature = MakeSignature(r, s)
	return nil
}

// Source returns the source hash associated with this block.
func (rb *ReceiveBlock) Source() BlockHash {
	return rb.Hashables.Source
}

// Timestamp returns the timestamp of when the block was created.
func (rb *ReceiveBlock) Timestamp() int64 {
	return rb.Hashables.Timestamp
}

// Type returns the type of this block.
func (rb *ReceiveBlock) Type() BlockType {
	return Receive
}

// Verify verifies whether this block was signed by the given public key owner.
func (rb *ReceiveBlock) Verify(pub *ecdsa.PublicKey) (bool, error) {
	hash, err := rb.Hash()

	if err != nil {
		return false, err
	}

	return crypto.Verify(hash[:], pub, rb.Signature.R, rb.Signature.S), nil
}

// Deserialize decodes byte data encoded by gob.
func (rb *ReceiveBlock) Deserialize(r io.Reader) error {
	decoder := gob.NewDecoder(r)
	return decoder.Decode(rb)
}

// DeserializeJSON decodes JSON data.
func (rb *ReceiveBlock) DeserializeJSON(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(rb)
}

// Serialize encodes to byte data using gob.
func (rb *ReceiveBlock) Serialize(w io.Writer) error {
	encoder := gob.NewEncoder(w)
	return encoder.Encode(rb)
}

// SerializeJSON encodes to JSON data.
func (rb *ReceiveBlock) SerializeJSON(w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(rb)
}

// String returns a json encoded string.
func (rb *ReceiveBlock) String() (string, error) {
	return rb.ToJSON()
}

// ToJSON returns a JSON encoded string.
func (rb *ReceiveBlock) ToJSON() (string, error) {
	bytes, err := json.Marshal(rb)

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// SendBlock represents the sending of a transaction.
type SendBlock struct {
	Hashables SendHashables `json:"hashables"`
	Signature Signature     `json:"signature"`
}

// NewSendBlock creates and initializes a SendBlock from the given arguments.
func NewSendBlock(amt Amount, dst IBAN, prev BlockHash) *SendBlock {
	return &SendBlock{
		Hashables: MakeSendHashables(amt, dst, prev),
	}
}

// Balance returns the balance associated with this block.
func (sb *SendBlock) Balance() Amount {
	return sb.Hashables.Balance
}

// Delegates returns the delegates associated with this block.
func (sb *SendBlock) Delegates() []IBAN {
	return nil
}

// Hash returns the SHA256 hash of the serialized bytes of Hashables.
func (sb *SendBlock) Hash() (BlockHash, error) {
	var buffer bytes.Buffer

	if err := sb.Hashables.Serialize(&buffer); err != nil {
		return BlockHashZero, err
	}

	return sha256.Sum256(buffer.Bytes()), nil
}

// Previous returns the previous hash associated with this block.
func (sb *SendBlock) Previous() BlockHash {
	return sb.Hashables.Previous
}

// Root returns the previous hash associated with this block.
func (sb *SendBlock) Root() BlockHash {
	return sb.Hashables.Previous
}

// Sign signs the block with the given private key.
func (sb *SendBlock) Sign(priv *ecdsa.PrivateKey) error {
	hash, err := sb.Hash()

	if err != nil {
		return err
	}

	r, s, err := crypto.Sign(hash[:], priv)

	if err != nil {
		return err
	}

	sb.Signature = MakeSignature(r, s)
	return nil
}

// Source returns the source hash associated with this block.
func (sb *SendBlock) Source() BlockHash {
	return BlockHashZero
}

// Timestamp returns the timestamp of when the block was created.
func (sb *SendBlock) Timestamp() int64 {
	return sb.Hashables.Timestamp
}

// Type returns the type of this block.
func (sb *SendBlock) Type() BlockType {
	return Send
}

// Verify verifies whether this block was signed by the given public key owner.
func (sb *SendBlock) Verify(pub *ecdsa.PublicKey) (bool, error) {
	hash, err := sb.Hash()

	if err != nil {
		return false, err
	}

	return crypto.Verify(hash[:], pub, sb.Signature.R, sb.Signature.S), nil
}

// Deserialize decodes byte data encoded by gob.
func (sb *SendBlock) Deserialize(r io.Reader) error {
	decoder := gob.NewDecoder(r)
	return decoder.Decode(sb)
}

// DeserializeJSON decodes JSON data.
func (sb *SendBlock) DeserializeJSON(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(sb)
}

// Serialize encodes to byte data using gob.
func (sb *SendBlock) Serialize(w io.Writer) error {
	encoder := gob.NewEncoder(w)
	return encoder.Encode(sb)
}

// SerializeJSON encodes to JSON data.
func (sb *SendBlock) SerializeJSON(w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(sb)
}

// String returns a JSON encoded string.
func (sb *SendBlock) String() (string, error) {
	return sb.ToJSON()
}

// ToJSON returns a JSON encoded string.
func (sb *SendBlock) ToJSON() (string, error) {
	bytes, err := json.Marshal(sb)

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
