# PoC Web App API

イベントの記録を行うサーバのRESTベースAPI（GoのWebアプリの習作）  

## 一覧

|API|機能|
| - | - |
|GET /api/events|✓イベントリストの取得|
|GET /api/events/{event_id}|✓イベントの取得（参加者リストや画像URL付）|
|POST /api/events|✓イベントの登録（内容をjsonで送信しidを自動発番）|
|DELETE /api/events/{event_id}|✓イベントの削除（参加者や画像は残る）※ToDo:バインド掃除|
|GET /api/persons|✓個人リストの取得|
|GET /api/persons/{person_id}|✓個人の取得|
|POST /api/persons|✓個人の登録（内容をjsonで送信しidを自動発番）|
|DELETE  /api/persons/{person_id}|✓個人の削除※ToDo:バインド掃除|
|GET /api/images|✓画像リストの取得|
|GET /api/images/{image_id}|✓画像情報の取得(画像本体はURLより取得)|
|POST /api/images|✓画像のアップロード(内容をフォームで送信しidを自動発番)|
|DELETE  /api/images/{image_id}|✓画像の削除※ToDo:バインド掃除|
|POST /api/events/{event_id}/persons|✓イベントへ参加者リストを追加|
|POST /api/events/{event_id}/images|✓イベントへ画像リスト追加|  
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

# イベント [events]

## イベントリスト取得 [GET]

```
GET /api/events HTTP/1.1
```

### Request

| パラメータ | 内容 | 必須 | デフォルト値 | 最大値 |
|  ---  |  ---  |  ---  |  ---  |  ---  |
| limit | 最大件数 | No | 50 | 500 |
| offset | 開始位置(件数ベース) | No | 0 | -|

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

### Test

```
$ curl localhost:1323/api/events
$ curl localhost:1323/api/events?limit=2\&offset=1
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
      "person_id": 1,
      "first_name": "岩佐",
      "last_name": ""
    },
    {
      "person_id": 3,
      "first_name": "武知",
      "last_name": ""
    }
  ],
  "images": [
    {
      "ImagePath": "/images/4.png",
      "image_id": 4,
      "image_name": "test-image",
      "content_type": "image/png"
    }
  ]
}
```

### Test

```
$ curl localhost:1323/api/events/2
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

### Response

```
HTTP/1.1 201 Created

{
  "EventID": 6
}
```

### Test

```
$ curl -X POST -H "Content-Type: application/json" localhost:1323/api/events -d '{"title": "追加ダイビング","description": "伊豆大島","event_date": "2022-08-15T00:00:00+09:00"}'
```

## イベントの削除 [DELETE]
イベントの削除（参加者や画像は残る）
```
DELETE /api/events/{event_id}
```

### Request

| パラメータ | 内容 | 必須 | デフォルト値 | 最大値 |
|  ---  |  ---  |  ---  |  ---  |  ---  |
|event_id|削除対象||||

### Response

```
HTTP/1.1 204 No Content
```

### Test

```
$ curl -X DELETE localhost:1323/api/events/4
```

# 個人 [persons]

## 個人リストの取得[GET]

```
GET /api/persons HTTP/1.1
```

### Request

| パラメータ | 内容 | 必須 | デフォルト値 | 最大値 |
|  ---  |  ---  |  ---  |  ---  |  ---  |
| limit | 最大人数 | No | 50 | 500 |
| offset | 開始位置(人ベース) | No | 0 | -|

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

### Test

```
$ curl localhost:1323/api/persons
$ curl localhost:1323/api/persons?limit=2\&offset=1
```

## 個人の取得[GET]

```
GET /api/persons/{person_id} HTTP/1.1
```

### Request

| パラメータ | 内容 | 必須 | デフォルト値 | 最大値 |
|  ---  |  ---  |  ---  |  ---  |  ---  |
| person_id | 取得対象個人 | Yes |  |  |

### Response

```
HTTP/1.1 200 OK

{
  "PersonId": 2,
  "FirstName": "高田",
  "LastName": ""
}
```

### Test

```
$ curl localhost:1323/api/persons/2
```
## 個人の登録 [POST]
個人を登録する。内容をjsonで送信しPersonIDが自動発番される。
```
POST /api/persons
```

### Request

| パラメータ | 内容 | 必須 | デフォルト値 | 最大値 |
|  ---  |  ---  |  ---  |  ---  |  ---  |
|first_name|名||||
|last_name|氏|||
```
{
  "first_name": "黒田",
  "last_name": ""
}
```

### Response

```
HTTP/1.1 201 Created

{
  "PersonID": 4
}
```

### Test

```
$ curl -X POST -H "Content-Type: application/json" localhost:1323/api/persons -d '{"first_name": "黒田","last_name": ""}'
```

## 個人の削除 [DELETE]
個人の削除
```
DELETE /api/persons/{person_id}
```

### Request

| パラメータ | 内容 | 必須 | デフォルト値 | 最大値 |
|  ---  |  ---  |  ---  |  ---  |  ---  |
|person_id|削除対象||||

### Response

```
HTTP/1.1 204 No Content
```

### Test

```
$ curl -X DELETE localhost:1323/api/persons/5
```
# 画像 [images]

## 画像リストの取得[GET]

```
GET /api/images HTTP/1.1
```

### Request

| パラメータ | 内容 | 必須 | デフォルト値 | 最大値 |
|  ---  |  ---  |  ---  |  ---  |  ---  |
| limit | 最大枚数 | No | 50 | 500 |
| offset | 開始位置(番号ベース) | No | 0 | -|

### Response

```
HTTP/1.1 200 OK

[
  {
    "ImagePath": "/images/3.png",
    "image_id": 3,
    "image_name": "test-image",
    "content_type": "image/png"
  },
  {
    "ImagePath": "/images/4.png",
    "image_id": 4,
    "image_name": "test-image",
    "content_type": "image/png"
  }
]
```

### Test

```
$ curl localhost:1323/api/images
$ curl localhost:1323/api/images?limit=2\&offset=1
```

## 画像の取得[GET]
画像情報を取得する。（画像本体はImagePathで返されるURLより取得する。）
```
GET /api/images/{image_id} HTTP/1.1
```

### Request

| パラメータ | 内容 | 必須 | デフォルト値 | 最大値 |
|  ---  |  ---  |  ---  |  ---  |  ---  |
| image_id | 取得対象画像 | Yes |  |  |

### Response

```
HTTP/1.1 200 OK

{
  "ImagePath": "/images/16.png",
  "image_id": 16,
  "image_name": "htmlicon.png",
  "content_type": "image/png"
}
```

### Test

```
$ curl localhost:1323/api/images/16
```

## 画像のアップロード [POST]
画像のアップロード。内容をフォームで送信しImageIDが自動発番される。
```
POST /api/images
```

### Request

| パラメータ | 内容 | 必須 | デフォルト値 | 最大値 |
|  ---  |  ---  |  ---  |  ---  |  ---  |
|form-data|formエンコードした画像データ|Yes||1MByte|
|filename|ファイル名||||
|Content-Type|image/png,image/jpegなど|||

### Response

```
HTTP/1.1 201 Created

{
  "ImageID": 19
}
```

### Test

```
$ curl -X POST -F file=@htmlicon.png localhost:1323/api/images
```

## 画像の削除 [DELETE]
画像の削除（※ファイル自体は残しDBから削除。ファイル削除は別途バッチ要。）
```
DELETE  /api/images/{image_id}
```

### Request

| パラメータ | 内容 | 必須 | デフォルト値 | 最大値 |
|  ---  |  ---  |  ---  |  ---  |  ---  |
|image_id|削除対象||||

### Response

```
HTTP/1.1 204 No Content
```

### Test

```
$ curl -X DELETE localhost:1323/api/images/17
```

## イベントへ参加者を追加 [POST]
イベントへ参加者リストを追加
```
POST /api/events/{event_id}/persons
```

### Request

| パラメータ | 内容 | 必須 | デフォルト値 | 最大値 |
|  ---  |  ---  |  ---  |  ---  |  ---  |
|event_id|追加する先||||
|person_id|追加する対象の配列||||
```
[
  {
    "person_id": 4
  },
  {
    "person_id": 7
  }
]
```

### Response

```
HTTP/1.1 204 No Content
```

### Test

```
$ curl -X POST -H "Content-Type: application/json" localhost:1323/api/events/3/persons -d '[{"person_id": 4}]'
$ curl -X POST -H "Content-Type: application/json" localhost:1323/api/events/4/persons -d '[{"person_id": 4},{"person_id": 5}]'
```

## イベントへ画像を追加 [POST]
イベントへ画像リストを追加
```
POST /api/events/{event_id}/images
```

### Request

| パラメータ | 内容 | 必須 | デフォルト値 | 最大値 |
|  ---  |  ---  |  ---  |  ---  |  ---  |
|event_id|追加する先||||
|image_id|追加する対象の配列||||

```
[
  {
    "image_id": 4
  },
  {
    "image_id": 5
  }
]
```

### Response

```
HTTP/1.1 204 No Content
```

### Test

```
curl -X POST -H "Content-Type: application/json" localhost:1323/api/events/7/images -d '[{"image_id": 4}]'
curl -X POST -H "Content-Type: application/json" localhost:1323/api/events/7/images -d '[{"image_id": 4},{"image_id": 5}]'
```