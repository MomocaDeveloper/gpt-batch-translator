import requests
import os,random
import string
import time
import json
import sys


def upload_file(fileName):
    #selected_file_name = "./" + fileName
 
    print("正在上传文件:", fileName)  
    myobj = {'type': 3}
    upload_url = 'http://35.219.174.249:9003/upload'
    files = {'file': open(fileName, 'rb')}
    # 发起Post请求上传文件
    print("发起POST请求:", fileName)  
    response = requests.post(upload_url, files = files, data = myobj)
    return fileName
    
if __name__ == '__main__':
    if len(sys.argv) < 1:
        print("没有输入文件名!")
    fileName = str(sys.argv[1])
    if os.path.exists(fileName):
        upload_file(fileName)
    else:
        print("文件不存在！")