package utils

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
	//	wu "webgo/utils"
)

const (
	// Domain operators.
	NOT_OPERATOR = "!"
	OR_OPERATOR  = "|"
	AND_OPERATOR = "&"

	TRUE_LEAF  = "(1, '=', 1)"
	FALSE_LEAF = "(0, '=', 1)"

	TRUE_DOMAIN  = "[" + TRUE_LEAF + "]"
	FALSE_DOMAIN = "[" + FALSE_LEAF + "]"
)

var (
	DOMAIN_OPERATORS = []string{NOT_OPERATOR, OR_OPERATOR, AND_OPERATOR}

	/*# List of available term operators. It is also possible to use the '<>'
	# operator, which is strictly the same as '!='; the later should be prefered
	# for consistency. This list doesn't contain '<>' as it is simpified to '!='
	# by the normalize_operator() function (so later part of the code deals with
	# only one representation).
	# Internals (i.e. not available to the user) 'inselect' and 'not inselect'
	# operators are also used. In this case its right operand has the form (subselect, params).
	*/
	TERM_OPERATORS = []string{"=", "!=", "<=", "<", ">", ">=", "=?", "=like", "=ilike",
		"like", "not like", "ilike", "not ilike", "in", "not in", "child_of"}
)

type (

	// lexer提供数据的扫描工作
	TLexer struct {
		length int
		stream string //即将被扫描的文件流 一般放在执行扫描的类里
		//leftdelim  string //开始的标志
		//rightdelim string //结束的标志
		pos int // 当前游标

		curline int
		char    byte
		isEof   bool
		/*
			begin     int //从这里开始
			end       int //结束的地方
			offset    int // current position in the input.
			ch        rune
			width     int // width of last rune read from input.
			insideTag bool
		*/
	}

	// Parser 提供解析的逻辑 和数据的变换
	TStringParser struct {
		//templete *Parser
		//parent   *Parser
		lexer   *TLexer
		strings *TStringList
	}

	TStringList struct {
		text string
		//type int // node, leaf
		quoteChar        string
		startPos, endPos int
		items            []*TStringList
		updated          bool // -- 是否更新修改
		//AutoUpdate       bool
	}
)

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

// 创建一个StringList
func NewStringList(strs ...string) (res_lst *TStringList) {
	res_lst = &TStringList{
		quoteChar: `'`,
		updated:   true,
	}

	lCnt := len(strs)
	if lCnt > 0 {
		if len(strs) == 1 {
			res_lst.text = strs[0] // 创建空白
		} else {
			res_lst.PushString(strs...)
		}
	}

	return
}

//主要-移动偏移
//
func (self *TLexer) next() byte {
	if self.pos >= self.length-1 { //如果大于Buf 则停止
		//selfl.width = 0
		self.isEof = true
		return 0
	}
	//utils.Dbg("next", self.pos, self.stream[self.pos])
	// 直接返回第一个字符
	if self.pos == 0 && self.char == 0 {
		self.char = self.stream[self.pos]
		//self.pos += 1
		return self.char
	}

	self.pos += 1
	self.char = self.stream[self.pos]
	return self.char
}

func (self *TLexer) backup() {
	if self.pos > 0 {
		self.pos--
	}
	self.char = self.stream[self.pos]
}

//主要-略过特殊字符移动
//保持上一个有效字符
func (self *TLexer) consume_whitespace() {
	lChar := self.char
	for {
		self.next()
		if self.isEof {
			break
		}

		//utils.Dbg("consume_whitespace", string(self.char))
		if self.char == ' ' || self.char == '\t' || self.char == '\n' || self.char == '\r' {
			continue
		} else {
			self.backup()
			//utils.Dbg("consume_whitespace", string(self.char))
			break
			//return self.char
		}
	}

	self.char = lChar
	//utils.Dbg("consume_whitespace", string(self.char))
}

// 判断扫描 {{
/*
	循环匹配 {{ 直到找到完全匹配的
*/

func (self *TLexer) scan_string(s string) string { // s={{
	var (
		lStr string
		//		lChar byte
		pos int

		//		match bool = true
		//		err error
	)

	if s == "" {
		return "" //, errors.New("invaid string")
	}

	// 如果查询的超出范围
	if self.pos+len(s) >= self.length { //是否已经到末端了??
		lStr = self.stream[self.pos:]
		self.pos = self.length // 移到末端
		return lStr            //, io.EOF    //如果是-->获取当前游标以后的数据
	}

	pos = self.pos + 1
	for { //循环直到匹配到"{{"返回
		self.consume_whitespace()
		self.next()
		if self.isEof {
			return ""
		}
		//if self.stream[cur_pos] == '\n' { //判断是否到行末
		//	newlines++
		//}

		// 直到找到 第一个 和S 要查询的字符一样
		if self.char != s[0] { //当前是不是 目标字符
			continue //重新循环
		}

		//pos = self.pos                    //偏移到包含 S之内的 游标地址
		text := self.stream[pos:self.pos] //读取{{这段数据
		//self.pos = pos                    //设置偏移为当前偏移					/* 2 offset n */

		//utils.Dbg("scan_string match", pos, self.pos)
		return text

		// 没找到的话继续循环,这是转向下一个字符的唯一代码
	}

	//should never be here
	return ""
}

// s Html文件流
/*
[('foo', '=', 'bar')]
foo = 'bar'

[('id', 'in', [1,2,3])]
id in (1, 2, 3)

[('field', '=', 'value'), ('field', '<>', 42)]
( field = 'value' AND field <> 42 )

[('&', ('field', '<', 'value'), ('field', '>', 'value'))]
( field < 'value' AND field > 'value' )

[('|', ('field', '=', 'value'), ('field', '=', 'value'))]
( field = 'value' OR field = 'value' )

[('&', ('field1', '=', 'value'), ('field2', '=', 'value'), ('|', ('field3', '<>', 'value'), ('field4', '=', 'value')))]
( field1 = 'value' AND field2 = 'value' AND ( field3 <> 'value' OR field4 = 'value' ) )

[('&', ('|', ('a', '=', 1), ('b', '=', 2)), ('|', ('c', '=', 3), ('d', '=', 4)))]
( ( a = 1 OR b = 2 ) AND ( c = 3 OR d = 4 ) )

[('|', (('a', '=', 1), ('b', '=', 2)), (('c', '=', 3), ('d', '=', 4)))]
( ( a = 1 AND b = 2 ) OR ( c = 3 AND d = 4 ) )
*/
// ['&', ('active', '=', True), ('value', '!=', 'foo')]
// ['|', ('active', '=', True), ('state', 'in', ['open', 'draft'])
// ['&', ('active', '=', True), '|', '!', ('state', '=', 'closed'), ('state', '=', 'draft')]
// ['|', '|', ('state', '=', 'open'), ('state', '=', 'closed'), ('state', '=', 'draft')]
// ['!', '&', '!', ('id', 'in', [42, 666]), ('active', '=', False)]
// ['!', ['=', 'company_id.name', ['&', ..., ...]]]
//	[('picking_id.picking_type_id.code', '=', 'incoming'), ('location_id.usage', '!=', 'internal'), ('location_dest_id.usage', '=', 'internal')]

func (self *TStringParser) _parselist(list *TStringList) error {
	// 处理规则
	// List 不直接跳出 [('field', '=', 'value'), ('field', '<>', 42)] 如果处理到第一个跳出 第二个则无法继续
	// 存储值的List 可以直接跳出 例如：('id', 'in', [42, 666])  100%确定不需要处理其他List

	var (
		lPos  int = 0
		lStr  string
		lVals []string
	)

	//utils.Dbg("begin for")
	// 不是文件末端
	for !self.lexer.isEof {

		// 忽略空白
		self.lexer.consume_whitespace()

		//utils.Dbg("switch", string(self.lexer.char))
		switch string(self.lexer.char) {
		case "(": //('active', '=', True)
			{
				// 新建链
				lLst := &TStringList{startPos: self.lexer.pos}

				// 移动游标到下一个非空格
				self.lexer.consume_whitespace()
				self.lexer.next()
				if self.lexer.isEof {
					goto exit
				}

				/*	if isAlphaNumeric(rune(self.lexer.char)) {
						lStr = self.lexer.scan_string(")")
						if self.lexer.isEof {
							goto exit
						}
						utils.Dbg("(", string(self.lexer.char), lStr)
						lVals = strings.Split(lStr, ",")
						for _, val := range lVals {
							lLst.items = append(lLst.items, &TStringList{text: val})
						}

					} else { // ['&', (
						utils.Dbg("(", string(self.lexer.char), lStr)
				*/
				//self.lexer.backup()
				err := self._parselist(lLst)
				if err != nil {

				}
				//}
				lLst.endPos = self.lexer.pos + 1
				lLst.text = self.lexer.stream[lLst.startPos:lLst.endPos]
				//utils.Dbg("()list ( text", lLst.text)
				list.items = append(list.items, lLst)

				//goto exit // NOTE:末尾不能直接跳出，有可能有下一个Listt要执行
			}
		case "[": //['&', ('active', '=', True), ('value', '!=', 'foo')]
			{
				// 新建链
				lLst := &TStringList{startPos: self.lexer.pos}

				self.lexer.consume_whitespace()
				self.lexer.next()
				if self.lexer.isEof {
					goto exit
				}

				if isAlphaNumeric(rune(self.lexer.char)) { // [1,2,3]
					self.lexer.backup()
					lStr = self.lexer.scan_string("]")
					if self.lexer.isEof {
						goto exit
					}

					//utils.Dbg("New val ", string(self.lexer.char), lStr)

					lVals = strings.Split(lStr, ",")
					for _, val := range lVals {
						lLst.items = append(lLst.items, &TStringList{text: val})
					}

					// 读写字符串传
					lLst.endPos = self.lexer.pos + 1
					lLst.text = self.lexer.stream[lLst.startPos:lLst.endPos]
					//utils.Dbg("list text", lLst.text)

					// 插入
					list.items = append(list.items, lLst)

					// 继续处理后面的")"  ...[42, 666])
					self.lexer.consume_whitespace()
					self.lexer.next()
					if self.lexer.isEof {
						goto exit
					}

					if self.lexer.char == ')' || self.lexer.char == ']' {
						goto exit // 列表末端 跳出
						break
					} else {
						self.lexer.backup() // 回滚 继续
					}
				} else { // ['&', (
					//utils.Dbg("[ into parselist", string(self.lexer.char), lStr)
					err := self._parselist(lLst)
					if err != nil {

					}

					//self.lexer.consume_whitespace()
					//self.lexer.next()
					//lStr = self.lexer.scan_string("]")
					//utils.Dbg("list rnd", lLst.text)
				}

				// 读写字符串传
				lLst.endPos = self.lexer.pos + 1
				lLst.text = self.lexer.stream[lLst.startPos:lLst.endPos]
				//utils.Dbg("[]list text", lLst.text)

				// 插入
				list.items = append(list.items, lLst)

				//goto exit // NOTE:末尾不能直接跳出，有可能有下一个Listt要执行
			}

		case "'":
			{

				//获取另一个引号 ('xx's', 可能字符有"’" 所以必须遇见","
				for {
					lStr = self.lexer.scan_string("'")
					if self.lexer.isEof {
						goto exit
					}
					//utils.Dbg("':", string(self.lexer.char), lStr)

					// 确定下一个字符是 "," 这是一个Block结束
					self.lexer.consume_whitespace()
					self.lexer.next()
					if self.lexer.isEof {
						goto exit
					}
					if self.lexer.char == ',' {
						//utils.Dbg("New text", lStr)
						list.items = append(list.items, &TStringList{text: lStr})
						self.lexer.backup()
						//utils.Dbg("New backup", string(self.lexer.char))
						break // 推出循环
					} else if self.lexer.char == ')' || self.lexer.char == ']' {
						//utils.Dbg("New text", lStr)
						list.items = append(list.items, &TStringList{text: lStr})
						// NOTE: 列表末尾不backup
						goto exit // 结束列表
					}
				}
			}

		case ",":
			{
				// 如果","的下一个字符是数字字母并“）]”结束 则取值
				// 否则 跳出循环
				self.lexer.consume_whitespace()
				lPos = self.lexer.pos + 1 // 减","
				for {
					self.lexer.consume_whitespace()
					self.lexer.next()
					if self.lexer.isEof {
						goto exit
					}

					// 如果是字母数字 向下
					if isAlphaNumeric(rune(self.lexer.char)) {
						continue
					} else {
						if self.lexer.char == ']' || self.lexer.char == ')' {
							lStr = self.lexer.stream[lPos:self.lexer.pos]
							list.items = append(list.items, &TStringList{text: lStr})
							//utils.Dbg("New text", lStr)
							goto exit

							// 以下代码不执行
							self.lexer.consume_whitespace()
							self.lexer.next()
							if self.lexer.isEof {
								goto exit
							}
							if self.lexer.char == ']' || self.lexer.char == ')' {
								goto exit
							}

						} else {
							// 回退1 并跳出For
							self.lexer.backup()
							break
						}
					}
				}

				//bracket = 1
			}
		default:
			{
				//bracket = 0
			}
		}
		//	utils.Dbg("end")
		self.lexer.consume_whitespace()
		self.lexer.next()
	}

exit:
	return nil
}

// 解析
func (self *TStringParser) parse() error {
	//utils.Dbg("domain", self.lexer.stream)

	// 检查合法性
	if !strings.HasPrefix(self.lexer.stream, "[") && !strings.HasSuffix(self.lexer.stream, "]") {
		//utils.Dbg("invaild domain")
		return errors.New("invaild domain")
	}

	return self._parselist(self.strings)
}

//---------------text
/*
func (self *TStringText) Item(idx int) IString {
	return self // 超范围返回自己
}

func (self *TStringText) Len() int {
	return 1
}
func (self *TStringText) Text() string {
	return self.text
}

func (self *TStringText) IsList() bool {
	return false
}

func (self *TStringText) String(idx int) string {
	return self.text
}
*/
//-----------list
func (self *TStringList) Item(idx int) *TStringList {
	self.Update()

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
			self.text = text
		}
	} else if len(self.items)-1 >= idx[0] {
		self.items[idx[0]] = &TStringList{text: text}
	} else {
		// 不是List才能直接修改TEXT 否者需要解析到Items 再生成Text
		if !self.IsList() {
			self.items[idx[0]].text = text
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
		self.text = ""
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
		self.text = ""
	}
	self.updated = false
	return lst
}

//栈方法Push：叠加元素
func (self *TStringList) Push(item ...*TStringList) *TStringList {
	// 当Self是一个值时必须添加自己到items里成为列表的一部分
	if self.text != "" && len(self.items) == 0 {
		self.items = append(self.items, NewStringList(self.text))
	}
	self.items = append(self.items, item...)

	self.updated = false
	return self
}

// #添加字符串到该 List
func (self *TStringList) PushString(strs ...string) *TStringList {
	// 当Self是一个值时必须添加自己到items里成为列表的一部分
	if self.text != "" && len(self.items) == 0 {
		self.items = append(self.items, NewStringList(self.text))
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

func (self *TStringList) Insert(idx int, item *TStringList) {
	// Grow the slice by one element.
	// make([]Token, len(self.Child)+1)
	// self.Child[0 : len(self.Child)+1]
	self.items = append(self.items, item)
	// Use copy to move the upper part of the slice out of the way and open a hole.

	copy(self.items[idx+1:], self.items[idx:])
	// Store the new value.
	self.items[idx] = item
	// Return the result.

	self.updated = false
	return
}

func (self *TStringList) InsertString(idx int, str string) {
	self.Insert(idx, &TStringList{text: str})
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
	if other == nil || other.Len() == 0 {
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
func (self *TStringList) Len() int {
	return len(self.items)
}

// retur the list length
func (self *TStringList) Count() int {
	return len(self.items)
}

func (self *TStringList) IsList() bool {
	return len(self.items) > 0
}

func (self *TStringList) In(strs ...interface{}) bool {
	for _, itr := range strs {
		// 处理字符串
		if str, ok := itr.(string); ok {
			if self.text == str {
				return true
			}
		} else
		// 处理*TStringList 类型
		if item, ok := itr.(*TStringList); ok {
			if self.text == item.text {
				return true
			}
		}

	}

	return false
}

//检查是否包含 一个或者多个内容
func (self *TStringList) Has(strs ...interface{}) bool {
	for _, itr := range strs {
		// 处理字符串
		if str, ok := itr.(string); ok {
			if strings.Index(self.text, str) == -1 {
				return false // 如果一个查不到就返回否
			}
		} else
		// 处理*TStringList 类型
		if item, ok := itr.(*TStringList); ok {
			if strings.Index(self.text, item.text) == -1 {
				return false
			}
		}

	}

	return true
}

// ('xx','xx')
func (self *TStringList) AsString(idx ...int) string {
	if self.IsList() {
		return self.String(idx...)
	}

	return "'" + self.String(idx...) + "'"
}

// return 'xx','xx'
func (self *TStringList) String(idx ...int) string {
	self.Update()

	if len(idx) > 0 {
		if self.IsList() {
			if idx[0] > -1 && idx[0] < len(self.items) {
				return self.items[idx[0]].text
			} else {
				return "" // 超范围返回空
			}
		}
	} else {
		return self.text
	}
	return ""
}

// 返回所有Items字符
func (self *TStringList) Strings(idx ...int) (result []string) {
	self.Update()

	lCnt := len(idx)

	if lCnt == 0 { // 返回所有
		for _, item := range self.items {
			result = append(result, item.text)
		}
	} else if lCnt == 1 {
		result = append(result, self.items[idx[0]].text)

	} else if lCnt > 1 {

		for _, item := range self.items[idx[0]:idx[1]] {
			result = append(result, item.Text())
		}
	}

	return
}

// 复制一个反转版
func (self *TStringList) Reversed() (result *TStringList) {
	self.Update()

	result = NewStringList()
	lCnt := self.Len()
	for i := lCnt - 1; i >= 0; i-- {
		result.Push(self.items[i]) //TODO: 复制
	}
	return
}

func (self *TStringList) Clear() {
	self.startPos = 0
	self.endPos = 0
	self.updated = true
	self.text = ""
	self.items = nil // make([]*TStringList, 0)
}

//在原有基础上克隆
func (self *TStringList) Clone(idx ...int) (result *TStringList) {
	self.Update()

	lCnt := len(idx)
	if lCnt == 0 {
		result = NewStringList()
		result.Push(self.Items(0, self.Len()-1)...)
	} else if lCnt == 1 && idx[0] < self.Len() { // idex 必须小于Self长度
		result = NewStringList()
		result.Push(self.Items(idx[0], self.Len()-1)...) //result.Push(self.items[idx[0]])
	} else if lCnt > 1 && idx[0] < self.Len() && idx[1] < self.Len() {
		result = NewStringList()
		if idx[1] == -1 {
			// 复制到end
			result.Push(self.items[idx[0]:self.Len()]...)
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
	self.Update()

	lCnt := len(idx)

	if !self.IsList() {
		// 如果是空白的对象 返回Nil
		if self.text == "" {
			return
		}

		return []*TStringList{self}
	}
	if lCnt == 0 {
		return self.items // 返回所有
	} else if lCnt == 1 && idx[0] < self.Len() { // idex 必须小于Self长度
		result = append(result, self.items[idx[0]])
	} else if lCnt > 1 && idx[0] < self.Len() && idx[1] < self.Len() {
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
	self.Update()

	// # 当StringList作为一个单字符串
	if self.text != "" && len(self.items) == 0 {
		result = []string{self.text}
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
func (self *TStringList) Text(idx ...int) string {
	self.Update()

	if len(idx) > 0 {
		return self.Item(idx[0]).text
	}
	return self.text
}

//废弃
func (self *TStringList) Type() int {
	if len(self.items) == 0 {
		return 0
	}
	return 0
}

// 添加''
func (self *TStringList) _quote(str string) string {
	return self.quoteChar + str + self.quoteChar
}

// 更新生成所有Text内容
func (self *TStringList) _update(node *TStringList) {
	//STEP  如果是Value Object 不处理
	if len(node.items) == 0 {
		return
	}

	// 处理有Child的Object
	lStr := ""
	lStrLst := make([]string, 0)
	IsList := false

	for _, item := range node.items {
		if is_leaf(item) {
			lStr = `(` + self._quote(item.String(0)) + `, ` + self._quote(item.String(1)) + `, `
			if item.Item(2).IsList() {
				lStr = lStr + `[` + item.String(2) + `])`
			} else {
				lStr = lStr + self._quote(item.String(2)) + `)`
			}

			item.text = lStr
			//utils.Dbg("_update leaf", lStr)
			IsList = true
			lStrLst = append(lStrLst, item.text)
		} else if item.IsList() {
			//utils.Dbg("IsList", item.text)
			self._update(item)
			lStrLst = append(lStrLst, item.text)
		} else {
			//utils.Dbg("_update val", item.text)
			lStrLst = append(lStrLst, self._quote(item.text))
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
	lStr = strings.Join(lStrLst, ",")

	//if lStr == "" {
	//	lStr = self._quote(node.text) // 如果是Val Node
	//}

	// 组合[XX,XX]
	if self == node && IsList {
		lStr = `[` + lStr + `]`
	}

	node.text = lStr
	//utils.Dbg("_update lst", lStr)
}

// 更新生成所有Text内容
func (self *TStringList) Update() {
	if !self.updated {
		self._update(self)
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

func (self *TStringList) AsBool() bool {
	if b, err := strconv.ParseBool(self.text); err == nil {
		return b //	fmt.Printf("%T, %v\n", s, s)
	}
	return false
}

func (self *TStringList) AsInt() int64 {
	if b, err := strconv.ParseInt(self.text, 10, 0); err == nil {
		return b //	fmt.Printf("%T, %v\n", s, s)
	}
	return -1
}

func (self *TStringList) IsBool() bool {
	if _, err := strconv.ParseBool(self.text); err == nil {
		return true //	fmt.Printf("%T, %v\n", s, s)
	}
	return false
}

func (self *TStringList) IsInt() bool {
	if _, err := strconv.ParseInt(self.text, 10, 0); err == nil {
		return true //	fmt.Printf("%T, %v\n", s, s)
	}
	return false
}

func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

/*""" Test whether an object is a valid domain term:
    - is a list or tuple
    - with 3 elements
    - second element if a valid op

    :param tuple element: a leaf in form (left, operator, right)
    :param boolean internal: allow or not the 'inselect' internal operator
        in the term. This should be always left to False.

    Note: OLD TODO change the share wizard to use this function.
"""*/
func is_leaf(element *TStringList, internal ...bool) bool {
	INTERNAL_OPS := append(TERM_OPERATORS, "<>")
	if internal != nil && internal[0] {
		INTERNAL_OPS = append(INTERNAL_OPS, "inselect")
		INTERNAL_OPS = append(INTERNAL_OPS, "not inselect")
	}

	//??? 出现过==Nil还是继续执行之下的代码
	return (element != nil && element.IsList()) &&
		(element.Len() == 3) &&
		InStrings(element.String(1), INTERNAL_OPS...) ||
		InStrings(element.String(), TRUE_LEAF, FALSE_LEAF)

	/*
	   def is_leaf(element, internal=False):

	       INTERNAL_OPS = TERM_OPERATORS + ('<>',)
	       if internal:
	           INTERNAL_OPS += ('inselect', 'not inselect')
	       return (isinstance(element, tuple) or isinstance(element, list)) \
	           and len(element) == 3 \
	           and element[1] in INTERNAL_OPS \
	           and ((isinstance(element[0], basestring) and element[0]) or tuple(element) in (TRUE_LEAF, FALSE_LEAF))
	*/
}
