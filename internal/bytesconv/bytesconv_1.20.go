//go:build go1.20

package bytesconv

// StringToBytes converts string to byte slice without a memory allocation.
// For more details, see https://github.com/golang/go/issues/53003#issuecomment-1140276077.
func StringToBytes(s string) []byte {
	return []byte(s)
	//return unsafe.Slice(unsafe.StringData(s), len(s))
}

// BytesToString converts byte slice to string without a memory allocation.
// For more details, see https://github.com/golang/go/issues/53003#issuecomment-1140276077.
//func BytesToString(b []byte) string {
//	return unsafe.String(unsafe.SliceData(b), len(b))
//}
