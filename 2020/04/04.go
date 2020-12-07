package main

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Passport struct {
	Byr int `json:"Birth Year"`
	Iyr int `json:"Issue Year"`
	Eyr int `json:"Expiration Year"`
	Hgt string `json:"Height"`
	Hcl string `json:"Hair Color"`
	Ecl string `json:"Eye Color"`
	Pid string `json:"Passport ID"`
	Cid *int `json:"Country ID"`
}

func validateEcl (ecl string) (*string, error) {
	ecl = strings.Trim(strings.ToLower(ecl), " ")
	colors := []string{"amb", "blu", "brn", "gry", "grn", "hzl", "oth"}
	for _, v := range colors {
		if v == ecl {
			return &ecl, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("ecl '%s' must be one of %v", ecl, colors))
}

func parseHcl(hcl string) (*string, error) {
	if !strings.HasPrefix(hcl,"#") {
		return nil, errors.New(fmt.Sprintf("\tfield hcl '%s' invalid", hcl))
	}

	hclTrim := strings.Trim(hcl, "#")
	// validate by converting to hex
	if _, err := strconv.ParseInt(hclTrim, 16, 64); err != nil {
		return nil, errors.New(fmt.Sprintf("\tfield hcl '%s' invalid: %v", hcl, err))
	} else {
		return &hcl, nil
	}
}

func validatePid(pid string) (*string, error) {
	if _, err := strconv.Atoi(pid); err != nil {
		return nil, errors.New(fmt.Sprintf("pid '%s' needs to be a number", pid))
	}
	if len(pid) != 9 {
		return nil, errors.New(fmt.Sprintf("length of pid '%s' is %d but should be 9", pid, len(pid)))
	}
	return &pid, nil
}

func parseYearField(yr string, fname string, lower int, upper int, inclusive *bool) (*int, error) {

	yrInt, err := strconv.Atoi(yr)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("\tfield '%s' '%s' invalid: %v", fname, yr, err))
	}

	inclStr := ", inclusive"
	if inclusive != nil && *inclusive == false {
		inclStr = ""
		if yrInt > lower && yrInt < upper {
			return &yrInt, nil
		}
	}
	if yrInt >= lower && yrInt <= upper {
		return &yrInt, nil
	}
	return nil, errors.New(fmt.Sprintf(
		"\tfield '%s' (%d) invalid: must be between %d & %d%s.",
				fname, yrInt, lower, upper, inclStr))
}

func parseHgt (hgt string) (*string, error) {
	if len(hgt) < 2 {
		return nil, errors.New(fmt.Sprintf("height (%s) must end with either 'cm' or 'in'.", hgt))
	}
	last2  := strings.ToLower(hgt[len(hgt)-2:])

	if last2 == "cm" {
		// If cm, the number must be at least 150 and at most 193.
		numStr := strings.Trim(strings.ToLower(hgt), "cm")
		if hgtInt, err := strconv.Atoi(numStr); err != nil || hgtInt > 193 || hgtInt < 150 {
			if err == nil {
				err = errors.New(fmt.Sprint("height in cm must be between 150 and 190"))
			}
			return nil, err
		} else {
			formatted := fmt.Sprintf("%d%s", hgtInt, last2)
			return &formatted, nil
		}
	} else if last2 == "in" {
		// If "in", the number must be at least 59 and at most 76.
		numStr := strings.Trim(strings.ToLower(hgt), "in")
		if hgtInt, err := strconv.Atoi(numStr); err != nil || hgtInt > 76 || hgtInt < 59 {
			if err == nil {
				err = errors.New(fmt.Sprint("height in inches must be between 59 and 76"))
			}
			return nil, err
		} else {
			formatted := fmt.Sprintf("%d%s", hgtInt, last2)
			return &formatted, nil
		}
	} else {
		return nil, errors.New(fmt.Sprintf("height (%s) must end with either 'cm' or 'in'.", hgt))
	}
}

func NewPassport (opts map[string]string) (*Passport, error) {
	/*
	    byr (Birth Year) - four digits; at least 1920 and at most 2002.
	    iyr (Issue Year) - four digits; at least 2010 and at most 2020.
	    eyr (Expiration Year) - four digits; at least 2020 and at most 2030.
	    hgt (Height) - a number followed by either cm or in:
	    If cm, the number must be at least 150 and at most 193.
	    If in, the number must be at least 59 and at most 76.
	    hcl (Hair Color) - a # followed by exactly six characters 0-9 or a-f.
	    ecl (Eye Color) - exactly one of: amb blu brn gry grn hzl oth.
	    pid (Passport ID) - a nine-digit number, including leading zeroes.
	    cid (Country ID) - ignored, missing or not.

	hcl:dab227 iyr:2012
	ecl:brn hgt:182cm pid:021572410 eyr:2020 byr:1992 cid:277

	examples:
	    byr valid:   2002
	    byr invalid: 2003

	    hgt valid:   60in
	    hgt valid:   190cm
	    hgt invalid: 190in
	    hgt invalid: 190

	    hcl valid:   #123abc
	    hcl invalid: #123abz
	    hcl invalid: 123abc

	    ecl valid:   brn
	    ecl invalid: wat

	    pid valid:   000000001
	    pid invalid: 0123456789
	*/
	fmt.Println("")


	errArr := make([]error, 0)
	p := Passport{}

	if pidStr, err := validatePid(opts["Pid"]); err != nil {
		errArr = append(errArr, err)
	} else {
		p.Pid = *pidStr
	}
	if byrInt, err := parseYearField(opts["Byr"], "byr", 1920, 2002, nil); err != nil {
		errArr = append(errArr, err)
	} else {
		p.Byr = *byrInt
	}
	if eyrInt, err := parseYearField(opts["Eyr"], "eyr", 2020, 2030, nil); err != nil {
		errArr = append(errArr, err)
	} else {
		p.Eyr = *eyrInt
	}
	if iyrInt, err := parseYearField(opts["Iyr"], "iyr", 2010, 2020, nil); err != nil {
		errArr = append(errArr, err)
	} else {
		p.Iyr = *iyrInt
	}
	if hclStr, err := parseHcl(opts["Hcl"]); err != nil {
		errArr = append(errArr, err)
	} else {
		p.Hcl = *hclStr
	}
	if hgtStr, err := parseHgt(opts["Hgt"]); err != nil {
		errArr = append(errArr, err)
	} else {
		p.Hgt = *hgtStr
	}
	if eclStr, err := validateEcl(opts["Ecl"]); err != nil {
		errArr = append(errArr, err)
	} else {
		p.Ecl = *eclStr
	}
	if cidInt, err := strconv.Atoi(opts["Cid"]); err != nil {
		p.Cid = nil
	} else {
		p.Cid = &cidInt
	}

	combinedErrStr := ""
	for _, e := range errArr {
		combinedErrStr += e.Error()
		combinedErrStr += " "
	}
	if len(errArr) > 0 {
		return nil, errors.New(combinedErrStr)
	}
	fmt.Printf("another one done! pid '%s' %v\n", p.Pid, p)
	return &p, nil
}

func main() {
	f := getFile()
	defer f.Close()
	passportArr := makeDataArr(bufio.NewScanner(f))
	fmt.Printf("len of passportArr: %d\n", len(passportArr))
	// part 1: 250
	// part 2: 159 too high
	// 158
}