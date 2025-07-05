package fn

import "github.com/jxs1211/keda-external-scaler/pkg/util/debug"

var dlog = debug.NewLogger()

func RunDebugLog() {
	dlog.Printf("my dlog output is %s", "awsome")
}
