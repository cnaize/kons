package kons

import "regexp"

func PathToRegs(path ...string) ([]*regexp.Regexp, error) {
	var regs []*regexp.Regexp
	for _, p := range path {
		reg, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		regs = append(regs, reg)
	}
	return regs, nil
}
