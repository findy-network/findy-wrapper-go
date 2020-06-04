package ctx

import (
	"fmt"
	"sync"

	"github.com/findy-network/findy-wrapper-go/dto"
	"github.com/golang/glog"
)

// Channel is channel type for findy API. Instead of callbacks findy API returns
// channels for async functions.
type Channel chan dto.Result

type chanInfo struct {
	Channel
	name string
}

// ChannelBufferSize is the channel size of the return channel of the API
// functions.
const ChannelBufferSize = 1

// We need buffered channels, indy callbacks CANNOT block or we cannot start
// many task simultaneously and concurrent, buffer size one is OK.

type cmdHandles struct {
	counter  uint32
	channels map[uint32]chanInfo
	lock     sync.Mutex
}

// With NamedPush() we can add string to cmd handle allocations. These names
// will be printed on pops.
func (c *cmdHandles) NamedPush(name string) (uint32, Channel) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.counter++
	ch := make(Channel, ChannelBufferSize)
	c.channels[c.counter] = chanInfo{
		Channel: ch,
		name:    name,
	}
	return c.counter, ch
}

func (c *cmdHandles) Push() (uint32, Channel) {
	return c.NamedPush("")
}

func (c *cmdHandles) Pop(handle uint32, s fmt.Stringer) Channel {
	c.lock.Lock()
	defer c.lock.Unlock()

	ch, found := c.channels[handle]
	if !found {
		panic("push/pop mismatch")
	}
	delete(c.channels, handle)

	if glog.V(10) {
		glog.Infof("%s(%d) -> %s\n", ch.name, handle, s)
	}
	return ch.Channel
}

// CmdContext is singleton handler for command context.
var CmdContext = &cmdHandles{
	channels: make(map[uint32]chanInfo),
}
