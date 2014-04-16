// Copyright 2014, The Serviced Authors. All rights reserved.
// Use of sc source code is governed by a
// license that can be found in the LICENSE file.

package web

import (
	"github.com/zenoss/go-json-rest"
	"github.com/zenoss/glog"
	"github.com/zenoss/serviced/domain/pool"

	"github.com/zenoss/serviced/dao"
	"net/url"
)

func getPoolRoutes(sc *ServiceConfig) []rest.Route {
	return []rest.Route{
		rest.Route{"POST", "/pools/add", sc.CheckAuth(RestAddPool)},
		rest.Route{"GET", "/pools/:poolId/hosts", sc.CheckAuth(RestGetHostsForResourcePool)},
		rest.Route{"DELETE", "/pools/:poolId", sc.CheckAuth(RestRemovePool)},
		rest.Route{"PUT", "/pools/:poolId", sc.CheckAuth(RestUpdatePool)},
		rest.Route{"GET", "/pools", sc.CheckAuth(RestGetPools)},
	}
}

func RestGetPools(w *rest.ResponseWriter, r *rest.Request, ctx *requestContext) {
	client, err := ctx.getMasterClient()
	if err != nil {
		RestServerError(w)
		return
	}

	pools, err := client.GetResourcePools()
	if err != nil {
		glog.Error("Could not get resource pools: ", err)
		RestServerError(w)
		return
	}
	var poolsMap map[string]*pool.ResourcePool
	for _, pool := range pools {
		poolsMap[pool.ID] = pool
	}
	w.WriteJson(&poolsMap)
}

func RestAddPool(w *rest.ResponseWriter, r *rest.Request, ctx *requestContext) {
	var payload pool.ResourcePool
	err := r.DecodeJsonPayload(&payload)
	if err != nil {
		glog.V(1).Info("Could not decode pool payload: ", err)
		RestBadRequest(w)
		return
	}
	client, err := ctx.getMasterClient()
	if err != nil {
		RestServerError(w)
		return
	}

	err = client.AddResourcePool(payload)
	if err != nil {
		glog.Error("Unable to add pool: ", err)
		RestServerError(w)
		return
	}
	glog.V(0).Info("Added pool ", payload.ID)
	w.WriteJson(&SimpleResponse{"Added resource pool", poolLinks(payload.ID)})
}

func RestUpdatePool(w *rest.ResponseWriter, r *rest.Request, ctx *requestContext) {
	poolId, err := url.QueryUnescape(r.PathParam("poolId"))
	if err != nil {
		RestBadRequest(w)
		return
	}
	var payload pool.ResourcePool
	err = r.DecodeJsonPayload(&payload)
	if err != nil {
		glog.V(1).Info("Could not decode pool payload: ", err)
		RestBadRequest(w)
		return
	}
	client, err := ctx.getMasterClient()
	if err != nil {
		RestServerError(w)
		return
	}
	err = client.UpdateResourcePool(payload)
	if err != nil {
		glog.Error("Unable to update pool: ", err)
		RestServerError(w)
		return
	}
	glog.V(1).Info("Updated pool ", poolId)
	w.WriteJson(&SimpleResponse{"Updated resource pool", poolLinks(poolId)})
}

func RestRemovePool(w *rest.ResponseWriter, r *rest.Request, ctx *requestContext) {
	poolId, err := url.QueryUnescape(r.PathParam("poolId"))
	if err != nil {
		RestBadRequest(w)
		return
	}
	client, err := ctx.getMasterClient()
	if err != nil {
		RestServerError(w)
		return
	}
	err = client.RemoveResourcePool(poolId)
	if err != nil {
		glog.Error("Could not remove resource pool: ", err)
		RestServerError(w)
		return
	}
	glog.V(0).Info("Removed pool ", poolId)
	w.WriteJson(&SimpleResponse{"Removed resource pool", poolsLinks()})
}

func RestGetHostsForResourcePool(w *rest.ResponseWriter, r *rest.Request, ctx *requestContext) {
	poolHosts := make([]*dao.PoolHost, 0)
	poolId, err := url.QueryUnescape(r.PathParam("poolId"))
	if err != nil {
		glog.V(1).Infof("Unable to acquire pool ID: %v", err)
		RestBadRequest(w)
		return
	}
	client, err := ctx.getMasterClient()
	if err != nil {
		RestServerError(w)
		return
	}
	hosts, err := client.FindHostsInPool(poolId)
	if err != nil {
		glog.Errorf("Could not get hosts: %v", err)
		RestServerError(w)
		return
	}
	for _, host := range hosts {
		ph := dao.PoolHost{
			HostId: host.ID,
			PoolId: poolId,
			HostIp: host.IPAddr,
		}
		poolHosts = append(poolHosts, &ph)
	}
	glog.V(2).Infof("Returning %d hosts for pool %s", len(poolHosts), poolId)
	w.WriteJson(&poolHosts)
}
