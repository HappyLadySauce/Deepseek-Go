package models

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

// 基于 Eino 的 DeepSeek Chat 模型
// 需要实现以下接口
// type ChatModel interface {
//     Generate(ctx context.Context, input []*schema.Message, opts ...Option) (*schema.Message, error)
//     Stream(ctx context.Context, input []*schema.Message, opts ...Option) (*schema.StreamReader[*schema.Message], error)
//     BindTools(tools []*schema.ToolInfo) error
// }

// Generate 方法
// 功能：生成完整的模型响应
// 参数：
// ctx：上下文对象，用于传递请求级别的信息，同时也用于传递 Callback Manager
// input：输入消息列表
// opts：可选参数，用于配置模型行为
// 返回值：
// *schema.Message：模型生成的响应消息
// error：生成过程中的错误信息

// Stream 方法
// 功能：以流式方式生成模型响应
// 参数：与 Generate 方法相同
// 返回值：
// *schema.StreamReader[*schema.Message]：模型响应的流式读取器
// error：生成过程中的错误信息

// BindTools 方法
// 功能：为模型绑定可用的工具
// 参数：
// tools：工具信息列表
// 返回值：
// error：绑定过程中的错误信息

// eino/schema/message.go
// 消息结构体
// type Message struct {
//     // Role 表示消息的角色（system/user/assistant/tool）
//     Role RoleType
//     // Content 是消息的文本内容
//     Content string
//     // MultiContent 是多模态内容，支持文本、图片、音频等
//     MultiContent []ChatMessagePart
//     // Name 是消息的发送者名称
//     Name string
//     // ToolCalls 是 assistant 消息中的工具调用信息
//     ToolCalls []ToolCall
//     // ToolCallID 是 tool 消息的工具调用 ID
//     ToolCallID string
//     // ResponseMeta 包含响应的元信息
//     ResponseMeta *ResponseMeta
//     // Extra 用于存储额外信息
//     Extra map[string]any
// }

// DeepSeekModel 是DeepSeek的模型
type DeepSeekModel struct {
	client     *http.Client
	apiKey     string
	baseURL    string
	model      string
	timeout    time.Duration
	retryCount int
}

// DeepSeekModelOptions 是DeepSeek的模型选项
type DeepSeekModelOptions struct {
	// eino 公共模型选项
	Options *model.Options
	// 重试次数
	RetryCount int
	// 超时时间
	Timeout time.Duration
}

// WithRetryCount 设置重试次数
func WithRetryCount(retryCount int) model.Option {
	return model.WrapImplSpecificOptFn(func(o *DeepSeekModelOptions) {
		o.RetryCount = retryCount
	})
}

// WithTimeout 设置超时时间
func WithTimeout(timeout time.Duration) model.Option {
	return model.WrapImplSpecificOptFn(func(o *DeepSeekModelOptions) {
		o.Timeout = timeout
	})
}

// DeepSeekModelConfig 是DeepSeek的模型配置
type DeepSeekModelConfig struct {
	APIKey     string
	BaseURL    string
	Model      string
	Timeout    time.Duration
	RetryCount int
}

// NewDeepSeekModel 创建一个DeepSeek模型
func NewDeepSeekModel(config DeepSeekModelConfig) (*DeepSeekModel, error) {

	// 检查APIKey是否设置
	if config.APIKey == "" {
		return nil, errors.New("NewDeepSeekModel: APIKey没有设置")
	}

	// 检查BaseURL是否设置
	if config.BaseURL == "" {
		config.BaseURL = "https://api.deepseek.com"
	}

	// 检查模型是否设置
	if config.Model == "" {
		// 默认使用deepseek-chat模型
		log.Println("NewDeepSeekModel: 未设置模型，使用默认模型deepseek-chat")
		config.Model = "deepseek-chat"
	}

	// 检查超时时间是否设置
	if config.Timeout <= 0 {
		// 默认超时时间为10秒
		log.Println("NewDeepSeekModel: 未设置超时时间，使用默认超时时间10秒")
		config.Timeout = 10 * time.Second
	}

	// 检查重试次数是否设置
	if config.RetryCount <= 0 {
		// 默认重试次数为3次
		log.Println("NewDeepSeekModel: 未设置重试次数，使用默认重试次数3次")
		config.RetryCount = 3
	}

	// 创建DeepSeek模型
	return &DeepSeekModel{
		client:     &http.Client{},
		apiKey:     config.APIKey,
		baseURL:    config.BaseURL,
		model:      config.Model,
		timeout:    config.Timeout,
		retryCount: config.RetryCount,
	}, nil
}

// model.Options 是模型选项
// type Options struct {
//     // Temperature is the temperature for the model, which controls the randomness of the model.
//     Temperature *float32
//     // MaxTokens is the max number of tokens, if reached the max tokens, the model will stop generating, and mostly return an finish reason of "length".
//     MaxTokens *int
//     // Model is the model name.
//     Model *string
//     // TopP is the top p for the model, which controls the diversity of the model.
//     TopP *float32
//     // Stop is the stop words for the model, which controls the stopping condition of the model.
//     Stop []string
//     // Tools is a list of tools the model may call.
//     Tools []*schema.ToolInfo
//     // ToolChoice controls which tool is called by the model.
//     ToolChoice *schema.ToolChoice
// }

// DeepSeekModel 的 Generate 方法
func (m *DeepSeekModel) Generate(ctx context.Context, messages []*schema.Message, opts ...model.Option) (*schema.Message, error) {
	// 1. 处理选项
	// 加载选项
	options := &DeepSeekModelOptions{
		Options: &model.Options{
			Model: &m.model,
		},
		RetryCount: m.retryCount,
		Timeout:    m.timeout,
	}
	// 加载通用选项
	options.Options = model.GetCommonOptions(options.Options, opts...)
	// 加载实现特定选项
	options = model.GetImplSpecificOptions(options, opts...)

	// 2. 开始生成前的回调
	ctx = callbacks.OnStart(ctx, &model.CallbackInput{
		Messages: messages,
		Config: &model.Config{
			Model: *options.Options.Model,
		},
	})

	// 3. 执行生成逻辑
	response, err := m.doGenerate(ctx, messages, options)

	// 4. 处理错误和完成回调
	if err != nil {
		ctx = callbacks.OnError(ctx, err)
		return nil, err
	}

	ctx = callbacks.OnEnd(ctx, &model.CallbackOutput{
		Message: response,
	})

	return response, nil
}

func (m *DeepSeekModel) Stream(ctx context.Context, messages []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	// 1. 处理选项
	// 加载选项
	options := &DeepSeekModelOptions{
		Options: &model.Options{
			Model: &m.model,
		},
		RetryCount: m.retryCount,
		Timeout:    m.timeout,
	}
	// 加载通用选项
	options.Options = model.GetCommonOptions(options.Options, opts...)
	// 加载实现特定选项
	options = model.GetImplSpecificOptions(options, opts...)

	// 2. 开始生成前的回调
	ctx = callbacks.OnStart(ctx, &model.CallbackInput{
		Messages: messages,
		Config: &model.Config{
			Model: *options.Options.Model,
		},
	})

    // 3. 创建流式响应
    // Pipe产生一个StreamReader和一个StreamWrite，向StreamWrite中写入可以从StreamReader中读到，二者并发安全。
    // 实现中异步向StreamWrite中写入生成内容，返回StreamReader作为返回值
    // ***StreamReader是一个数据流，仅可读一次，组件自行实现Callback时，既需要通过OnEndWithCallbackOutput向callback传递数据流，也需要向返回一个数据流，需要对数据流进行一次拷贝
    // 考虑到此种情形总是需要拷贝数据流，OnEndWithCallbackOutput函数会在内部拷贝并返回一个未被读取的流
	// 以下代码演示了一种流处理方式，处理方式不唯一
	sr, sw := schema.Pipe[*model.CallbackOutput](1)

	// 4. 启动异步生成
	go func() {
		defer sw.Close()
		// 流式写入
		m.doStream(ctx, messages, options, sw)
	}()

	// 5. 完成回调
	_, nsr := callbacks.OnEndWithStreamOutput(ctx, sr)

	return schema.StreamReaderWithConvert(nsr, func(t *model.CallbackOutput) (*schema.Message, error) {
		return t.Message, nil
	}), nil
}

func (m *DeepSeekModel) BindTools(tools []*schema.ToolInfo) error {
	// 实现工具绑定逻辑
	return nil
}

func (m *DeepSeekModel) doGenerate(ctx context.Context, messages []*schema.Message, opts *DeepSeekModelOptions) (*schema.Message, error) {
	// 实现生成逻辑
	return nil, nil
}

func (m *DeepSeekModel) doStream(ctx context.Context, messages []*schema.Message, opts *DeepSeekModelOptions, sr *schema.StreamWriter[*model.CallbackOutput]) {
	// 流式生成文本写入sr中
	return
}
