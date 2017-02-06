// primitives.go
package main

import (
	"fmt"
	"math"
)

// CUERPOS
type Object interface {
	Type() string
	Material() int
	Intersect(r *Ray, i int) bool
	getNormal(point Vector) Vector
}

// ESFERA
type Sphere struct {
	material int
	position Vector
	radius   float64
}

func (e Sphere) Type() string {
	return "sphere"
}

func (e Sphere) Material() int {
	return e.material
}

func (e Sphere) Intersect(r *Ray, i int) bool {
	sphereRay := e.position.Sub(r.origin)
	v := sphereRay.Dot(r.direction)
	if v-e.radius > r.interDist {
		return false
	}
	collisionDist := e.radius*e.radius + v*v - sphereRay.x*sphereRay.x - sphereRay.y*sphereRay.y - sphereRay.z*sphereRay.z
	if collisionDist < 0.0 {
		return false
	}
	collisionDist = v - math.Sqrt(collisionDist)
	if collisionDist > r.interDist || collisionDist < 0.0 {
		return false
	}
	r.interDist = collisionDist
	r.interObj = i
	return true
}

func (e Sphere) getNormal(point Vector) Vector {
	normal := point.Sub(e.position)
	return normal.Normalize()
}

func (e Sphere) String() string {
	return fmt.Sprintf("<Esf: %d %s %.2f>", e.material, e.position.String(), e.radius)
}

// PLANO
type Plane struct {
	material  int
	normal    Vector
	distancia float64
}

func (p Plane) Type() string {
	return "plane"
}

func (p Plane) Material() int {
	return p.material
}

func (p Plane) Intersect(r *Ray, i int) bool {
	v := p.normal.Dot(r.direction)
	if v == 0 {
		return false
	}
	collisionDist := -(p.normal.Dot(r.origin) + p.distancia) / v
	if collisionDist < 0.0 {
		return false
	}
	if collisionDist > r.interDist {
		return false
	}
	r.interDist = collisionDist
	r.interObj = i
	return true
}

func (p Plane) getNormal(point Vector) Vector {
	return p.normal
}

func (p Plane) String() string {
	return fmt.Sprintf("<Pla: %d %s %.2f>", p.material, p.normal.String(), p.distancia)
}
