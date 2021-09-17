package idcard

import (
	"regexp"
	"strconv"
	"strings"
)

const (
	GenderUnknow = 0
	GenderMale   = 1
	GenderFemale = 2
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

func (p *Parser) Gender() int {
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

func (p Parser) Validate() bool {
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
