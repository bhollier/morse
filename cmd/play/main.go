package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/bhollier/morse"
	"github.com/bhollier/morse/play"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"io"
	"os"
	"strings"
	"time"
)

var inputModeStr = flag.String("mode", "", "Required. The input mode of the program, either [t]ext or [m]orse")

var sampleRate = flag.Int("sampleRate", 44100, "The speaker sample rate")
var freq = flag.Int("freq", 800, "The tone frequency")
var wpm = flag.Uint("wpm", 20, "The words per minute (standard word PARIS)")
var farnsworthWPM = flag.Uint("fwpm", 15, "The farnsworth words per minute")

var interactive = flag.Bool("i", false, "Input interactively")

var printMorse = flag.Bool("p", false, "Print the morse code as it's played")

func main() {
	flag.Parse()

	sr := beep.SampleRate(*sampleRate)

	inputMode, err := ParseInputMode(*inputModeStr)
	if err != nil {
		fmt.Fprintln(flag.CommandLine.Output(), err.Error())
		os.Exit(2)
	}

	signalChannel := make(chan morse.Signal)
	morseWriter := morse.WriterFromChan(signalChannel, true)
	morseReader := morse.ReaderFromChan(signalChannel, false)
	if *printMorse {
		morseReader = morse.PrintWrapReader(morseReader)
	}

	streamer, err := play.MorseStreamer(sr, *freq, *wpm, *farnsworthWPM, morseReader)
	if err != nil {
		fmt.Fprintln(flag.CommandLine.Output(), err.Error())
		os.Exit(2)
	}

	err = speaker.Init(sr, sr.N(time.Second/10))
	if err != nil {
		panic(err)
	}
	defer speaker.Close()

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	// Output the initial input from the program arguments
	{
		input := strings.Join(flag.Args(), " ")
		_, err = morseWriter.Write(inputMode.ConvertInput(input))
		if err != nil {
			panic(err)
		}
	}

	// Output input from the user if in interactive mode
	if *interactive {
		fmt.Println("Press Ctrl+C to exit")
		inputScanner := bufio.NewScanner(os.Stdin)
		for inputScanner.Scan() {
			input := inputScanner.Text()
			_, _ = morseWriter.Write(inputMode.ConvertInput(input))
		}
		if inputScanner.Err() != io.EOF {
			panic(inputScanner.Err())
		}
	}

	close(signalChannel)
	<-done
}
