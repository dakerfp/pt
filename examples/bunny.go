package main

import (
	"flag"
	. "github.com/fogleman/pt/pt"
)

var width = flag.Int("w", 500, "width")
var height = flag.Int("h", 300, "height")
var spp = flag.Int("spp", 1, "spp")
var interactions = flag.Int("interactions", 1024, "interactions")
var pathTemplate = flag.String("path", "bunny-%04d.npy", "")

func main() {
	flag.Parse()

	scene := Scene{}
	material := GlossyMaterial(HexColor(0xF2EBC7), 1.5, Radians(0))
	mesh, err := LoadOBJ("examples/bunny.obj", material)
	if err != nil {
		panic(err)
	}
	mesh.SmoothNormals()
	mesh.FitInside(Box{V(-1, 0, -1), V(1, 2, 1)}, V(0.5, 0, 0.5))
	scene.Add(mesh)
	floor := GlossyMaterial(HexColor(0x33332D), 1.2, Radians(20))
	scene.Add(NewCube(V(-10000, -10000, -10000), V(10000, 0, 10000), floor))
	scene.Add(NewSphere(V(0, 5, 0), 1, LightMaterial(White, 20)))
	scene.Add(NewSphere(V(4, 5, 4), 1, LightMaterial(White, 20)))
	camera := LookAt(V(-1, 2, 3), V(0, 0.75, 0), V(0, 1, 0), 50)
	sampler := NewSampler(4, 4)
	sampler.SpecularMode = SpecularModeFirst
	renderer := NewRenderer(&scene, &camera, sampler, *width, *height, *spp)
	renderer.FireflySamples = 128
	//renderer.IterativeRender("out%03d.png", 1000)
	renderer.ExportFeatures(*pathTemplate, *interactions)
}
