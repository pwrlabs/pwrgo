package encode

func txnBaseBytes(txType int, nonce int) ([]byte, error){
	typeByte := decToBytes(txType, 1)
	chainByte := decToBytes(0, 1)
	nonceBytes := decToBytes(nonce, 4)

	paddedNonce := make([]byte, 4)
	copy(paddedNonce[4-len(nonceBytes):], nonceBytes)

	var txnBytes []byte
	txnBytes = append(txnBytes, typeByte...)
	txnBytes = append(txnBytes, chainByte...)
	txnBytes = append(txnBytes, paddedNonce...)

	return txnBytes, nil
}

func decToBytes(value, length int) []byte {
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
