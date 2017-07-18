import tensorflow as tf
import numpy as np
import create_batches
import lbf

if __name__ == '__main__':
    kwidth=11
    ninstances = 1

    import sys
    dataset = None
    if sys.argv[1].endswith(".zip"):
        dataset = create_batches.ZipDataset(sys.argv[1:], prefixes=[1, 2, 3, 4, 5, 6], lowres=16, hires=1024, kernel_size=kwidth)
    else:
        npys = zip(sys.argv[1::2], sys.argv[2::2])
        dataset = create_batches.Dataset(npys, 25)

    depth = dataset.next()[0].shape[2]

    flt = lbf.LearningBasedFilter(width=kwidth,depth=depth)
    saver = tf.train.Saver()
    with tf.Session() as sess:
        sess.run(tf.global_variables_initializer())
        saver.restore(sess, 'lbf-basic-4000')

        import matplotlib.pyplot as plt
        fig = plt.figure()
        for i in range(ninstances):
            x = dataset.next_scene()

            fig.add_subplot(ninstances, 2, 2 * i + 1)
            plt.imshow(x.arr[:,:,:3])

            y = flt.filter_scene(x)
            fig.add_subplot(ninstances, 2, 2 * i + 2)
            plt.imshow(y)

        plt.show()

