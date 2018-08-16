package sJson

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Env struct {
	Index string
}

var e Env
var BaseStructureName = "Base"

func getMaxStruct(in []map[string]interface{}) map[string]interface{} {
	var multi = map[string]interface{}{}
	var out []string
	var spe []interface{}
	for i := 0; i < len(in); i++ {
		var flag int
		for j, _ := range in[i] {
			flag = 0
			if j == e.Index && e.Index != "" {
				spe = append(spe, in[i][j])
			}
			for k := 0; k < len(out); k++ {
				if j == out[k] {
					flag = 1
					break
				}
			}
			if flag == 0 {
				out = append(out, j)
				multi[j] = in[i][j]
			}
		}
	}
	if e.Index != "" && len(spe) != 0 {
		// 最も上位の []struct 内直下のデータなら拾ってこれる。
		fmt.Println(spe, len(spe))
	}
	return multi
}

func typeCheck(input interface{}) string {
	if input == nil {
		return "interface{}"
	}
	return reflect.TypeOf(input).String()
}

func jsonOrNot(input string) bool {
	if input == "map[string]interface {}" {
		return true
	}
	return false
}

func arrayOrNot(input string) bool {
	if input == "[]interface {}" {
		return true
	}
	return false
}

func returnType(input interface{}, count *int, Map *[]map[string]interface{}, n *[]string, inner int) string {
	var temp = typeCheck(input)
	var flag int
	var ff = 0
	if temp == "float64" {
		return checkFloat64(input)
	} else if jsonOrNot(temp) {
		*count = *count + 1
		tempMap, _ := input.(map[string]interface{})
		if inner == 0 {
			*Map = append(*Map, tempMap)
		}
		for i, _ := range tempMap {
			flag = 0
			for j := 0; j < len(*n); j++ {
				if i == (*n)[j] {
					// 同じ名前の要素に関しては一番初めに出てきたものでしか評価していないため、後で下位にもっと情報を持っている構造が出てきても反映されない。
					flag = 1
					break
				}
			}
			if flag == 0 {
				*n = append(*n, i)
				if inner == 1 {
					if ff == 0 {
						*Map = append(*Map, map[string]interface{}{})
						ff = 1
					}
					(*Map)[len(*Map)-1][i] = tempMap[i]
				}
			}
		}
		return ""
	} else if arrayOrNot(temp) {
		return forArray(input, count, Map, n, 1)
	}
	return temp
}

func forArray(input interface{}, count *int, Map *[]map[string]interface{}, n *[]string, inner int) string {
	k, _ := input.([]interface{})
	var te []map[string]interface{}
	var tempMap map[string]interface{}
	if len(k) == 0 {
		return "[]interface{}"
	}
	var temp = typeCheck(k[0])
	if jsonOrNot(temp) {
		for l := 0; l < len(k); l++ {
			tempMap, _ = k[l].(map[string]interface{})
			te = append(te, tempMap)
		}
		tempMap = getMaxStruct(te)
		return "[]" + returnType(tempMap, count, Map, n, 1)
	} else if arrayOrNot(temp) {
		kk, _ := k[0].([]interface{})
		if len(kk) > 0 {
			if jsonOrNot(typeCheck(kk[0])) {
				for i := 0; i < len(kk); i++ {
					tempMap, _ = kk[i].(map[string]interface{})
					te = append(te, tempMap)
				}
				tempMap = getMaxStruct(te)
				return "[][]" + returnType(tempMap, count, Map, n, 1)
			}
		}
		return "[]" + forArray(k[0], count, Map, n, 1)
	}
	if temp == "float64" {
		return "[]" + checkFloat64(k[0])
	}
	return "[]" + temp
}

func checkFloat64(input interface{}) string {
	s, _ := input.(float64)
	index := strings.Index(strconv.FormatFloat(s, 'f', -1, 64), ".")
	if index == -1 {
		return "int"
	}
	return "float64"
}

func convert(input string) string {
	if len(input) > 0 {
		if string(input[0]) == "_" {
			return "UuU" + input[1:]
		}
		return firstUpper(input)
	}
	return ""
	//	firstUpper(replace.Replace(input, "_", "UuU"))
}

func firstUpper(input string) string {
	if len(input) > 0 {
		return strings.ToUpper(input[0:1]) + input[1:]
	}
	return ""
}

func ParseCont(body string, s ...string) {
	if len(s) == 1 {
		e.Index = s[0]
	} else if len(os.Args) == 2 {
		e.Index = os.Args[1]
	} else {
		e.Index = ""
	}
	parseJson(string(body))
}

func Parse(filename string, s ...string) {
	file, _ := os.Open(filename)
	defer file.Close()
	body, _ := ioutil.ReadAll(file)
	if len(s) == 1 {
		e.Index = s[0]
	} else if len(os.Args) == 2 {
		e.Index = os.Args[1]
	} else {
		e.Index = ""
	}
	parseJson(string(body))
}

func parseJson(input string) {
	var d map[string]interface{}
	var err error
	err = json.Unmarshal([]byte(input), &d)
	if err != nil {
		d = map[string]interface{}{}
		json.Unmarshal([]byte(`{"Data":`+input+"}"), &d)
	}
	var out [][]string
	var outb = [][]string{{BaseStructureName}}
	var cont [][]interface{}
	var count int
	var dd []map[string]interface{}
	dd = append(dd, d)
	var layer = -1
	for {
		layer = layer + 1
		count = 0
		out = [][]string{}
		cont = [][]interface{}{}
		for i := 0; i < len(dd); i++ {
			out = append(out, []string{})
			cont = append(cont, []interface{}{})
			for j, _ := range dd[i] {
				out[i] = append(out[i], j)
				cont[i] = append(cont[i], dd[i][j])
			}
		}
		printType(out, &outb, cont, &count, &dd)
		if len(outb[layer+1]) == 0 {
			break
		}
	}
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func trim(input string, length int) string {
	var L = len(input)
	for i := 0; i < length-L; i++ {
		input = input + " "
	}
	return input
}

func printType(out [][]string, outb *[][]string, cont [][]interface{}, count *int, dd *[]map[string]interface{}) {
	var cb int
	newb := []string{}
	var temp string
	var newbb []string
	var nowLayer = len(*outb)
	var printCont []string
	var printType []string
	var contMax int
	var typeMax int
	*dd = []map[string]interface{}{}
	if e.Index == "" {
		fmt.Println("")
	}
	for i := 0; i < len(out); i++ {
		printCont = []string{}
		printType = []string{}
		contMax = 0
		typeMax = 0
		for j := 0; j < len(out[i]); j++ {
			cb = *count
			temp = returnType(cont[i][j], count, dd, &[]string{}, 0)
			printCont = append(printCont, convert(out[i][j]))
			contMax = max(contMax, len(convert(out[i][j])))
			if cb != *count {
				printType = append(printType, temp+convert(out[i][j]))
				typeMax = max(typeMax, len(temp+convert(out[i][j])))
				newb = append(newb, out[i][j])
				newbb = append(newbb, (*outb)[nowLayer-1][i]+"."+out[i][j]+temp)
			} else {
				printType = append(printType, temp)
				typeMax = max(typeMax, len(temp))
			}
		}
		if e.Index == "" {
			fmt.Println("type " + convert((*outb)[nowLayer-1][i]) + " struct {")
			for j := 0; j < len(printCont); j++ {
				fmt.Println("  " + trim(printCont[j], contMax) + " " + trim(printType[j], typeMax) + " `json:\"" + out[i][j] + "\"`")
			}
			fmt.Println("}")
			fmt.Println("")
		}
	}
	*outb = append(*outb, newb)
	// ひとつ前の outb の中に過去の json 構造の情報が入っている。
	// 後はここで保存した形に応じて最初の map[string]interface{} から []map[string]interface{} を作り
	// getMaxStruct を用いて map[string]interface{} に落とし込めば、全てを包含する構造体を定義できる。
	(*outb)[nowLayer-1] = newbb
	var series string
	var t []string
	var tt []string
	if nowLayer > 1 {
		for i := 0; i < len((*outb)[nowLayer-1]); i++ {
			t = strings.Split((*outb)[nowLayer-1][i], ".")
			for j := 0; j < len((*outb)[nowLayer-2]); j++ {
				tt = strings.Split((*outb)[nowLayer-2][j], ".")
				if t[0] == tt[len(tt)-1] || t[0]+"[]" == tt[len(tt)-1] {
					series = ""
					for k := 1; k < len(t); k++ {
						series = series + "." + t[k]
					}
					(*outb)[nowLayer-1][i] = strings.TrimLeft((*outb)[nowLayer-2][j], ".") + series
					break
				}
			}
		}
	}
	/*else {
		for i := 0; i < len((*outb)[nowLayer-1]); i++ {
			(*outb)[nowLayer-1][i] = ""
		}
	}*/
	//	fmt.Println(*outb)
}
