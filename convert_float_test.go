package utils

import (
	"encoding/json"
	"math"
	"testing"
	"time"
)

// ★ 这个文件盯的是同一个洞：类型 switch 认精确类型，漏掉 `case int:` 时 Go 里
// 最常见的整数形态会掉进 default，打一行 "unable to cast 0 of type int to
// float64" 然后返回 0。0 是合法金额，报表上分辨不出来是"真的零"还是"转换失败"。
func TestToFloat64_NativeTypes(t *testing.T) {
	type namedInt64 int64
	type namedFloat float64

	f32 := float32(1.5)
	i := 42

	cases := []struct {
		in   any
		want float64
	}{
		{int(0), 0}, // ← 就是日志里那一条
		{int(1500), 1500},
		{int(-7), -7},
		{int8(8), 8},
		{int16(16), 16},
		{int32(32), 32},
		{int64(64), 64},
		{uint(1), 1},
		{uint8(2), 2},
		{uint16(3), 3},
		{uint32(4), 4},
		{uint64(5), 5},
		{float32(1.5), 1.5},
		{float64(2.25), 2.25},
		{"3.5", 3.5},
		{json.Number("4.75"), 4.75},
		{true, 1},
		{false, 0},
		{nil, 0},
		{&i, 42}, // 指针要先解引用
		{&f32, 1.5},
		// 命名数值类型：精确类型对不上任何 case，靠 Kind 兜底。
		{namedInt64(9), 9},
		{namedFloat(0.5), 0.5},
		{time.Month(3), 3},
		{time.Weekday(2), 2},
		{3 * time.Second, float64(3 * time.Second)},
	}
	for _, c := range cases {
		got, err := ToFloat64E(c.in)
		if err != nil {
			t.Errorf("ToFloat64E(%#v) 报错: %v", c.in, err)
			continue
		}
		if got != c.want {
			t.Errorf("ToFloat64E(%#v) = %v，期望 %v", c.in, got, c.want)
		}
		if got := ToFloat64(c.in); got != c.want {
			t.Errorf("ToFloat64(%#v) = %v，期望 %v", c.in, got, c.want)
		}
	}
}

// 转换失败必须能被调用方看见——非 E 版本返回 0 是刻意的兼容行为，但 E 版本
// 必须回错，否则"失败"和"真的是 0"永远分不开。
func TestToFloat64E_ReportsFailure(t *testing.T) {
	for _, in := range []any{"abc", []string{"1"}, struct{}{}, map[string]int{}} {
		if v, err := ToFloat64E(in); err == nil {
			t.Errorf("ToFloat64E(%#v) 应当报错，实得 %v", in, v)
		}
		if v := ToFloat64(in); v != 0 {
			t.Errorf("ToFloat64(%#v) 失败时应回 0，实得 %v", in, v)
		}
	}
}

func TestToFloat32_NativeInt(t *testing.T) {
	if got, err := ToFloat32E(1500); err != nil || got != 1500 {
		t.Errorf("ToFloat32E(1500) = %v, %v", got, err)
	}
	if got := ToFloat32(0); got != 0 {
		t.Errorf("ToFloat32(0) = %v", got)
	}
	// 字符串走 32 位解析，不是先解成 float64 再截。
	if got, err := ToFloat32E("0.1"); err != nil || got != float32(0.1) {
		t.Errorf("ToFloat32E(\"0.1\") = %v, %v", got, err)
	}
	if _, err := ToFloat32E("abc"); err == nil {
		t.Error("ToFloat32E(\"abc\") 应当报错")
	}
}

// 命名整数类型走反射兜底时不许经 float64 中转：雪花 id 超过 2^53，中转一趟就被
// 舍成末位带 0 的假 id。
func TestToInt64_NamedTypeKeepsSnowflakePrecision(t *testing.T) {
	type recordID int64
	const id = int64(1806159265236451329) // > 2^53

	if got := ToInt64(recordID(id)); got != id {
		t.Errorf("ToInt64(recordID) = %d，期望 %d（差 %d）", got, id, got-id)
	}
	if got := ToInt(recordID(id)); int64(got) != id {
		t.Errorf("ToInt(recordID) = %d，期望 %d", got, id)
	}
	// float64 确实存不下，这里只是标明代价的边界。
	if int64(float64(id)) == id {
		t.Skip("这台机器上 float64 恰好无损，本用例失去意义")
	}
	if got := ToFloat64(recordID(id)); int64(got) == id {
		t.Errorf("ToFloat64 不该无损，实得 %v", got)
	}
	_ = math.MaxInt64
}
