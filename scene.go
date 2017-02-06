package main

import (
	"bufio"
	"errors"
	"image"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

// SCENE
type Scene struct {
	imgWidth     int
	imgHeight    int
	traceDepth   int
	oversampling int
	visionField  float64
	startline    int
	endline      int
	gridWidth    int
	gridHeight   int
	cameraPos    Vector
	cameraLook   Vector
	cameraUp     Vector
	look         Vector
	Vhor         Vector
	Vver         Vector
	Vp           Vector
	image        *image.RGBA
	objectList   []Object
	lightList    []Light
	materialList []Material
}

func NewScene(sceneFilename string) *Scene {
	scn := new(Scene)
	// defaults
	scn.imgWidth = 320
	scn.imgHeight = 200

	scn.traceDepth = 3   // bounces
	scn.oversampling = 1 // no oversampling
	scn.visionField = 60

	scn.startline = 0 // Start rendering line
	scn.endline = scn.imgHeight - 1

	//scn.objectList = append(scn.objectList, Sphere{0,0.0,0.0,0.0,0.0})

	f, err := os.Open(sceneFilename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := bufio.NewReaderSize(f, 4*1024)
	if err != nil {
		panic(err)
	}
	line, isPrefix, err := r.ReadLine()

	for err == nil && !isPrefix {

		s := string(line)
		if len(s) == 0 {
			line, isPrefix, err = r.ReadLine()
			continue
		}

		if s[0:1] == "#" {
			line, isPrefix, err = r.ReadLine()
			continue
		}

		sline := strings.Split(s, " ")
		word := sline[0]
		untrimmed := sline[1:]
		var data []string

		for _, item := range untrimmed {
			if item == "" || item == " " {
				continue
			}
			data = append(data, strings.Trim(item, " "))
		}

		switch word {
			case "image_size":
				scn.imgWidth, _ = strconv.Atoi(data[0])
				scn.imgHeight, _ = strconv.Atoi(data[1])
				scn.endline = scn.imgHeight - 1 // End rendering line
			case "depth":
				scn.traceDepth, _ = strconv.Atoi(data[0]) // n. bounces
			case "oversampling":
				scn.oversampling, _ = strconv.Atoi(data[0])
			case "field_of_view":
				scn.visionField, _ = strconv.ParseFloat(data[0], 64)
			case "renderslice":
				scn.startline, _ = strconv.Atoi(data[0])
				scn.endline, _ = strconv.Atoi(data[1])
			case "camera_position":
				scn.cameraPos = ParseVector(data)
			case "camera_look":
				scn.cameraLook = ParseVector(data)
			case "camera_up":
				scn.cameraUp = ParseVector(data)
			case "sphere":
				mat, _ := strconv.Atoi(data[0])
				rad, _ := strconv.ParseFloat(data[4], 64)
				scn.objectList = append(scn.objectList, Sphere{mat, ParseVector(data[1:4]), rad})
			case "plane":
				mat, _ := strconv.Atoi(data[0])
				dis, _ := strconv.ParseFloat(data[4], 64)
				scn.objectList = append(scn.objectList, Plane{mat, ParseVector(data[1:4]), dis})
			case "light":
				light := Light{ParseVector(data[1:4]), ParseColor(data[4:7]), data[0]}
				scn.lightList = append(scn.lightList, light)
			case "material":
				mat := ParseMaterial(data)
				scn.materialList = append(scn.materialList, mat)
		}
		line, isPrefix, err = r.ReadLine()
	}

	if isPrefix {
		panic(errors.New("buffer size to small"))
	}

	if err != io.EOF {
		panic(err)
	}

	scn.image = image.NewRGBA(image.Rect(0, 0, scn.imgWidth, scn.imgHeight))

	scn.gridWidth = scn.imgWidth * scn.oversampling
	scn.gridHeight = scn.imgHeight * scn.oversampling

	scn.look = scn.cameraLook.Sub(scn.cameraPos)
	scn.Vhor = scn.look.Cross(scn.cameraUp)
	scn.Vhor = scn.Vhor.Normalize()

	scn.Vver = scn.look.Cross(scn.Vhor)
	scn.Vver = scn.Vver.Normalize()

	fl := float64(scn.gridWidth) / (2 * math.Tan((0.5*scn.visionField)*PI_180))

	Vp := scn.look.Normalize()

	Vp.x = Vp.x*fl - 0.5*(float64(scn.gridWidth)*scn.Vhor.x+float64(scn.gridHeight)*scn.Vver.x)
	Vp.y = Vp.y*fl - 0.5*(float64(scn.gridWidth)*scn.Vhor.y+float64(scn.gridHeight)*scn.Vver.y)
	Vp.z = Vp.z*fl - 0.5*(float64(scn.gridWidth)*scn.Vhor.z+float64(scn.gridHeight)*scn.Vver.z)

	scn.Vp = Vp

	return scn
}

// Auxiliary Methods
func ParseVector(line []string) Vector {
	x, _ := strconv.ParseFloat(line[0], 64)
	y, _ := strconv.ParseFloat(line[1], 64)
	z, _ := strconv.ParseFloat(line[2], 64)
	return Vector{x, y, z}
}

func ParseColor(line []string) Color {
	r, _ := strconv.ParseFloat(line[0], 64)
	g, _ := strconv.ParseFloat(line[1], 64)
	b, _ := strconv.ParseFloat(line[2], 64)
	return Color{r, g, b}
}

func ParseMaterial(line []string) Material {
	var f [6]float64
	for i, item := range line[3:] {
		f[i], _ = strconv.ParseFloat(item, 64)
	}
	return Material{ParseColor(line[0:3]), f[0], f[1], f[2], f[3], f[4], f[5]}
}
