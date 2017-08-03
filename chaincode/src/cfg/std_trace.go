package cfg

import (
	"fmt"
	"os"
)

func StdError(v ...interface{}) {
	fmt.Fprint(os.Stderr, v)
}

func StdErrorf(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format, v)
}

func StdOut(v ...interface{}) {
	fmt.Fprint(os.Stdout, v)
}

func StdOutf(format string, v ...interface{}) {
	fmt.Fprintf(os.Stdout, format, v)
}
