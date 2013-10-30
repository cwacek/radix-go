package radix

import "testing"
import "math/rand"
import "strconv"
import "fmt"
import "bytes"
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

func  Assert(t *testing.T, actual interface{}, expected interface{}, msg string) {
  if expected != actual {
    t.Errorf("%s. %v != %v", msg, expected, actual)
  }
}

type test_entry struct{
  Key []byte
  Val interface{}
}

func (t test_entry) RadixKey() []byte {
  return t.Key
}

func TestInsert(t *testing.T) {
  SetupTestLogging()

  var (
    val interface{}
    found bool
  )


  radix := NewTrie()

  Assert(t, radix.Len(), 0, "Length mismatch")

  radix.Insert(test_entry{[]byte("james"), 4})

  Assert(t, radix.Len(), 1, "Length mismatch")
  val, found = radix.Find([]byte("james"))
  Assert(t, found, true, "Couldn't find 'james'")
  Assert(t, val.(test_entry).Val, 4, "Found incorrect value for 'james'")

  radix.Insert(test_entry{ []byte("janice"), 4 })

  Assert(t, radix.Len(), 2, "Length mismatch")

  val, found = radix.Find([]byte("janice"))
  Assert(t, found, true, "Couldn't find 'janice'")
  Assert(t, val.(test_entry).Val, 4, "Found incorrect value for 'janice'")

  val, found = radix.Find([]byte("james"))
  Assert(t, found, true, "Couldn't find 'james'")
  Assert(t, val.(test_entry).Val, 4, "Found incorrect value for 'james'")

  radix.Insert(test_entry{ []byte("james"), "different" })

  val, found = radix.Find([]byte("james"))
  Assert(t, found, true, "Couldn't find 'james'")
  Assert(t, val.(test_entry).Val, "different", "Found incorrect value for 'james'")

  val, found = radix.Find([]byte("a"))
  Assert(t, found, false, "Couldn't find 'a'")
  radix.Insert(test_entry{ []byte("a"), "different" })

  radix.Insert(test_entry{ []byte("freddie"), "kruger" })
  val, found = radix.Find([]byte("fredie"))
  Assert(t, found, false, "Found 'fredie' but shouldn't have")

  Assert(t, radix.Len(), 4, "Length mismatch")

  radix.Insert(test_entry{ []byte("jimbo"), "kruger" })
  val, found = radix.Find([]byte("jimbo"))
  Assert(t, found, true, "Couldn't find 'jimbo'")
  Assert(t, val.(test_entry).Val, "kruger", "Found incorrect value for 'jimbo'")

  expected := []test_entry{
    test_entry{[]byte("a"), "different"},
    test_entry{[]byte("freddie"), "kruger"},
    test_entry{[]byte("james"), "different"},
    test_entry{[]byte("janice"), 4},
    test_entry{[]byte("jimbo"), "kruger"},
  }

  entries := radix.Walk()
  for i, entry := range entries {
    if bytes.Compare(entry.RadixKey(), expected[i].Key) != 0 {
      t.Errorf("Error when walking: %s != %s at position %d", entry, expected[i], i)
    }
  }

}

func BenchmarkRapidIteration(b *testing.B) {

  var T test_entry

  radix := NewTrie()

  for i := 0; i < b.N; i++ {

    T = test_entry{[]byte(strconv.Itoa(rand.Intn(200000000000))), "What"}

    radix.Insert(T)
  }

}

