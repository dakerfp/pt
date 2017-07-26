from sys import argv
import zipfile as zf

z = zf.ZipFile(argv[1], mode='a')
z.write(argv[2], argv[3])
z.close()