package radix

import "testing"
import "fmt"
import log "github.com/cihub/seelog"

func SetupTestLogging() {
	var appConfig = `
  <seelog type="sync" minlevel='%s'>
  <outputs formatid="scanner">
    <filter levels="critical,error,warn,info">
      <console formatid="scanner" />
    </filter>
    <filter levels="debug">
      <console formatid="debug" />
    </filter>
  </outputs>
  <formats>
  <format id="scanner" format="test: [%%LEV] %%Msg%%n" />
  <format id="debug" format="test: [%%LEV] %%Func :: %%Msg%%n" />
  </formats>
  </seelog>
`

var config string
if testing.Verbose() {
  config =  fmt.Sprintf(appConfig, "debug")
} else {
  config =  fmt.Sprintf(appConfig, "info")
}

	logger, err := log.LoggerFromConfigAsBytes([]byte(config))

	if err != nil {
		fmt.Println(err)
		return
	}

	log.ReplaceLogger(logger)
}

func  Assert(t *testing.T, expected interface{}, actual interface{}, msg string) {
  if expected != actual {
    t.Errorf("%s. %v != %v", msg, expected, actual)
  }
}

func TestInsert(t *testing.T) {
  SetupTestLogging()

  var (
    val interface{}
    found bool
  )


  radix := NewRadixTree()

  Assert(t, radix.Len(), 0, "Length mismatch")

  radix.Insert([]byte("james"), 4)

  Assert(t, radix.Len(), 1, "Length mismatch")
  val, found = radix.Find([]byte("james"))
  Assert(t, found, true, "Couldn't find 'james'")
  Assert(t, val, 4, "Found incorrect value for 'james'")

  radix.Insert([]byte("janice"), 4)

  Assert(t, radix.Len(), 2, "Length mismatch")

  val, found = radix.Find([]byte("janice"))
  Assert(t, found, true, "Couldn't find 'janice'")
  Assert(t, val, 4, "Found incorrect value for 'janice'")

  val, found = radix.Find([]byte("james"))
  Assert(t, found, true, "Couldn't find 'james'")
  Assert(t, val, 4, "Found incorrect value for 'james'")

  radix.Insert([]byte("james"), "different")

  val, found = radix.Find([]byte("james"))
  Assert(t, found, true, "Couldn't find 'james'")
  Assert(t, val, "different", "Found incorrect value for 'james'")

}

