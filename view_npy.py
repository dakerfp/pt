
import sys
import numpy as np

arr = np.load(sys.argv[1])

# primary features
col = arr[:,:,:3]
norm = arr[:,:,3:6]
dist = arr[:,:,6]
alb = arr[:,:,7:10]

# sec features
col_var = arr[:,:,10:13]
norm_var = arr[:,:,13:16]
dist_var = arr[:,:,16]

import matplotlib.pyplot as plt

fig = plt.figure()
fig.add_subplot(4,3,1)
plt.imshow(col)
fig.add_subplot(4,3,2)
plt.imshow(norm)
fig.add_subplot(4,3,3)
plt.imshow(dist)
fig.add_subplot(4,3,4)
plt.imshow(alb)
fig.add_subplot(4,3,5)
plt.imshow(col_var)
fig.add_subplot(4,3,6)
plt.imshow(norm_var)
fig.add_subplot(4,3,7)
plt.imshow(dist_var)
fig.add_subplot(4,3,8)
plt.imshow(np.clip(col ** (1 / 2.2), 0, 255))
plt.show()

