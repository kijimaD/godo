package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"

	"golang.org/x/sync/errgroup"
)

// HTTPサーバが起動しているか、テストコードから意図通りに終了するかチェック
func TestRun(t *testing.T) {
	t.Skip("wip")

	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen port %v", err)
	}

	// キャンセル可能なcontext.Contextオブジェクトを作成する
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)

	// 別ゴルーチンでテスト対象のrun関数を実行する
	eg.Go(func() error {
		return run(ctx)
	})

	// エンドポイントに対してGETリクエストを送信する
	in := "message"
	url := fmt.Sprintf("http://%s/%s", l.Addr().String(), in)

	// 使用ポートを確認する
	// t.Logf("try request to %q", url)

	rsp, err := http.Get(url)
	if err != nil {
		t.Errorf("failed to get: %+v", err)
	}
	defer rsp.Body.Close()
	got, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	// HTTPサーバの戻り値を検証する
	want := fmt.Sprintf("Hello, %s!", in)
	if string(got) != want {
		t.Errorf("want %q, but got %q", want, got)
	}

	// run関数に終了通知を送信する
	cancel()
	// run関数の戻り値を検証する
	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}
}
