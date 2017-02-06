// ray.go
package main

// RAY
type Ray struct {
	origin    Vector
	direction Vector
	interDist float64 // MAX_DIST
	interObj  int
}
