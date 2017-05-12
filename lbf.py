import sys
import tensorflow as tf
import numpy as np

import scipy.stats as st

def gaussian_kernel(kernlen=11, nsig=3):
    """Returns a 2D Gaussian kernel array."""
    interval = (2*nsig+1.)/(kernlen)
    x = np.linspace(-nsig-interval/2., nsig+interval/2., kernlen+1)
    kern1d = np.diff(st.norm.cdf(x))
    kernel_raw = np.sqrt(np.outer(kern1d, kern1d))
    kernel = kernel_raw/kernel_raw.sum()
    return kernel

tf.InteractiveSession()

def conv2d(x, W):
  return tf.nn.conv2d(x, W, strides=[1, 1, 1, 1], padding='SAME')

samples = np.load(sys.argv[1])
n_instances, height, width, feat_size = samples.shape

x = tf.placeholder(tf.float32, [None, height, width, feat_size])
y_ = tf.placeholder(tf.float32, [3])

def network(x):
    w_conv1 = tf.Variable(tf.random_uniform((5, 5, feat_size, 32)))
    b_conv1 = tf.Variable(tf.random_uniform((32,)))

    w_conv2 = tf.Variable(tf.random_uniform((5, 5, 32, 32)))
    b_conv2 = tf.Variable(tf.random_uniform((32,)))

    w_conv3 = tf.Variable(tf.random_uniform((5, 5, 32, 3)))
    b_conv3 = tf.Variable(tf.random_uniform((3,)))

    layer1 = tf.sigmoid(conv2d(x, w_conv1) + b_conv1)
    layer2 = tf.sigmoid(conv2d(layer1, w_conv2) + b_conv2)
    layer3 = tf.sigmoid(conv2d(layer2, w_conv3) + b_conv3)

    return layer3

def bilateral_filter_window(i, w):
    g = gaussian_kernel(i.shape[0])
    g = np.stack((g,g,g), axis=2)
    w = g * w
    return tf.reduce_sum(i * w, axis=[0,1]) / tf.reduce_sum(w)

w = network(x)
y = bilateral_filter_window(x, w)

mse = tf.reduce_sum(tf.square(y - y_))
train_step = tf.train.AdamOptimizer(1e-4).minimize(mse)

sess = tf.InteractiveSession()
sess.run(tf.initialize_all_variables())

bsize = 16
frame = np.average(samples[:,:,:,:3], axis=0)
expected = np.repeat(np.expand_dims(frame, axis=0), bsize, axis=0)

print(expected.shape)
for i in range(n_instances / bsize):
    begin = i * bsize
    end = begin + bsize
    train_step.run(feed_dict={x: samples[begin:end,:,:,:], y_: expected})
