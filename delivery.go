package main

import "context"

type AssetJob struct {
	AssetID   string `json:"asset_id"`
	AccountID string `json:"account_id"`
	Ready     bool   `json:"ready"`
}

func shouldNotify(job AssetJob) bool { return job.Ready && job.AssetID != "" && job.AccountID != "" }

func deliver(ctx context.Context, c *Client, job AssetJob) error {
	if !shouldNotify(job) {
		return nil
	}
	channel := "creator:" + job.AccountID
	if err := c.CreateChannel(ctx, channel, "private"); err != nil {
		return err
	}
	return c.Publish(ctx, channel, "asset.ready", map[string]string{"asset_id": job.AssetID}, job.AccountID)
}
