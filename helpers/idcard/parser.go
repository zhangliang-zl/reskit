package idcard

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	GenderUnknow = "未知"
	GenderMale   = "男"
	GenderFemale = "女"
)

var (
	salt     = []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	checksum = []string{"1", "0", "X", "9", "8", "7", "6", "5", "4", "3", "2"}
)

type Parser struct {
	id        string
	checked   bool
	validated bool
}

func NewParser(id string) *Parser {
	return &Parser{
		id: id,
	}
}

func (p *Parser) Birthday() string {
	if len(p.id) == 18 {
		return p.id[6:14]
	} else {
		return "19" + p.id[6:12]
	}
}

func (p *Parser) Age() int {
	birthday := p.Birthday()

	birYear, _ := strconv.Atoi(birthday[0:4])

	birMonth, _ := strconv.Atoi(strings.TrimRight(birthday[4:6], "0"))

	age := time.Now().Year() - birYear

	if int(time.Now().Month()) < birMonth {
		age--
	}
	return age
}

func (p *Parser) Gender() string {
	var g string
	if len(p.id) == 18 {
		g = p.id[16:17]
	} else if len(p.id) == 15 {
		g = p.id[14:15]
	} else {
		return GenderUnknow
	}

	intG, _ := strconv.Atoi(g)
	if intG%2 == 0 {
		return GenderFemale
	}
	return GenderMale
}

func (p *Parser) Validate() bool {
	if len(p.id) != 15 && len(p.id) != 18 {
		return false
	}

	if p.checked {
		return p.validated
	}
	p.validated = p.checkFormat() && p.checkBirthday() && p.checkLastCode()
	p.checked = true
	return p.validated
}

func (p Parser) checkFormat() bool {
	pattern := `^([\d]{17}[xX\d]|[\d]{15})$`
	matched, _ := regexp.MatchString(pattern, p.id)
	return matched
}

func (p Parser) checkBirthday() bool {
	pattern := `^(([0-9]{3}[1-9]|[0-9]{2}[1-9][0-9]{1}|[0-9]{1}[1-9][0-9]{2}|[1-9][0-9]{3})(((0[13578]|1[02])(0[1-9]|[12][0-9]|3[01]))|((0[469]|11)(0[1-9]|[12][0-9]|30))|(02(0[1-9]|[1][0-9]|2[0-8]))))|((([0-9]{2})(0[48]|[2468][048]|[13579][26])|((0[48]|[2468][048]|[3579][26])00))0229)$`
	matched, _ := regexp.MatchString(pattern, p.Birthday())
	return matched
}

func (p Parser) Province() string {
	provinces := make(map[string]string)
	provinces["11"] = "北京市"
	provinces["12"] = "天津市"
	provinces["13"] = "河北省"
	provinces["14"] = "山西省"
	provinces["15"] = "内蒙古自治区"
	provinces["21"] = "辽宁省"
	provinces["22"] = "吉林省"
	provinces["23"] = "黑龙江省"
	provinces["31"] = "上海市"
	provinces["32"] = "江苏省"
	provinces["33"] = "浙江省"
	provinces["34"] = "安徽省"
	provinces["35"] = "福建省"
	provinces["36"] = "江西省"
	provinces["37"] = "山东省"
	provinces["41"] = "河南省"
	provinces["42"] = "湖北省"
	provinces["43"] = "湖南省"
	provinces["44"] = "广东省"
	provinces["45"] = "广西壮族自治区"
	provinces["46"] = "海南省"
	provinces["50"] = "重庆市"
	provinces["51"] = "四川省"
	provinces["52"] = "贵州省"
	provinces["53"] = "云南省"
	provinces["54"] = "西藏自治区"
	provinces["61"] = "陕西省"
	provinces["62"] = "甘肃省"
	provinces["63"] = "青海省"
	provinces["64"] = "宁夏回族自治区"
	provinces["65"] = "新疆维吾尔自治区"
	provinces["71"] = "台湾省"
	provinces["81"] = "香港特别行政区"
	provinces["82"] = "澳门特别行政区"

	left := p.id[0:2]
	pro, ok := provinces[left]
	if !ok {
		return "未知"
	}

	return pro
}

func (p Parser) checkLastCode() bool {
	if len(p.id) == 15 {
		return true
	}

	if len(p.id) != 18 {
		return false
	}

	sum := 0
	for i := 0; i < 17; i++ {
		char, _ := strconv.Atoi(p.id[i : i+1])
		sum += char * salt[i]
	}

	seek := sum % 11
	return checksum[seek] == strings.ToUpper(p.id[17:18])
}
