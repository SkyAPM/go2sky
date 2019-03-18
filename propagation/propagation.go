package propagation

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

func NewSW3CarrierItem() CarrierItem {
	item := new(sw3CarrierItem)

	return item
}

// Context is a data carrier of tracing context,
// it holds a snapshot for across process propagation.
type ContextCarrier struct {
	items []CarrierItem
}

func (c *ContextCarrier) GetAllItems() []CarrierItem {
	return c.items
}

func NewContextCarrier() *ContextCarrier {
	carrier := ContextCarrier{items: []CarrierItem{
		NewSW3CarrierItem(),
	}}
	return &carrier
}

type Extractor func() (ContextCarrier, error)

type Injector func(carrier *ContextCarrier) error
