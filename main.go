package main

import (
	"github.com/neo4l/x/tool"
	"log"
	"fmt"
	"os"
	"bufio"
	"bytes"
	"io"
	"strings"
	"strconv"
	"encoding/xml"
)

var (
	AndroidFormat      = "<string name=\"%s\">%s</string>"
	AndroidFormatFalse = "<string name=\"%s\" formatted=\"false\">%s</string>"
	IosFormat          = "\"%s\" = \"%s\";"
)

func main() {
	//convertToAndroidAndIOS()
	// "^((?!@).)*$"
	// %[^@]|
	//text := `确认中%@%`
	//reg := regexp.MustCompile(`%[^@]|%$`)
	//log.Println("match:", reg.MatchString(text))
	//fmt.Printf("%q\n", reg.FindAllString(text, -1))
	convertToAndroidAndIOS()
}

type AndroidXML struct {
	Name    string   `xml:"name,attr"`
	InnerText  string   `xml:",innerxml"`
}

func convertToCVS() {
	xmlLines, err := tool.ReadLines("./android_xml.txt")
	if err != nil {
		log.Printf("read android_xml error, %s", err)
		return
	}
	cvsLines := make([]string, 0)
	vincentCvsLines := getCVSKeyList()
	vincentEnCvsLines := getCVSEnList()
	for _, xmlLine := range xmlLines {
		androidXml := &AndroidXML{}
		xml.Unmarshal([]byte(xmlLine), androidXml)
		log.Println("name:", androidXml.Name, ";value:", androidXml.InnerText)
		value := getCVSStr(androidXml.InnerText)
		keys := strings.Split(androidXml.Name, ",")
		cvsLine := ""
		if value == "" {
			continue
		}
		for i := 0; i < len(keys); i++ {
			cvsLine = cvsLine + keys[i] + ","
		}
		cvsLine = cvsLine + "\"" + value + "\""
		isInVincent := false
		for i := 0; i < len(vincentCvsLines); i++ {
			log.Println("vincentCvsLines[i]:", vincentCvsLines[i], ";", androidXml.Name)
			if vincentCvsLines[i] == androidXml.Name {
				cvsLine = cvsLine + ",\"" + vincentEnCvsLines[i] + "\","
				isInVincent = true
				break
			}
		}
		if !isInVincent {
			cvsLine = cvsLine + ", ,新增"
		}
		cvsLines = append(cvsLines, cvsLine)
	}

	WriteToFile("./cvs.txt", cvsLines)
}

func convertToAndroidAndIOS() {
	pageLines, err := tool.ReadLines("./source/page.txt")
	elementLines, err := tool.ReadLines("./source/element.txt")
	elementKeyLines, err := tool.ReadLines("./source/element_key.txt")
	cnLines, err := tool.ReadLines("./source/cn.txt")
	enLines, err := tool.ReadLines("./source/en.txt")
	riBenLines, err := tool.ReadLines("./source/riben.txt")
	hanGuoLines, err := tool.ReadLines("./source/hanguo.txt")
	if err != nil {
		log.Printf("read file error, %s", err)
		return
	}
	androidChineseLines := make([]string, 0)
	iosChineseLines := make([]string, 0)
	androidEnglishLines := make([]string, 0)
	iosEnglishLines := make([]string, 0)
	androidRiBenLines := make([]string, 0)
	iosRiBenLines := make([]string, 0)
	androidHanGuoLines := make([]string, 0)
	iosHanGuoLines := make([]string, 0)
	for i := 1; i < len(pageLines); i++ {
		key := strings.TrimSpace(pageLines[i]) + "." + strings.TrimSpace(elementLines[i]) + "." + strings.TrimSpace(elementKeyLines[i])
		keyIos := strings.TrimSpace(pageLines[i]) + "_" + strings.TrimSpace(elementLines[i]) + "_" + strings.TrimSpace(elementKeyLines[i])
		chineseValue := strings.TrimSpace(cnLines[i])
		englishValue := strings.TrimSpace(enLines[i])
		riBenValue := strings.TrimSpace(riBenLines[i])
		hanGuoValue := strings.TrimSpace(hanGuoLines[i])
		log.Println("key:", key)
		//if tools.ContrainPercentSign(chineseValue) || tools.ContrainPercentSign(englishValue) {
		//	androidChineseLines = append(androidChineseLines, fmt.Sprintf(AndroidFormatFalse, key, getAndroidStr(chineseValue)))
		//	androidEnglishLines = append(androidEnglishLines, fmt.Sprintf(AndroidFormatFalse, key, getAndroidStr(englishValue)))
		//} else {
		//	androidChineseLines = append(androidChineseLines, fmt.Sprintf(AndroidFormat, key, getAndroidStr(chineseValue)))
		//	androidEnglishLines = append(androidEnglishLines, fmt.Sprintf(AndroidFormat, key, getAndroidStr(englishValue)))
		//}
		if getAndroidStr(chineseValue) != "" {
			androidChineseLines = append(androidChineseLines, fmt.Sprintf(AndroidFormat, key, getAndroidStr(chineseValue)))
		}
		if getAndroidStr(englishValue) != "" {
			androidEnglishLines = append(androidEnglishLines, fmt.Sprintf(AndroidFormat, key, getAndroidStr(englishValue)))
		}
		if getAndroidStr(riBenValue) != "" {
			androidRiBenLines = append(androidRiBenLines, fmt.Sprintf(AndroidFormat, key, getAndroidStr(riBenValue)))
		}
		if getAndroidStr(hanGuoValue) != "" {
			androidHanGuoLines = append(androidHanGuoLines, fmt.Sprintf(AndroidFormat, key, getAndroidStr(hanGuoValue)))
		}

		iosChineseLines = append(iosChineseLines, fmt.Sprintf(IosFormat, keyIos, getIosStr(chineseValue)))
		iosEnglishLines = append(iosEnglishLines, fmt.Sprintf(IosFormat, keyIos, getIosStr(englishValue)))
		iosRiBenLines = append(iosRiBenLines, fmt.Sprintf(IosFormat, keyIos, getIosStr(riBenValue)))
		iosHanGuoLines = append(iosHanGuoLines, fmt.Sprintf(IosFormat, keyIos, getIosStr(hanGuoValue)))
	}
	os.MkdirAll("./output/", os.ModePerm)
	WriteToFile("./output/android_cn.txt", androidChineseLines)
	WriteToFile("./output/android_en.txt", androidEnglishLines)
	WriteToFile("./output/android_riben.txt", androidRiBenLines)
	WriteToFile("./output/android_hanguo.txt", androidHanGuoLines)
	WriteToFile("./output/ios_cn.txt", iosChineseLines)
	WriteToFile("./output/ios_en.txt", iosEnglishLines)
	WriteToFile("./output/ios_riben.txt", iosRiBenLines)
	WriteToFile("./output/ios_hanguo.txt", iosHanGuoLines)
}

func getCVSEnList() []string {
	enLines, err := tool.ReadLines("./en.txt")
	if err != nil {
		log.Printf("read file error, %s", err)
		return nil
	}
	cvsEnLines := make([]string, 0)
	for i := 1; i < len(enLines); i++ {
		cvsEnLines = append(cvsEnLines, enLines[i])
	}
	return cvsEnLines
}

func getCVSKeyList() []string {
	pageLines, err := tool.ReadLines("./page.txt")
	elementLines, err := tool.ReadLines("./element.txt")
	elementKeyLines, err := tool.ReadLines("./element_key.txt")
	if err != nil {
		log.Printf("read file error, %s", err)
		return nil
	}
	cvsLines := make([]string, 0)
	for i := 1; i < len(pageLines); i++ {
		key := strings.TrimSpace(pageLines[i]) + "." + strings.TrimSpace(elementLines[i]) + "." + strings.TrimSpace(elementKeyLines[i])
		cvsLines = append(cvsLines, key)
	}
	return cvsLines
}

func getAndroidStr(value string) string {
	// 单引号需要转义
	if strings.Contains(value, "'") {
		value = strings.Replace(value, "'", "\\'", 100)
	}
	// 中文单引号需要转英文的
	if strings.Contains(value, "’") {
		value = strings.Replace(value, "’", "\\'", 100)
	}
	// 两个&&符号作为百分号
	if strings.Contains(value, "&&") {
		// %@是%1$s的含义，此处处理比较特殊，如果字符中有%1$s那么百分号的使用只需要%%，如果没有%1$s那么必须给%转义
		if strings.Contains(value, "%@") {
			value = strings.Replace(value, "&&", "%%", 100)
		} else {
			value = strings.Replace(value, "&&", "\\%%", 100)
		}
	}
	if strings.Contains(value, "&#160;") {
		// 先替换成一个特殊符号
		value = strings.Replace(value, "&#160;", "####%%", 100)
	}
	if strings.Contains(value, "&#8230;") {
		// 先替换成一个特殊符号
		value = strings.Replace(value, "&#8230;", "####-%%", 100)
	}
	if strings.Contains(value, "&") {
		value = strings.Replace(value, "&", "&#38;", 100)
	}
	if strings.Contains(value, "####%%") {
		// 再替换回来
		value = strings.Replace(value, "####%%", "&#160;", 100)
	}
	if strings.Contains(value, "####-%%") {
		// 再替换回来
		value = strings.Replace(value, "####-%%", "&#8230;", 100)
	}
	// 尖括号转义
	if strings.Contains(value, "<") {
		value = strings.Replace(value, "<", "&lt;", 100)
	}
	if strings.Contains(value, ">") {
		value = strings.Replace(value, ">", "&gt;", 100)
	}
	if !strings.Contains(value, "%@") {
		return value
	}

	for i := 1; strings.Index(value, "%@") >= 0; i++ {
		value = strings.Replace(value, "%@", "%"+strconv.Itoa(i)+"$s", 1)
	}
	return value
}

func getIosStr(value string) string {
	if strings.Contains(value, "&&") {
		if strings.Contains(value, "%@") {
			value = strings.Replace(value, "&&", "%", 100)
		} else {
			value = strings.Replace(value, "&&", "%", 100)
		}
	}
	if strings.Contains(value, "&#160;") {
		// 先替换成一个特殊符号
		value = strings.Replace(value, "&#160;", "####%%", 100)
	}
	if strings.Contains(value, "&#8230;") {
		// 先替换成一个特殊符号
		value = strings.Replace(value, "&#8230;", "####-%%", 100)
	}
	if strings.Contains(value, "####%%") {
		// 再替换回来
		value = strings.Replace(value, "####%%", " ", 100)
	}
	if strings.Contains(value, "####-%%") {
		// 再替换回来
		value = strings.Replace(value, "####-%%", "...", 100)
	}
	return value
}

func getCVSStr(androidXML string) string {
	if !strings.Contains(androidXML, "%1$s") {
		return androidXML
	}
	for i := 1; strings.Index(androidXML, "%" + strconv.Itoa(i) + "$s") > 0; i++ {
		androidXML = strings.Replace(androidXML, "%" + strconv.Itoa(i) + "$s", "%@", 1)
	}
	return androidXML
}

func ReadLines(path string) (lines []string, err error) {
	var (
		file   *os.File
		part   []byte
		prefix bool
	)

	if file, err = os.Open(path); err != nil {
		return
	}

	reader := bufio.NewReader(file)
	buffer := bytes.NewBuffer(make([]byte, 1024))

	for {
		if part, prefix, err = reader.ReadLine(); err != nil {
			break
		}
		buffer.Write(part)
		if !prefix {
			lines = append(lines, buffer.String())
			buffer.Reset()
		}
	}
	if err == io.EOF {
		err = nil
	}
	return
}

func WriteToFile(path string, values []string) {
	file, err := os.Create(path)
	if err != nil {
		log.Println("Open error:", err.Error())
		return
	}
	writer := bufio.NewWriter(file)
	for _, value := range values {
		n, err := writer.WriteString(value + "\n")
		if err != nil {
			log.Println("WriteString error:", err.Error())
			return
		}
		log.Println("n:", n)
	}
	writer.Flush()
	file.Close()
	if err == io.EOF {
		err = nil
	}
	return
}
