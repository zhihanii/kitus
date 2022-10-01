package v1

import (
    "context"
    "github.com/zhihanii/kitus"
)

type VideoClient interface {
    Feed(context.Context, *FeedRequest) (*FeedResponse, error)
    PublishAction(context.Context, *PublishActionRequest) (*PublishActionResponse, error)
    PublishList(context.Context, *PublishListRequest) (*PublishListResponse, error)
}

type videoClient struct {
    kc kitus.Client
}

func (c *videoClient) Feed(ctx context.Context, req *FeedRequest) (*FeedResponse, error) {
    resp := new(FeedResponse)
    err := c.kc.Call(ctx, "/examples.api.video.v1.Video/Feed", req, resp)
    if err != nil {
        return nil, err
    }
    return resp, nil
}

func (c *videoClient) PublishAction(ctx context.Context, req *PublishActionRequest) (*PublishActionResponse, error) {
    resp := new(PublishActionResponse)
    err := c.kc.Call(ctx, "/examples.api.video.v1.Video/PublishAction", req, resp)
    if err != nil {
        return nil, err
    }
    return resp, nil
}

func (c *videoClient) PublishList(ctx context.Context, req *PublishListRequest) (*PublishListResponse, error) {
    resp := new(PublishListResponse)
    err := c.kc.Call(ctx, "/examples.api.video.v1.Video/PublishList", req, resp)
    if err != nil {
        return nil, err
    }
    return resp, nil
}

type VideoServer interface {
    Feed(context.Context, *FeedRequest) (*FeedResponse, error)
    PublishAction(context.Context, *PublishActionRequest) (*PublishActionResponse, error)
    PublishList(context.Context, *PublishListRequest) (*PublishListResponse, error)
}

func RegisterVideoServer(s kitus.Server, srv VideoServer) {
    s.RegisterService(&Video_ServiceInfo, srv)
}

func _Video_Feed_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
    req := new(FeedRequest)
    if err := dec(req); err != nil {
        return nil, err
    }
    return srv.(VideoServer).Feed(ctx, req)
}

func _Video_PublishAction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
    req := new(PublishActionRequest)
    if err := dec(req); err != nil {
        return nil, err
    }
    return srv.(VideoServer).PublishAction(ctx, req)
}

func _Video_PublishList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
    req := new(PublishListRequest)
    if err := dec(req); err != nil {
        return nil, err
    }
    return srv.(VideoServer).PublishList(ctx, req)
}

var Video_ServiceInfo = kitus.ServiceInfo{
    ServiceName: "examples.api.video.v1.Video",
    HandlerType: (*VideoServer)(nil),
    Methods: map[string]*kitus.MethodInfo{
        "Feed": {
            MethodName: "Feed",
            Handler: _Video_Feed_Handler,
        },
        "PublishAction": {
            MethodName: "PublishAction",
            Handler: _Video_PublishAction_Handler,
        },
        "PublishList": {
            MethodName: "PublishList",
            Handler: _Video_PublishList_Handler,
        },
    },
}

