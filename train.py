import tensorflow as tf
import numpy as np
import create_batches
import lbf
import random


def smooth(xs, w):
    return [sum(xs[i:i+w]) / w for i in range(len(xs)-w)]

if __name__ == '__main__':
    import sys
    kwidth=11
    dataset = None
    random.seed(0)

    if any(fn.endswith(".zip") for fn in sys.argv[1:]):
        dataset = create_batches.ZipDataset(sys.argv[1:], prefixes=[3, 4, 5, 6], lowres=32, hires=1024, kernel_size=kwidth, depth=37)
    else:
        npys = zip(sys.argv[1::2], sys.argv[2::2])
        dataset = create_batches.Dataset(npys, 25)

    flt = lbf.LearningBasedFilter(width=kwidth,depth=dataset.depth)
    saver = tf.train.Saver()
    sess = tf.InteractiveSession()
    sess.run(tf.global_variables_initializer())

    errs = []
    berrs = []
    for epoch in range(1001):
        flt.run_epoch(dataset, 1000)
        e = flt.test_model(sess, dataset)
        print("epoch:", epoch, e)
        errs.append(e)
        if epoch % 100 == 0 and epoch > 10:
            print("save")
            saver.save(sess, 'lbf-new-v')

    np.save('./errs.txt', errs)

    import matplotlib.pyplot as plt

    fig = plt.figure()
    plt.plot(smooth(errs, 20))
    # plt.plot(smooth(berrs, 20))
    plt.plot(smooth(errs, 100))
    # plt.plot(smooth(berrs, 100))
    scene = dataset.scenes[1]
    plt.show()
    # plt.save("a.png")
