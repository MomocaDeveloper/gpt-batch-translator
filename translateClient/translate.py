import requests
import os,random
import string
import time
import json
import math

httpUrl = 'http://{your_IP}:{your_port}/'

def upload_file():
    file_list = os.listdir("./upload")
    for i, file_name in enumerate(file_list):
        print(f"{i+1}. {file_name}")
    selected_file_index = int(input("Please enter the file number to upload："))
    fileName = str(file_list[selected_file_index - 1])
    selected_file_name = "./upload/" + str(file_list[selected_file_index - 1])
    
    print("uoloading:", selected_file_name)  
    
    upload_url = httpUrl + 'upload'
    files = {'file': open(selected_file_name, 'rb')}
    # Post update file
    print("start POST:", selected_file_name)  
    response = requests.post(upload_url, files=files)
    return fileName
    
#fileName="test1.xlsx"
def download_file(fileName):
    download_url = httpUrl + 'download?type=download&file='
    download_url = download_url + 'tr_' + fileName
    response = requests.get(download_url)
    
    with open(('./download/' + 'tr_' + fileName), 'wb') as f:
        f.write(response.content)

def translate_file(fileName):
    # can edit your temperature here
    translate_url = httpUrl + 'translate?temperature=0.7&file='
    translate_url = translate_url + fileName
    response = requests.get(translate_url)
    #print(response.content)
    
def progress_file(fileRealName, result):
    progress_url = httpUrl + 'progress'
    response = requests.get(progress_url)
    try:
        resp_content = response.json()
    except ValueError:
        resp_content = None
        
    if resp_content:
        dataList = response.json()
        return int(dataList[fileRealName]), 0
    else:
        if result > 0:
            return 0, 999
        else:
            return 0, 0
    
def timer_api(fileName):
    index = 1
    startTick = False
    _result = 0
    _lastProgress = 0
    fileRealName = fileName.split('.')[0]
    fanyiText = "translate progress:"
    badTokenCount = 0
    print("translating，wait please，now " + fanyiText + "0%")
    while (_result <= 99 and index < 10000):
        _result, errcode = progress_file(fileRealName, _result)
        if errcode == 999:
            badTokenCount += 1
            if badTokenCount >= 5:
                break
        if _result != _lastProgress:
            if _result < _lastProgress:
                _result = _lastProgress
            print(fanyiText + str(_result) + "%")
            _lastProgress = _result
        if _result > 0:
            startTick = True
        if startTick == True:
            index += 1
        time.sleep(0.1)
    
if __name__ == '__main__':
    fileName = upload_file()
    print("uoload sucess, start translate:", fileName)
    translate_file(fileName)
    start = time.time()
    timer_api(fileName)
    end = time.time()
    time.sleep(6)
    runtime = math.floor(end - start)
    print("translate progress:100%")
    print("Translation total duration:", str(runtime) + "sec")
    download_file(fileName)
