package logger

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type testLogger struct {
	strings.Builder
}

func (tl *testLogger) AnyLevelln(args []interface{}) {
	defer tl.Builder.WriteByte('\n')

	for i := range args {
		strArg := args[i].(string)
		_, _ = tl.Builder.WriteString(strArg)

		if i != len(args)-1 {
			tl.Builder.WriteString(" ")
		}
	}
}

func (tl *testLogger) String() string {
	return tl.Builder.String()
}

func (tl *testLogger) Debugln(args ...interface{}) {
	var argsForAnyLevel []interface{} = args
	tl.AnyLevelln(argsForAnyLevel)
}

func (tl *testLogger) Infoln(args ...interface{}) {
	var argsForAnyLevel []interface{} = args
	tl.AnyLevelln(argsForAnyLevel)
}

func (tl *testLogger) Warnln(args ...interface{}) {
	var argsForAnyLevel []interface{} = args
	tl.AnyLevelln(argsForAnyLevel)
}

func (tl *testLogger) Errorln(args ...interface{}) {
	var argsForAnyLevel []interface{} = args
	tl.AnyLevelln(argsForAnyLevel)
}

func (tl *testLogger) DPanicln(args ...interface{}) {
	var argsForAnyLevel []interface{} = args
	tl.AnyLevelln(argsForAnyLevel)
}

func (tl *testLogger) Panicln(args ...interface{}) {
	var argsForAnyLevel []interface{} = args
	tl.AnyLevelln(argsForAnyLevel)
}

func (tl *testLogger) Fatalln(args ...interface{}) {
	var argsForAnyLevel []interface{} = args
	tl.AnyLevelln(argsForAnyLevel)
}

func TestLogger(t *testing.T) {
	Logger().Debugln("")

	tl := &testLogger{}
	log.logger = tl

	Logger().Debugln("a")
	Logger().Infoln("b")
	Logger().Warnln("c")
	Logger().Errorln("d")
	Logger().DPanicln("e")
	Logger().Panicln("f")
	Logger().Fatalln("g")

	require.Equal(t, "a\nb\nc\nd\ne\nf\ng\n", tl.String())
}
