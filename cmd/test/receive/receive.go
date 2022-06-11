package receive

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/bhollier/morse"
	"github.com/bhollier/morse/play"
	"github.com/bhollier/morse/prosign"
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

var attention = SubCmd.Bool("attention", false,
	"Whether to send the 'attention' prosign (－・－・－) before each test code")

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
		if *attention {
			// Prepend the attention prosign to indicate the start of the message
			randCode = morse.JoinWords(prosign.Attention, randCode)
		}

		// Keep asking for input until the user gives something
		ok := true
		var input string
		for ok && input == "" {
			cancel := make(chan bool, 1)
			finished := make(chan bool)
			// Write the morse code in another go routine
			// so the user can write their input as it plays
			go func() {
				i := 0
				for i < len(randCode) {
					// Can only cancel on a space
					if !randCode[i].Audible() {
						select {
						// If the user wants to cancel
						case <-cancel:
							// Send a rune space
							signalChannel <- morse.RuneSpace
							// Exit the loop early
							break

							// Otherwise just try to write to the signal channel
						case signalChannel <- randCode[i]:
							i++
						}

						// Otherwise just write like normal
					} else {
						signalChannel <- randCode[i]
						i++
					}
				}
				// Let the main thread know we're finished,
				// so it can (possibly) replay the code
				finished <- true
			}()

			ok = inputScanner.Scan()
			// The user has given their input, cancel playing the morse code
			// (if it hasn't finished already)
			cancel <- true
			// Wait for the go routine to finish writing at an appropriate spot
			<-finished
			if ok {
				input = strings.ToLower(inputScanner.Text())
			}
		}
		if !ok || inputScanner.Err() != nil {
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
