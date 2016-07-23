package pubsub

var cs channelStore

type Consumer interface {
	Consume(data interface{})
}

type Publisher struct {
	c *Channel
}

func NewPublisher(name string) Publisher {
	return Publisher{c: NewChannel(name)}
}

func (p Publisher)Publish(data interface{})  {
	for _, s := range (p.c.s){
		go s.Consume(data)
	}
}

type Channel struct {
	name string
	s []Consumer
}

func NewChannel(name string) *Channel {
	if cs.channels[name] != nil {
		return cs.channels[name]
	}
	c := &Channel{
		name:name,
	}
	cs.channels[name] = c
	return c
}

func (ch Channel)Subscribe(c Consumer) {
	ch.s = append(ch.s, c)
}

type channelStore struct {
	channels map[string]*Channel
}

func init() {
	cs = channelStore{channels: make(map[string]*Channel)}
}