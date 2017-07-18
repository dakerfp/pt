
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
            layer1_size = input_size * 8
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
            layer3 = tf.sigmoid(tf.matmul(layer2, W3) + b3)

            return layer3
            # return tf.reshape(layer3, shape=(None,window_width * window_width * 3))

        # def bilateral_filter_window(img, w):
        #     return tf.reduce_sum(img * w, (0,1)) / tf.reduce_sum(w)

        x = tf.placeholder(tf.float32, [None, width * width * depth])
        xcol = tf.placeholder(tf.float32, [None, width * width * 3])
        y_ = tf.placeholder(tf.float32, [None, 3])
        w = filter_weights(x, depth=depth, window_width=width)
        # y = bilateral_filter_window(xcol, w)
        W4 = tf.Variable(tf.random_uniform((width * width * 3, 3)))
        b4 = tf.Variable(tf.random_uniform((3,)))
        y = tf.matmul(w, W4) + b4 # no sigmoid on last layer

        eps = tf.constant(1e-8)
        # relmse = tf.reduce_sum(tf.square(y - y_)/(tf.square(y_) + eps))
        mse = tf.reduce_sum(tf.square(y - y_))
        # train_step = tf.train.AdamOptimizer(1e-4).minimize(relmse)
        train_step = tf.train.AdamOptimizer(1e-3).minimize(mse)

        self.x = x
        self.xcol = xcol
        self.y_ = y_
        self.y = y
        self.train_step = train_step
        self.mse = mse

    def run_epoch(self, dataset, epoch_size=100):
        xs, ys_ = zip(*dataset.next_batch(epoch_size))
        flat_xs, flat_xcols, ys_ = flatten_xy(xs, ys_)
        self.train_step.run(feed_dict={
            self.x: flat_xs,
            self.xcol: flat_xcols,
            self.y_: ys_})

    def test_model(self, dataset, instances=10):
        xs, xcols, ys_ = get_batch(dataset, instances)
        errsum = self.mse.eval(feed_dict={
            self.x: xs,
            self.xcol: xcols,
            self.y_: ys_})
        return errsum / instances


    def filter_scene(self, scene):
        wins = scene.flat_windows()
        print(len(wins), wins[0].shape)

        xs, xcols = flatten_x(wins)
        pixels = self.y.eval(feed_dict={
            self.x: xs,
            self.xcol: xcols})


        # reassemble image
        print(len(pixels), pixels.shape)
        w, h = scene.shape_windows()
        print(w, h)
        img = []
        i = 0
        for _ in range(h):
            img.append([pixels[i:i+w] for _ in range(w)])
            i += w


        arr = np.stack(img)
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


if __name__ == '__main__':
    kwidth=11

    import sys
    dataset = None
    if sys.argv[1].endswith(".zip"):
        dataset = create_batches.ZipDataset(sys.argv[1:], prefixes=[1, 2, 3, 4, 5, 6], lowres=16, hires=1024, kernel_size=kwidth)
    else:
        npys = zip(sys.argv[1::2], sys.argv[2::2])
        dataset = create_batches.Dataset(npys, 25)

    depth = dataset.next()[0].shape[2]

    x, xcol, y_, y, train_step, err = create_network(width=kwidth,depth=depth)
    saver = tf.train.Saver()

    sess = tf.InteractiveSession()

    sess.run(tf.global_variables_initializer())

    errs = []
    for epoch in range(101):
        run_epoch(train_step, dataset)
        e = test_model(x, xcol, y_, err, dataset)
        print("epoch:", epoch, e)
        errs.append(e)
        if epoch % 100 == 0:
            saver.save(sess, 'lbf-basic', global_step=epoch)

    import matplotlib.pyplot as plt
    fig = plt.figure()
    smooth = lambda xs, w: [sum(xs[i:i+w]) / w for i in range(len(xs)-w)]
    fig.add_subplot(2,2,1)
    # plt.plot(errs)
    plt.plot(smooth(errs, 20))
    plt.plot(smooth(errs, 100))
    plt.plot(smooth(errs, 500))
    plt.plot(smooth(errs, 1000))
    scene = dataset.scenes[1]
    img = filter_scene(y, scene)
    fig.add_subplot(2,2,2)
    plt.imshow(np.clip(img, 0, 1))
    fig.add_subplot(2,2,3)
    plt.imshow(np.clip(scene.color(), 0, 1))
    fig.add_subplot(2,2,4)
    plt.imshow(np.clip(scene.gt_color(), 0, 1))
    plt.show()
