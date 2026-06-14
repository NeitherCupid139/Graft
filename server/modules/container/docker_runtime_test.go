// Copyright (c) 2025-2026 GeWuYou
// SPDX-License-Identifier: Apache-2.0

package container

import (
	"bytes"
	"context"
	"errors"
	"io"
	"reflect"
	"sync/atomic"
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/pkg/stdcopy"
)

func TestReadDockerLogLinesDoesNotReportExactTailAsTruncated(t *testing.T) {
	t.Parallel()

	lines, truncated, err := readDockerLogLines(dockerLogStream(t, "one\n", "two\n"), 2)
	if err != nil {
		t.Fatalf("read log lines: %v", err)
	}
	if truncated {
		t.Fatalf("expected exactly tail lines to avoid truncation")
	}
	if !reflect.DeepEqual(lines, []string{"one", "two"}) {
		t.Fatalf("unexpected lines %#v", lines)
	}
}

func TestReadDockerLogLinesTruncatesOnlyAfterDiscardingLines(t *testing.T) {
	t.Parallel()

	lines, truncated, err := readDockerLogLines(dockerLogStream(t, "one\n", "two\n", "three\n"), 2)
	if err != nil {
		t.Fatalf("read log lines: %v", err)
	}
	if !truncated {
		t.Fatalf("expected more than tail lines to report truncation")
	}
	if !reflect.DeepEqual(lines, []string{"two", "three"}) {
		t.Fatalf("unexpected lines %#v", lines)
	}
}

func TestReadDockerLogLinesAvoidsUserSizedPreallocation(t *testing.T) {
	t.Parallel()

	const excessiveTail = int(^uint(0) >> 1)
	lines, truncated, err := readDockerLogLines(dockerLogStream(t, "one\n"), excessiveTail)
	if err != nil {
		t.Fatalf("read log lines: %v", err)
	}
	if truncated {
		t.Fatalf("expected one line to avoid truncation")
	}
	if !reflect.DeepEqual(lines, []string{"one"}) {
		t.Fatalf("unexpected lines %#v", lines)
	}
}

func TestDockerRuntimeLogsAvoidsRuntimeInfoCall(t *testing.T) {
	t.Parallel()

	client := &countingDockerClient{
		logReader: dockerLogReadCloser(t, "one\n", "two\n"),
		inspect: container.InspectResponse{
			ContainerJSONBase: &container.ContainerJSONBase{
				ID:   "full-id",
				Name: "/web",
			},
		},
	}
	runtime := &DockerRuntime{client: client, endpoint: "unix:///var/run/docker.sock"}

	logs, err := runtime.Logs(context.Background(), Ref{Value: "web"}, LogQuery{
		Tail:   2,
		Stdout: true,
	})
	if err != nil {
		t.Fatalf("logs: %v", err)
	}
	if logs.ID != "full-id" || logs.Name != "web" {
		t.Fatalf("unexpected log metadata %#v", logs)
	}
	if calls := client.infoCalls.Load(); calls != 0 {
		t.Fatalf("expected logs to avoid Info calls, got %d", calls)
	}
	if calls := client.inspectCalls.Load(); calls != 1 {
		t.Fatalf("expected one inspect call for log metadata, got %d", calls)
	}
	if calls := client.logCalls.Load(); calls != 1 {
		t.Fatalf("expected one log call, got %d", calls)
	}
}

func TestDockerRuntimeLogsReturnsInvalidLogQueryError(t *testing.T) {
	t.Parallel()

	runtime := &DockerRuntime{client: &countingDockerClient{}, endpoint: "unix:///var/run/docker.sock"}

	_, err := runtime.Logs(context.Background(), Ref{Value: "web"}, LogQuery{
		Tail:   1,
		Since:  "not-a-time",
		Stdout: true,
	})
	if !errors.Is(err, errInvalidLogQuery) {
		t.Fatalf("expected invalid log query error, got %v", err)
	}
}

func TestDockerRuntimeListUsesCheapSummaryFields(t *testing.T) {
	t.Parallel()

	client := &countingDockerClient{
		list: []container.Summary{
			{
				ID:      "1234567890abcdef",
				Names:   []string{"/graft-web"},
				Image:   "graft/web:latest",
				ImageID: "sha256:web",
				Labels: map[string]string{
					composeProjectLabel: "graft",
					composeServiceLabel: "web",
				},
				Ports:  []container.Port{{IP: "0.0.0.0", PrivatePort: 80, PublicPort: 8080, Type: "tcp"}},
				State:  container.StateRunning,
				Status: "Up 10 minutes",
				NetworkSettings: &container.NetworkSettingsSummary{
					Networks: map[string]*network.EndpointSettings{
						"bridge": {
							NetworkID:  "net-1",
							EndpointID: "endpoint-1",
							Gateway:    "172.18.0.1",
							IPAddress:  "172.18.0.2",
							MacAddress: "02:42:ac:12:00:02",
						},
					},
				},
				Created: 1781409600,
			},
		},
	}
	runtime := &DockerRuntime{client: client, endpoint: "unix:///var/run/docker.sock"}

	items, err := runtime.List(context.Background(), ListQuery{})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected one item, got %#v", items)
	}
	item := items[0]
	if item.ShortID != "1234567890ab" || item.Name != "graft-web" {
		t.Fatalf("unexpected identity fields %#v", item)
	}
	if item.PrimaryIP != "172.18.0.2" || item.NetworkSummary != "bridge" {
		t.Fatalf("unexpected network fields %#v", item)
	}
	if item.Health != containerHealthUnavailable || item.Resource.Available {
		t.Fatalf("expected unavailable health/resource semantics, got %#v", item)
	}
	if item.ComposeProject != "graft" || item.ComposeService != "web" {
		t.Fatalf("unexpected compose fields %#v", item)
	}
	if !item.CanStop || !item.CanRestart || item.CanStart {
		t.Fatalf("unexpected action availability %#v", item)
	}
	if calls := client.inspectCalls.Load(); calls != 0 {
		t.Fatalf("expected list to avoid inspect calls, got %d", calls)
	}
}

func dockerLogStream(t *testing.T, chunks ...string) io.Reader {
	t.Helper()

	var output bytes.Buffer
	writer := stdcopy.NewStdWriter(&output, stdcopy.Stdout)
	for _, chunk := range chunks {
		if _, err := writer.Write([]byte(chunk)); err != nil {
			t.Fatalf("write log chunk: %v", err)
		}
	}
	return bytes.NewReader(output.Bytes())
}

func dockerLogReadCloser(t *testing.T, chunks ...string) io.ReadCloser {
	t.Helper()

	return io.NopCloser(dockerLogStream(t, chunks...))
}

type countingDockerClient struct {
	infoCalls    atomic.Int64
	inspectCalls atomic.Int64
	logCalls     atomic.Int64
	listCalls    atomic.Int64
	logReader    io.ReadCloser
	inspect      container.InspectResponse
	list         []container.Summary
}

func (c *countingDockerClient) Info(context.Context) (systemInfo, error) {
	c.infoCalls.Add(1)
	return dockerClientSystemInfo{}, nil
}

func (c *countingDockerClient) ContainerList(context.Context, container.ListOptions) ([]container.Summary, error) {
	c.listCalls.Add(1)
	return c.list, nil
}

func (c *countingDockerClient) ContainerInspect(context.Context, string) (container.InspectResponse, error) {
	c.inspectCalls.Add(1)
	return c.inspect, nil
}

func (c *countingDockerClient) ContainerLogs(context.Context, string, container.LogsOptions) (io.ReadCloser, error) {
	c.logCalls.Add(1)
	if c.logReader == nil {
		return io.NopCloser(bytes.NewReader(nil)), nil
	}
	return c.logReader, nil
}

func (c *countingDockerClient) ContainerStart(context.Context, string, container.StartOptions) error {
	return nil
}

func (c *countingDockerClient) ContainerStop(context.Context, string, container.StopOptions) error {
	return nil
}

func (c *countingDockerClient) ContainerRestart(context.Context, string, container.StopOptions) error {
	return nil
}

func (c *countingDockerClient) Close() error {
	return nil
}
