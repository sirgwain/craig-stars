package cs

// Compute cloak percent (as an int) based on total cloakUnits
//
// https://wiki.starsautohost.org/wiki/Guts_of_Cloaking
func getCloakPercentForCloakUnits(cloakUnits int) (cloakPercent int) {
	// This is just a piecewise function
	if cloakUnits <= 100 {
		return cloakUnits/2 + cloakUnits%2 // round units up
	}

	cloakUnits -= 100
	if cloakUnits <= 200 {
		return 50 + cloakUnits/8
	}

	cloakUnits -= 200
	if cloakUnits < 312 {
		return 75 + cloakUnits/24
	}

	cloakUnits -= 312
	if cloakUnits <= 512 {
		return 88 + cloakUnits/64
	} else if cloakUnits < 768 {
		return 96
	} else if cloakUnits < 1000 {
		return 97
	}
	return 98
}

// get the factor to reduce scan ranges by based on a cloakPercent and
// a cloak reduction factor (i.e. tachyons)
func getCloakFactor(cloakPercent int, cloakReductionFactor float64) float64 {
	if cloakPercent <= 0 {
		return 1
	}
	return 1 - float64(cloakPercent)/100*cloakReductionFactor
}
