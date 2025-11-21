# 材积重计算

一个用于计算体积重量的 Go 包，支持多种单位和高精度计算。

## 特性

- 计算包裹的体积重量。
- 支持多种尺寸单位：`cm`、`mm`、`in`、`m`。
- 支持多种重量单位：`g`、`kg`、`lb`。
- 使用 `decimal` 进行高精度计算，以避免浮点数错误。
- 允许为计算指定自定义的系数（factor）。
- 允许为输出结果指定所需的精度。

## 安装

```bash
go get github.com/hiscaler/volumetric-weight
```

## 使用方法

这是一个如何使用该包的简单示例：

```go
package main

import (
	"fmt"
	"log"

	"github.com/hiscaler/volumetric-weight"
)

func main() {
	// 为一个 30x20x10 厘米的包裹创建一个新的计算器实例。
	// 初始单位：尺寸为 cm，重量为 g。
	vw, err := volumetricweight.New(30, 20, 10, "cm", "g")
	if err != nil {
		log.Fatalf("创建新的计算器失败: %v", err)
	}

	// 计算体积重量。
	// 在计算时将尺寸单位转换为 'in'。
	// 使用 5000 作为系数。
	// 将最终重量转换为 'kg'。
	// 将结果四舍五入到 4 位小数。
	precision := int32(4)
	volumetricWeight, err := vw.Calc("in", 5000, "kg", precision)
	if err != nil {
		log.Fatalf("计算体积重量失败: %v", err)
	}

	fmt.Printf("体积重量是: %v kg\n", volumetricWeight)
	// 预期输出: 体积重量是: 0.0073 kg
}
```

## API

### `New(length, width, height float64, sizeUnit, weightUnit string) (*VolumetricWeight, error)`

创建一个新的 `VolumetricWeight` 计算器实例。

- `length`, `width`, `height`: 包裹的尺寸。
- `sizeUnit`: 尺寸的单位（`cm`、`mm`、`in`、`m`）。
- `weightUnit`: 重量计算的基础单位（`g`、`kg`、`lb`）。

### `(vw *VolumetricWeight) Calc(toSizeUnit string, factor float64, toWeightUnit string, precision int32) (float64, error)`

计算体积重量。

- `toSizeUnit`: 在计算体积时，尺寸将转换到的单位。
- `factor`: 用于除以体积的体积系数。
- `toWeightUnit`: 计算出的重量的最终单位。如果为空，则使用原始的 `weightUnit`。
- `precision`: 最终结果要四舍五入到的小数位数。

## 许可证

该项目根据 BSD 3-Clause 许可证授权 - 详细信息请参阅 [LICENSE](LICENSE) 文件。
