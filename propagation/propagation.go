package propagation

//CarrierItem is sub entity of propagation specification
type CarrierItem interface {
	HeadKey() string
	HeadValue() string
	SetValue(t string)
	IsValid() bool
}

type sw3CarrierItem struct {
}

func (s *sw3CarrierItem) HeadKey() string {
	return "sw3"
}

func (s *sw3CarrierItem) HeadValue() string {
	return ""
}

func (s *sw3CarrierItem) SetValue(t string) {
}

func (s *sw3CarrierItem) IsValid() bool {
	return true
}

//NewSW3CarrierItem create a new SkyWalking v3 propagation protocol carrier object
func NewSW3CarrierItem() CarrierItem {
	item := new(sw3CarrierItem)

	return item
}

// ContextCarrier is a data carrier of tracing context,
// it holds a snapshot for across process propagation.
type ContextCarrier struct {
	items []CarrierItem
}

// GetAllItems gets all data from ContextCarrier
func (c *ContextCarrier) GetAllItems() []CarrierItem {
	return c.items
}

// NewContextCarrier create a new ContextCarrier object
func NewContextCarrier() *ContextCarrier {
	carrier := ContextCarrier{items: []CarrierItem{
		NewSW3CarrierItem(),
	}}
	return &carrier
}

// Extractor is a tool specification which define how to
// extract trace parent context from propagation context
type Extractor func() (ContextCarrier, error)

// Injector is a tool specification which define how to
// inject trace context into propagation context
type Injector func(carrier *ContextCarrier) error
