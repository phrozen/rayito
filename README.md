# Rayito

**Rayito** is a simple but educational ray tracer written in Go meant as a project to learn programming languages and optimization in Computer Graphics.

*"In computer graphics, ray tracing is a technique for generating an image by tracing the path of light through pixels in an image plane and simulating the effects of its encounters with virtual objects." - Wikipedia*

**Rayito** is meant to be simple, but implements a full set of features like:

+ Multi Threading
+ Oversampling (anti aliasing)
+ Ambient and Point Lights
+ Custom Materials
+ Specular and Difuse shading (Phong)
+ Reflection (mirror)
+ Refraction (transparency)
+ Soft Shadows
+ Scene description and parsing
+ PNG Export
+ and more...

Code is meant to be as clear as possible (no cutting names on variables) so it's meant to be read directly. In order to try **Rayito** either clone the git repository or do:

```
go get github.com/phrozen/rayito
```

After building the binary you can start to render by using some scenes in the *samples* directory or write your own.

### Usage
Set workers to the number of threads your CPU can handle (defaults to 4).
```
rayito -workers=8 -file=samples/test.txt
```

### Samples

![Helix](https://raw.githubusercontent.com/phrozen/rayito/master/samples/helix.txt.png)

![Spheres](https://raw.githubusercontent.com/phrozen/rayito/master/samples/test.txt.png)
