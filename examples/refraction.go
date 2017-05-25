package main

import (
	"flag"
	. "github.com/fogleman/pt/pt"
)

var width = flag.Int("w", 500, "width")
var height = flag.Int("h", 300, "height")
var spp = flag.Int("spp", 1, "spp")
var interactions = flag.Int("interactions", 1024, "interactions")
var pathTemplate = flag.String("path", "refraction-%04d.npy", "")

func main() {
	flag.Parse()

	scene := Scene{}

	glass := ClearMaterial(1.5, 0)

	// add a sphere primitive
	scene.Add(NewSphere(V(-1.5, 0, 0.5), 1, glass))

	// add a mesh sphere
	mesh, err := LoadSTL("examples/sphere.stl", glass)
	if err != nil {
		panic(err)
	}
	mesh.SaveSTL("examples/sphere2.stl")
	mesh.SmoothNormals()
	mesh.Transform(Translate(V(1.5, 0, 0.5)))
	scene.Add(mesh)

	// add the floor
	scene.Add(NewPlane(V(0, 0, -1), V(0, 0, 1), DiffuseMaterial(White)))

	// add the light
	scene.Add(NewSphere(V(0, 0, 5), 1, LightMaterial(White, 30)))

	camera := LookAt(V(0, -5, 5), V(0, 0, 0), V(0, 0, 1), 50)
	sampler := NewSampler(16, 8)
	sampler.SpecularMode = SpecularModeAll
	renderer := NewRenderer(&scene, &camera, sampler, *width, *height, *spp)
	
	//renderer.IterativeRender("out%03d.png", 1000)
	renderer.ExportFeatures(*pathTemplate, *interactions)
}
