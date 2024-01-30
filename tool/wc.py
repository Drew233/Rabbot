import sys
from PIL import Image
from wordcloud import WordCloud
from matplotlib.font_manager import FontProperties
import matplotlib.pyplot as plt
import numpy as np
import jieba
from datetime import datetime

def chinese_jieba(text):
    words = [word for word in jieba.cut(text) if len(word) >= 2]
    wordst = " ".join(words)
    print(wordst)
    return wordst

def read_file(file_path):
    with open(file_path, 'r', encoding='utf-8') as file:
        content = file.read()
    return content

def main():
    arguments = sys.argv
    origin_file = arguments[1]
    target_file = arguments[2]
    if len(arguments) <= 3:
        # 获取当前日期
        current_date = datetime.now()

        # 格式化日期为指定格式
        formatted_date = current_date.strftime("%Y.%m.%d")

        data = formatted_date + " 词云"
    else:
        data = arguments[3] + " 词云"

    wc = WordCloud(font_path="./tool/cus.otf",
                   width=800,
                   height=400,
                   background_color="white",
                   max_words=200,
                   mask=np.array(Image.open("./tool/rab.png")),
                   )


    file_path = origin_file
    file_content = read_file(file_path)

    wordcloud = wc.generate(chinese_jieba(file_content))
    plt.figure(figsize=(8, 5))
    plt.imshow(wordcloud, interpolation='bilinear')
    plt.axis("off")
    # 设置标题字体
    title_font = FontProperties(fname="./tool/cus.otf", size=24)  # 替换为你想要使用的字体文件路径和字体大小
    plt.title(data, fontproperties=title_font, color='black')

    # 保存词云图（可以根据需要修改文件名和格式）
    plt.savefig(target_file)
    # plt.show()

main()
