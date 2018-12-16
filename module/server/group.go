// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
    "fmt"

    "github.com/tinystack/goweb"
    "github.com/tinystack/goutil"
    "github.com/tinystack/syncd"
    "github.com/tinystack/syncd/route"
    groupModel "github.com/tinystack/syncd/model/server/group"
    baseModel "github.com/tinystack/syncd/model"

    serverService "github.com/tinystack/syncd/service/server"
)

func init() {
    route.Register(route.API_SERVER_GROUP_UPDATE, updateServerGroup)
    route.Register(route.API_SERVER_GROUP_LIST, listServerGroup)
    route.Register(route.API_SERVER_GROUP_DETAIL, detailServerGroup)
    route.Register(route.API_SERVER_GROUP_DELETE, deleteServerGroup)
}

func updateServerGroup(c *goweb.Context) error {
    groupId := c.PostFormInt("id")
    groupName := c.PostForm("name")
    if groupName == "" {
        syncd.RenderParamError(c, "group name can not empty")
        return nil
    }
    var ok bool
    g := groupModel.ServerGroup{
        Name: groupName,
    }
    if groupId > 0 {
        ok = groupModel.Update(groupId, g)
    } else {
        ok = groupModel.Create(&g)
    }
    if !ok {
        syncd.RenderAppError(c, "server group data update failed")
        return nil
    }
    syncd.RenderJson(c, nil)
    return nil
}

func listServerGroup(c *goweb.Context) error {
    var (
        total, offset, limit, groupId int
        ok bool
        keyword string
        where []baseModel.WhereParam
    )
    offset, limit = c.QueryInt("offset"), c.QueryInt("limit")
    keyword = c.Query("keyword")
    if keyword != "" {
        if goutil.IsInteger(keyword) {
            groupId = c.QueryInt("keyword")
            if groupId > 0 {
                where = append(where, baseModel.WhereParam{
                    Field: "id",
                    Prepare: groupId,
                })
            }
        } else {
            where = append(where, baseModel.WhereParam{
                Field: "name",
                Tag: "LIKE",
                Prepare: fmt.Sprintf("%%%s%%", keyword),
            })
        }
    }

    list, ok := groupModel.List(baseModel.QueryParam{
        Fields: "id, name",
        Offset: offset,
        Limit: limit,
        Order: "id DESC",
        Where: where,
    })
    if !ok {
        syncd.RenderAppError(c, "get server group list data failed")
        return nil
    }

    total, ok = groupModel.Total(baseModel.QueryParam{
        Where: where,
    })
    if !ok {
        syncd.RenderAppError(c, "get server group total count failed")
        return nil
    }

    syncd.RenderJson(c, goweb.JSON{
        "list": list,
        "total": total,
    })
    return nil
}

func detailServerGroup(c *goweb.Context) error {
    id := c.QueryInt("id")
    if id == 0 {
        syncd.RenderParamError(c, "id can not be empty")
        return nil
    }
    detail, ok := groupModel.Get(id)
    if !ok {
        syncd.RenderAppError(c, "get server group detail data failed")
        return nil
    }
    syncd.RenderJson(c, goweb.JSON{
        "detail": detail,
    })
    return nil
}

func deleteServerGroup(c *goweb.Context) error {
    id := c.PostFormInt("id")
    if id == 0 {
        syncd.RenderParamError(c, "id can not be empty")
        return nil
    }
    ok := groupModel.Delete(id)
    if !ok {
        syncd.RenderAppError(c, "delete server group data failed")
        return nil
    }
    syncd.RenderJson(c, nil)
    return nil
}