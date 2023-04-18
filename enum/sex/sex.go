package sex

import "database/sql/driver"

type Sex string

const (
	FEMALE Sex = "FEMALE"
	MALE   Sex = "MALE"
	OTHER  Sex = "OTHER"
)

func (s Sex) String() string {
	switch s {
	case FEMALE:
		return "FEMALE"
	case MALE:
		return "MALE"
	case OTHER:
		return "OTHER"
	}
	return "OTHER"
}

func (s *Sex) Scan(value interface{}) error {
	*s = Sex(value.([]byte))
	return nil
}

func (s Sex) Value() (driver.Value, error) {
	return string(s), nil
}
