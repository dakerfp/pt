
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

        w_conv1 = tf.Variable(tf.random_uniform((input_size, 32)))
        b_conv1 = tf.Variable(tf.random_uniform((32,)))

        w_conv2 = tf.Variable(tf.random_uniform((32, 32)))
        b_conv2 = tf.Variable(tf.random_uniform((32,)))

        w_conv3 = tf.Variable(tf.random_uniform((32, output_size)))
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
    w = network(x, feat_size=feat_size, window_width=11)
    y = bilateral_filter_window(x[:,:,:3], w)
    mse = tf.reduce_sum(tf.square(y - y_))
    train_step = tf.train.AdamOptimizer(1e-4).minimize(mse)

    return x, y_, y, train_step

if __name__ == '__main__':
    import sys
    sess = tf.InteractiveSession()
    x, y_, y, train_step = create_network()
    sess.run(tf.initialize_all_variables())

    dataset = create_batches.Dataset(sys.argv[1:], kernel_size=11)
    for i in range(100):
        sample = dataset.next()
        train_step.run(feed_dict={x: sample, y_: sample[11/2,11/2,:3]})


