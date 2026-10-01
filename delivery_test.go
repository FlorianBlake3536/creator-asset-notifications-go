package main

import "testing"

func TestDeliveryOnlyPublishesReadyAssets(t *testing.T) {
	for _, tc := range []struct {
		name                    string
		ready                   bool
		wantCreate, wantPublish bool
	}{
		{"queued", false, false, false}, {"processed", true, true, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := shouldNotify(AssetJob{AssetID: "asset", AccountID: "creator", Ready: tc.ready})
			if got != tc.wantPublish {
				t.Fatalf("ready=%v notify=%v", tc.ready, got)
			}
		})
	}
}
