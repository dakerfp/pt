package main

import (
	"flag"
	"log"
	"math/rand"

	. "github.com/fogleman/pt/pt"
)

var width = flag.Int("w", 500, "")
var height = flag.Int("h", 300, "")
var spp = flag.Int("spp", 1024, "")
var pathTemplate = flag.String("path", "cube-%04d.npy", "")

func createMesh(material Material) Shape {
	mesh, err := LoadSTL("examples/cube.stl", material)
	if err != nil {
		log.Fatalln("LoadSTL error:", err)
	}
	mesh.FitInside(Box{Vector{0, 0, 0}, Vector{1, 1, 1}}, Vector{0.5, 0.5, 0.5})
	return mesh
}

func main() {
	flag.Parse()
	scene := Scene{}
	meshes := []Shape{
		createMesh(GlossyMaterial(HexColor(0x3B596A), 1.5, Radians(20))),
		createMesh(GlossyMaterial(HexColor(0x427676), 1.5, Radians(20))),
		createMesh(GlossyMaterial(HexColor(0x3F9A82), 1.5, Radians(20))),
		createMesh(GlossyMaterial(HexColor(0xA1CD73), 1.5, Radians(20))),
		createMesh(GlossyMaterial(HexColor(0xECDB60), 1.5, Radians(20))),
	}
	for x := -8; x <= 8; x++ {
		for z := -12; z <= 12; z++ {
			fx := float64(x)
			fy := rand.Float64() * 2
			fz := float64(z)
			scene.Add(NewTransformedShape(meshes[rand.Intn(len(meshes))], Translate(Vector{fx, fy, fz})))
			scene.Add(NewTransformedShape(meshes[rand.Intn(len(meshes))], Translate(Vector{fx, fy - 1, fz})))
		}
	}
	scene.Add(NewSphere(Vector{8, 10, 0}, 3, LightMaterial(White, 30)))
	camera := LookAt(Vector{-10, 10, 0}, Vector{-2, 0, 0}, Vector{0, 1, 0}, 45)
	sampler := NewSampler(4, 4)
	renderer := NewRenderer(&scene, &camera, sampler, *width, *height)
	// renderer.IterativeRender("cube-out%03d.png", 100)
	renderer.ExportFeatures(*pathTemplate, *spp)
}
