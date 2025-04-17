package enum

import (
	"encoding/json"
)

type Gender int

const (
	GENDER_UNKNOWN Gender = iota
	GENDER_FEMALE
	GENDER_MALE
)

var (
	genderName = map[Gender]string{
		GENDER_UNKNOWN: "ไม่ระบุ",
		GENDER_FEMALE:  "หญิง",
		GENDER_MALE:    "ชาย",
	}
	intToGender = map[int]Gender{
		0: GENDER_UNKNOWN,
		1: GENDER_FEMALE,
		2: GENDER_MALE,
	}
)

func (s Gender) List() map[Gender]string {
	return genderName
}

func (s Gender) String() string {
	return genderName[s]
}

func GetGender(value int) Gender {
	if gender, ok := intToGender[value]; ok {
		return gender
	}
	return GENDER_UNKNOWN
}

func (s *Gender) UnmarshalJSON(data []byte) error {
	var value int
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*s = GetGender(value)
	return nil
}

// GetGenderFromString converts a string to a Gender enum value
func GetGenderFromInt(i int) Gender {
	switch i {
	case 0:
		return GENDER_MALE
	case 1:
		return GENDER_FEMALE
	default:
		return GENDER_UNKNOWN
	}
}

func (s Gender) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}
