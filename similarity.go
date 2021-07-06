package cilin

import (
	"bufio"
	"errors"
	"io"
	"math"
	"os"
	"strconv"
	"strings"

	// 使用go 1.16以上版本的embed功能
	_ "embed"
)

//go:embed cilin.txt
var cilinData string

var (
	DefaultLayerTable = map[int]int{
		1: 1,
		2: 2,
		4: 3,
		5: 4,
		7: 5,
	}
	DEGREE float64 = 180
)

// Layer 编码按层次结构化
type Layer [6]string

// NewLayer 将编码按层次结构化
// Aa01A01=
// 第三层和第五层是两个数字表示；
// 第一、二、四层分别是一个字母；
// 最后一个字符用来区分所有字符相同的情况。
func NewLayer(c string) Layer {
	return Layer{
		string(c[0]), string(c[1]), string(c[2:4]), string(c[4]), string(c[5:7]), string(c[7]),
	}
}

// Param 为Similarity相似性计算参数
type Param struct {
	a float64
	b float64
	c float64
	d float64
	e float64
	f float64
}

// NewDefaultParam 默认计算参数
func NewDefaultParam() *Param {
	return &Param{
		a: 0.65,
		b: 0.8,
		c: 0.9,
		d: 0.96,
		e: 0.5,
		f: 0.1,
	}
}

// Similarity 结构体
type Similarity struct {
	codeWords map[string][]string // 以编码为key，单词list为value的dict，一个编码有多个单词
	wordCodes map[string][]string // 以单词为key，编码为value的dict，一个单词可能有多个编码
	param     *Param
}

// NewSimilarity 新建Simility
func NewSimilarity() *Similarity {
	s := &Similarity{
		codeWords: make(map[string][]string),
		wordCodes: make(map[string][]string),
		param:     NewDefaultParam(),
	}
	s.ParseCilin()
	return s
}

// Load 读入同义词词林，编码为key，词群为value，保存在codeWord
// 单词为key，编码为value，保存在wordCode
func (s *Similarity) Load(filename string) error {
	fd, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fd.Close()
	buf := bufio.NewReader(fd)
	for {
		line, err := buf.ReadString('\n')
		if err != nil && errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return err
		}
		line = strings.TrimSpace(line)
		parts := strings.Split(line, " ")
		code := parts[0]
		words := make([]string, len(parts)-1)
		copy(words, parts[1:])
		s.codeWords[code] = words
		for _, w := range words {
			s.wordCodes[w] = append(s.wordCodes[w], code)
		}
	}
	return nil
}

func (s *Similarity) ParseCilin() {
	lines := strings.Split(cilinData, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		parts := strings.Split(line, " ")
		code := parts[0]
		words := make([]string, len(parts)-1)
		copy(words, parts[1:])
		s.codeWords[code] = words
		for _, w := range words {
			s.wordCodes[w] = append(s.wordCodes[w], code)
		}
	}
}

// commonCode 获取两个编码字符的公共部分
func commonCode(c1 string, c2 string) string {
	length := len(c1)
	if length > len(c2) {
		length = len(c2)
	}
	var idx int
	for idx < length {
		if c1[idx] != c2[idx] {
			break
		}
		idx += 1
	}
	if idx == 3 || idx == 6 {
		return c1[0 : idx-1]
	}
	return c1[0:idx]
}

// getLayer 根据commonCodeStr返回两个编码所在的层数。
// 如果没有共同的str，则位于第一层，用0表示；
// 如果第1个字符相同，则位于第二层，用1表示；
// 这里第一层用0表示。
func getLayer(commonCodeStr string) int {
	length := len(commonCodeStr)
	if length == 0 {
		return 0
	}
	if layer, found := DefaultLayerTable[length]; found {
		return layer
	}
	return 0
}

// getK 返回两个编码对应分支的距离，相邻距离为1
func getK(l1 Layer, l2 Layer) int {
	var (
		idx    int
		layers int = 5
	)
	for idx < layers {
		if (idx == 2 || idx == 4) && l1[idx] != l2[idx] {
			n1, _ := strconv.Atoi(strings.TrimPrefix(l1[idx], "0"))
			n2, _ := strconv.Atoi(strings.TrimPrefix(l2[idx], "0"))
			return abs(n1 - n2)
		} else if l1[idx] != l2[idx] {
			return abs(int(l1[idx][0]) - int(l2[idx][0]))
		}
		idx += 1
	}
	return 0
}

// getN 计算所在分支层的分支数
// 即计算分支的父节点总共有多少个子节点
// 两个编码的common_str决定了它们共同处于哪一层
// 例如，它们的common_str为前两层，则它们共同处于第三层，则我们统计前两层为common_str的第三层编码个数就好了
func (s *Similarity) getN(commonCodeStr string) int {
	if len(commonCodeStr) == 0 {
		return 0
	}
	layerIdx := getLayer(commonCodeStr)
	mp := make(map[string]struct{})
	var total int
	for code := range s.codeWords {
		if strings.HasPrefix(code, commonCodeStr) {
			layer := NewLayer(code)
			layerValue := layer[layerIdx]
			if _, found := mp[layerValue]; found {
				continue
			}
			total += 1
			mp[layer[layerIdx]] = struct{}{}
		}
	}
	return total
}

// GetCode 获取字符串编码
func (s *Similarity) GetCodes(w string) []string {
	codes, _ := s.wordCodes[w]
	return codes
}

// Calculate 计算相似度
// 根据以下论文提出的改进方法计算：
// 《基于知网与词林的词语语义相似度计算》，朱新华，马润聪， 孙柳，陈宏朝（ 广西师范大学 计算机科学与信息工程学院，广西 桂林541004）
func (s *Similarity) Calculate(w1 string, w2 string) float64 {
	var sim float64
	// 如果有一个词不在词林中，则相似度为0
	c1s := s.GetCodes(w1)
	c2s := s.GetCodes(w2)
	if len(c1s) == 0 || len(c2s) == 0 {
		return 0
	}
	for _, c1 := range c1s {
		for _, c2 := range c2s {
			if curSim := s.sim(c1, c2); curSim > sim {
				sim = curSim
			}
		}
	}
	return sim
}

// sim 根据编码计算相似度
func (s *Similarity) sim(c1 string, c2 string) float64 {
	// 先把code的层级信息提取出来
	c1layer := NewLayer(c1)
	c2layer := NewLayer(c2)
	commonStr := commonCode(c1, c2)
	commonLen := len(commonStr)

	// 如果有一个编码以'@'结尾，那么表示自我封闭，这个编码中只有一个词，直接返回f
	if commonLen == 0 || strings.HasSuffix(c1, "@") || strings.HasSuffix(c2, "@") {
		return s.param.f
	}

	var sim float64
	if commonLen >= 7 {
		// 如果前面七个字符相同，则第八个字符也相同，要么同为'='，要么同为'#''
		if strings.HasSuffix(c1, "=") && strings.HasSuffix(c2, "=") {
			sim = 1
		} else if strings.HasSuffix(c1, "#") && strings.HasSuffix(c2, "#") {
			sim = s.param.e
		}
	} else {
		k := float64(getK(c1layer, c2layer))
		n := float64(s.getN(commonStr))
		switch commonLen {
		case 1:
			sim = simFormula(s.param.a, n, k)
		case 2:
			sim = simFormula(s.param.b, n, k)
		case 4:
			sim = simFormula(s.param.c, n, k)
		case 5:
			sim = simFormula(s.param.d, n, k)
		}
	}
	return sim
}

// simFormula 计算相似度的公式，不同的层系数不同
func simFormula(coef float64, n float64, k float64) float64 {
	return coef * math.Cos(n*math.Pi/DEGREE) * (n - k + 1) / n
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
