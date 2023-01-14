# PoC Web App API

イベントの記録を行うサーバのRESTベースAPI（GoのWebアプリの習作）  

## 一覧

|API|機能|
| - | - |
|GET /api/events|✓イベントリストの取得|
|GET /api/events/{event_id}|✓イベントの取得（参加者リストや画像URL付）|
|POST /api/events|✓イベントの登録（内容をjsonで送信しidを自動発番）|
|DELETE /api/events/{event_id}|イベントの削除（参加者や画像は残る）※バインド掃除|
|GET /api/persons|✓個人リストの取得|
|GET /api/persons/{person_id}|個人の取得|
|POST /api/persons|個人の登録（内容をjsonで送信しidを自動発番）|
|DELETE  /api/persons/{person_id}|個人の削除|
|POST /api/images|画像のアップロード(内容をフォームで送信しidを自動発番)|
|DELETE  /api/images/{image_id}|画像の削除|
|POST /api/events/{event_id}/persons|イベントへ参加者リストを追加|
|POST /api/events/{event_id}/images|イベントへ画像リスト追加|  
※アップロードされた画像は画像URLから取得できる

|ホスト|プロトコル|データ形式|
| - | - | - |
|lolalhost:1323|http|JSON|

## ステータスコード

| ステータスコード |net/http定義| 説明 |
| - | - | - |
| 200 |StatusOK| リクエスト成功 |
| 201 |StatusCreated|登録成功| 
| 204 |StatusNoContent| リクエストに成功したが返却するbodyが存在しない |
| 400 |StatusBadRequest| 不正なリクエストパラメータを指定している |
| 401 |StatusUnauthorized| APIアクセストークンが不正、または権限不正 |
| 404 |StatusNotFound| 存在しないURLにアクセス |
| 429 |StatusTooManyRequests| リクエスト制限を超えている |
| 500 |StatusInternalServerError| 不明なエラー |

# イベント [/events]

## イベントリスト取得 [GET]

```
GET /api/events HTTP/1.1
```

### Request

| パラメータ | 内容 | 必須 | デフォルト値 | 最大値 |
|  ---  |  ---  |  ---  |  ---  |  ---  |
| limit | 最大件数 | No | 50 | 500 |
| offset | 開始位置(件数ベース) | No | 0 | -|


```
$ curl localhost:1323/api/events?limit=2\&offset=1
```

### Response

```
HTTP/1.1 200 OK

[
  {
    "EventID": 2,
    "AccountID": 0,
    "Title": "フィリピンダイビング",
    "Description": "パラワン島",
    "EventDate": "2020-05-03T00:00:00+09:00"
  },
  {
    "EventID": 3,
    "AccountID": 0,
    "Title": "伊豆ダイビング",
    "Description": "熱海",
    "EventDate": "2022-08-10T00:00:00+09:00"
  }
]
```
## イベントの取得 [GET]

参加者リストや画像URLリストを含むイベント情報を取得する
```
GET /api/events/{event_id}
```

### Request

| パラメータ | 内容 | 必須 | デフォルト値 | 最大値 |
|  ---  |  ---  |  ---  |  ---  |  ---  |
| event_id | 取得対象イベント | Yes |  |  |


```
$ curl localhost:1323/api/events/2
```

### Response

```
HTTP/1.1 200 OK

{
  "event": {
    "EventID": 2,
    "AccountID": 1,
    "Title": "フィリピンダイビング",
    "Description": "パラワン島",
    "EventDate": "2020-05-03T00:00:00+09:00"
  },
  "persons": [
    {
      "PersonId": 1,
      "FirstName": "岩佐",
      "LastName": ""
    },
    {
      "PersonId": 3,
      "FirstName": "武知",
      "LastName": ""
    }
  ]
}
```
## イベントの登録 [POST]
イベントを登録する。EventIDが自動発番される。
```
POST /api/events
```

### Request

| パラメータ | 内容 | 必須 | デフォルト値 | 最大値 |
|  ---  |  ---  |  ---  |  ---  |  ---  |
|title|題名||||
|description|説明|||
|event_date|日時|||
```
{
  "title": "追加ダイビング",
  "description": "伊豆大島",
  "event_date": "2022-08-15T00:00:00+09:00"
}
```
```
$ curl -X POST -H "Content-Type: application/json" localhost:1323/api/events -d '{"title": "追加ダイビング","description": "伊豆大島","event_date": "2022-08-15T00:00:00+09:00"}'
```

### Response

```
HTTP/1.1 201 Created

{
  "EventID": 6
}
```

## イベントの削除 [DELETE]
イベントの削除（参加者や画像は残る）
```
DELETE /api/events/{event_id}
```

### Request

| パラメータ | 内容 | 必須 | デフォルト値 | 最大値 |
|  ---  |  ---  |  ---  |  ---  |  ---  |
|event_id}|削除対象||||

```
$ curl -X DELETE localhost:1323/api/events/4
```

### Response

```
HTTP/1.1 204 No Content
```
## 個人リストの取得[GET]

```
GET /api/persons HTTP/1.1
```

### Request

| パラメータ | 内容 | 必須 | デフォルト値 | 最大値 |
|  ---  |  ---  |  ---  |  ---  |  ---  |
| limit | 最大人数 | No | 50 | 500 |
| offset | 開始位置(人ベース) | No | 0 | -|

```
$ curl localhost:1323/api/persons?limit=2\&offset=1
```

### Response

```
HTTP/1.1 200 OK

[
  {
    "PersonId": 2,
    "FirstName": "高田",
    "LastName": ""
  },
  {
    "PersonId": 3,
    "FirstName": "武知",
    "LastName": ""
  }
]
```

