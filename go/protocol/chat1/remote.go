// Auto-generated by avdl-compiler v1.3.7 (https://github.com/keybase/node-avdl-compiler)
//   Input file: avdl/chat1/remote.avdl

package chat1

import (
	rpc "github.com/keybase/go-framed-msgpack-rpc"
	context "golang.org/x/net/context"
)

type MessageBoxed struct {
	ServerHeader     *MessageServerHeader `codec:"serverHeader,omitempty" json:"serverHeader,omitempty"`
	SupersededBy     *MessageBoxed        `codec:"supersededBy,omitempty" json:"supersededBy,omitempty"`
	ClientHeader     MessageClientHeader  `codec:"clientHeader" json:"clientHeader"`
	HeaderCiphertext EncryptedData        `codec:"headerCiphertext" json:"headerCiphertext"`
	BodyCiphertext   EncryptedData        `codec:"bodyCiphertext" json:"bodyCiphertext"`
	KeyGeneration    int                  `codec:"keyGeneration" json:"keyGeneration"`
}

type ThreadViewBoxed struct {
	Messages   []MessageBoxed `codec:"messages" json:"messages"`
	Pagination *Pagination    `codec:"pagination,omitempty" json:"pagination,omitempty"`
}

type GetInboxRemoteRes struct {
	Inbox     InboxView  `codec:"inbox" json:"inbox"`
	RateLimit *RateLimit `codec:"rateLimit,omitempty" json:"rateLimit,omitempty"`
}

type GetInboxByTLFIDRemoteRes struct {
	Convs     []Conversation `codec:"convs" json:"convs"`
	RateLimit *RateLimit     `codec:"rateLimit,omitempty" json:"rateLimit,omitempty"`
}

type GetThreadRemoteRes struct {
	Thread    ThreadViewBoxed `codec:"thread" json:"thread"`
	RateLimit *RateLimit      `codec:"rateLimit,omitempty" json:"rateLimit,omitempty"`
}

type GetConversationMetadataRemoteRes struct {
	Conv      Conversation `codec:"conv" json:"conv"`
	RateLimit *RateLimit   `codec:"rateLimit,omitempty" json:"rateLimit,omitempty"`
}

type PostRemoteRes struct {
	MsgID     MessageID  `codec:"msgID" json:"msgID"`
	RateLimit *RateLimit `codec:"rateLimit,omitempty" json:"rateLimit,omitempty"`
}

type NewConversationRemoteRes struct {
	ConvID    ConversationID `codec:"convID" json:"convID"`
	RateLimit *RateLimit     `codec:"rateLimit,omitempty" json:"rateLimit,omitempty"`
}

type GetMessagesRemoteRes struct {
	Msgs      []MessageBoxed `codec:"msgs" json:"msgs"`
	RateLimit *RateLimit     `codec:"rateLimit,omitempty" json:"rateLimit,omitempty"`
}

type MarkAsReadRes struct {
	RateLimit *RateLimit `codec:"rateLimit,omitempty" json:"rateLimit,omitempty"`
}

type GetInboxRemoteArg struct {
	Query      *GetInboxQuery `codec:"query,omitempty" json:"query,omitempty"`
	Pagination *Pagination    `codec:"pagination,omitempty" json:"pagination,omitempty"`
}

type GetThreadRemoteArg struct {
	ConversationID ConversationID  `codec:"conversationID" json:"conversationID"`
	Query          *GetThreadQuery `codec:"query,omitempty" json:"query,omitempty"`
	Pagination     *Pagination     `codec:"pagination,omitempty" json:"pagination,omitempty"`
}

type PostRemoteArg struct {
	ConversationID ConversationID `codec:"conversationID" json:"conversationID"`
	MessageBoxed   MessageBoxed   `codec:"messageBoxed" json:"messageBoxed"`
}

type NewConversationRemoteArg struct {
	IdTriple ConversationIDTriple `codec:"idTriple" json:"idTriple"`
}

type NewConversationRemote2Arg struct {
	IdTriple   ConversationIDTriple `codec:"idTriple" json:"idTriple"`
	TLFMessage MessageBoxed         `codec:"TLFMessage" json:"TLFMessage"`
}

type GetMessagesRemoteArg struct {
	ConversationID ConversationID `codec:"conversationID" json:"conversationID"`
	MessageIDs     []MessageID    `codec:"messageIDs" json:"messageIDs"`
}

type MarkAsReadArg struct {
	ConversationID ConversationID `codec:"conversationID" json:"conversationID"`
	MsgID          MessageID      `codec:"msgID" json:"msgID"`
}

type TlfFinalizeArg struct {
	TlfID TLFID `codec:"tlfID" json:"tlfID"`
}

type RemoteInterface interface {
	GetInboxRemote(context.Context, GetInboxRemoteArg) (GetInboxRemoteRes, error)
	GetThreadRemote(context.Context, GetThreadRemoteArg) (GetThreadRemoteRes, error)
	PostRemote(context.Context, PostRemoteArg) (PostRemoteRes, error)
	NewConversationRemote(context.Context, ConversationIDTriple) (NewConversationRemoteRes, error)
	NewConversationRemote2(context.Context, NewConversationRemote2Arg) (NewConversationRemoteRes, error)
	GetMessagesRemote(context.Context, GetMessagesRemoteArg) (GetMessagesRemoteRes, error)
	MarkAsRead(context.Context, MarkAsReadArg) (MarkAsReadRes, error)
	TlfFinalize(context.Context, TLFID) error
}

func RemoteProtocol(i RemoteInterface) rpc.Protocol {
	return rpc.Protocol{
		Name: "chat.1.remote",
		Methods: map[string]rpc.ServeHandlerDescription{
			"getInboxRemote": {
				MakeArg: func() interface{} {
					ret := make([]GetInboxRemoteArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]GetInboxRemoteArg)
					if !ok {
						err = rpc.NewTypeError((*[]GetInboxRemoteArg)(nil), args)
						return
					}
					ret, err = i.GetInboxRemote(ctx, (*typedArgs)[0])
					return
				},
				MethodType: rpc.MethodCall,
			},
			"getThreadRemote": {
				MakeArg: func() interface{} {
					ret := make([]GetThreadRemoteArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]GetThreadRemoteArg)
					if !ok {
						err = rpc.NewTypeError((*[]GetThreadRemoteArg)(nil), args)
						return
					}
					ret, err = i.GetThreadRemote(ctx, (*typedArgs)[0])
					return
				},
				MethodType: rpc.MethodCall,
			},
			"postRemote": {
				MakeArg: func() interface{} {
					ret := make([]PostRemoteArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]PostRemoteArg)
					if !ok {
						err = rpc.NewTypeError((*[]PostRemoteArg)(nil), args)
						return
					}
					ret, err = i.PostRemote(ctx, (*typedArgs)[0])
					return
				},
				MethodType: rpc.MethodCall,
			},
			"newConversationRemote": {
				MakeArg: func() interface{} {
					ret := make([]NewConversationRemoteArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]NewConversationRemoteArg)
					if !ok {
						err = rpc.NewTypeError((*[]NewConversationRemoteArg)(nil), args)
						return
					}
					ret, err = i.NewConversationRemote(ctx, (*typedArgs)[0].IdTriple)
					return
				},
				MethodType: rpc.MethodCall,
			},
			"newConversationRemote2": {
				MakeArg: func() interface{} {
					ret := make([]NewConversationRemote2Arg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]NewConversationRemote2Arg)
					if !ok {
						err = rpc.NewTypeError((*[]NewConversationRemote2Arg)(nil), args)
						return
					}
					ret, err = i.NewConversationRemote2(ctx, (*typedArgs)[0])
					return
				},
				MethodType: rpc.MethodCall,
			},
			"getMessagesRemote": {
				MakeArg: func() interface{} {
					ret := make([]GetMessagesRemoteArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]GetMessagesRemoteArg)
					if !ok {
						err = rpc.NewTypeError((*[]GetMessagesRemoteArg)(nil), args)
						return
					}
					ret, err = i.GetMessagesRemote(ctx, (*typedArgs)[0])
					return
				},
				MethodType: rpc.MethodCall,
			},
			"markAsRead": {
				MakeArg: func() interface{} {
					ret := make([]MarkAsReadArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]MarkAsReadArg)
					if !ok {
						err = rpc.NewTypeError((*[]MarkAsReadArg)(nil), args)
						return
					}
					ret, err = i.MarkAsRead(ctx, (*typedArgs)[0])
					return
				},
				MethodType: rpc.MethodCall,
			},
			"tlfFinalize": {
				MakeArg: func() interface{} {
					ret := make([]TlfFinalizeArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]TlfFinalizeArg)
					if !ok {
						err = rpc.NewTypeError((*[]TlfFinalizeArg)(nil), args)
						return
					}
					err = i.TlfFinalize(ctx, (*typedArgs)[0].TlfID)
					return
				},
				MethodType: rpc.MethodCall,
			},
		},
	}
}

type RemoteClient struct {
	Cli rpc.GenericClient
}

func (c RemoteClient) GetInboxRemote(ctx context.Context, __arg GetInboxRemoteArg) (res GetInboxRemoteRes, err error) {
	err = c.Cli.Call(ctx, "chat.1.remote.getInboxRemote", []interface{}{__arg}, &res)
	return
}

func (c RemoteClient) GetThreadRemote(ctx context.Context, __arg GetThreadRemoteArg) (res GetThreadRemoteRes, err error) {
	err = c.Cli.Call(ctx, "chat.1.remote.getThreadRemote", []interface{}{__arg}, &res)
	return
}

func (c RemoteClient) PostRemote(ctx context.Context, __arg PostRemoteArg) (res PostRemoteRes, err error) {
	err = c.Cli.Call(ctx, "chat.1.remote.postRemote", []interface{}{__arg}, &res)
	return
}

func (c RemoteClient) NewConversationRemote(ctx context.Context, idTriple ConversationIDTriple) (res NewConversationRemoteRes, err error) {
	__arg := NewConversationRemoteArg{IdTriple: idTriple}
	err = c.Cli.Call(ctx, "chat.1.remote.newConversationRemote", []interface{}{__arg}, &res)
	return
}

func (c RemoteClient) NewConversationRemote2(ctx context.Context, __arg NewConversationRemote2Arg) (res NewConversationRemoteRes, err error) {
	err = c.Cli.Call(ctx, "chat.1.remote.newConversationRemote2", []interface{}{__arg}, &res)
	return
}

func (c RemoteClient) GetMessagesRemote(ctx context.Context, __arg GetMessagesRemoteArg) (res GetMessagesRemoteRes, err error) {
	err = c.Cli.Call(ctx, "chat.1.remote.getMessagesRemote", []interface{}{__arg}, &res)
	return
}

func (c RemoteClient) MarkAsRead(ctx context.Context, __arg MarkAsReadArg) (res MarkAsReadRes, err error) {
	err = c.Cli.Call(ctx, "chat.1.remote.markAsRead", []interface{}{__arg}, &res)
	return
}

func (c RemoteClient) TlfFinalize(ctx context.Context, tlfID TLFID) (err error) {
	__arg := TlfFinalizeArg{TlfID: tlfID}
	err = c.Cli.Call(ctx, "chat.1.remote.tlfFinalize", []interface{}{__arg}, nil)
	return
}
