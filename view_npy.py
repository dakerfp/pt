
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

# indirect
col2 = arr[:,:,38:41]
wasCast2 = arr[:,:,41]
dist2 = arr[:,:,42]
spec2 = arr[:,:,43]
diff2 = arr[:,:,44]

col2_var = arr[:,:,45:48]
wasCast2_var = arr[:,:,48]
dist2_var = arr[:,:,49]
spec2_var = arr[:,:,50]
diff2_var = arr[:,:,51]

import matplotlib.pyplot as plt

def imgshow_all(*args):
	fig = plt.figure()
	cols=6
	rows=6
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

	col2,
	wasCast2,
	dist2,
	spec2,
	diff2,

	col2_var,
	wasCast2_var,
	dist2_var,
	spec2_var,
	diff2_var
)

plt.show()

