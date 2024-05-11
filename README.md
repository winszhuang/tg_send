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
> 花費時間 ? 15:14
- [] 部屬到vercel

## api規格  
### req  
- message string

### resp
- success bool
- message string

