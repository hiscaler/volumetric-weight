package volumetricweight

import (
	"fmt"
	"slices"
)

type VolumetricWeight struct {
	length   float64
	width    float64
	height   float64
	fromUnit string
}

func New(length, width, height float64, unit string) (*VolumetricWeight, error) {
	return &VolumetricWeight{
		length:   length,
		width:    width,
		height:   height,
		fromUnit: unit,
	}, nil
}

func (vw *VolumetricWeight) Calc(toUnit string, divisor float64) (float64, error) {
	if divisor <= 0 {
		return 0, fmt.Errorf("divisor must be greater than 0")
	}

	type caster struct {
		from, to string
		divisor  float64
	}
	var casters = []caster{
		{"cm", "in", 2.54},
		{"in", "cm", 1.0 / 2.54},
		{"mm", "cm", 10},
		{"cm", "mm", 1.0 / 10},
		{"in", "m", 1.0 / 39.3701},
		{"m", "in", 39.3701},
	}
	index := slices.IndexFunc(casters, func(c caster) bool {
		return c.from == vw.fromUnit && c.to == toUnit
	})
	if index == -1 {
		return 0, fmt.Errorf("unsupported unit conversion %s to %s", vw.fromUnit, toUnit)
	}
	value := casters[index].divisor
	weight := (vw.length * value * vw.width * value * vw.height * value) / divisor
	return weight, nil
}
