package spotify

import (
	"errors"
	"github.com/longbridgeapp/opencc"
	"strings"
	"unicode"
)

var t2s *opencc.OpenCC

// InitT2S 初始化简繁转化器 t2s。
// 默认懒汉调用，但暴露出来提供饿汉支持。
func InitT2S() error {
	var err error
	t2s, err = opencc.New("t2s")
	if err != nil {
		return err
	}
	return nil
}

// T2S converts Traditional Chinese to Simplified Chinese
// 繁 => 简
func T2S(t string) (string, error) {
	var err error
	if t2s == nil {
		err = InitT2S()
	}
	if err != nil || t2s == nil {
		err = errors.New("init t2s failed:" + err.Error())
		//fmt.Println("[Debug] T2S Error:", err)
		return t, err
	}
	s, err := t2s.Convert(t)
	if err != nil {
		//fmt.Println("[Debug] T2S Error:", err)
		return t, err
	}
	//fmt.Printf("[Debug] T2S: \n\tin=%v\n\tout=%v\n", t, s)
	return s, nil
}

// ChineseStringEqual checks if string a == b **in Chinese**.
// It converts traditional Chinese text into simplified Chinese if necessary.
// lower/upper and Punct chars will be ignored.
// 繁简混合的情况下 a 繁、b 简 比较快。
// For example:
//      ChineseStringEqual("自然語言處理", "自然语言处理") == true
// See Also: cc_benchmark_test.go
func ChineseStringEqual(a, b string) bool {
	a, b = strings.ToLower(a), strings.ToLower(b)
	a, b = RemovePunct(a), RemovePunct(b)
	if len(a) != len(b) {
		return false
	}
	if a == b {
		return true
	}
	// opencc: t2s
	sa, _ := T2S(a)
	if sa == b {
		return true
	}
	sb, _ := T2S(b)
	if sa == sb {
		return true
	}
	return false
}

func RemovePunct(s string) string {
	r := []rune(s)
	var rr []rune
	for i := 0; i < len(r); i++ {
		if !unicode.IsPunct(r[i]) {
			rr = append(rr, r[i])
		}
	}
	return string(rr)
}
