# gpt-batch-translator
大家可以使用我们开发的基于 ChatGPT API 的角色剧情翻译工具，可以对关键词进行限定，在character中对角色语言性格进行限定使翻译出的话更贴近角色，对于一段文字进行整体prompt优化等，欢迎使用，也希望大家一起完善修改。<br>
You can use our role-playing translation tool based on the ChatGPT API, which allows you to limit keywords and the language and personality of characters to make the translated dialogue more true to the character. Additionally, it optimizes the entire text prompt. We welcome your use and hope that everyone can work together to improve and modify it.

# translateClient
简单的python脚本，从translateClient/upload上传excel文件，上传成功后自动开始翻译并显示翻译进度，结束后下载文件到translateClient/download文件夹下。运行:python ./translate.py<br>
A simple Python script that uploads an Excel file from translateClient/upload, automatically starts the translation process and displays the progress, and downloads the file to translateClient/download upon completion. To run, use: python ./translate.py.

# translator
go语言用来连接gptAPI的服务端部分，原本是打算上传下载等功能都写在页面的，奈何我们webUI太糟糕，只能将这部分当做服务器功能来做。请在config.yaml配置自己的openai的APIkey，然后执行go run main.go<br>
The server side of the Go language is used to connect to the GPT API. Originally, it was planned to write all the functions such as uploading and downloading on the webpage. However, due to our poor web UI, we can only treat this part as a server function. Please configure your own OpenAI API key in config.yaml and then execute "go run main.go".
