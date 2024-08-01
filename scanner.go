package main

func scanNumber(l *lexer) stateFunc {
	// Optional leading sign.
	l.accept("+-")
	digits := "0123456789_"
	// Is it hex?
	if l.accept("0") {
		// Note: Leading 0 does not mean octal in floats.
		if l.accept("xX") {
			digits = "0123456789abcdefABCDEF_"
		} else if l.accept("oO") {
			digits = "01234567_"
		} else if l.accept("bB") {
			digits = "01_"
		}
	}
	l.acceptRun(digits)
	if l.accept(".") {
		l.acceptRun(digits)
	}
	if len(digits) == 10+1 && l.accept("eE") {
		l.accept("+-")
		l.acceptRun("0123456789_")
	}
	if len(digits) == 16+6+1 && l.accept("pP") {
		l.accept("+-")
		l.acceptRun("0123456789_")
	}
	// Is it imaginary?
	l.accept("i")
	r := l.peek()
	// Next thing mustn't be alphanumeric.
	if isAlphaNumeric(r) {
		l.next()
		return l.errorf("bad syntax for number")
	}
	return nil
}

func scanString(l *lexer) {
	valid := func(r rune) bool {
		return !isSpace(r) && r != eof
	}
	l.acceptFunc(valid)
}
