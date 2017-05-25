package main

import (
	"flag"
	. "github.com/fogleman/pt/pt"
)

var width = flag.Int("w", 500, "width")
var height = flag.Int("h", 300, "height")
var spp = flag.Int("spp", 1, "spp")
var interactions = flag.Int("interactions", 1024, "interactions")
var pathTemplate = flag.String("path", "hdri-%04d.npy", "")

func main() {
	flag.Parse()

	scene := Scene{}
	scene.Color = White
	scene.Texture = GetTexture("examples/courtyard_ccby/courtyard_8k.png")
	material := GlossyMaterial(White, 2, Radians(0))
	material.Texture = GetTexture("examples/checker.png")
	scene.Add(NewSphere(V(0, 0, 0), 1, material))
	scene.Add(NewSphere(V(-2.5, 0, 0), 1, material))
	scene.Add(NewSphere(V(2.5, 0, 0), 1, material))
	scene.Add(NewSphere(V(0, 0, -2.5), 1, material))
	scene.Add(NewSphere(V(0, 0, 2.5), 1, material))
	material = GlossyMaterial(HexColor(0xEFECCA), 1.1, Radians(45))
	scene.Add(NewCube(V(-100, -100, -100), V(100, -1, 100), material))
	camera := LookAt(V(2, 3, 4), V(0, 0, 0), V(0, 1, 0), 40)
	sampler := NewSampler(4, 4)
	renderer := NewRenderer(&scene, &camera, sampler, *width, *height, *spp)
	
	//renderer.IterativeRender("out%03d.png", 1000)
	renderer.ExportFeatures(*pathTemplate, *interactions)
}
