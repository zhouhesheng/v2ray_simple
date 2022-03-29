// Package direct provies direct proxy support for proxy.Client
package direct

import (
	"io"
	"net"
	"net/url"

	"github.com/hahahrfool/v2ray_simple/netLayer"
	"github.com/hahahrfool/v2ray_simple/proxy"
)

const name = "direct"

func init() {
	proxy.RegisterClient(name, &ClientCreator{})
}

//实现了 proxy.Client, netLayer.UDP_Extractor, netLayer.UDP_Putter
type Direct struct {
	proxy.ProxyCommonStruct
	*netLayer.UDP_Pipe

	targetAddr *netLayer.Addr
	addrStr    string
}

type ClientCreator struct{}

func NewClient() (proxy.Client, error) {
	d := &Direct{
		UDP_Pipe: netLayer.NewUDP_Pipe(),
	}
	//单单的pipe是无法做到转发的，它就像一个缓存一样;
	// 一方是未知的, 将 Direct 视为 UDP_Putter, 放入请求,
	// 然后我们这边就要通过一个 goroutine 来不断提取请求然后转发到direct.

	go netLayer.RelayUDP_to_Direct(d.UDP_Pipe)
	return d, nil
}

func (_ ClientCreator) NewClientFromURL(*url.URL) (proxy.Client, error) {
	return NewClient()
}

func (_ ClientCreator) NewClient(*proxy.DialConf) (proxy.Client, error) {
	return NewClient()
}

func (d *Direct) Name() string { return name }

//若 underlay 为nil，则我们会自动对target进行拨号。
func (d *Direct) Handshake(underlay net.Conn, target *netLayer.Addr) (io.ReadWriter, error) {

	if underlay == nil {
		d.targetAddr = target
		d.SetAddrStr(d.targetAddr.String())
		return target.Dial()
	}

	return underlay, nil

}
