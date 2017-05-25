package main

import (
	"flag"
	. "github.com/fogleman/pt/pt"
)

var width = flag.Int("w", 500, "width")
var height = flag.Int("h", 300, "height")
var spp = flag.Int("spp", 1, "spp")
var interactions = flag.Int("interactions", 1024, "interactions")
var pathTemplate = flag.String("path", "exemple1-%04d.npy", "")

func main() {
	flag.Parse()
	scene := Scene{}
	scene.Add(NewSphere(V(1.5, 1.25, 0), 1.25, SpecularMaterial(HexColor(0x004358), 1.3)))
	scene.Add(NewSphere(V(-1, 1, 2), 1, SpecularMaterial(HexColor(0xFFE11A), 1.3)))
	scene.Add(NewSphere(V(-2.5, 0.75, 0), 0.75, SpecularMaterial(HexColor(0xFD7400), 1.3)))
	scene.Add(NewSphere(V(-0.75, 0.5, -1), 0.5, ClearMaterial(1.5, 0)))
	scene.Add(NewCube(V(-10, -1, -10), V(10, 0, 10), GlossyMaterial(White, 1.1, Radians(10))))
	scene.Add(NewSphere(V(-1.5, 4, 0), 0.5, LightMaterial(White, 30)))
	camera := LookAt(V(0, 2, -5), V(0, 0.25, 3), V(0, 1, 0), 45)
	camera.SetFocus(V(-0.75, 1, -1), 0.1)
	sampler := NewSampler(4, 8)
	sampler.SpecularMode = SpecularModeFirst
	renderer := NewRenderer(&scene, &camera, sampler, *width, *height, *spp)
	renderer.AdaptiveSamples = 32
	renderer.FireflySamples = 256

	//renderer.IterativeRender("out%03d.png", 1000)
	renderer.ExportFeatures(*pathTemplate, *interactions)
}
