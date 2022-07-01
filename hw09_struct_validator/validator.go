package hw09structvalidator

import (
	"reflect"
	"regexp"
	"strconv"
)

const TAGV = "validate"

var (
	reConjunction = regexp.MustCompile(`\|`)
	reRule        = regexp.MustCompile(`:`)
	reOccurrence  = regexp.MustCompile(`,`)
)

// Validate validates structure s according rules in the structure tags.
func Validate(s interface{}) error {
	vldErrs := make(ValidationErrors, 0)

	sType := reflect.TypeOf(s)
	if sType.Kind() != reflect.Struct {
		return ErrExpectedStruct
	}
	sValue := reflect.ValueOf(s)

	for i := 0; i < sType.NumField(); i++ {
		sField := sType.Field(i)
		fieldName := sField.Name
		fieldValue := sValue.Field(i)

		if alias, ok := sField.Tag.Lookup(TAGV); ok {
			if alias == "" {
				vldErrs.Add(fieldName, ErrRuleIncorrect)
				continue
			}
			rules, err := splitRules(alias)
			if err != nil {
				vldErrs.Add(fieldName, err)
				continue
			}
			vldFieldErrs := validateField(fieldName, fieldValue, rules)
			if vldFieldErrs != nil {
				vldErrs.AddList(vldFieldErrs)
			}
		}
	}
	if len(vldErrs) > 0 {
		return vldErrs
	}

	return nil
}

type Ruler struct {
	Rule  string
	Value string
}

// splitRules reads rules in validating structure tags.
func splitRules(str string) ([]Ruler, error) {
	ruleList := make([]Ruler, 0)
	for _, rl := range reConjunction.Split(str, -1) {
		fr := reRule.Split(rl, -1)
		if len(fr) == 2 && fr[0] != "" && fr[1] != "" {
			ruleList = append(ruleList, Ruler{Rule: fr[0], Value: fr[1]})
		} else {
			return nil, ErrRuleIncorrect
		}
	}
	return ruleList, nil
}

// validateField is the base method for the field validation.
func validateField(fieldName string, value reflect.Value, rules []Ruler) ValidationErrors {
	vldErrs := make(ValidationErrors, 0)
	switch value.Kind() { //nolint:exhaustive
	case reflect.String:
		for _, rule := range rules {
			err := validateString(value.String(), rule)
			if err != nil {
				vldErrs.Add(fieldName, err)
			}
		}
	case reflect.Int:
		for _, rule := range rules {
			err := validateInt(int(value.Int()), rule)
			if err != nil {
				vldErrs.Add(fieldName, err)
			}
		}
	case reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			err := validateField(fieldName, value.Index(i), rules)
			if err != nil {
				vldErrs.AddList(err)
			}
		}
	}
	return vldErrs
}

// validateString validates string fields.
func validateString(value string, rule Ruler) error {
	switch rule.Rule {
	case "len":
		rv, err := strconv.Atoi(rule.Value)
		if err != nil {
			return ErrRuleValueIncorrect
		}
		if len(value) != rv {
			return ErrValidationStringLength
		}
	case "regexp":
		matched, err := regexp.MatchString(rule.Value, value)
		if err != nil {
			return ErrRuleValueIncorrect
		}
		if !matched {
			return ErrValidationStringRegexp
		}
	case "in":
		patterns := reOccurrence.Split(rule.Value, -1)
		for _, p := range patterns {
			if p == value {
				return nil
			}
		}
		return ErrValidationOccurrence
	default:
		return ErrRuleUnsupported
	}
	return nil
}

// validateString validates int fields.
func validateInt(value int, rule Ruler) error {
	switch rule.Rule {
	case "max":
		rv, err := strconv.Atoi(rule.Value)
		if err != nil {
			return ErrRuleValueIncorrect
		}
		if value > rv {
			return ErrValidationIntMax
		}
	case "min":
		rv, err := strconv.Atoi(rule.Value)
		if err != nil {
			return ErrRuleValueIncorrect
		}
		if value < rv {
			return ErrValidationIntMin
		}
	case "in":
		patterns := reOccurrence.Split(rule.Value, -1)
		v := strconv.Itoa(value)
		for _, p := range patterns {
			if p == v {
				return nil
			}
		}
		return ErrValidationOccurrence
	default:
		return ErrRuleUnsupported
	}
	return nil
}
