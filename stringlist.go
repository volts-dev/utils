package utils

import (
	"fmt"
	//	"errors"
	"strconv"
	//	"strings"
	//	"unicode"
)

type (
	IStringList interface {
		AddSubList(strs ...string) IStringList
		Clear()
		Count() int
		Clone(idx ...int) (result IStringList)
		Find(itr interface{}) IStringList
		Flatten() (result []string)
		Has(strs ...interface{}) bool
		In(strs ...interface{}) bool
		Insert(idx int, item interface{})
		IsBaseInt() bool
		IsBaseString() bool
		IsBool() bool
		IsInt() bool
		IsList() bool
		Item(idx int) IStringList
		Items(idx ...int) (result []IStringList)
		Pop() (lst IStringList)
		Push(item ...IStringList) IStringList
		PushString(strs ...string) IStringList
		Quote(str string) string
		Remove(idx int) IStringList
		Reversed() (result IStringList)
		SetItem(idx int, item IStringList)
		SetText(text string, idx ...int)
		Shift() IStringList
		String(idx ...int) string
		Strings(idx ...int) (result []string)
		Union(other IStringList)
		Update()
	}

	TStringList struct {
		Text string
		//type int // node, leaf
		quoteChar        string
		StartPos, EndPos int
		items            []*TStringList
		updated          bool // -- 是否更新修改
		format_func      func(node *TStringList)
		//AutoUpdate       bool
	}
)

/*
func NewStringListParser(domain ...string) *TStringList {
	return NewDomain(domain...)
}

// [""]
func NewDomain(domain ...string) *TStringList {
	if len(domain) == 0 || (len(domain) == 1 && domain[0] == "") {
		return NewStringList() // 创建空白
	}

	lLexer := &TLexer{
		stream: domain[0],
		length: len(domain[0]),
		pos:    0,
		isEof:  false,
	}

	lParse := &TStringParser{
		lexer:   lLexer,
		strings: NewStringList(),
	}

	err := lParse.parse()
	if err != nil {
		return nil
	}
	return lParse.strings.Item(0)
}
*/
// 创建一个StringList
func NewStringList(strs ...string) (res_lst *TStringList) {
	res_lst = &TStringList{
		quoteChar:   `'`,
		updated:     true,
		format_func: _update,
	}

	lCnt := len(strs)
	if lCnt > 0 {
		if len(strs) == 1 {
			res_lst.Text = strs[0] // 创建空白
		} else {
			res_lst.PushString(strs...)
		}
	}

	return
}

func PrintStringList(list *TStringList) {
	fmt.Printf(
		"[Root]  Count:%d  Text:%v  IsList:%v ",
		list.Count(),
		list.String(),
		list.IsList(),
	)
	fmt.Println()

	_printStringList(1, list)
}
func _printStringList(idx int, list *TStringList) {
	for i, Item := range list.Items() {
		for j := 0; j < idx; j++ { // 空格距离ss
			fmt.Print("  ")
		}
		if idx > 0 {
			fmt.Print("┗", "  ")
		}
		fmt.Printf(
			"[%d]  Count:%d  Text:%v  IsList:%v ",
			i,
			Item.Count(),
			Item.String(),
			Item.IsList(),
		)

		fmt.Println()
		_printStringList(idx+1, Item)
	}
}

//-----------list
func (self *TStringList) Item(idx int) *TStringList {
	//self.Update()

	if self.IsList() {
		if idx < len(self.items) {
			return self.items[idx]
		}
	}

	return self // 超范围返回自己
}

// 添加 term
func (self *TStringList) SetItem(idx int, item *TStringList) {
	self.items[idx] = item

	self.updated = false
}

// 修改非List 的Text
// idx 只会自动添加当新的Items=0时
// TODO:自动生成Items 可修改List
func (self *TStringList) SetText(text string, idx ...int) {
	if len(idx) == 0 {
		// 不是List才能直接修改TEXT 否者需要解析到Items 再生成Text
		if !self.IsList() {
			self.Text = text
		}
	} else if len(self.items)-1 >= idx[0] {
		self.items[idx[0]] = NewStringList(text) //
	} else {
		// 不是List才能直接修改TEXT 否者需要解析到Items 再生成Text
		if !self.IsList() {
			self.items[idx[0]].Text = text
		}
	}

	self.updated = false

}

//栈方法Pop :取栈方式出栈<最后一个>元素 即最后一个添加进列的元素
func (self *TStringList) Shift() *TStringList {
	var lst *TStringList
	if len(self.items) > 0 {
		lst = self.items[0]
		self.items = self.items[1:]
	}

	// # 正式清空所有字符
	if len(self.items) == 0 {
		self.Text = ""
	}

	self.updated = false
	return lst
}

//""" Pop a leaf to process. """
//栈方法Pop :取栈方式出栈<最后一个>元素 即最后一个添加进列的元素
func (self *TStringList) Pop() (lst *TStringList) {
	cnt := len(self.items)
	if cnt == 0 {
		return
	}

	lst = self.items[cnt-1]
	self.items = self.items[:cnt-1]

	//# 正式清空所有字符 避免Push时添加Text为新item
	if len(self.items) == 0 {
		self.Text = ""
	}
	self.updated = false
	return lst
}

//栈方法Push：叠加元素
func (self *TStringList) __Push(items ...interface{}) *TStringList {
	if len(items) > 0 {
		// 当Self是一个值时必须添加自己到items里成为列表的一部分
		add_self := self.Text != "" && len(self.items) == 0
		for _, item := range items {
			switch item.(type) {
			case string:
				if add_self {
					self.items = append(self.items, NewStringList(self.Text))
					add_self = false
				}
				self.items = append(self.items, NewStringList(item.(string)))
				self.updated = false
			case *TStringList:
				if add_self {
					self.items = append(self.items, NewStringList(self.Text))
					add_self = false
				}
				self.items = append(self.items, item.(*TStringList))
				self.updated = false
			}
		}
	}

	return self
}

//栈方法Push：叠加元素
func (self *TStringList) Push(item ...*TStringList) *TStringList {
	// 当Self是一个值时必须添加自己到items里成为列表的一部分
	if self.Text != "" && len(self.items) == 0 {
		self.items = append(self.items, NewStringList(self.Text))
	}
	self.items = append(self.items, item...)
	self.updated = false
	return self
}

// #添加字符串到该 List
func (self *TStringList) PushString(strs ...string) *TStringList {
	// 当Self是一个值时必须添加自己到items里成为列表的一部分
	if self.Text != "" && len(self.items) == 0 {
		self.items = append(self.items, NewStringList(self.Text))
	}

	for _, str := range strs {
		//logger.Dbg("PushString", str, self)
		self.items = append(self.items, NewStringList(str))
	}

	self.updated = false

	return self
}

// # all the strs will push in a list and insert as a item to current list
func (self *TStringList) AddSubList(strs ...string) *TStringList {
	self.Push(NewStringList(strs...))
	return self
}

func (self *TStringList) Insert(idx int, item interface{}) {
	var must_move bool = false
	var new_item *TStringList

	switch item.(type) {
	case string:
		new_item = NewStringList(item.(string))
		self.items = append(self.items, new_item)
		must_move = true
	case *TStringList:
		new_item = item.(*TStringList)
		// Grow the slice by one element.
		// make([]Token, len(self.Child)+1)
		// self.Child[0 : len(self.Child)+1]
		self.items = append(self.items, new_item)
		// Use copy to move the upper part of the slice out of the way and open a hole.
		must_move = true
	default:
		// 不支持
	}

	if must_move {
		copy(self.items[idx+1:], self.items[idx:])
		// Store the new value.
		self.items[idx] = new_item
		// Return the result.
	}
	self.updated = false
	return
}

func (self *TStringList) __InsertString(idx int, str string) {
	self.Insert(idx, &TStringList{Text: str})
}

// TODO: 为避免错乱,移除后复制一个新的返回结果
func (self *TStringList) Remove(idx int) *TStringList {
	self.items = append(self.items[:idx], self.items[idx+1:]...)

	self.updated = false
	return self
}

// 所有Item 都是String
func (self *TStringList) IsBaseString() bool {
	for _, item := range self.items {
		if item.IsList() {
			return false
		}
	}

	return true
}

// 生成集合 one 和集合 other 的并集
func (self *TStringList) Union(other *TStringList) {
	if other == nil || other.Count() == 0 {
		return
	}

	for _, v := range other.items {
		self.Push(v)
	}

	self.updated = false
	return
}

func (self *TStringList) IsBaseInt() bool {
	for _, item := range self.items {
		if item.IsList() || !item.IsInt() {
			return false
		}
	}

	return true
}

// 废弃retur the list length
func (self *TStringList) __Len() int {
	return len(self.items)
}

// return the list length
func (self *TStringList) Count() int {
	return len(self.items)
}

func (self *TStringList) IsList() bool {
	return len(self.items) > 0
}

func (self *TStringList) In(strs ...interface{}) bool {
	for _, itr := range strs {
		switch itr.(type) {
		case string: // 处理字符串
			if self.Text == itr.(string) {
				return true
			}
		case *TStringList: // 处理*TStringList 类型
			if self.Text == itr.(*TStringList).Text {
				return true
			}

		}
	}

	return false
}

// find by string or object
func (self *TStringList) Find(itr interface{}) *TStringList {
	for _, itm := range self.items {
		if itm.IsList() {
			return itm.Find(itr)
		} else if str, ok := itr.(string); ok {
			if itm.Text == str {
				return itm
			}
		} else if item, ok := itr.(*TStringList); ok {
			if itm.Text == item.Text {
				return itm
			}
		}
	}

	return nil
}

//检查是否包含 一个或者多个内容
func (self *TStringList) Has(strs ...interface{}) bool {
	for _, itr := range strs {
		if self.Find(itr) != nil {
			return true
		}
		/*
			// 处理字符串
			if str, ok := itr.(string); ok {
				if strings.Index(self.String(), str) == -1 {
					return false // 如果一个查不到就返回否
				}
			} else if item, ok := itr.(*TStringList); ok {
				// 处理*TStringList 类型
				if strings.Index(self.String(), item.text) == -1 {
					return false
				}
			}
		*/
	}

	return false
}

// ('xx','xx')
func (self *TStringList) __AsString(idx ...int) string {
	if self.IsList() {
		return self.String(idx...)
	}

	return "'" + self.String(idx...) + "'"
}

// return 'xx','xx'
// 当self 时列表时idx有效,反则返回Text
func (self *TStringList) String(idx ...int) string {
	if len(idx) > 0 {
		if self.IsList() {
			if idx[0] > -1 && idx[0] < len(self.items) {
				return self.items[idx[0]].Text
			} else {
				return "" // 超范围返回空
			}
		}
	}

	//self.Update()
	return self.Text
}

// 返回所有Items字符
func (self *TStringList) Strings(idx ...int) (result []string) {
	//self.Update()

	lCnt := len(idx)

	if lCnt == 0 { // 返回所有
		if len(self.items) == 0 {
			result = append(result, self.Text)
		} else {
			for _, item := range self.items {
				result = append(result, item.Text)
			}
		}

	} else if lCnt == 1 {
		result = append(result, self.items[idx[0]].Text)

	} else if lCnt > 1 {

		for _, item := range self.items[idx[0]:idx[1]] {
			result = append(result, item.Text)
		}
	}

	return
}

// 复制一个反转版
func (self *TStringList) Reversed() (result *TStringList) {
	//self.Update()

	result = NewStringList()
	lCnt := self.Count()
	for i := lCnt - 1; i >= 0; i-- {
		result.Push(self.items[i]) //TODO: 复制
	}
	return
}

func (self *TStringList) Clear() {
	self.StartPos = 0
	self.EndPos = 0
	self.updated = true
	self.Text = ""
	self.items = nil // make([]*TStringList, 0)
}

//在原有基础上克隆
// len(idx)==0:返回所有
// len(idx)==1:返回Idx 指定item
// len(idx)>1:返回Slice 范围的items
func (self *TStringList) Clone(idx ...int) (result *TStringList) {
	//self.Update()

	lCnt := len(idx)
	if lCnt == 0 {
		result = NewStringList()
		result.Push(self.Items(0, self.Count()-1)...)
	} else if lCnt == 1 && idx[0] < self.Count() { // idex 必须小于Self长度
		result = NewStringList()
		result.Push(self.Items(idx[0], self.Count()-1)...) //result.Push(self.items[idx[0]])
	} else if lCnt > 1 && idx[0] < self.Count() && idx[1] < self.Count() {
		result = NewStringList()
		if idx[1] == -1 {
			// 复制到end
			result.Push(self.items[idx[0]:self.Count()]...)
		} else if idx[0] < idx[1] {
			// 复制到Offset
			result.Push(self.items[idx[0]:idx[1]]...)
		}

	}

	return
}

// len(idx)==0:返回所有
// len(idx)==1:返回Idx 指定item
// len(idx)>1:返回Slice 范围的items
func (self *TStringList) Items(idx ...int) (result []*TStringList) {
	//self.Update()

	lCnt := len(idx)

	/*  ??? 待考究
	if !self.IsList() {
		// 如果是空白的对象 返回Nil
		if self.Text == "" {
			return nil
		}

		return []*TStringList{self}
	}
	*/

	if lCnt == 0 {
		return self.items // 返回所有
	} else if lCnt == 1 && idx[0] < self.Count() { // idex 必须小于Self长度
		result = append(result, self.items[idx[0]])
	} else if lCnt > 1 && idx[0] < self.Count() && idx[1] < self.Count() {
		//if idx[0]>idx[1]{
		//
		//}

		for _, item := range self.items[idx[0]:idx[1]] {
			result = append(result, item)
		}
	}

	return
}

/*"Flatten a list of elements into a uniqu list
Author: Christophe Simonis (christophe@tinyerp.com)

Examples::
>>> flatten(['a'])
['a']
>>> flatten('b')
['b']
>>> flatten( [] )
[]
>>> flatten( [[], [[]]] )
[]
>>> flatten( [[['a','b'], 'c'], 'd', ['e', [], 'f']] )
['a', 'b', 'c', 'd', 'e', 'f']
>>> t = (1,2,(3,), [4, 5, [6, [7], (8, 9), ([10, 11, (12, 13)]), [14, [], (15,)], []]])
>>> flatten(t)
[1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15]
"*/
// 返回列表中的所有值
func (self *TStringList) Flatten() (result []string) {
	//self.Update()

	// # 当StringList作为一个单字符串
	if self.Text != "" && len(self.items) == 0 {
		result = []string{self.Text}
		return
	}

	// # 当StringList作为字符串组
	for _, lst := range self.items {
		if lst.IsList() {
			result = append(result, lst.Strings()...)
		} else {
			result = append(result, lst.String())
		}
	}
	return
}

// 废弃
func (self *TStringList) __Text(idx ...int) string {
	//self.Update()

	if len(idx) > 0 {
		return self.Item(idx[0]).Text
	}
	return self.Text
}

//废弃
func (self *TStringList) __Type() int {
	if len(self.items) == 0 {
		return 0
	}
	return 0
}

// 添加''
func (self *TStringList) Quote(str string) string {
	return self.quoteChar + str + self.quoteChar
}

// 更新生成所有Text内容
func _update(node *TStringList) {
	/*	//STEP  如果是Value Object 不处理
		IsList := false
		if len(node.items) == 0 {
			return
		} else {
			IsList = true
		}

		// 处理有Child的Object
		lStr := ""
		lStrLst := make([]string, 0)

		for _, item := range node.items {
			if is_leaf(item) {
				lStr = `(` + node._quote(item.String(0)) + `, ` + node._quote(item.String(1)) + `, `
				if item.Item(2).IsList() {
					lStr = lStr + `[` + item.String(2) + `])`
				} else {
					lStr = lStr + node._quote(item.String(2)) + `)`
				}

				item.text = lStr
				//utils.Dbg("_update leaf", lStr)
				IsList = true
				lStrLst = append(lStrLst, item.text)
			} else if item.IsList() {
				//utils.Dbg("IsList", item.text)
				_update(item)
				lStrLst = append(lStrLst, item.text)
			} else {
				//utils.Dbg("_update val", item.text)
				lStrLst = append(lStrLst, node._quote(item.text))
			}
		}

		/*	// 组合 XX,XX
			lStrLst := make([]string, 0)
			for _, item := range node.items {
				if item.IsList() {
					lStrLst = append(lStrLst, item.text)
				} else {
					lStrLst = append(lStrLst, self._quote(item.text))
				}
			}*/
	/*	lStr = strings.Join(lStrLst, ",")

		//if lStr == "" {
		//	lStr = self._quote(node.text) // 如果是Val Node
		//}

		// 组合[XX,XX]
		//if self == node && IsList {
		if IsList {
			lStr = `[` + lStr + `]`
		}

		node.text = lStr
		//utils.Dbg("_update lst", lStr)
	*/
}

// 更新生成所有Text内容
func (self *TStringList) Update() {
	if !self.updated {
		self.format_func(self)
		self.updated = true
	}

	/*for _, item := range self.items {

		if len(item.items) > 0 {
			//utils.Dbg(method)
			printNode(1, self.Root[method])
			utils.Dbg()
		}
	}*/
}

func (self *TStringList) FormatFunc(f func(self *TStringList)) {
	self.format_func = f
}

func (self *TStringList) AsBool() bool {
	if b, err := strconv.ParseBool(self.Text); err == nil {
		return b //	fmt.Printf("%T, %v\n", s, s)
	}
	return false
}

func (self *TStringList) AsInt() int64 {
	if b, err := strconv.ParseInt(self.Text, 10, 0); err == nil {
		return b //	fmt.Printf("%T, %v\n", s, s)
	}
	return -1
}

func (self *TStringList) IsBool() bool {
	if _, err := strconv.ParseBool(self.Text); err == nil {
		return true //	fmt.Printf("%T, %v\n", s, s)
	}
	return false
}

func (self *TStringList) IsInt() bool {
	if _, err := strconv.ParseInt(self.Text, 10, 0); err == nil {
		return true //	fmt.Printf("%T, %v\n", s, s)
	}
	return false
}
