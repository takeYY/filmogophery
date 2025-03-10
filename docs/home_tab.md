```mermaid
---
title: Homeタブの遷移
---

flowchart TB
    subgraph Home
        direction TB

        H001@{ img: "https://raw.githubusercontent.com/takeYY/filmogophery/main/docs/images/H001_home.png", pos: "t",label: "Home", w: 720,constraint: "on" }
        H002@{ img: "https://raw.githubusercontent.com/takeYY/filmogophery/main/docs/images//H002_detail.png", pos: "t",label: "Detail", w: 720,constraint: "on" }
        H003@{ img: "https://raw.githubusercontent.com/takeYY/filmogophery/main/docs/images//H003_edit_impression.png", pos: "t",label: "Edit Impression", w: 720,constraint: "on" }
        H004@{ img: "https://raw.githubusercontent.com/takeYY/filmogophery/main/docs/images//H004_create_record.png", pos: "t",label: "Create Record", w: 720,constraint: "on" }

        H001 -- cardを選択 --> H002

        H002 -- 感想を編集 --> H003
        H002 -- 視聴履歴を作成 --> H004

        H003 -- 更新 --> H001
        H004 -- 作成 --> H001
    end

    H001 -- Watch List ----> WL001
    H001 -- Watch Calendar ----> WC001
    H001 -- 検索 ----> S001
```
