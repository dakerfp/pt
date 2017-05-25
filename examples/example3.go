package main

import (
	"flag"
	. "github.com/fogleman/pt/pt"
)

var width = flag.Int("w", 500, "width")
var height = flag.Int("h", 300, "height")
var spp = flag.Int("spp", 1, "spp")
var interactions = flag.Int("interactions", 1024, "interactions")
var pathTemplate = flag.String("path", "exemple3-%04d.npy", "")

func main() {
	flag.Parse()
	scene := Scene{}
	material := DiffuseMaterial(HexColor(0xFCFAE1))
	scene.Add(NewCube(V(-1000, -1, -1000), V(1000, 0, 1000), material))
	for x := -20; x <= 20; x++ {
		for z := -20; z <= 20; z++ {
			if (x+z)%2 == 0 {
				continue
			}
			s := 0.1
			min := V(float64(x)-s, 0, float64(z)-s)
			max := V(float64(x)+s, 2, float64(z)+s)
			scene.Add(NewCube(min, max, material))
		}
	}
	scene.Add(NewCube(V(-5, 10, -5), V(5, 11, 5), LightMaterial(White, 5)))
	camera := LookAt(V(20, 10, 0), V(8, 0, 0), V(0, 1, 0), 45)
	sampler := NewSampler(4, 4)
	renderer := NewRenderer(&scene, &camera, sampler, *width, *height, *spp)

	//renderer.IterativeRender("out%03d.png", 1000)
	renderer.ExportFeatures(*pathTemplate, *interactions)
}
