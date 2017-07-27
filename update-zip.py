import sys
import zipfile as zf

z = zf.ZipFile(sys.argv[1], mode='a', allowZip64=True)
spp = int(sys.argv[4])
filename = sys.argv[2] % spp
to = sys.argv[3] % spp
print(filename, to)
z.write(filename, to)
z.close()
