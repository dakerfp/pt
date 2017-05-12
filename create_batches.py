
import random
import numpy as np


class Scene(object):
	def __init__(self, filename, kernel_size=11):
		self.arr = np.load(filename)
		self.kernel_size = kernel_size

	def next(self):
		w, h, d = self.arr.shape
		k = self.kernel_size
		x = random.randint(0, w - k - 1)
		y = random.randint(0, h - k - 1)
		return self.arr[y:y+k,x:x+k,:]


class Dataset(object):
	def __init__(self, filenames, kernel_size=11):
		self.scenes = [Scene(fn, kernel_size) for fn in filenames]

	def next(self):
		s = random.choice(self.scenes)
		return s.next()

	def next_batch(self, n):
		return [self.next() for _ in range(n)]


if __name__ == '__main__':
	import matplotlib.pyplot as plt
	import sys

	npys = sys.argv[1:]

	dset = Dataset(sys.argv[1:], 25)

	fig = plt.figure()
	for i in range(1,11):
		arr = dset.next()
		col = arr[:,:,:3]
		col = 255 * np.clip(col, 0, 1.0)
		print(arr.shape)
		norm = arr[:,:,3:6]
		fig.add_subplot(2,10,i)
		plt.imshow(col)
		fig.add_subplot(1,10,i)
		plt.imshow(norm)

	plt.show()
	sys.exit()

"""
# primary features
dist = arr[:,:,6]
alb = arr[:,:,7:10]

# sec features
col_var = arr[:,:,10:13]
norm_var = arr[:,:,13:16]
dist_var = arr[:,:,16]
"""
