package cilin

import (
	"strings"
	"testing"
)

var s = NewSimilarity()

func TestSimilarity(t *testing.T) {
	w1 := "抄袭"
	w2 := "克隆"
	sim := s.Calculate(w1, w2)
	if sim < 0.585643 && sim > 0.585644 {
		t.Errorf("%s %s 返回值:%f, 期待值:0.585643\n", w1, w2, sim)
	}
	t.Logf("%s %s, 相似度:%f\n", w1, w2, sim)
	w1 = "人民"
	lst := []string{"国民", "群众", "党群", "良民", "同志", "成年人", "市民", "亲属", "志愿者", "先锋"}
	for _, w2 := range lst {
		sim := s.Calculate(w1, w2)
		t.Logf("%s %s, 相似度:%f\n", w1, w2, sim)
	}
}

func TestGetCodes(t *testing.T) {
	w := "抄袭"
	expects := "Hb08B04=,Hn10C01="
	codes := s.GetCodes(w)
	result := strings.Join(codes, ",")
	if result != expects {
		t.Errorf("%s 编码错误，返回值:%s, 期待值:%+v\n", w, result, expects)
	}
}
