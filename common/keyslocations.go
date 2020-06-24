package common

import (
	"errors"
	"strings"
)

type KeysLocations string

const (
	ALL_FLORIDA_KEYS KeysLocations = "All Florida Keys"
	KEY_LARGO                      = "Key Largo"
	ISLAMORADA                     = "Islamorada"
	MARATHON                       = "Marathon"
	THE_LOWER_KEYS                 = "The Lower Keys"
	KEY_WEST                       = "Key West"
)

func (kl *KeysLocations) UnmarshalJSON(b []byte) error {
	keyLocation := KeysLocations(strings.Trim(string(b), `"`))
	switch keyLocation {
	case ALL_FLORIDA_KEYS, KEY_LARGO, ISLAMORADA, MARATHON, THE_LOWER_KEYS, KEY_WEST:
		*kl = keyLocation
		return nil
	}
	return errors.New("Invalid Location type")
}
