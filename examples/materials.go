package main

import (
	"flag"
	. "github.com/fogleman/pt/pt"
)

var width = flag.Int("w", 500, "width")
var height = flag.Int("h", 300, "height")
var spp = flag.Int("spp", 1, "spp")
var interactions = flag.Int("interactions", 1024, "interactions")
var pathTemplate = flag.String("path", "hits-%04d.npy", "")

func main() {
	flag.Parse()

	scene := Scene{}
	r := 0.4
	var material Material

	material = DiffuseMaterial(HexColor(0x334D5C))
	scene.Add(NewSphere(V(-2, r, 0), r, material))

	material = SpecularMaterial(HexColor(0x334D5C), 2)
	scene.Add(NewSphere(V(-1, r, 0), r, material))

	material = GlossyMaterial(HexColor(0x334D5C), 2, Radians(50))
	scene.Add(NewSphere(V(0, r, 0), r, material))

	material = TransparentMaterial(HexColor(0x334D5C), 2, Radians(20), 1)
	scene.Add(NewSphere(V(1, r, 0), r, material))

	material = ClearMaterial(2, 0)
	scene.Add(NewSphere(V(2, r, 0), r, material))

	material = MetallicMaterial(HexColor(0xFFFFFF), 0, 1)
	scene.Add(NewSphere(V(0, 1.5, -4), 1.5, material))

	scene.Add(NewCube(V(-1000, -1, -1000), V(1000, 0, 1000), GlossyMaterial(HexColor(0xFFFFFF), 1.4, Radians(20))))
	scene.Add(NewSphere(V(0, 5, 0), 1, LightMaterial(White, 25)))
	camera := LookAt(V(0, 3, 6), V(0, 1, 0), V(0, 1, 0), 30)
	sampler := NewSampler(16, 16)
	renderer := NewRenderer(&scene, &camera, sampler, *width, *height, *spp)
	
	//renderer.IterativeRender("out%03d.png", 1000)
	renderer.ExportFeatures(*pathTemplate, *interactions)
}
