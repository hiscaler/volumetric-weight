package volumetricweight

import (
	"math"
	"testing"

	"github.com/shopspring/decimal"
)

func TestNew(t *testing.T) {
	// 测试正常创建
	vw, err := New(10, 20, 30, "cm", "g")
	if err != nil {
		t.Fatalf("Failed to create VolumetricWeight: %v", err)
	}
	// 使用 decimal 比较结构体字段
	if !vw.length.Equal(decimal.NewFromFloat(10)) || !vw.width.Equal(decimal.NewFromFloat(20)) || !vw.height.Equal(decimal.NewFromFloat(30)) || vw.fromSizeUnit != "cm" || vw.fromWeightUnit != "g" {
		t.Errorf("New() failed to initialize fields correctly")
	}

	// 测试无效尺寸
	_, err = New(-10, 20, 30, "cm", "g")
	if err == nil {
		t.Error("Expected error for negative length, but got none")
	}

	// 测试无效尺寸单位
	_, err = New(10, 20, 30, "km", "g")
	if err == nil {
		t.Error("Expected error for unsupported size unit, but got none")
	}

	// 测试无效重量单位
	_, err = New(10, 20, 30, "cm", "ton")
	if err == nil {
		t.Error("Expected error for unsupported weight unit, but got none")
	}
}

func TestVolumetricWeight_Calc(t *testing.T) {
	tests := []struct {
		name           string
		length         float64
		width          float64
		height         float64
		fromSizeUnit   string
		fromWeightUnit string
		toSizeUnit     string
		toWeightUnit   string
		factor         float64
		precision      int32
		expected       float64
		expectedError  bool
	}{
		{
			name:           "cm to in, g to kg, factor 5000",
			length:         10,
			width:          10,
			height:         10,
			fromSizeUnit:   "cm",
			fromWeightUnit: "g",
			toSizeUnit:     "in",
			toWeightUnit:   "kg",
			factor:         5000,
			precision:      8,
			expected:       0.00001220, // Corrected value
			expectedError:  false,
		},
		{
			name:           "in to cm, lb to kg, factor 6000",
			length:         5,
			width:          5,
			height:         5,
			fromSizeUnit:   "in",
			fromWeightUnit: "lb",
			toSizeUnit:     "cm",
			toWeightUnit:   "kg",
			factor:         6000,
			precision:      4,
			expected:       0.1549, // Corrected value
			expectedError:  false,
		},
		{
			name:           "mm to cm, kg to g",
			length:         100,
			width:          100,
			height:         100,
			fromSizeUnit:   "mm",
			fromWeightUnit: "kg",
			toSizeUnit:     "cm",
			toWeightUnit:   "g",
			factor:         1,
			precision:      0,
			expected:       1000000,
			expectedError:  false,
		},
		{
			name:           "cm to mm, g to lb",
			length:         10,
			width:          10,
			height:         10,
			fromSizeUnit:   "cm",
			fromWeightUnit: "g",
			toSizeUnit:     "mm",
			toWeightUnit:   "lb",
			factor:         1,
			precision:      8,
			expected:       2204.6226218,
			expectedError:  false,
		},
		{
			name:           "in to m, lb to kg",
			length:         39.3701,
			width:          39.3701,
			height:         39.3701,
			fromSizeUnit:   "in",
			fromWeightUnit: "lb",
			toSizeUnit:     "m",
			toWeightUnit:   "kg",
			factor:         1,
			precision:      8,
			expected:       0.4535931,
			expectedError:  false,
		},
		{
			name:           "m to in, kg to lb",
			length:         1,
			width:          1,
			height:         1,
			fromSizeUnit:   "m",
			fromWeightUnit: "kg",
			toSizeUnit:     "in",
			toWeightUnit:   "lb",
			factor:         1,
			precision:      2,
			expected:       134534.33,
			expectedError:  false,
		},
		{
			name:           "same units",
			length:         10,
			width:          10,
			height:         10,
			fromSizeUnit:   "cm",
			fromWeightUnit: "g",
			toSizeUnit:     "cm",
			toWeightUnit:   "g",
			factor:         1,
			precision:      0,
			expected:       1000,
			expectedError:  false,
		},
		{
			name:           "unsupported size conversion",
			length:         10,
			width:          10,
			height:         10,
			fromSizeUnit:   "cm",
			fromWeightUnit: "g",
			toSizeUnit:     "km",
			toWeightUnit:   "g",
			factor:         1,
			precision:      2,
			expected:       0,
			expectedError:  true,
		},
		{
			name:           "unsupported weight conversion",
			length:         10,
			width:          10,
			height:         10,
			fromSizeUnit:   "cm",
			fromWeightUnit: "g",
			toSizeUnit:     "cm",
			toWeightUnit:   "ton",
			factor:         1,
			precision:      2,
			expected:       0,
			expectedError:  true,
		},
		{
			name:           "invalid factor",
			length:         10,
			width:          10,
			height:         10,
			fromSizeUnit:   "cm",
			fromWeightUnit: "g",
			toSizeUnit:     "in",
			toWeightUnit:   "g",
			factor:         0,
			precision:      2,
			expected:       0,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vw, err := New(tt.length, tt.width, tt.height, tt.fromSizeUnit, tt.fromWeightUnit)
			if err != nil {
				if !tt.expectedError {
					t.Fatalf("Failed to create VolumetricWeight: %v", err)
				} else {
					return
				}
			}

			got, err := vw.Calc(tt.toSizeUnit, tt.factor, tt.toWeightUnit, tt.precision)
			if (err != nil) != tt.expectedError {
				t.Errorf("VolumetricWeight.Calc(%s) error = %v, expectedError %v", tt.name, err, tt.expectedError)
				return
			}

			if !tt.expectedError {
				// Allow a small tolerance for float64 comparison
				if math.Abs(got-tt.expected) > 1e-9 {
					t.Errorf("VolumetricWeight.Calc(%s) = %v, want %v", tt.name, got, tt.expected)
				}
			}
		})
	}
}

// 测试仅指定尺寸单位，重量单位使用默认值的情况
func TestVolumetricWeight_CalcWithDefaultWeightUnit(t *testing.T) {
	vw, err := New(10, 10, 10, "cm", "g")
	if err != nil {
		t.Fatalf("Failed to create VolumetricWeight: %v", err)
	}

	// 调用Calc但不指定目标重量单位，应该使用原始重量单位
	result, err := vw.Calc("in", 1, "", 8)
	if err != nil {
		t.Errorf("VolumetricWeight.Calc() error = %v", err)
		return
	}

	expected := 61.02374409
	if math.Abs(result-expected) > 1e-9 {
		t.Errorf("VolumetricWeight.Calc() = %v, want %v", result, expected)
	}
}
