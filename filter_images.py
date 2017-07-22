import tensorflow as tf
import numpy as np
import create_batches
import lbf


def show_images(flt, dataset, ninstances=1):
    imgs = []
    for i in range(ninstances):
        print("filter scene: ", i)
        x = dataset.next_scene()
        y = flt.filter_scene(x)
        imgs.append((x.arr[:,:,:3], y))


    import matplotlib.pyplot as plt
    fig = plt.figure()
    for i in range(ninstances):
        print("plotting ", i)
        fig.add_subplot(ninstances, 2, 2 * i + 1)
        plt.imshow(imgs[i][0])
        fig.add_subplot(ninstances, 2, 2 * i + 2)
        plt.imshow(imgs[i][1])

    plt.show()

if __name__ == '__main__':
    kwidth=11

    import sys
    dataset = None
    if sys.argv[1].endswith(".zip"):
        dataset = create_batches.ZipDataset(sys.argv[1:], prefixes=[3, 4, 5, 6], lowres=16, hires=1024, kernel_size=kwidth)
    else:
        npys = zip(sys.argv[1::2], sys.argv[2::2])
        dataset = create_batches.Dataset(npys, 25)

    depth = dataset.next()[0].shape[2]

    flt = lbf.LearningBasedFilter(width=kwidth,depth=depth)
    saver = tf.train.Saver()
    with tf.Session() as sess:
        sess.run(tf.global_variables_initializer())
        saver.restore(sess, 'lbf-g-v')
        show_images(flt, dataset, ninstances=1)
