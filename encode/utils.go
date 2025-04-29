package encode

func txnBaseBytes(txIdentifier int, nonce int, address []byte, feePerByte int) ([]byte, error) {
	identifierByte := DecToBytes(txIdentifier, 4)
	chainIdByte := DecToBytes(0, 1)
	nonceBytes := DecToBytes(nonce, 4)
	feePerByteBytes := decToBytes64(int64(feePerByte), 8)

	paddedNonce := make([]byte, 4)
	copy(paddedNonce[4-len(nonceBytes):], nonceBytes)

	paddedAddress := make([]byte, 20)
	copy(paddedAddress[20-len(address):], address)

	var txnBytes []byte
	txnBytes = append(txnBytes, identifierByte...)
	txnBytes = append(txnBytes, chainIdByte...)
	txnBytes = append(txnBytes, paddedNonce...)
	txnBytes = append(txnBytes, feePerByteBytes...)
	txnBytes = append(txnBytes, paddedAddress...)

	return txnBytes, nil
}

func DecToBytes(value, length int) []byte {
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[length-1-i] = byte(value >> (8 * i))
	}
	return result
}

func decToBytes64(value int64, length int) []byte {
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[length-1-i] = byte(value >> (8 * i))
	}
	return result
}
