package volumetricweight

import (
	"fmt"
	"slices"
	"strings"

	"github.com/shopspring/decimal"
)

type VolumetricWeight struct {
	length         decimal.Decimal
	width          decimal.Decimal
	height         decimal.Decimal
	fromSizeUnit   string
	fromWeightUnit string
}

func New(length, width, height float64, sizeUnit, weightUnit string) (*VolumetricWeight, error) {
	if length <= 0 || width <= 0 || height <= 0 {
		return nil, fmt.Errorf("长度、宽度和高度必须大于 0")
	}
	if !slices.Contains([]string{"cm", "mm", "in", "m"}, sizeUnit) {
		return nil, fmt.Errorf("不支持的单位 %s", sizeUnit)
	}
	if !slices.Contains([]string{"g", "kg", "lb"}, weightUnit) {
		return nil, fmt.Errorf("不支持的单位 %s", weightUnit)
	}
	return &VolumetricWeight{
		length:         decimal.NewFromFloat(length),
		width:          decimal.NewFromFloat(width),
		height:         decimal.NewFromFloat(height),
		fromSizeUnit:   sizeUnit,
		fromWeightUnit: weightUnit,
	}, nil
}

// unitConverter 存储两个单位之间的转换信息。
type unitConverter struct {
	from, to string
	factor   string
}

// getConversionFactor 查找两个单位之间的转换因子。
func getConversionFactor(fromUnit, toUnit string, converters []unitConverter) (decimal.Decimal, error) {
	if strings.EqualFold(fromUnit, toUnit) {
		return decimal.NewFromInt(1), nil
	}
	index := slices.IndexFunc(converters, func(c unitConverter) bool {
		return c.from == fromUnit && c.to == toUnit
	})
	if index == -1 {
		return decimal.Zero, fmt.Errorf("不支持的单位转换 %s => %s", fromUnit, toUnit)
	}
	factorDec, err := decimal.NewFromString(converters[index].factor)
	if err != nil {
		return decimal.Zero, fmt.Errorf("转换因子 %s 无效: %w", converters[index].factor, err)
	}
	return factorDec, nil
}

func (vw *VolumetricWeight) Calc(toSizeUnit string, factor float64, toWeightUnit string, precision int32) (float64, error) {
	if factor <= 0 {
		return 0.0, fmt.Errorf("系数必须大于 0")
	}
	if toWeightUnit == "" {
		toWeightUnit = vw.fromWeightUnit
	}
	// factor = 将“from”单位乘以该值以获得“to”单位
	var sizeConverters = []unitConverter{
		{from: "cm", to: "in", factor: "0.3937007874015748031496062992126"}, // 1.0 / 2.54
		{from: "in", to: "cm", factor: "2.54"},
		{from: "mm", to: "cm", factor: "0.1"},
		{from: "cm", to: "mm", factor: "10"},
		{from: "in", to: "m", factor: "0.0254"},
		{from: "m", to: "in", factor: "39.37007874015748031496062992126"}, // 1 / 0.0254
		{from: "cm", to: "m", factor: "0.01"},
		{from: "m", to: "cm", factor: "100"},
		{from: "mm", to: "in", factor: "0.03937007874015748031496062992126"}, // 1.0 / 25.4
		{from: "in", to: "mm", factor: "25.4"},
		{from: "mm", to: "m", factor: "0.001"},
		{from: "m", to: "mm", factor: "1000"},
	}
	sizeConversionFactorDec, err := getConversionFactor(vw.fromSizeUnit, toSizeUnit, sizeConverters)
	if err != nil {
		return 0.0, err
	}

	volume := vw.length.Mul(vw.width).Mul(vw.height)
	volumeInToSizeUnit := volume.Mul(sizeConversionFactorDec.Pow(decimal.NewFromInt(3)))
	weightInFromUnit := volumeInToSizeUnit.Div(decimal.NewFromFloat(factor))
	var weightConverters = []unitConverter{
		{from: "g", to: "kg", factor: "0.001"},
		{from: "kg", to: "g", factor: "1000"},
		{from: "lb", to: "kg", factor: "0.45359237"},
		{from: "kg", to: "lb", factor: "2.2046226218"},
		{from: "lb", to: "g", factor: "453.59237"},
		{from: "g", to: "lb", factor: "0.0022046226218"},
	}
	weightConversionFactorDec, err := getConversionFactor(vw.fromWeightUnit, toWeightUnit, weightConverters)
	if err != nil {
		return 0.0, err
	}

	result := weightInFromUnit.Mul(weightConversionFactorDec).Round(precision)
	f, _ := result.Float64()
	return f, nil
}
