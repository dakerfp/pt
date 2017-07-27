import tensorflow as tf
import numpy as np
import create_batches
import lbf


def prod(xs):
    p = 1
    for x in xs:
        if x is not None:
            p *= x
    return p

def flatten_x(xs):
    xcols = [xi[:,:,10:13] for xi in xs]

    xsize = int(prod(xs[0].shape))
    flat_xs = np.stack([np.reshape(xi, (xsize,)) for xi in xs], axis=0)

    xcsize = prod(xcols[0].shape)
    flat_xcols = np.stack([np.reshape(xi, (xcsize,)) for xi in xcols])

    return flat_xs, flat_xcols

def fake_filter(flt, scene, sess):
    h, w = scene.shape_windows()
    arr = np.zeros((h,w,3))

    for i, j, win in scene.iter_windows():
        xs, xcols = flatten_x([win])
        # pixels = xs # flt.x.eval(feed_dict={flt.x: xs})
        pixel = flt.y.eval(session=sess, feed_dict={flt.x: xs, flt.xcol: xcols})
        color = pixel[0,:3]
        arr[i,j,:] = color
        print(h, i, j)

    print("HERE")
    # print(flt.y)

    # reassemble image
    # print(flt.x)
    # print(len(pixels), pixels.shape)
    # print(pixels)
    # arr = np.reshape(pixels, (h, w, 3))
    # arr = np.array(arr)
    print(arr.shape)
    return arr[:,:,:3]

def show_images(flt, dataset, sess):
    imgs = []
     
    scene = dataset.next_scene()
    y = fake_filter(flt, scene, sess)
    imgs.append((scene.arr[:,:,:3], y))

    import matplotlib.pyplot as plt
    fig = plt.figure()
    print("plotting")
    fig.add_subplot(1, 2, 1)
    plt.imshow(imgs[0][0])
    fig.add_subplot(1, 2, 2)
    plt.imshow(imgs[0][1])
    plt.show()


if __name__ == '__main__':
    kwidth=11

    import sys
    dataset = None
    if sys.argv[1].endswith(".zip"):
        dataset = create_batches.ZipDataset(sys.argv[1:], prefixes=[3], lowres=16, hires=1024, kernel_size=kwidth)
    else:
        npys = zip(sys.argv[1::2], sys.argv[2::2])
        dataset = create_batches.Dataset(npys, 25)

    depth = dataset.depth()
    flt = lbf.LearningBasedFilter(width=kwidth,depth=depth)
    saver = tf.train.Saver()
    with tf.Session() as sess:
        sess.run(tf.global_variables_initializer())
        saver.restore(sess, 'lbf-new-v')
        show_images(flt, dataset, sess)
