
import tensorflow as tf
import numpy as np
import create_batches
import scipy.stats as st    


def create_network(width=11, feat_size=17):
    def shape(tensor):
        return [d.value for d in tensor.get_shape()]

    def prod(xs):
        p = 1
        for x in xs:
            p *= x
        return p

    def gaussian_kernel(kernlen=11, nsig=3):
        """Returns a 2D Gaussian kernel array."""
        interval = (2*nsig+1.)/(kernlen)
        x = np.linspace(-nsig-interval/2., nsig+interval/2., kernlen+1)
        kern1d = np.diff(st.norm.cdf(x))
        kernel_raw = np.sqrt(np.outer(kern1d, kern1d))
        kernel = kernel_raw/kernel_raw.sum()
        return kernel

    def network(x, feat_size, window_width):
        input_size = prod(shape(x))
        output_size = window_width * window_width * 3

        w_conv1 = tf.Variable(tf.random_uniform((input_size, input_size)))
        b_conv1 = tf.Variable(tf.random_uniform((input_size,)))

        w_conv2 = tf.Variable(tf.random_uniform((input_size, input_size)))
        b_conv2 = tf.Variable(tf.random_uniform((input_size,)))

        w_conv3 = tf.Variable(tf.random_uniform((input_size, output_size)))
        b_conv3 = tf.Variable(tf.random_uniform((output_size,)))

        flat_x = tf.reshape(x, shape=(1, input_size))
        layer1 = tf.sigmoid(tf.matmul(flat_x, w_conv1) + b_conv1)
        layer2 = tf.sigmoid(tf.matmul(layer1, w_conv2) + b_conv2)
        layer3 = tf.sigmoid(tf.matmul(layer2, w_conv3) + b_conv3)

        return tf.reshape(layer3, shape=(window_width, window_width, 3))

    def bilateral_filter_window(img, w):
        g = gaussian_kernel(shape(img)[0])
        g = np.stack((g,g,g), axis=2)
        w = g * w
        return tf.reduce_sum(img * w, (0,1)) / tf.reduce_sum(w)

    x = tf.placeholder(tf.float32, [width, width, feat_size])
    y_ = tf.placeholder(tf.float32, [3])
    w = network(x, feat_size=feat_size, window_width=width)
    y = bilateral_filter_window(x[:,:,:3], w)
    # mse = tf.reduce_sum(tf.square(y - y_))
    eps = tf.constant(1e-8)
    relmse = tf.clip_by_value(tf.reduce_sum(tf.square(y - y_)/(tf.square(y_) + eps)), 0.0, 1.0)
    train_step = tf.train.AdamOptimizer(1e-4).minimize(relmse)

    return x, y_, y, train_step, relmse

def run_epoch(train_step, dataset, epoch_size=100):
    for _ in range(epoch_size):
        sample, label = dataset.next()
        train_step.run(feed_dict={x: sample, y_: label})

def test_model(err, dataset, instances=10):
    errsum = 0
    for _ in range(instances):
        sample, label = dataset.next()
        errsum += err.eval(feed_dict={x: sample, y_: label})
    return errsum / instances


def filter_scene(y, scene):
    wins = scene.windows()
    return np.array([[y.eval(feed_dict={x: w}) for w in row] for row in wins])

if __name__ == '__main__':
    import sys
    sess = tf.InteractiveSession()
    x, y_, y, train_step, err = create_network(width=5,feat_size=20)
    sess.run(tf.global_variables_initializer())

    npys = zip(sys.argv[1::2], sys.argv[2::2])
    dataset = create_batches.Dataset(npys, kernel_size=5)
    errs = []
    for epoch in range(1000):
        run_epoch(train_step, dataset)
        e = test_model(err, dataset)
        print("epoch:", epoch, e)
        errs.append(e)

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





