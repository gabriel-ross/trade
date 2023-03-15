package trade

type Currency int

const (
	Dollars Currency = iota
	Apples
)

func (c Currency) String() string {
	switch c {
	case Dollars:
		return "dollars"
	case Apples:
		return "apples"
	}

	return "unknown"
}
