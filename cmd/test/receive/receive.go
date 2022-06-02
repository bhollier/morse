package receive

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/bhollier/morse"
	"github.com/bhollier/morse/play"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

var flagSet = flag.NewFlagSet("receive", flag.ExitOnError)

var sampleRate = flagSet.Int("sampleRate", 44100, "The speaker sample rate")
var freq = flagSet.Int("freq", 800, "The tone frequency")
var wpm = flagSet.Uint("wpm", 20, "The words per minute (standard word PARIS)")
var farnsworthWPM = flagSet.Uint("fwpm", 15, "The farnsworth words per minute. "+
	"Only applicable if group is equal to '[w]ords' or '[s]entences")

var groupingStr = flagSet.String("group", "", "Required. How many morse code signals to send for each test, "+
	"either individual [c]haraters, [w]ords or [s]entences")
var characters = flagSet.String("characters", "abcdefghijklmnopqrstuvwxyz", "The characters to use as input in tests")
var maxWordLength = flagSet.Uint("maxWord", 4, "The maximum length of a word. "+
	"Only applicable if group is equal to '[w]ords' or '[s]entences'")
var minSentenceLength = flagSet.Uint("minSentence", 2, "The minimum words in a sentence. "+
	"Only applicable if group is equal to '[s]entences'")
var maxSentenceLength = flagSet.Uint("maxSentence", 6, "The maximum words in a sentence. "+
	"Only applicable if group is equal to '[s]entences'")

func Main(args []string) {
	randSrc := rand.NewSource(time.Now().UnixNano())
	r := rand.New(randSrc)

	_ = flagSet.Parse(args)

	sr := beep.SampleRate(*sampleRate)

	grouping, err := ParseGrouping(*groupingStr)
	if err != nil {
		fmt.Fprintln(flag.CommandLine.Output(), err.Error())
		os.Exit(2)
	}

	signalChannel := make(chan morse.Signal)
	morseReader := morse.ReaderFromChan(signalChannel, false)

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

	fmt.Println("Press Ctrl+C to exit. Enter to repeat")
	inputScanner := bufio.NewScanner(os.Stdin)

	for true {
		// Generate a random string
		randString := strings.ToLower(grouping.GenerateString(r))
		randCode := morse.FromText(randString)

		// Keep asking for input until the user gives something
		var input string
		for input == "" {
			go func() {
				for _, s := range randCode {
					// todo cancel if input is received or if repeat is requested
					signalChannel <- s
				}
			}()

			if !inputScanner.Scan() {
				break
			}
			input = strings.ToLower(inputScanner.Text())
		}
		if input == "" || inputScanner.Err() != nil {
			break
		}

		if randString == input {
			fmt.Println("Correct!")
		} else {
			fmt.Println("Incorrect! Was actually " + randString)
		}
	}
	if inputScanner.Err() != io.EOF {
		panic(inputScanner.Err())
	}

	close(signalChannel)
	<-done
}
