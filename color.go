package main

import (
	"fmt"
	"strconv"
)

// https://shipengqi.github.io/posts/2019-09-18-go-color/

// 示例：\033[1;31;40m%s\033[0m\n
// \033：\ 表示转义，\033 表示设置颜色。
// [1;31;40m：定义颜色，[ 表示开始颜色设置，m 为颜色设置结束，以 ; 号分隔。1 代码，表示显示方式，31 表示前景颜色（文字的 颜色），40 表示背景颜色。
// \033[0m：表示恢复终端默认样式。

// 代码 意义
// -------------------------
//  0  终端默认设置
//  1  高亮显示
//  4  使用下划线
//  5  闪烁
//  7  反白显示
//  8  不可见

// 3 位前景色, 4 位背景色

// 前景 背景 颜色
// ---------------------------------------
// 30  40  黑色
// 31  41  红色
// 32  42  绿色
// 33  43  黄色
// 34  44  蓝色
// 35  45  紫红色
// 36  46  青蓝色
// 37  47  白色
// 38  48  终端默认

// Color defines a single SGR Code
type Color int

// Foreground text colors
const (
	FgBlack Color = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
	FgDefault
)

// Background text colors
const (
	BgBlack Color = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
	BgDefault
)

// // Foreground Hi-Intensity text colors
// const (
// 	FgHiBlack Color = iota + 90
// 	FgHiRed
// 	FgHiGreen
// 	FgHiYellow
// 	FgHiBlue
// 	FgHiMagenta
// 	FgHiCyan
// 	FgHiWhite
// )

// Colorize a string based on given color.
func Colorize(s string, fgc Color, bgc Color) string {
	return fmt.Sprintf("\033[1;%s;%sm%s\033[0m", strconv.Itoa(int(fgc)), strconv.Itoa(int(bgc)), s)
}
