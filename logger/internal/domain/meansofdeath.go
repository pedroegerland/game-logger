package domain

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
)

type MeansOfDeath int

var (
	_ sql.Scanner   = (*MeansOfDeath)(nil)
	_ driver.Valuer = Undefined

	mapToString = map[MeansOfDeath]string{
		Undefined:        "",
		ModUnknown:       "MOD_UNKNOWN",
		ModShotgun:       "MOD_SHOTGUN",
		ModGauntlet:      "MOD_GAUNTLET",
		ModMachineGun:    "MOD_MACHINEGUN",
		ModGrenade:       "MOD_GRENADE",
		ModGrenadeSplash: "MOD_GRENADE_SPLASH",
		ModRocket:        "MOD_ROCKET",
		ModRocketSplash:  "MOD_ROCKET_SPLASH",
		ModPlasma:        "MOD_PLASMA",
		ModPlasmaSplash:  "MOD_PLASMA_SPLASH",
		ModRailGun:       "MOD_RAILGUN",
		ModLightning:     "MOD_LIGHTNING",
		ModBfg:           "MOD_BFG",
		ModBfgSplash:     "MOD_BFG_SPLASH",
		ModWater:         "MOD_WATER",
		ModSlime:         "MOD_SLIME",
		ModLava:          "MOD_LAVA",
		ModCrush:         "MOD_CRUSH",
		ModTelefrag:      "MOD_TELEFRAG",
		ModFalling:       "MOD_FALLING",
		ModSuicide:       "MOD_SUICIDE",
		ModTargetLaser:   "MOD_TARGET_LASER",
		ModTriggerHurt:   "MOD_TRIGGER_HURT",
	}

	mapStringTo = map[string]MeansOfDeath{
		"":                   Undefined,
		"MOD_UNKNOWN":        ModUnknown,
		"MOD_SHOTGUN":        ModShotgun,
		"MOD_GAUNTLET":       ModGauntlet,
		"MOD_MACHINEGUN":     ModMachineGun,
		"MOD_GRENADE":        ModGrenade,
		"MOD_GRENADE_SPLASH": ModGrenadeSplash,
		"MOD_ROCKET":         ModRocket,
		"MOD_ROCKET_SPLASH":  ModRocketSplash,
		"MOD_PLASMA":         ModPlasma,
		"MOD_PLASMA_SPLASH":  ModPlasmaSplash,
		"MOD_RAILGUN":        ModRailGun,
		"MOD_LIGHTNING":      ModLightning,
		"MOD_BFG":            ModBfg,
		"MOD_BFG_SPLASH":     ModBfgSplash,
		"MOD_WATER":          ModWater,
		"MOD_SLIME":          ModSlime,
		"MOD_LAVA":           ModLava,
		"MOD_CRUSH":          ModCrush,
		"MOD_TELEFRAG":       ModTelefrag,
		"MOD_FALLING":        ModFalling,
		"MOD_SUICIDE":        ModSuicide,
		"MOD_TARGET_LASER":   ModTargetLaser,
		"MOD_TRIGGER_HURT":   ModTriggerHurt,
	}
)

const (
	Undefined MeansOfDeath = iota
	ModUnknown
	ModShotgun
	ModGauntlet
	ModMachineGun
	ModGrenade
	ModGrenadeSplash
	ModRocket
	ModRocketSplash
	ModPlasma
	ModPlasmaSplash
	ModRailGun
	ModLightning
	ModBfg
	ModBfgSplash
	ModWater
	ModSlime
	ModLava
	ModCrush
	ModTelefrag
	ModFalling
	ModSuicide
	ModTargetLaser
	ModTriggerHurt
)

func (m MeansOfDeath) String() string {
	if mod, ok := mapToString[m]; ok {
		return mod
	}
	return mapToString[Undefined]
}

func (m MeansOfDeath) Value() (driver.Value, error) {
	return m.String(), nil
}

func (m *MeansOfDeath) Scan(src interface{}) (err error) {
	switch src := src.(type) {
	case string:
		if mod, ok := mapStringTo[src]; ok {
			*m = mod
			return nil
		}

		return fmt.Errorf("unrecognized means of death value: %s", src)
	case nil:
		*m = Undefined
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", src)
	}
}

func Parse(str string) (MeansOfDeath, error) {
	if mod, ok := mapStringTo[str]; ok {
		return mod, nil
	}

	return Undefined, fmt.Errorf("unrecognized means of death value: %s", str)
}
