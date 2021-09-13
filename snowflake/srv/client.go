package srv

import "github.com/zhangliang-zl/reskit/snowflake/srv/protocols"

type Client struct {
}

func (*Client) NextID() {

	protocols.New("")
}

func (*Client) Quit() {

}
