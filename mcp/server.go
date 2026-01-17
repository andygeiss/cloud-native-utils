package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"sync"
)

// Server is an MCP server that handles JSON-RPC requests over STDIO.
type Server struct {
	tools       map[string]Tool
	reader      *bufio.Reader
	writer      io.Writer
	name        string
	version     string
	mu          sync.RWMutex
	initialized bool
}

// NewServer creates a new MCP server with the given name and version.
func NewServer(name, version string) *Server {
	return &Server{
		name:    name,
		version: version,
		tools:   make(map[string]Tool),
		reader:  bufio.NewReader(os.Stdin),
		writer:  os.Stdout,
	}
}

// NewServerWithIO creates a new MCP server with custom IO (for testing).
func NewServerWithIO(name, version string, reader io.Reader, writer io.Writer) *Server {
	return &Server{
		name:    name,
		version: version,
		tools:   make(map[string]Tool),
		reader:  bufio.NewReader(reader),
		writer:  writer,
	}
}

// RegisterTool registers a tool with the server.
func (a *Server) RegisterTool(tool Tool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.tools[tool.Definition.Name] = tool
}

// Serve starts the server and processes requests until context is canceled.
func (a *Server) Serve(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := a.handleRequest(ctx); err != nil {
				if errors.Is(err, io.EOF) {
					return nil
				}
			}
		}
	}
}

// handleRequest reads a single request from the reader.
func (a *Server) handleRequest(ctx context.Context) error {
	line, err := a.reader.ReadBytes('\n')
	if err != nil {
		return err
	}

	var req Request
	if err := json.Unmarshal(line, &req); err != nil {
		return a.writeResponse(NewErrorResponse(nil, ErrorCodeParse, "Parse error"))
	}

	resp := a.routeRequest(ctx, req)
	return a.writeResponse(resp)
}

// routeRequest routes the request to the appropriate handler.
func (a *Server) routeRequest(ctx context.Context, req Request) Response {
	switch req.Method {
	case "initialize":
		return a.handleInitialize(req)
	case "initialized":
		return a.handleInitialized()
	case "tools/list":
		return a.handleToolsList(req)
	case "tools/call":
		return a.handleToolsCall(ctx, req)
	default:
		return NewErrorResponse(req.ID, ErrorCodeMethodNotFound, "Method not found")
	}
}

// handleInitialize handles the initialize request.
func (a *Server) handleInitialize(req Request) Response {
	a.mu.Lock()
	a.initialized = true
	a.mu.Unlock()

	result := InitializeResult{
		ProtocolVersion: "2024-11-05",
		Capabilities: Capabilities{
			Tools: &ToolsCapability{},
		},
		ServerInfo: Implementation{
			Name:    a.name,
			Version: a.version,
		},
	}

	return NewResponse(req.ID, result)
}

// handleInitialized handles the initialized notification.
func (a *Server) handleInitialized() Response {
	return Response{}
}

// handleToolsList handles the tools/list request.
func (a *Server) handleToolsList(req Request) Response {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if !a.initialized {
		return NewErrorResponse(req.ID, ErrorCodeInternal, ErrMCPNotInitialized.Error())
	}

	tools := make([]ToolDefinition, 0, len(a.tools))
	for _, tool := range a.tools {
		tools = append(tools, tool.Definition)
	}

	return NewResponse(req.ID, ToolsListResult{Tools: tools})
}

// handleToolsCall handles the tools/call request.
func (a *Server) handleToolsCall(ctx context.Context, req Request) Response {
	a.mu.RLock()
	initialized := a.initialized
	a.mu.RUnlock()

	if !initialized {
		return NewErrorResponse(req.ID, ErrorCodeInternal, ErrMCPNotInitialized.Error())
	}

	var params ToolsCallParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return NewErrorResponse(req.ID, ErrorCodeInvalidParams, "Invalid params")
	}

	a.mu.RLock()
	tool, ok := a.tools[params.Name]
	a.mu.RUnlock()

	if !ok {
		return NewErrorResponse(req.ID, ErrorCodeInvalidParams, ErrMCPToolNotFound.Error())
	}

	result, err := tool.Handler(ctx, params)
	if err != nil {
		return NewResponse(req.ID, ToolsCallResult{
			Content: []ContentBlock{NewTextContent(err.Error())},
			IsError: true,
		})
	}

	return NewResponse(req.ID, result)
}

// writeResponse writes a response to the output writer.
func (a *Server) writeResponse(resp Response) error {
	if resp.JSONRPC == "" {
		return nil
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	_, err = a.writer.Write(append(data, '\n'))
	return err
}
