package commands

import (
	"fmt"
	"regexp"
)

type ValidString struct {
	Input string
	name  string
	regex string
}

func (v *ValidString) Validate() error {
	var err error

	if v.regex == "." && v.Input == "" {
		err = (fmt.Errorf("%s can not be empty", v.name))
	} else {
		_, regerr := regexp.MatchString(v.regex, v.Input)
		if regerr != nil {
			err = fmt.Errorf("%s is not accepted by parameter %s (regex %s)", v.Input, v.name, v.regex)
		}
	}

	return err
}

func ValidateArray(vss []ValidString) error {
	var err error

	for i := 0; i < len(vss) && err == nil; i = i + 1 {
		err = vss[i].Validate()
	}

	return err
}
func ValidateOrArray(vsso [][]ValidString) error {
	var err error

	for i := 0; i < len(vsso) && err == nil; i = i + 1 {
		var ferr error
		foundnil := false
		vss := vsso[i]
		for j := 0; j < len(vss) && !foundnil; j++ {
			v := vss[j].Validate()
			if v != nil {
				if ferr == nil {
					ferr = v
				}
			} else {
				foundnil = true
			}
		}
		if !foundnil {
			err = ferr
		}
	}

	return err

}
