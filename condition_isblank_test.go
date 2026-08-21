package utils

import (
	"testing"
	"time"
)

// 底层是 int64 的具名类型——vectors 的 json.ID 就是这一类。
type namedID int64

type namedStr string

type someStruct struct {
	A int
	B string
}

// TestIsBlank_NamedTypesDoNotPanic 是这次修复的主判据。
//
// 原来 default 分支无条件调 reflect.ValueOf(v).Len()，而 Len() 只对
// Array/Slice/Map/String/Chan 合法。具名 int64 掉进 default 就是
// `reflect: call of reflect.Value.Len on int64 Value` —— 真栈上表现为
// **这台机器上建不了新租户**（/api/setup 在写 sys_tenant 那一步 panic），
// 而返回体带着 done:true，只看成功标志会误判成"卡住"。
func TestIsBlank_NamedTypesDoNotPanic(t *testing.T) {
	cases := []struct {
		name string
		val  any
		want bool
	}{
		{"具名 int64 零值", namedID(0), true},
		{"具名 int64 非零", namedID(2079533084918681600), false},
		{"具名 string 空", namedStr(""), true},
		{"具名 string 非空", namedStr("x"), false},
		{"结构体零值", someStruct{}, true},
		{"结构体非零", someStruct{A: 1}, false},
		{"nil 指针", (*someStruct)(nil), true},
		{"非 nil 指针", &someStruct{}, false},
		{"nil map", map[string]int(nil), true},
		{"空 map", map[string]int{}, true},
		{"非空 map", map[string]int{"a": 1}, false},
		{"nil 切片", []int(nil), true},
		{"空切片", []int{}, true},
		{"非空切片", []int{1}, false},
		{"数组零长", [0]int{}, true},
		{"数组非零长", [1]int{0}, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("IsBlank(%#v) panic 了: %v", c.val, r)
				}
			}()
			if got := IsBlank(c.val); got != c.want {
				t.Fatalf("IsBlank(%#v) = %v, want %v", c.val, got, c.want)
			}
		})
	}
}

// TestIsBlank_ExistingBehaviourUnchanged 钉住既有分支：这一处是全仓共用的，
// 只补 default 的判断，别的语义一个都不许动。
func TestIsBlank_ExistingBehaviourUnchanged(t *testing.T) {
	cases := []struct {
		val  any
		want bool
	}{
		{nil, true},
		{false, true}, {true, false},
		{0, true}, {1, false},
		{int64(0), true}, {int64(1), false},
		{uint(0), true}, {uint(1), false},
		{float64(0), true}, {float64(1.5), false},
		{"", true}, {"x", false},
		{time.Time{}, true},
		{[]byte(nil), true}, {[]byte{}, true}, {[]byte("a"), false},
	}
	for _, c := range cases {
		if got := IsBlank(c.val); got != c.want {
			t.Fatalf("IsBlank(%#v) = %v, want %v", c.val, got, c.want)
		}
	}
	// ★ *time.Time 那个 case 是**恒 false** 的：它写的是
	//     case *time.Time: blank = new(time.Time)
	//   然后 `return v == blank` —— 比的是两个**指针**，new() 出来的那个永远不等于
	//   传进来的那个。看意图应该是"指向零值时间的指针算空"，实际一个都判不出来。
	//
	//   这里只**钉住现状**，不顺手改：utils 是全仓共用的独立仓，把它改成 true 会让
	//   某些 orm 写入路径突然把零值时间当成"没填"而跳过写库，影响面远超本次要修的
	//   panic。留给专门的一次改动，连同调用方一起过。
	var zero time.Time
	if IsBlank(&zero) {
		t.Fatal("现状变了：*time.Time 分支原来恒为 false（见上面的说明），改动前请先查调用方")
	}
}
