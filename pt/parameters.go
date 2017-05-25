package pt

import "time"

func ParameterTest(scene *Scene, camera *Camera, w, h, spp int, duration time.Duration) {
	var sampler *DefaultSampler
	var renderer *Renderer

	sampler = NewSampler(1, 4)
	renderer = NewRenderer(scene, camera, sampler, w, h, spp)
	SavePNG("1.Default.png", renderer.TimedRender(duration))

	sampler = NewSampler(4, 4)
	renderer = NewRenderer(scene, camera, sampler, w, h, spp)
	SavePNG("2.StratifiedFirstHit.png", renderer.TimedRender(duration))

	sampler = NewSampler(4, 4)
	sampler.LightMode = LightModeAll
	renderer = NewRenderer(scene, camera, sampler, w, h, spp)
	SavePNG("3.LightModeAll.png", renderer.TimedRender(duration))

	sampler = NewSampler(4, 4)
	sampler.SpecularMode = SpecularModeFirst
	renderer = NewRenderer(scene, camera, sampler, w, h, spp)
	SavePNG("4.SpecularModeFirst.png", renderer.TimedRender(duration))

	sampler = NewSampler(4, 4)
	sampler.SpecularMode = SpecularModeAll
	renderer = NewRenderer(scene, camera, sampler, w, h, spp)
	SavePNG("5.SpecularModeAll.png", renderer.TimedRender(duration))

	sampler = NewSampler(4, 4)
	renderer = NewRenderer(scene, camera, sampler, w, h, spp)
	renderer.AdaptiveSamples = 16
	SavePNG("6.AdaptiveSamples.png", renderer.TimedRender(duration))

	sampler = NewSampler(4, 4)
	renderer = NewRenderer(scene, camera, sampler, w, h, spp)
	renderer.FireflySamples = 256
	SavePNG("7.FireflySamples.png", renderer.TimedRender(duration))

	sampler = NewSampler(4, 4)
	sampler.LightMode = LightModeAll
	sampler.SpecularMode = SpecularModeFirst
	renderer = NewRenderer(scene, camera, sampler, w, h, spp)
	renderer.AdaptiveSamples = 16
	renderer.FireflySamples = 256
	SavePNG("8.Everything.png", renderer.TimedRender(duration))
}
