# 同义词词林扩展版

[![Go Reference](https://pkg.go.dev/badge/github.com/bububa/cilin.svg)](https://pkg.go.dev/github.com/bububa/cilin)
[![Go](https://github.com/bububa/cilin/actions/workflows/go.yml/badge.svg)](https://github.com/bububa/cilin/actions/workflows/go.yml)
[![goreleaser](https://github.com/bububa/cilin/actions/workflows/goreleaser.yml/badge.svg)](https://github.com/bububa/cilin/actions/workflows/goreleaser.yml)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/bububa/cilin.svg)](https://github.com/bububa/cilin)
[![GoReportCard](https://goreportcard.com/badge/github.com/bububa/cilin)](https://goreportcard.com/report/github.com/bububa/cilin)
[![GitHub license](https://img.shields.io/github/license/bububa/cilin.svg)](https://github.com/bububa/cilin/blob/master/LICENSE)
[![GitHub release](https://img.shields.io/github/release/bububa/cilin.svg)](https://GitHub.com/bububa/cilin/releases/)


Word similarity computation based on Tongyici Cilin.
这是一个基于哈工大同义词词林扩展版的单词相似度计算方法的golang实现，参考论文如下：
2010 田久乐等，吉林大学学报（信息科学版），基于同义词词林的词语相似度计算方法。

## Install
go get -u github.com/bububa/cilin

## Usage 
```golang
import "github.com/bububa/cilin"

cs := cilin.NewSimilarity()
w1 := "抄袭"
w2 := "克隆"
sim := cs.Calculate(w1, w2)
fmt.Printf("%s %s 相似度为 %f\n", w1, w2, sim)
// 抄袭 克隆 相似度为 0.585642777645155

w1 = "人民"
lst := []string{"国民", "群众", "党群", "良民", "同志", "成年人", "市民", "亲属", "志愿者", "先锋"}
for _, w2 := range lst {
    sim := s.Calculate(w1, w2)
    fmt.Printf("%s %s, 相似度:%f\n", w1, w2, sim)
}

// 人民 国民, 相似度:1.000000
// 人民 群众, 相似度:0.957661
// 人民 党群, 相似度:0.897808
// 人民 良民, 相似度:0.718246
// 人民 同志, 相似度:0.663015
// 人民 成年人, 相似度:0.630692
// 人民 市民, 相似度:0.540593
// 人民 亲属, 相似度:0.360396
// 人民 志愿者, 相似度:0.225247
// 人民 先锋, 相似度:0.180198
```

## 同类项目

- 
- https://github.com/ashengtx/CilinSimilarity  实现了三种计算方法。
- https://github.com/Xls1994/Cilin
- http://www.codepub.cn/2015/08/04/Based-on-the-extended-version-of-synonyms-Cilin-word-similarity-computing/  Java实现

## 致谢

本代码的实现要感谢下面几位作者：
* 哈工大信息检索研究室所著的《哈工大信息检索研究室同义词词林扩展版》
* 田久乐  赵 蔚在2010年所发表论文<基于同义词词林的词语相似度计算方法>
* http://codepub.cn/2015/08/04/Based-on-the-extended-version-of-synonyms-Cilin-word-similarity-computing/
* http://www.cnblogs.com/einyboy/archive/2012/09/09/2677265.html?from=singlemessage&isappinstalled=0
* 本代码参考了https://github.com/ashengtx/CilinSimilarity 部分实现
