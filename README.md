## 対象のデータ：
東京都の分譲マンション推移
https://catalog.data.metro.tokyo.lg.jp/dataset/t000008d0000000036/resource/f4279188-839d-401c-94d8-f3d404f4bf90

## APIのエンドポイント
 https://lgwhkdlifi.execute-api.ap-northeast-1.amazonaws.com/Prod/


## APIの仕様
https://ito-kota.github.io/techboost/api_document.html


## データベーステーブル
テーブル名：ApartmentData<br>
|カラム名|タイプ|
| ---- | ---- |
| year | integer not null primary |
| number | integer not null |
| area | double not null |
| price | integer not null |

