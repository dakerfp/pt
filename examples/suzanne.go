package main

import (
	"flag"
	. "github.com/fogleman/pt/pt"
)

var width = flag.Int("w", 500, "width")
var height = flag.Int("h", 300, "height")
var spp = flag.Int("spp", 1, "spp")
var interactions = flag.Int("interactions", 1024, "interactions")
var pathTemplate = flag.String("path", "suzanne-%04d.npy", "")

func main() {
	flag.Parse()

	scene := Scene{}
	material := DiffuseMaterial(HexColor(0x334D5C))
	scene.Add(NewSphere(V(0.5, 1, 3), 1, LightMaterial(White, 4)))
	scene.Add(NewSphere(V(1.5, 1, 3), 1, LightMaterial(White, 4)))
	scene.Add(NewCube(V(-5, -5, -2), V(5, 5, -1), material))
	mesh, err := LoadOBJ("examples/suzanne.obj", SpecularMaterial(HexColor(0xEFC94C), 1.3))
	if err != nil {
		panic(err)
	}
	scene.Add(mesh)
	camera := LookAt(V(1, -0.45, 4), V(1, -0.6, 0.4), V(0, 1, 0), 40)
	sampler := NewSampler(16, 8)
	renderer := NewRenderer(&scene, &camera, sampler, *width, *height, *spp)
	
	//renderer.IterativeRender("out%03d.png", 1000)
	renderer.ExportFeatures(*pathTemplate, *interactions)
}
