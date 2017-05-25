package main

import (
	"flag"
	. "github.com/fogleman/pt/pt"
)

var width = flag.Int("w", 500, "width")
var height = flag.Int("h", 300, "height")
var spp = flag.Int("spp", 1, "spp")
var interactions = flag.Int("interactions", 1024, "interactions")
var pathTemplate = flag.String("path", "exemple2-%04d.npy", "")

func main() {
	flag.Parse()
	scene := Scene{}
	material := GlossyMaterial(HexColor(0xEFC94C), 3, Radians(30))
	whiteMat := GlossyMaterial(White, 3, Radians(30))
	for x := 0; x < 40; x++ {
		for z := 0; z < 40; z++ {
			center := V(float64(x)-19.5, 0, float64(z)-19.5)
			scene.Add(NewSphere(center, 0.4, material))
		}
	}
	scene.Add(NewCube(V(-100, -1, -100), V(100, 0, 100), whiteMat))
	scene.Add(NewSphere(V(-1, 4, -1), 1, LightMaterial(White, 30)))
	camera := LookAt(V(0, 4, -8), V(0, 0, -2), V(0, 1, 0), 45)
	sampler := NewSampler(4, 4)
	renderer := NewRenderer(&scene, &camera, sampler, *width, *height, *spp)
	
	//renderer.IterativeRender("out%03d.png", 1000)
	renderer.ExportFeatures(*pathTemplate, *interactions)
}
