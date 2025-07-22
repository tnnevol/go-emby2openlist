package jsons

import (
	"encoding/json"
	"strings"
)

// Struct 将 item 转换为结构体对象
func (i *Item) Struct() any {
	switch i.jType {
	case JsonTypeVal:
		return i.val
	case JsonTypeObj:
		m := make(map[string]any)
		for key, value := range i.obj {
			m[key] = value.Struct()
		}
		return m
	case JsonTypeArr:
		a := make([]any, len(i.arr))
		for idx, value := range i.arr {
			a[idx] = value.Struct()
		}
		return a
	default:
		return "null"
	}
}

// String 将 item 转换为 json 字符串
func (i *Item) String() string {
	var buf strings.Builder
	writeJSON(&buf, i)
	return buf.String()
}

func writeJSON(buf *strings.Builder, i *Item) {
	switch i.jType {
	case JsonTypeVal:
		writeValueJSON(buf, i.val)
	case JsonTypeObj:
		buf.WriteByte('{')
		first := true
		for k, v := range i.obj {
			if !first {
				buf.WriteByte(',')
			}
			first = false

			// Key 编码（禁用 HTML 转义）
			var keyBuf strings.Builder
			enc := json.NewEncoder(&keyBuf)
			enc.SetEscapeHTML(false)
			_ = enc.Encode(k)
			key := strings.TrimSpace(keyBuf.String())
			buf.WriteString(key)

			buf.WriteByte(':')
			writeJSON(buf, v)
		}
		buf.WriteByte('}')
	case JsonTypeArr:
		buf.WriteByte('[')
		for idx, elem := range i.arr {
			if idx > 0 {
				buf.WriteByte(',')
			}
			writeJSON(buf, elem)
		}
		buf.WriteByte(']')
	default:
		buf.WriteString("null")
	}
}

func writeValueJSON(buf *strings.Builder, val any) {
	var tmp strings.Builder
	enc := json.NewEncoder(&tmp)
	enc.SetEscapeHTML(false)
	err := enc.Encode(val)
	if err != nil {
		buf.WriteString("null")
		return
	}
	s := strings.TrimSpace(tmp.String())
	buf.WriteString(s)
}