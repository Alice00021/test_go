package entity

type ReagentType string

const (
	ReagentTypeVER     ReagentType = "VER"
	ReagentTypeCAL     ReagentType = "CAL"
	ReagentTypeALCOHOL ReagentType = "ALCOHOL"
	ReagentTypeBLEACH  ReagentType = "BLEACH"
	ReagentTypeWATER   ReagentType = "WATER"
	ReagentTypeALKALI  ReagentType = "ALKALI"
	ReagentTypeDF      ReagentType = "DF"
)

type Address string

const (
	AddressRA Address = "RA"
	AddressRB Address = "RB"
	AddressRC Address = "RC"
	AddressRD Address = "RD"
	AddressSA Address = "SA"
	AddressSB Address = "SB"
	AddressSC Address = "SC"
	AddressSD Address = "SD"
	AddressSE Address = "SE"
	AddressSF Address = "SF"
	AddressSG Address = "SG"
	AddressSH Address = "SH"
)

type Command struct {
	Entity
	Name             string
	SystemName       string
	Reagent          ReagentType
	AverageTime      int64
	VolumeWaste      int64
	VolumeDriveFluid int64
	VolumeContainer  int64
	DefaultAddress   Address
}
