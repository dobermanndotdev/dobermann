package monitor

var (
	// Americas
	RegionUsWest                = Region{value: "us_west"}
	RegionUsEast                = Region{value: "us_east"}
	RegionNorthAmericaNorthEast = Region{value: "north_america_north_east"}
	RegionSouthAmericaWest      = Region{value: "south_america_west"}
	RegionSouthAmericaEast      = Region{value: "south_america_east"}

	// Europe
	RegionEuWest  = Region{value: "eu_west"}
	RegionEuEast  = Region{value: "eu_east"}
	RegionEuSouth = Region{value: "eu_south"}
	RegionEuNorth = Region{value: "eu_north"}

	// Asia
	RegionAsiaNorth     = Region{value: "asia_north"}
	RegionAsiaNorthEast = Region{value: "asia_north_east"}
	RegionAsiaSouth     = Region{value: "asia_south"}
	RegionAsiaEast      = Region{value: "asia_east"}
	RegionAsiaSouthEast = Region{value: "asia_south_east"}

	// Pacific
	RegionAustraliaSouthEast = Region{value: "australia_south_east"}

	// Middle East
	RegionMiddleEastWest    = Region{value: "middle_east_west"}
	RegionMiddleEastCentral = Region{value: "middle_east_central"}

	// Africa
	RegionAfricaSouth = Region{value: "africa_south"}
)

type Region struct {
	value string
}

func NewDefaultRegions() []Region {
	return []Region{
		RegionUsWest,
		RegionUsEast,
		RegionNorthAmericaNorthEast,
		RegionSouthAmericaWest,
		RegionSouthAmericaEast,
		RegionEuWest,
		RegionEuEast,
		RegionEuSouth,
		RegionEuNorth,
		RegionAsiaNorth,
		RegionAsiaNorthEast,
		RegionAsiaSouth,
		RegionAsiaEast,
		RegionAsiaSouthEast,
		RegionAustraliaSouthEast,
		RegionMiddleEastWest,
		RegionMiddleEastCentral,
		RegionAfricaSouth,
	}
}
