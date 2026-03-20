# Razer Battery Monitor

Razer Basilisk Ultimate のバッテリー残量を Windows タスクバー（システムトレイ）に表示する常駐アプリ。

**Python 不要。単体 exe で動作します。**

---

## インストールと使い方

1. [Releases](https://github.com/4quarius-crate/razer-battery-monitor/releases) から `RazerBatteryMonitor.exe` と `data/config.toml` をダウンロード
2. 同じフォルダに配置
3. `RazerBatteryMonitor.exe` を起動
4. タスクバー右下のトレイアイコンにバッテリー残量が表示される

---

## 設定

`data/config.toml` を編集して動作をカスタマイズできます。

```toml
[monitor]
poll_interval_normal = 300   # ポーリング間隔（秒）

[alert]
low_battery_threshold = 20   # 低バッテリー通知の閾値（%）

[fps_guard]
enabled = true
game_processes = ["cs2.exe", "valorant.exe"]   # ゲーム中はHID通信をスキップ

[debug]
dump_raw_bytes = false   # trueにすると生バイトをログ出力（デバッグ用）
```

---

## 自分でビルドする場合

### 必要なもの

- [Go 1.21+](https://go.dev/dl/)
- [MinGW-w64](https://www.mingw-w64.org/)（CGO に必要）

### ビルド手順

```bat
go mod download
build.bat
```

または GitHub にプッシュすると Actions が自動でビルドして exe を生成します。

---

## トラブルシューティング

| 症状 | 対処 |
|---|---|
| トレイに表示されない | Razer USB レシーバーが接続されているか確認 |
| バッテリーが取得できない | Razer Synapse を終了してから再起動 |
| アクセス拒否エラー | 右クリック→「管理者として実行」で起動 |

---

## 対応マウス

- Razer Basilisk Ultimate
- Razer DeathAdder V2 Pro
- Razer Viper Ultimate
- Razer Naga Pro
