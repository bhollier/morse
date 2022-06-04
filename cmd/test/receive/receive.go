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

type subCmd struct {
	*flag.FlagSet
}

var SubCmd = subCmd{
	FlagSet: flag.NewFlagSet("receive", flag.ExitOnError),
}

var sampleRate = SubCmd.Int("sampleRate", 44100, "The speaker sample rate")
var freq = SubCmd.Int("freq", 800, "The tone frequency")
var wpm = SubCmd.Uint("wpm", 20, "The words per minute (standard word PARIS)")
var farnsworthWPM = SubCmd.Uint("fwpm", 15, "The farnsworth words per minute. "+
	"Only applicable if group is equal to '[w]ords' or '[s]entences")

var groupingStr = SubCmd.String("group", "", "Required. How many morse code signals to send for each test, "+
	"either individual [c]haraters, [w]ords or [s]entences")
var characters = SubCmd.String("characters", "abcdefghijklmnopqrstuvwxyz", "The characters to use as input in tests")
var maxWordLength = SubCmd.Uint("maxWord", 4, "The maximum length of a word. "+
	"Only applicable if group is equal to '[w]ords' or '[s]entences'")
var minSentenceLength = SubCmd.Uint("minSentence", 2, "The minimum words in a sentence. "+
	"Only applicable if group is equal to '[s]entences'")
var maxSentenceLength = SubCmd.Uint("maxSentence", 6, "The maximum words in a sentence. "+
	"Only applicable if group is equal to '[s]entences'")

func (s subCmd) Run(args []string) {
	randSrc := rand.NewSource(time.Now().UnixNano())
	s.Name()
	r := rand.New(randSrc)

	_ = s.Parse(args)

	sr := beep.SampleRate(*sampleRate)

	grouping, err := ParseGrouping(*groupingStr)
	if err != nil {
		fmt.Fprintln(s.Output(), err.Error())
		os.Exit(2)
	}

	signalChannel := make(chan morse.Signal)
	morseReader := morse.ReaderFromChan(signalChannel, false)

	streamer, err := play.MorseStreamer(sr, *freq, *wpm, *farnsworthWPM, morseReader)
	if err != nil {
		fmt.Fprintln(s.Output(), err.Error())
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
