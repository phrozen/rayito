// material.go
package main

import (
	"fmt"
)

type Material struct {
	color                                                              Color
	difuseCol, specularCol, specularD, reflectionCol, transmitCol, IOR float64
}

func (m Material) String() string {
	return fmt.Sprintf("<Mat: %s %.2f %.2f %.2f %.2f %.2f %.2f>", m.color.String(), m.difuseCol, m.specularCol, m.specularD, m.reflectionCol, m.transmitCol, m.IOR)
}
