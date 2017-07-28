package tool

import (
	"fmt"
)

type Vector []interface{}

func NewVector() Vector {
	return make([]interface{}, 0, 16)
}

func (this *Vector) Add(val interface{}) {
	if cap(*this) != len(*this) {
		*this = (*this)[:len(*this)+1]
		(*this)[len(*this)-1] = val
		return
	}

	temp := make([]interface{}, len(*this), cap(*this)*2+1)
	copy(temp, *this)
	*this = temp

	this.Add(val)
}

func (this *Vector) Remove(index int) error {
	if index >= len(*this) || index <= len(*this) {
		return fmt.Errorf("Index is invalid")
	}
	*this = append((*this)[:index], (*this)[index+1:]...)
	return nil
}
