
import tensorflow as tf
import numpy as np
import create_batches

def prod(xs):
    p = 1
    for x in xs:
        if x is not None:
            p *= x
    return p

class LearningBasedFilter(object):

    def __init__(self, width=11, depth=37):
        def shape(tensor):
            return [d.value for d in tensor.get_shape()]

        def filter_weights(x, depth, window_width):
            input_size = prod(shape(x))
            layer1_size = input_size / 2
            layer2_size = input_size / 4
            output_size = window_width * window_width * 3

            W1 = tf.Variable(tf.random_uniform((input_size, layer1_size)))
            b1 = tf.Variable(tf.random_uniform((layer1_size,)))
            layer1 = tf.sigmoid(tf.matmul(x, W1) + b1)

            W2 = tf.Variable(tf.random_uniform((layer1_size, layer2_size)))
            b2 = tf.Variable(tf.random_uniform((layer2_size,)))
            layer2 = tf.sigmoid(tf.matmul(layer1, W2) + b2)

            W3 = tf.Variable(tf.random_uniform((layer2_size, output_size)))
            b3 = tf.Variable(tf.random_uniform((output_size,)))
            layer3 = tf.nn.softplus(tf.matmul(layer2, W3) + b3)
            w = layer3
            w = w / tf.reduce_sum(w)
            return w
            # return tf.reshape(layer3, shape=(None,window_width * window_width * 3))

        x = tf.placeholder(tf.float32, [None, width * width * depth])
        xcol = tf.placeholder(tf.float32, [None, width * width * 3])
        y_ = tf.placeholder(tf.float32, [None, 3])

        w = filter_weights(x, depth, width)
        # g = tf.Variable(tf.random_uniform((width * width * 3, 3), minval=0, maxval=1)) # learn geometric relationship and how they add to final color
        # relu to ensure it is positive
        # y = tf.matmul(xcol * w, g) # bilateral filter
        # print(w, g, y)
        sumcol = xcol * w
        sw = tf.reduce_sum(w)
        r = tf.reduce_sum(sumcol[:,0::3], axis=1)
        g = tf.reduce_sum(sumcol[:,1::3], axis=1)
        b = tf.reduce_sum(sumcol[:,2::3], axis=1)
        y = tf.stack((r,g,b), axis=1) / sw
        print(xcol, w, r, g, b, y)

        def relmse(y, y_):
            eps = tf.constant(1e-8)
            return tf.reduce_sum(
                tf.reduce_sum(tf.square(y - y_)/
                    (tf.square(y_) + eps), axis=1)
            )

        err = relmse(y, y_)
        train_step = tf.train.RMSPropOptimizer(1e-3, use_locking=True, centered=True).minimize(err)

        self.x = x
        self.xcol = xcol
        self.y_ = y_
        self.y = y
        self.train_step = train_step
        self.mse = err

    def run_epoch(self, dataset, epoch_size=100):
        xs, ys_ = zip(*dataset.next_batch(epoch_size))
        flat_xs, flat_xcols, ys_ = flatten_xy(xs, ys_)
        self.train_step.run(feed_dict={
            self.x: flat_xs,
            self.xcol: flat_xcols,
            self.y_: ys_})

    def test_model(self, sess, dataset, instances=100):
        xs, xcols, ys_ = get_batch(dataset, instances)
        err = sess.run(self.mse, feed_dict={
            self.x: xs,
            self.xcol: xcols,
            self.y_: ys_
        })
        return err / instances


    def filter_scene(self, scene):
        wins = scene.flat_windows()
        print(len(wins), wins[0].shape)

        xs, xcols = flatten_x(wins)
        pixels = self.y.eval(feed_dict={
            self.x: xs,
            self.xcol: xcols})


        # reassemble image
        print(self.y)
        print(len(pixels), pixels.shape)
        print(pixels)

        w, h = scene.shape_windows()
        arr = np.reshape(pixels, (h,w,3))
        print(arr.shape)
        return arr


def get_batch(dataset, instances=10):
    xs, ys_ = zip(*dataset.next_batch(instances))
    xcols = [xi[:,:,10:13] for xi in xs]

    xsize = prod(xs[0].shape)
    flat_xs = np.stack([np.reshape(xi, (xsize,)) for xi in xs], axis=0)

    xcsize = prod(xcols[0].shape)
    flat_xcols = np.stack([np.reshape(xi, (xcsize,)) for xi in xcols])

    return flat_xs, flat_xcols, ys_

def flatten_x(xs):
    xcols = [xi[:,:,10:13] for xi in xs]

    xsize = int(prod(xs[0].shape))
    flat_xs = np.stack([np.reshape(xi, (xsize,)) for xi in xs], axis=0)

    xcsize = prod(xcols[0].shape)
    flat_xcols = np.stack([np.reshape(xi, (xcsize,)) for xi in xcols])

    return flat_xs, flat_xcols


def flatten_xy(xs, ys_):
    xcols = [xi[:,:,10:13] for xi in xs]

    xsize = int(prod(xs[0].shape))
    flat_xs = np.stack([np.reshape(xi, (xsize,)) for xi in xs], axis=0)

    xcsize = prod(xcols[0].shape)
    flat_xcols = np.stack([np.reshape(xi, (xcsize,)) for xi in xcols])

    return flat_xs, flat_xcols, ys_


# if __name__ == '__main__':
#     kwidth=11

#     import sys
#     dataset = None
#     if sys.argv[1].endswith(".zip"):
#         dataset = create_batches.ZipDataset(sys.argv[1:], prefixes=[1, 2, 3, 4, 5, 6], lowres=16, hires=1024, kernel_size=kwidth)
#     else:
#         npys = zip(sys.argv[1::2], sys.argv[2::2])
#         dataset = create_batches.Dataset(npys, 25)

#     depth = dataset.next()[0].shape[2]

#     x, xcol, y_, y, train_step, err = create_network(width=kwidth,depth=depth)
#     saver = tf.train.Saver()

#     sess = tf.InteractiveSession()

#     sess.run(tf.global_variables_initializer())

#     errs = []
#     for epoch in range(201):
#         run_epoch(train_step, dataset, epoch_size=500)
#         e = test_model(x, xcol, y_, err, dataset)
#         print("epoch:", epoch, e)
#         errs.append(e)
#         if epoch > 10 and epoch % 1000 == 0:
#             saver.save(sess, 'lbf-basic', global_step=epoch)

#     import matplotlib.pyplot as plt
#     fig = plt.figure()
#     smooth = lambda xs, w: [sum(xs[i:i+w]) / w for i in range(len(xs)-w)]
#     fig.add_subplot(2,2,1)
#     # plt.plot(errs)
#     plt.plot(smooth(errs, 20))
#     plt.plot(smooth(errs, 100))
#     plt.plot(smooth(errs, 500))
#     plt.plot(smooth(errs, 1000))
#     scene = dataset.scenes[1]
#     img = filter_scene(y, scene)
#     fig.add_subplot(2,2,2)
#     plt.imshow(np.clip(img, 0, 1))
#     fig.add_subplot(2,2,3)
#     plt.imshow(np.clip(scene.color(), 0, 1))
#     fig.add_subplot(2,2,4)
#     plt.imshow(np.clip(scene.gt_color(), 0, 1))
#     plt.show()
