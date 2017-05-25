
import sys
import numpy as np

arr = np.load(sys.argv[1])

# primary features
col = arr[:,:,:3]
norm = arr[:,:,3:6]
dist = arr[:,:,6]
alb = arr[:,:,7:10]
gamma = arr[:,:,10:13]
bump = arr[:,:,13]
emit = arr[:,:,14]
idx = arr[:,:,15]
gloss = arr[:,:,16]
tint = arr[:,:,17]
refl = arr[:,:,18]
transp = arr[:,:,19]

# sec features
col_var = arr[:,:,20:23]
norm_var = arr[:,:,23:26]
dist_var = arr[:,:,26]
alb_var = arr[:,:,27:30]
bump_var = arr[:,:,30]
emit_var = arr[:,:,31]
idx_var = arr[:,:,32]
gloss_var = arr[:,:,33]
tint_var = arr[:,:,34]
refl_var = arr[:,:,35]
transp_var = arr[:,:,36]

# other
hits = arr[:,:,37]

import matplotlib.pyplot as plt

def imgshow_all(*args):
	fig = plt.figure()
	cols=5
	rows=5
	for i, img in enumerate(args):
		fig.add_subplot(rows, cols, i+1)
		plt.imshow(img, interpolation='nearest')

# imgshow_all(col, norm, dist)
imgshow_all(
	col,
	norm,
	dist,
	alb,
	gamma,
	bump,
	emit,
	idx,
	gloss,
	tint,
	refl,
	transp,

	col_var,
	norm_var,
	dist_var,
	bump_var,
	emit_var,
	idx_var,
	gloss_var,
	tint_var,
	refl_var,
	transp_var,
	hits,
)

plt.show()

