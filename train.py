import tensorflow as tf
import create_batches
import lbf

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

    x, xcol, y_, y, train_step, err = lbf.create_network(width=kwidth,depth=depth)
    saver = tf.train.Saver()

    sess = tf.InteractiveSession()

    sess.run(tf.global_variables_initializer())

    errs = []
    for epoch in range(101):
        lbf.run_epoch(x, xcol, y_, train_step, dataset)
        e = lbf.test_model(x, xcol, y_, err, dataset)
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
    # scene = dataset.scenes[1]
    # img = lbf.filter_scene(y, scene)
    # fig.add_subplot(2,2,2)
    # plt.imshow(np.clip(img, 0, 1))
    fig.add_subplot(2,2,3)
    plt.imshow(np.clip(scene.color(), 0, 1))
    fig.add_subplot(2,2,4)
    plt.imshow(np.clip(scene.gt_color(), 0, 1))
    plt.show()
