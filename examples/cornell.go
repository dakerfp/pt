package main

import (
	"flag"
	. "github.com/fogleman/pt/pt"
)

var width = flag.Int("w", 500, "width")
var height = flag.Int("h", 300, "height")
var spp = flag.Int("spp", 1, "spp")
var interactions = flag.Int("interactions", 1024, "interactions")
var pathTemplate = flag.String("path", "cornell-%04d.npy", "")

func main() {
	flag.Parse()

	scene := Scene{}
	material := DiffuseMaterial(Black)
	mesh, err := LoadOBJ("examples/CornellBox-Original.obj", material)
	if err != nil {
		panic(err)
	}
	for _, t := range mesh.Triangles {
		scene.Add(t)
	}
	camera := LookAt(V(0, 1, 3), V(0, 1, 0), V(0, 1, 0), 50)
	sampler := NewSampler(4, 8)
	renderer := NewRenderer(&scene, &camera, sampler, *width, *height, *spp)
	
	//renderer.IterativeRender("out%03d.png", 1000)
	renderer.ExportFeatures(*pathTemplate, *interactions)
}
