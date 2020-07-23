package meta

import (
	"context"
	"fmt"
	"time"

	"github.com/influxdata/influxdb/services/meta"
	v2 "github.com/influxdata/influxdb/servicesv2"
)

type Client struct {
	BucketService      v2.BucketService
	DBRPMappingService v2.DBRPMappingServiceV2
	ShardGroupService  v2.ShardGroupService
}

func (c *Client) Database(db string) *meta.DatabaseInfo {
	dbrps, count, err := c.DBRPMappingService.FindMany(
		context.Background(),
		v2.DBRPMappingFilterV2{Database: &db},
	)
	if err != nil {
		return nil
	} else if count != 1 {
		return nil
	}
	dbrp := dbrps[0]
	dbinfo := meta.DatabaseInfo{
		Name: db,
	}

	rp := dbrp.RetentionPolicy
	if dbrp.Default {
		dbinfo.DefaultRetentionPolicy = rp
	}

	rpi, err := c.RetentionPolicy(db, rp)
	if err != nil {
		return nil
	}
	dbinfo.RetentionPolicies = append(dbinfo.RetentionPolicies, *rpi)

	return &dbinfo
}

func (c *Client) RetentionPolicy(db, rp string) (*meta.RetentionPolicyInfo, error) {
	dbrps, count, err := c.DBRPMappingService.FindMany(context.Background(), v2.DBRPMappingFilterV2{
		Database:        &db,
		RetentionPolicy: &rp,
	})
	if err != nil {
		return nil, err
	}
	if count != 1 {
		return nil, fmt.Errorf("expected 1 value - got %d", count)
	}

	dbrp := dbrps[0]
	rpi := meta.RetentionPolicyInfo{
		Name:     dbrp.RetentionPolicy,
		ReplicaN: 1,
		// We need to populate the duration values here
		// Duration
		// ShardGroupDuration
		// ShardGroups
		// Subscriptions - this one may be unnecessary
	}

	bucket, err := c.BucketService.FindBucket(context.Background(), v2.BucketFilter{
		ID: &dbrp.BucketID,
	})
	if err != nil {
		return nil, err
	}
	rpi.Duration = bucket.RetentionPeriod
	rpi.ShardGroupDuration = bucket.RetentionPeriod
	return &rpi, nil
}

func (c *Client) CreateShardGroup(db, rp string, timestamp time.Time) (*meta.ShardGroupInfo, error) {
	dbrps, count, err := c.DBRPMappingService.FindMany(context.Background(), v2.DBRPMappingFilterV2{
		Database:        &db,
		RetentionPolicy: &rp,
	})
	if err != nil {
		return nil, err
	} else if count != 1 {
		return nil, fmt.Errorf("expected 1 DBRP - got %d", count)
	}
	dbrp := dbrps[0]
	sgi, err := c.ShardGroupService.CreateShardGroup(context.Background(), dbrp.BucketID, timestamp)
	if err != nil {
		return nil, err
	}
	return sgi, nil
}
