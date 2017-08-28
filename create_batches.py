
import random
import numpy as np
import zipfile
import io
import itertools as itt

class Scene(object):
	def __init__(self, arr, gt, kernel_size=11, depth=52):
		self.arr = arr
		self.gt = gt
		self.kernel_size = kernel_size
		self.depth = depth

	def next(self):
		h, w, d = self.arr.shape
		k = self.kernel_size
		x = random.randint(0, w - k - 1)
		y = random.randint(0, h - k - 1)
		return self.arr[y:y+k,x:x+k,:self.depth], self.gt[y+k/2,x+k/2,10:13]

	def windows(self):
		h, w, d = self.arr.shape
		k = self.kernel_size
		return [[self.arr[y:y+k,x:x+k,:] for x in range(w-k)]
			for y in range(h-k)]

	def flat_windows(self):
		h, w, d = self.arr.shape
		k = self.kernel_size
		return [self.arr[y:y+k,x:x+k,:] for x, y in itt.product(range(w-k),range(h-k))]

	def iter_windows(self):
		h, w, d = self.arr.shape
		k = self.kernel_size
		return ((y, x, self.arr[y:y+k,x:x+k,:]) for x, y in itt.product(range(w-k),range(h-k)))

	def shape_windows(self):
		h, w, _ = self.arr.shape
		k = self.kernel_size
		return h-k, w-k

	def color(self):
		return self.arr[:,:,10:3]

	def gt_color(self):
		# col = self.gt[:,:,:3]
		col = self.gt[:,:,10:13]
		return np.clip(col, 0, 1)

class Dataset(object):
	def __init__(self, filetuples, kernel_size=11, depth=52):
		self.scenes = [Scene(np.load(fn), np.load(gt), kernel_size, depth) for (fn, gt) in filetuples]
		self.depth = depth

	def next(self):
		s = random.choice(self.scenes)
		return s.next()

	def next_batch(self, n):
		return [self.next() for _ in xrange(n)]


class ZipDataset(object):
	def __init__(self, zipfilenames, prefixes=[1], lowres=16, hires=1024, kernel_size=11, depth=52):
		def load(zf, prefix, spp):
			template = '%d-%04d.npy'
			buff = io.BufferedReader(zf.open(template % (prefix, spp)))
			return np.load(buff)

		zipfiles = [zipfile.ZipFile(fn) for fn in zipfilenames]

		filetuples = []
		for pre, zf in itt.product(prefixes, zipfiles):
			try:
		 		tupl = load(zf, pre, lowres), load(zf, pre, hires)
		 		filetuples.append(tupl)
		 	except Exception as e:
		 		print(e, "with", zf)
		self.scenes = [Scene(fn, gt, kernel_size, depth) for (fn, gt) in filetuples]
		self.depth = depth

	def next(self):
		s = random.choice(self.scenes)
		return s.next()

	def next_batch(self, n):
		return [self.next() for _ in xrange(n)]

	def next_scene(self):
		return random.choice(self.scenes)


if __name__ == '__main__':
	import matplotlib.pyplot as plt
	import sys

	dset = None
	if sys.argv[1].endswith(".zip"):
		dset = ZipDataset(sys.argv[1:], prefixes=[1, 2, 3, 4, 5, 6])
	else:
		npys = zip(sys.argv[1::2], sys.argv[2::2])
		dset = Dataset(npys, 25)

	fig = plt.figure()
	for i in range(1,11):
		arr, _ = dset.next()
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
