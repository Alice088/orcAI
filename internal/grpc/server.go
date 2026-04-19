package grpc

import (
	"OrcAI/internal/broker/task"
	"OrcAI/internal/core"
	"OrcAI/internal/grpc/pb"
	"OrcAI/pkg/str"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/mailru/easyjson"
	"github.com/nats-io/nats.go"
)

type Event struct {
	Type   string
	Log    string
	Level  string
	Status string
}

type BridgeServer struct {
	NC     *nats.Conn
	Memory sync.Map
	Router core.AgentRouter
	pb.UnimplementedBridgeServiceServer
}

func (s *BridgeServer) StartSubscribers() error {
	_, err := s.NC.Subscribe("task.*.*", func(msg *nats.Msg) {
		var taskID string

		parts := strings.Split(msg.Subject, ".")
		if len(parts) < 3 {
			return
		}
		taskID = parts[2]

		chRaw, _ := s.Memory.LoadOrStore(taskID, make(chan Event, 100))
		ch := chRaw.(chan Event)

		switch {
		case strings.HasPrefix(msg.Subject, "task.logs"):
			var log task.Log
			if json.Unmarshal(msg.Data, &log) != nil {
				return
			}
			ch <- Event{
				Type:  "log",
				Log:   log.Message,
				Level: log.Level,
			}

		case strings.HasPrefix(msg.Subject, "task.progress"):
			var p task.Progress
			if json.Unmarshal(msg.Data, &p) != nil {
				return
			}
			ch <- Event{
				Type:   "progress",
				Status: p.Status,
			}

		case strings.HasPrefix(msg.Subject, "task.result"):
			ch <- Event{Type: "result"}
		}
	})

	return err
}

func (s *BridgeServer) Ask(_ context.Context, req *pb.AskRequest) (*pb.AskResponse, error) {
	id := uuid.NewString()

	routeType, err := s.Router.Route(req.Input)
	if err != nil {
		return nil, err
	}

	t := task.Created{
		TaskID: id,
		Type:   routeType,
		Payload: task.CreatedPayload{
			Msg: req.Input,
		},
	}

	raw, err := easyjson.Marshal(t)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal task.created.%s: %w", str.TrimChars(id, 10), err)
	}

	err = s.NC.Publish("task.created."+id, raw)
	if err != nil {
		return nil, fmt.Errorf("failed to publish -> [task.created.%s]: %w", str.TrimChars(id, 10), err)
	}

	if err = s.NC.Flush(); err != nil {
		return nil, fmt.Errorf("failed to flush nats: %w", err)
	}

	return &pb.AskResponse{TaskId: id}, nil
}

func (s *BridgeServer) Logs(req *pb.LogRequest, stream pb.BridgeService_LogsServer) error {
	ctx := stream.Context()

	chRaw, ok := s.Memory.Load(req.TaskId)
	if !ok {
		return fmt.Errorf("task not found")
	}

	ch := chRaw.(chan Event)

	for {
		select {
		case <-ctx.Done():
			return nil

		case e := <-ch:
			if e.Type == "log" {
				if err := stream.Send(&pb.LogResponse{
					Log:    e.Log,
					TaskId: req.TaskId,
					Level:  e.Level,
				}); err != nil {
					return err
				}
			}

			if e.Type == "result" {
				return nil
			}
		}
	}
}

func (s *BridgeServer) Status(req *pb.StatusRequest, stream pb.BridgeService_StatusServer) error {
	ctx := stream.Context()

	chRaw, ok := s.Memory.Load(req.TaskId)
	if !ok {
		return fmt.Errorf("task not found")
	}

	ch := chRaw.(chan Event)

	for {
		select {
		case <-ctx.Done():
			return nil

		case e := <-ch:
			if e.Type == "progress" {
				if err := stream.Send(&pb.StatusResponse{
					TaskId: req.TaskId,
					Status: e.Status,
				}); err != nil {
					return err
				}
			}

			if e.Type == "result" {
				return nil
			}
		}
	}
}
