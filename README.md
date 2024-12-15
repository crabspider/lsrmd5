# lsrmd5

指定されたディレクトリ内のファイル一覧をMD5ハッシュ付きで出力する。

写真や古いファイル等、更新や削除が行われないデータを格納したディレクトリに対して実行し、意図せぬ変更が生じていないか確認する事が目的。

# 使用例

優先度を下げて実行する事を推奨。Windowsの場合は`start /belownormal /b /wait .\lsrmd5.exe`等。

## ファイルリスト作成

    .\lsrmd5.exe Documents/archive Dropbox/archive > archive.md5

## フォルダ構造無視モード

    .\lsrmd5.exe --flat Pictures > pictures.md5

通常モードとの違い
- パスはファイル名のみを出力する
- 出力結果はファイル名でソートされる
