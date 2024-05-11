# Telegram Bot 訊息發送服務 api

## 目標
只要某一方呼叫這個api, 指定頻道和訊息就會發到對應telegram聊天室

## 代辦  

### 階段一  
> 花費時間 40分
- [x] 研究telegram bot價錢
- [x] 看telegram bot怎麼用go串接  
- [x] 使用gin實作一隻發送api(先不要區分不同聊天室, 那是第二階段)  

### 階段二  
> 花費時間 ? 
- [x] 部屬到render.com 15:14 ~ 3:35
- [x] 調整請求格式 可以直接貼整段文字 ~ 3:46  

### 階段三
> 花費時間 ~ 5:00
- [x] 確認其他人是不是也可以收到這個bot通知
- [x] 權限, 不是任何人都可以用這個api
- [x] bot只要加入群組, 該群組就可以收到bot通知

## api規格  
### req  
- message string

### resp
- success bool
- message string

