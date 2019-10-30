// asynchronous network,support tcp,kcp,websocket, etc.
package anet

import (
	"net"

	"github.com/jeckbjy/gsk/util/buffer"
)

var Default NewTranFunc
var gTranFuncMap = make(map[string]NewTranFunc)

func Add(name string, fn NewTranFunc) {
	gTranFuncMap[name] = fn
	Default = fn
}

func New(name string) Tran {
	if fn, ok := gTranFuncMap[name]; ok {
		return fn()
	}

	return nil
}

// NewDefault 新建一个默认的Transport
func NewDefault() Tran {
	return Default()
}

type NewTranFunc func() Tran

// Tran 创建Conn,可以是tcp,websocket等协议
// 不同的Tran可以配置不同的FilterChain
type Tran interface {
	String() string
	GetChain() FilterChain
	SetChain(chain FilterChain)
	AddFilters(filters ...Filter)
	Dial(addr string, opts ...DialOption) (Conn, error)
	Listen(addr string, opts ...ListenOption) (Listener, error)
	Close() error
}

type Status int

const (
	DISCONNECTED = Status(iota)
	CONNECTED
	CLOSED
	RECONNECTING
	CONNECTING
)

// Conn 异步收发消息
type Conn interface {
	Tag() string                     // 额外标识类型
	Get(key string) interface{}      // 获取自定义数据
	Set(key string, val interface{}) // 设置自定义数据
	Status() Status                  // Socket状态
	LocalAddr() net.Addr             // 本地地址
	RemoteAddr() net.Addr            // 远程地址
	Read() *buffer.Buffer            // 异步读缓存,非线程安全,通常在一个线程中解析消息,在分发到其他线程处理消息
	Write(data *buffer.Buffer) error // 异步写数据,线程安全
	Send(msg interface{}) error      // 异步发消息,会调用HandleWrite,没有连接成功时也可以发送,当连接成功后会自动发送缓存数据
	Close() error
}

type Listener interface {
	Close() error
	Addr() net.Addr
}

// Filter 用于链式处理IConn各种回调
type Filter interface {
	Name() string
	HandleRead(ctx FilterCtx) error
	HandleWrite(ctx FilterCtx) error
	HandleOpen(ctx FilterCtx) error
	HandleClose(ctx FilterCtx) error
	HandleError(ctx FilterCtx) error
}

// FilterCtx Filter上下文，默认会自动调用Next,如需终止，需要主动调用Abort
type FilterCtx interface {
	Conn() Conn               // Socket Connection
	Data() interface{}        // 获取数据
	SetData(data interface{}) // 设置数据
	Error() error             // 错误信息
	SetError(err error)       // 设置错误信息
	IsAbort() bool            // 是否已经强制终止
	Abort()                   // 终止调用
	Next()                    // 调用下一个
	Jump(index int) error     // 跳转到指定位置,可以是负索引
	JumpBy(name string) error // 通过名字跳转
	Clone() FilterCtx         // 拷贝当前状态,可用于转移到其他协程中继续执行
	Call()                    // 从当前位置开始执行
}

// FilterChain 管理Filter,并链式调用所有Filter
// Filter分为Inbound和Outbound
// InBound: 从前向后执行,包括Read,Open,Error
// OutBound:从后向前执行,包括Write,Close
type FilterChain interface {
	Len() int                   // 长度
	Front() Filter              // 第一个
	Back() Filter               // 最后一个
	Get(index int) Filter       // 通过索引获取filter
	Index(name string) int      // 通过名字查询索引
	AddFirst(filters ...Filter) // 在前边插入
	AddLast(filters ...Filter)  // 在末尾插入
	HandleOpen(conn Conn)
	HandleClose(conn Conn)
	HandleRead(conn Conn, msg interface{})
	HandleWrite(conn Conn, msg interface{})
	HandleError(conn Conn, err error)
}
