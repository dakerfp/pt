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
	floor := GlossyMaterial(HexColor(0xFCFFF5), 1.2, Radians(30))
	cube := GlossyMaterial(HexColor(0xFF8C00), 1.3, Radians(20))
	ball := GlossyMaterial(HexColor(0xD90000), 1.4, Radians(10))
	n := 7
	fn := float64(n)
	for z := 0; z < n; z++ {
		for x := 0; x < n-z; x++ {
			for y := 0; y < n-z-x; y++ {
				fx, fy, fz := float64(x), float64(y), float64(z)
				scene.Add(NewCube(V(fx, fy, fz), V(fx+1, fy+1, fz+1), cube))
				if x+y == n-z-1 {
					if rand.Float64() > 0.75 {
						scene.Add(NewSphere(V(fx+0.5, fy+0.5, fz+1.5), 0.35, ball))
					}
				}
			}
		}
	}
	scene.Add(NewCube(V(-1000, -1000, -1), V(1000, 1000, 0), floor))
	scene.Add(NewSphere(V(fn, fn/3, fn*2), 1, LightMaterial(White, 100)))
	camera := LookAt(V(fn*2, fn*2, fn*2), V(0, 0, fn/4), V(0, 0, 1), 35)
	sampler := NewSampler(4, 4)
	renderer := NewRenderer(&scene, &camera, sampler, *width, *height, *spp)
	
	//renderer.IterativeRender("out%03d.png", 1000)
	renderer.ExportFeatures(*pathTemplate, *interactions)
}
