# nostr-relay

## Usage

* Rename .env.example to .env
* Set value of AUTOMIGRATE=1 in .env to migrate Database table at first time
* Set dababase addres and username/password in .env
* Run project at  http://localhost:8100/ for default

* Default watcher http://127.0.0.1:8100/watcher

 ## Questions

 #### **Why did you choose this database?**
 - 我使用postgresql，因為時間上的關係使用比較熟悉的DB

 #### **If the number of events to be stored will be huge, what would you do to scale the database?**
 - 將數據儲存在不同的表，如果系統和硬碟效能已經達到瓶頸，則考慮分庫或其他分散式資料庫。
