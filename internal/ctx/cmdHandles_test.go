package ctx

import (
	"reflect"
	"testing"
	"time"

	"github.com/findy-network/findy-wrapper-go/dto"
)

func fillChannel(cmd uint32, ch Channel) {
	r := dto.Result{}
	r.SetHandle(int(cmd))
	time.Sleep(time.Duration(cmd) * 10 * time.Millisecond)
	//	fmt.Printf("%d -> ch\n", cmd)
	ch <- r
}

func readChannel(ch Channel) uint32 {
	r := <-ch
	//	fmt.Printf("%d <- ch\n", r.Handle())
	return uint32(r.Handle())
}

func Test_cmdHandles_Push(t *testing.T) {
	// allocate channels and _cmd_handle_values_
	cmd1, ch1 := CmdContext.Push()
	cmd2, ch2 := CmdContext.NamedPush("named push test")
	cmd3, ch3 := CmdContext.Push()

	// write _cmd_values_ to these same channels
	go fillChannel(cmd1, ch1)
	go fillChannel(cmd2, ch2)
	go fillChannel(cmd3, ch3)

	// read cmds from channels
	r1 := readChannel(ch1)
	r2 := readChannel(ch2)
	r3 := readChannel(ch3)

	// check that we did read same values from same channel
	// order of values is taken care in fillChannel()
	tests := []struct {
		name string
		got  uint32
		want uint32
	}{
		{"1st ch", r1, cmd1},
		{"2nd ch", r2, cmd2},
		{"3rd ch", r3, cmd3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("add() = %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestGetCmdHandles(t *testing.T) {
	instance := CmdContext
	tests := []struct {
		name string
		want *cmdHandles
	}{
		{"startup", instance},
		{"first", instance},
		{"second", instance},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CmdContext; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CmdContext = %v, want %v", got, tt.want)
			}
		})
	}
}
