package morse

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSignal_Duration(t *testing.T) {
	a := assert.New(t)

	standardWordDuration := append(FromText(StandardWord), WordSpace).Duration(20, 0)
	a.Equal(time.Minute, (standardWordDuration * 20).Round(time.Second))

	standardWordFarnsworthDuration := append(FromText(StandardWord), WordSpace).Duration(20, 15)
	a.Equal(time.Minute, (standardWordFarnsworthDuration * 15).Round(time.Second))
}
