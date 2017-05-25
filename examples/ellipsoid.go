package main

import (
	"flag"
	. "github.com/fogleman/pt/pt"
)

var width = flag.Int("w", 500, "width")
var height = flag.Int("h", 300, "height")
var spp = flag.Int("spp", 1, "spp")
var interactions = flag.Int("interactions", 1024, "interactions")
var pathTemplate = flag.String("path", "ellipsoid-%04d.npy", "")

func main() {
	flag.Parse()
	scene := Scene{}
	wall := GlossyMaterial(HexColor(0xFCFAE1), 1.333, Radians(30))
	scene.Add(NewSphere(Vector{10, 10, 10}, 2, LightMaterial(White, 50)))
	scene.Add(NewCube(Vector{-100, -100, -100}, Vector{-12, 100, 100}, wall))
	scene.Add(NewCube(Vector{-100, -100, -100}, Vector{100, -1, 100}, wall))
	material := GlossyMaterial(HexColor(0x167F39), 1.333, Radians(30))
	sphere := NewSphere(Vector{}, 1, material)
	for i := 0; i < 180; i += 30 {
		m := Identity()
		m = m.Scale(Vector{0.3, 1, 5})
		m = m.Rotate(Vector{0, 1, 0}, Radians(float64(i)))
		shape := NewTransformedShape(sphere, m)
		scene.Add(shape)
	}
	camera := LookAt(Vector{8, 8, 0}, Vector{1, 0, 0}, Vector{0, 1, 0}, 45)
	sampler := NewSampler(4, 4)
	renderer := NewRenderer(&scene, &camera, sampler, *width, *height, *spp)
	//renderer.IterativeRender("out%03d.png", 1000)
	renderer.ExportFeatures(*pathTemplate, *interactions)
}
