package monitor

var (
	RegionAfrica       = Region{"africa"}
	RegionEurope       = Region{"europe"}
	RegionNorthAmerica = Region{"north-america"}
	RegionSouthAmerica = Region{"south-america"}
	RegionMiddleEast   = Region{"middle-east"}
	RegionAsia         = Region{"asia"}
	RegionAustralia    = Region{"australia"}
)

type Region struct {
	value string
}

func NewRegion(value string) (Region, error) {
	switch value {
	case RegionAfrica.value:
		return RegionAfrica, nil
	case RegionAsia.value:
		return RegionAsia, nil
	case RegionEurope.value:
		return RegionEurope, nil
	case RegionAustralia.value:
		return RegionAustralia, nil
	case RegionMiddleEast.value:
		return RegionMiddleEast, nil
	case RegionNorthAmerica.value:
		return RegionNorthAmerica, nil
	case RegionSouthAmerica.value:
		return RegionSouthAmerica, nil
	}

	return Region{}, nil
}

func (r Region) String() string {
	return r.value
}
