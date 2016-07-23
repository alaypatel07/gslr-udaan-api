package pubsub

//TODO :
//	1) improve the design of data structures using interfaces
//	2) Add comments


type Client struct {
	C map[string] chan interface{}
}

func NewClient() *Client {
	return &Client{make(map[string]chan interface{})}
}

func (c *Client)Subscribe(channel *Channel) {
	(*channel).Clients = append((*channel).Clients, c)
	c.C[(*channel).Name] = make(chan interface{})
}

type Channel struct {
	Name string
	Clients []*Client
}

func NewChannel(name string) *Channel {
	return &Channel{
		Name: name,
	}
}

func (c *Channel)Subscribe(cli *Client) {
	cli.C[c.Name] = make(chan interface{})
	c.Clients = append(c.Clients, cli)
}

func (c *Channel)Publish(v interface{}) {
	for _, client := range c.Clients {
		//fmt.Println("Published", v, "to", *client)
		(*client).C[c.Name] <- v
	}
}