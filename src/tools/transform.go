package tools

import (
	"github.com/360EntSecGroup-Skylar/excelize"
	"io/ioutil"
	"strings"
)

//把文件内容按行列处理成表
func TextToTable(fileName string) (out [][]string, err error) {
	d, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	rows := strings.Split(string(d), "\n")
	for _, row := range rows {
		if row == "" {
			continue
		}

		var line []string
		fields := strings.Split(row, "\t")
		for _, field := range fields {
			line = append(line, field)
		}
		out = append(out, line)
	}
	return
}

//把文件内容按行列处理成表
//以第一行为map的键
func ReadTextToTableMap(fileName string) (out []map[string]string, err error) {
	d, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	rows := strings.Split(string(d), "\n")
	keys := make(map[int]string)
	for i, row := range rows {
		if i == 0 {
			if row == "" {
				return
			}
			fields := strings.Split(row, "\t")
			for j, field := range fields {
				//如果有中横线转为下划线
				//如果有空格转为下划线
				//大写转小写
				field = strings.Replace(field, "-", "_", -1)
				field = strings.Replace(field, " ", "_", -1)
				field = strings.ToLower(field)
				keys[j] = field
			}
			continue
		}

		if row == "" {
			continue
		}

		line := make(map[string]string)
		fields := strings.Split(row, "\t")
		for j, field := range fields {
			if key, ok := keys[j]; ok {
				line[key] = field
			}
		}
		out = append(out, line)
	}
	return
}

//忽略Bom头
func IgnoreBom(fileName string) []*string {
	dat, err := ioutil.ReadFile(fileName)
	if err != nil {

	}
	if dat[0] == 0xef || dat[1] == 0xbb || dat[2] == 0xbf {
		dat = dat[3:]
	}
	var cleaned = strings.Replace(string(dat), "\r", "", -1)
	var lines = strings.Split(cleaned, "\n")
	n := len(lines)
	var r []*string
	for i := 0; i < n; i++ {
		if lines[i] != "" {
			r = append(r, &lines[i])
		}
	}
	return r
}

//读取xls把内容转成table
func ReadXlsToTableMap(fileName string) (out []map[string]string, err error) {
	var f *excelize.File
	f, err = excelize.OpenFile(fileName)
	if err != nil {
		return nil, err
	}

	//兼容名称被修改
	index := f.GetActiveSheetIndex()
	sheetMap := f.GetSheetMap()
	var rows [][]string
	keys := make(map[int]string)
	rows = f.GetRows(sheetMap[index])
	for i, row := range rows {
		if i == 0 { //标题
			for j, colCell := range row {
				keys[j] = colCell
			}
			continue
		}

		line := make(map[string]string)
		for j, colCell := range row {
			if key, ok := keys[j]; ok {
				line[key] = colCell
			}
		}
		out = append(out, line)
	}
	return out, nil
}
