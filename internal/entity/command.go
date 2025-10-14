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

func ValidateUniqueReagentAddress(commands []Command) error {
	addressMap := make(map[Address]ReagentType)

	for _, cmd := range commands {
		if cmd.DefaultAddress == "" {
			continue
		}

		if existingReagentType, exists := addressMap[cmd.DefaultAddress]; exists {
			if existingReagentType != cmd.Reagent {
				return ErrCommandDuplicateAddress
			}
		}
		addressMap[cmd.DefaultAddress] = cmd.Reagent
	}

	return nil
}

func ValidateMaxVolumeAddress(commands []Command) error {
	addressVolumeMap := make(map[Address]int64)

	for _, cmd := range commands {
		if cmd.DefaultAddress == "" {
			continue
		}
		addressVolumeMap[cmd.DefaultAddress] += cmd.VolumeContainer
	}

	for address, totalVolume := range addressVolumeMap {
		if address == "SA" || address == "SB" || address == "SC" ||
			address == "SE" || address == "SF" || address == "SG" || address == "SH" {
			if totalVolume > 200 {
				return ErrCommandVolumeExceeded
			}
		} else if address == "RA" || address == "RB" || address == "RC" || address == "RD" {
			if totalVolume > 5000 {
				return ErrCommandVolumeExceeded
			}
		}
	}

	return nil
}
