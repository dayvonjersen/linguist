// go port of github.com/github/linguist/lib/linguist/tokenizer.rb
package tokenizer

import (
	"bufio"
	"bytes"
	"log"
	"regexp"
)

var (
	ByteLimit int = 100000

	StartLineComments []string = []string{
		"\"", // Vim
		"%",  // Tex
	}

	SingleLineComments []string = []string{
		"//", // C
		"--", // Ada, Haskell, AppleScript
		"#",  // Perl, Bash, Ruby
	}

	MultiLineComments [][]string = [][]string{
		[]string{"/*", "*/"},         // C
		[]string{"<!--", "-->"},      // XML
		[]string{"{-", "-}"},         // Haskell
		[]string{"(*", "*)"},         // Coq
		[]string{"\"\"\"", "\"\"\""}, // Python
		[]string{"'''", "'''"},       // Python
		[]string{"#`(", ")"},         // Perl6
	}

	Strings []string = []string{
		`"`,
		`'`,
		"`",
	}

	Shebang          *regexp.Regexp   = regexp.MustCompile(`#!.*$`)
	Number           *regexp.Regexp   = regexp.MustCompile(`(0x[0-9a-f]([0-9a-f]|\.)*|\d(\d|\.)*)([uU][lL]{0,2}|([eE][-+]\d*)?[fFlL]*)`)
	StartLineComment []*regexp.Regexp = func() []*regexp.Regexp {
		ree := []*regexp.Regexp{}
		for _, normies := range append(StartLineComments, SingleLineComments...) {
			getout := regexp.MustCompile(`^\s*` + normies)
			ree = append(ree, getout)
		}
		return ree
	}()
)

func checkErr(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

func Tokenize(b []byte) (tokens []string) {
	if len(b) == 0 {
		return tokens
	}
	if len(b) >= ByteLimit {
		b = b[:ByteLimit]
	}

	ml_in := false
	ml_end := ""
	str_in := false
	str_end := ""

	buf := bytes.NewBuffer(b)
	scanlines := bufio.NewScanner(buf)
	scanlines.Split(bufio.ScanLines)
line:
	for scanlines.Scan() {
		ln := scanlines.Bytes()
		//fmt.Println(scanlines.Text)

		for _, re := range StartLineComment {
			if c := re.Find(ln); c != nil {
				goto line
			}
		}

		ln_buf := bytes.NewBuffer(ln)
		scanwords := bufio.NewScanner(ln_buf)
		scanwords.Split(bufio.ScanWords)
	word:
		for scanwords.Scan() {
			tk_b := scanwords.Bytes()
			tk := scanwords.Text()

			//fmt.Println(tk)
			if ml_in {
				if tk == ml_end {
					ml_in = false
					ml_end = ""
				}
				goto word
			}

			if str_in {
				if tk == str_end {
					str_in = false
					str_end = ""
				}
				goto word
			}

			for _, c := range SingleLineComments {
				if tk == c {
					goto line
				}
			}

			for _, m := range MultiLineComments {
				if tk == m[0] {
					ml_in = true
					ml_end = m[1]
					goto word
				}
			}

			for _, s := range Strings {
				if tk == s {
					str_in = true
					str_end = s
					goto word
				}
			}

			if n := Number.Find(tk_b); n != nil {
				goto word
			}

			tokens = append(tokens, tk)

		}
	}
	return tokens
}
