package linguist

// IsBinary checks contents for known character escape codes which
// frequently show up in binary files but rarely (if ever) in text.
//
// Use this check before using DetectFromContents to reduce likelihood
// of passing binary data into it.
//
// NOTE(tso): preliminary testing on this method of checking for binary
// contents were promising, having fed a document consisting of all
// utf-8 codepoints from 0000 to FFFF with satisfactory results. Thanks
// to robpike.io/cmd/unicode:
// ```
// unicode -c $(seq 0 65535 | xargs printf "%04x ") | tr -d '\n' > unicode_test
// ```
//
// However, the intentional presence of character escape codes to throw
// this function off is entirely possible, as is, potentially, a binary
// file consisting entirely of the 4 exceptions to the rule for the first
// 512 bytes. It is also possible that more character escape codes need
// to be added.
//
// Further analysis and real world testing of this is required.
func IsBinary(contents []byte) (probably bool) {
	for _, b := range contents[:512] {
		if b < 32 {
			switch b {
			case 0:
				fallthrough
			case 9:
				fallthrough
			case 10:
				fallthrough
			case 13:
				continue
			default:
				return true
			}
		}
	}
	return false
}
