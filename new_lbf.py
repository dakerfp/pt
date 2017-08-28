
import tensorflow as tf
from keras import layers, models
import create_batches
import numpy as np
import itertools as it


class Model(object):
    def __init__(self, window_shape=(11,11,51), output_depth=3):
        with tf.name_scope("model") as scope:
            assert window_shape[0] == window_shape[1]
            assert len(window_shape) == 3

            width = window_shape[0]
            depth = window_shape[2]

            x = tf.placeholder(tf.float32, shape=(None, width, width, depth), name="x")
            y_ = tf.placeholder(tf.float32, shape=(None, 3), name="y_")

            xnorm = None
            with tf.name_scope("normalize"):
                xmax = tf.reduce_max(x, axis=(1,2), name="xmax")
                xmin = tf.reduce_min(x, axis=(1,2), name="xmin")
                xnorm = (x - xmin) / (xmax - xmin)
            ccolor = xnorm[:,(width-1)/2,(width-1)/2,:3]

            # xnorm = None
            model = models.Sequential()
            model.add(layers.InputLayer(input_tensor=xnorm))
            model.add(layers.Reshape((width * width * depth,)))
            model.add(layers.Dense(64, activation='sigmoid'))
            model.add(layers.Dense(16, activation='sigmoid'))
            model.add(layers.Dense(2, activation='softplus'))

            def pdiff(width, alpha):
                diff2 = np.zeros((width,width))
                c = (width-1)/2
                for i, j in it.product(range(width), range(width)):
                    diff2[i,j] = (i-c)**2 + (j-c)**2

                return tf.exp(-diff2/(2*alpha**2))

            def D(color, center_color, beta):
                cdiff2 = tf.reduce_sum((center_color - color) ** 2, axis=3) ## FIX COLOR DIFF WITH VARIANCE
                return -cdiff2 / (2 * beta ** 2)

            def Dk(f, f_center, gamma):
                tdiff2 = tf.reduce_sum((f - f_center) ** 2, axis=3) ## FIX COLOR DIFF WITH VARIANCE
                return -tdiff2 / (2 * gamma ** 2)

            alpha = model.output[:,0]
            beta = model.output[:,1]

            c = xnorm[:,:,:,:3]
            d = tf.exp(
                pdiff(width, alpha)
                + D(c, ccolor, beta)
            )
            d3 = tf.stack((d,d,d), axis=3)

            y = tf.reduce_sum(d3 * c, axis=(1,2)) \
                / tf.reduce_sum(d, axis=(1,2))

            mse = tf.reduce_sum((y - y_)**2)
            opt = tf.train.RMSPropOptimizer(0.001)
            train = opt.minimize(mse)

            self.x = x
            self.y_ = y_
            self.y = y
            self.max = xmax
            self.model = model
            self.err = mse
            self.train = train

def prod(xs):
    p = 1
    for x in xs:
        if x is not None:
            p *= x
    return p


def apply_filter(flt, scene, sess):
    h, w = scene.shape_windows()
    arr = np.zeros((h,w,3))
    for i, j, win in scene.iter_windows():
        # pixels = xs # flt.x.eval(feed_dict={flt.x: xs})
        win = np.array([win])
        pixel = flt.y.eval(session=sess, feed_dict={flt.x: win})
        color = pixel[0]
        arr[i,j,0] = color[0]
        arr[i,j,1] = color[1]
        arr[i,j,2] = color[2]

    return arr

def show_images(flt, dataset, sess):
    scene = dataset.next_scene()
    return apply_filter(flt, scene, sess)


def train(sess, model, dataset, epochs=100, batch_size=100):
    errs = []
    sess.run(tf.global_variables_initializer())
    for i in xrange(epochs):
        xs, ys = zip(*dataset.next_batch(batch_size))
        xs = np.array(xs)
        ys = np.array(ys)
        print(xs.shape, ys.shape)
        _, err = sess.run((model.train, model.err), feed_dict={model.x: xs, model.y_: ys})
        errs.append(err)
        if i % 10 == 0:
            print(i)
    return err


if __name__ == '__main__':
    kwidth=11
    depth=52
    m = Model((kwidth, kwidth, depth))

    import sys
    dataset = None
    if sys.argv[1].endswith(".zip"):
        dataset = create_batches.ZipDataset(sys.argv[1:], prefixes=[3, 4, 5, 6], lowres=32, hires=1024, kernel_size=kwidth, depth=52)
    else:
        npys = zip(sys.argv[1::2], sys.argv[2::2])
        dataset = create_batches.Dataset(npys, 25)

    errs = []
    import matplotlib.pyplot as plt
    scene = dataset.next_scene()
    plt.subplot(2, 2, 1)
    plt.imshow(np.clip(scene.arr[:,:,:3], 0, 1))
    with tf.Session() as sess:
        errs = train(sess, m, dataset)
        plt.subplot(2,2,2)
        img = apply_filter(m, dataset.next_scene(), sess)
        plt.imshow(img)
    plt.subplot(2, 2, 3)
    plt.imshow(np.clip(scene.gt[:,:,:3], 0, 1))
    plt.subplot(2, 2, 4)
    plt.plot(errs, 0, 4)
    plt.show()

