package main

import (
	"flag"
	"fmt"
	"image/png"
	"math"
	"os"
	"runtime"
)

const (
	MAX_DIST = 1999999999
	PI_180   = 0.017453292
	SMALL    = 0.000000001
)

func calcShadow(r *Ray, collisionObj int) float64 {
	shadow := 1.0 //starts with no shadow
	for i, obj := range scene.objectList {
		r.interObj = -1
		r.interDist = MAX_DIST

		if obj.Intersect(r, i) && i != collisionObj {
			shadow *= scene.materialList[obj.Material()].transmitCol
		}
	}
	return shadow
}

func trace(r *Ray, depth int) (c Color) {

	for i, obj := range scene.objectList {
		obj.Intersect(r, i)
	}

	if r.interObj >= 0 {
		matIndex := scene.objectList[r.interObj].Material()
		interPoint := r.origin.Add(r.direction.Mul(r.interDist))
		incidentV := interPoint.Sub(r.origin)
		originBackV := r.direction.Mul(-1.0)
		originBackV = originBackV.Normalize()
		vNormal := scene.objectList[r.interObj].getNormal(interPoint)
		for _, light := range scene.lightList {
			switch light.kind {
			case "ambient":
				c = c.Add(light.color)
			case "point":
				lightDir := light.position.Sub(interPoint)
				lightDir = lightDir.Normalize()
				lightRay := Ray{interPoint, lightDir, MAX_DIST, -1}
				shadow := calcShadow(&lightRay, r.interObj)
				NL := vNormal.Dot(lightDir)

				if NL > 0.0 {
					if scene.materialList[matIndex].difuseCol > 0.0 { // ------- Difuso
						difuseColor := light.color.Mul(scene.materialList[matIndex].difuseCol).Mul(NL)
						difuseColor.r *= scene.materialList[matIndex].color.r * shadow
						difuseColor.g *= scene.materialList[matIndex].color.g * shadow
						difuseColor.b *= scene.materialList[matIndex].color.b * shadow
						c = c.Add(difuseColor)
					}
					if scene.materialList[matIndex].specularCol > 0.0 { // ----- Especular
						R := (vNormal.Mul(2).Mul(NL)).Sub(lightDir)
						spec := originBackV.Dot(R)
						if spec > 0.0 {
							spec = scene.materialList[matIndex].specularCol * math.Pow(spec, scene.materialList[matIndex].specularD)
							specularColor := light.color.Mul(spec).Mul(shadow)
							c = c.Add(specularColor)
						}
					}
				}
			}
		}
		if depth < scene.traceDepth {
			if scene.materialList[matIndex].reflectionCol > 0.0 { // -------- Reflexion
				T := originBackV.Dot(vNormal)
				if T > 0.0 {
					vDirRef := (vNormal.Mul(2).Mul(T)).Sub(originBackV)
					vOffsetInter := interPoint.Add(vDirRef.Mul(SMALL))
					rayoRef := Ray{vOffsetInter, vDirRef, MAX_DIST, -1}
					c = c.Add(trace(&rayoRef, depth+1.0).Mul(scene.materialList[matIndex].reflectionCol))
				}
			}
			if scene.materialList[matIndex].transmitCol > 0.0 { // ---- Refraccion
				RN := vNormal.Dot(incidentV.Mul(-1.0))
				incidentV = incidentV.Normalize()
				var n1, n2 float64
				if vNormal.Dot(incidentV) > 0.0 {
					vNormal = vNormal.Mul(-1.0)
					RN = -RN
					n1 = scene.materialList[matIndex].IOR
					n2 = 1.0
				} else {
					n2 = scene.materialList[matIndex].IOR
					n1 = 1.0
				}
				if n1 != 0.0 && n2 != 0.0 {
					par_sqrt := math.Sqrt(1 - (n1*n1/n2*n2)*(1-RN*RN))
					refactDirV := incidentV.Add(vNormal.Mul(RN).Mul(n1 / n2)).Sub(vNormal.Mul(par_sqrt))
					vOffsetInter := interPoint.Add(refactDirV.Mul(SMALL))
					refractRay := Ray{vOffsetInter, refactDirV, MAX_DIST, -1}
					c = c.Add(trace(&refractRay, depth+1.0).Mul(scene.materialList[matIndex].transmitCol))
				}
			}
		}
	}
	return c
}

func renderPixel(line chan int, done chan bool) {
	for y := range line { // 1: 1, 5: 2, 8: 3,
		for x := 0; x < scene.imgWidth; x++ {
			var c Color
			yo := y * scene.oversampling
			xo := x * scene.oversampling
			for i := 0; i < scene.oversampling; i++ {
				for j := 0; j < scene.oversampling; j++ {
					var dir Vector
					dir.x = float64(xo)*scene.Vhor.x + float64(yo)*scene.Vver.x + scene.Vp.x
					dir.y = float64(xo)*scene.Vhor.y + float64(yo)*scene.Vver.y + scene.Vp.y
					dir.z = float64(xo)*scene.Vhor.z + float64(yo)*scene.Vver.z + scene.Vp.z
					dir = dir.Normalize()
					r := Ray{scene.cameraPos, dir, MAX_DIST, -1}
					c = c.Add(trace(&r, 1.0))
					yo += 1
				}
				xo += 1
			}
			srq_oversampling := float64(scene.oversampling * scene.oversampling)
			c.r /= srq_oversampling
			c.g /= srq_oversampling
			c.b /= srq_oversampling
			scene.image.SetRGBA(x, y, c.ToPixel())
			//fmt.Println("check")
		}
		if y%10 == 0 {
			fmt.Printf("%d ", y)
		}
	}
	done <- true
}

var scene *Scene

func main() {
	var sceneFilename string
	var numWorkers int
	flag.StringVar(&sceneFilename, "file", "samples/scene.txt", "Scene file to render.")
	flag.IntVar(&numWorkers, "workers", runtime.NumCPU(), "Number of worker threads.")
	flag.Parse()

	scene = NewScene(sceneFilename)
	done := make(chan bool, numWorkers)
	line := make(chan int)

	// launch the workers
	for i := 0; i < numWorkers; i++ {
		go renderPixel(line, done)
	}

	fmt.Println("Rendering: ", sceneFilename)
	fmt.Printf("Line (from %d to %d): ", scene.startline, scene.endline)
	for y := scene.startline; y < scene.endline; y++ {
		line <- y
	}
	close(line)

	// wait for all workers to finish
	for i := 0; i < numWorkers; i++ {
		<-done
	}

	output, err := os.Create(sceneFilename[0:len(sceneFilename)-4] + ".png")
	if err != nil {
		panic(err)
	}

	if err = png.Encode(output, scene.image); err != nil {
		panic(err)
	}
	fmt.Println(" DONE!")
}
