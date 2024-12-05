# 変換ロジックは https://qiita.com/Cartelet/items/5c1c012c132be3aa9608 を参考にした

from PIL import Image

import matplotlib.pyplot as plt
import numpy as np
import cv2
import sys

# (PIC: 元画像、R: 正方形領域の一辺)
def kuwahara(pic,r=5):
    h,w,_=pic.shape
    pic=np.pad(pic,((r,r),(r,r),(0,0)),"edge")
    ave,var=cv2.integral2(pic)
    # 平均値の一括計算
    ave=(ave[:-r-1,:-r-1]+ave[r+1:,r+1:]-ave[r+1:,:-r-1]-ave[:-r-1,r+1:])/(r+1)**2
    # 分散の一括計算
    var=((var[:-r-1,:-r-1]+var[r+1:,r+1:]-var[r+1:,:-r-1]-var[:-r-1,r+1:])/(r+1)**2-ave**2).sum(axis=2)

    def filt(i,j):
        return np.array([ave[i,j],ave[i+r,j],ave[i,j+r],ave[i+r,j+r]])[(np.array([var[i,j],var[i+r,j],var[i,j+r],var[i+r,j+r]]).argmin(axis=0).flatten(),j.flatten(),i.flatten())].reshape(w,h,_).transpose(1,0,2)
    # 色の決定
    filtered_pic = filt(*np.meshgrid(np.arange(h),np.arange(w))).astype(pic.dtype)
    return filtered_pic

args = sys.argv
input = args[1]

pic = np.array(plt.imread(input))
filtered_pic = kuwahara(pic,5)
filtered_pic = np.clip(filtered_pic, 0, 1)  # 値を0〜1にクリップ
filtered_pic = (filtered_pic * 255).astype(np.uint8)
filtered_image = Image.fromarray(filtered_pic)
filtered_image.save(input)
