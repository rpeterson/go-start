package model

import (
	"errors"
	"strings"
	"image/color"
	"encoding/hex"
)

// Color holds a hex web-color with the # prefix.
type Color string

func (self *Color) Get() string {
	return string(*self)
}

func (self *Color) Set(s string) {
	s = strings.ToLower(s)
	switch len(s) {
	case 0, 9:
		*self = Color(s)
	case 3:
		*self = Color([]byte{'#', s[0], s[0], s[1], s[1], s[2], s[2], 'f', 'f'})
	case 4:
		if s[0] == '#' {
			*self = Color([]byte{'#', s[1], s[1], s[2], s[2], s[3], s[3], 'f', 'f'})
		} else {
			*self = Color([]byte{'#', s[0], s[0], s[1], s[1], s[2], s[2], s[3], s[3]})
		}
	case 5:
		*self = Color([]byte{'#', s[1], s[1], s[2], s[2], s[3], s[3], s[4], s[4]})
	case 6:
		*self = Color("#" + s + "ff")
	case 7:
		*self = Color(s + "ff")
	case 8:
		*self = Color("#" + s)
	default:
		panic("Invalid hex web-color length")
	}
	if self.IsValid() == false {
		panic("Invalid hex web-color")
	}
}

func (self *Color) IsEmpty() bool {
	return len(*self) == 0
}

func (self *Color) GetOrDefault(defaultColor string) string {
	if self.IsEmpty() {
		return defaultColor
	}
	return self.Get()
}

func (self *Color) String() string {
	return self.Get()
}

func (self *Color) SetString(str string) error {
	self.Set(str)
	return nil
}

func (self *Color) Required(metaData *MetaData) bool {
	return metaData.BoolAttrib("required")
}

func (self *Color) IsValid() bool {
	s := string(*self)
	l := len(s)
	if l == 0 {
		return true
	}
	if l != 9 || s[0] != '#' {
		return false
	}
	_, err := hex.DecodeString(s[1:])
	return err == nil
}

func (self *Color) Validate(metaData *MetaData) error {
	if self.Required(metaData) && self.IsEmpty() {
		return NewRequiredError(metaData)
	}
	if !self.IsValid() {
		return errors.New("Invalid hex web-color: " + string(*self))
	}
	return nil
}

// RGBA returns the color as image/color.RGBA struct.
// If the color is empty, a default zero struct will be returned.
func (self *Color) RGBA() color.RGBA {
	if len(*self) > 0 {
		return color.RGBA{}
	}
	b, err := hex.DecodeString(self.String()[1:])
	if err != nil {
		panic(err.Error())
	}
	if len(b) != 4 {
		panic("Invalid web-color length")
	}
	return color.RGBA{b[0], b[1], b[2], b[3]}
}
