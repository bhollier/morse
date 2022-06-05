package morse

import (
	"fmt"
	"io"
)

type Printer interface {
	Print(a ...any)
	Println()
}

// PrintReaderWrapper is a simple wrapper around a Reader to
// print any signals that are read. Also prints newlines on EOF
type PrintReaderWrapper struct {
	Reader
	Printer
}

// PrintWrapReader wraps a Reader with a PrintReaderWrapper.
// All signals are printed with fmt.Print
func PrintWrapReader(r Reader) PrintReaderWrapper {
	return PrintWrapReaderWithPrinter(r, nil)
}

// PrintWrapReaderWithPrinter wraps a Reader with a custom
// Printer. All signals are printed with Printer.Print
func PrintWrapReaderWithPrinter(r Reader, p Printer) PrintReaderWrapper {
	return PrintReaderWrapper{
		Reader:  r,
		Printer: p,
	}
}

func (r PrintReaderWrapper) Read(p []Signal) (n int, err error) {
	n, err = r.Reader.Read(p)
	for _, s := range p[:n] {
		if r.Printer == nil {
			fmt.Print(s)
		} else {
			r.Print(s)
		}
	}
	if err == io.EOF {
		if r.Printer == nil {
			fmt.Println()
		} else {
			r.Println()
		}
	}
	return
}

// PrintWriterWrapper is a simple wrapper around a Writer to
// print any signals that are written
type PrintWriterWrapper struct {
	Writer
	Printer
}

// PrintWrapWriter wraps a Writer with a PrintWriterWrapper.
// All signals are printed with fmt.Print
func PrintWrapWriter(w Writer) PrintWriterWrapper {
	return PrintWrapWriterWithPrinter(w, nil)
}

// PrintWrapWriterWithPrinter wraps a Writer with a custom
// Printer. All signals are printed with Printer.Print
func PrintWrapWriterWithPrinter(w Writer, p Printer) PrintWriterWrapper {
	return PrintWriterWrapper{
		Writer:  w,
		Printer: p,
	}
}

func (w PrintWriterWrapper) Write(p []Signal) (n int, err error) {
	n, err = w.Writer.Write(p)
	for _, s := range p[:n] {
		if w.Printer == nil {
			fmt.Print(s)
		} else {
			w.Print(s)
		}
	}
	return
}
