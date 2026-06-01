package plugin

import (
	"reflect"
	"strings"
	"testing"
)

type testPlugin struct {
	name      string
	version   string
	dependsOn []string
}

func (p testPlugin) Name() string { return p.name }

func (p testPlugin) Version() string { return p.version }

func (p testPlugin) DependsOn() []string { return append([]string(nil), p.dependsOn...) }

func (p testPlugin) Register(_ *Context) error { return nil }

func (p testPlugin) Boot(_ *Context) error { return nil }

func (p testPlugin) Shutdown(_ *Context) error { return nil }

// TestManagerOrderedUsesDependencyOrderAndAlphabeticalTieBreak 验证同一批插件在
// 不同注册顺序下仍会按依赖和字母序得到稳定的运行时顺序。
func TestManagerOrderedUsesDependencyOrderAndAlphabeticalTieBreak(t *testing.T) {
	manager := NewManager()
	input := []Plugin{
		testPlugin{name: "user", version: "0.1.0"},
		testPlugin{name: "scheduler", version: "0.1.0"},
		testPlugin{name: "rbac", version: "0.1.0", dependsOn: []string{"user"}},
		testPlugin{name: "audit", version: "0.1.0"},
	}

	for _, current := range input {
		if err := manager.RegisterPlugin(current); err != nil {
			t.Fatalf("register plugin %s: %v", current.Name(), err)
		}
	}

	ordered, err := manager.Ordered()
	if err != nil {
		t.Fatalf("order plugins: %v", err)
	}

	got := make([]string, 0, len(ordered))
	for _, current := range ordered {
		got = append(got, current.Name())
	}

	expected := []string{"audit", "scheduler", "user", "rbac"}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}

// TestManagerOrderedRejectsMissingDependency 验证缺失依赖会在排序阶段直接阻断。
func TestManagerOrderedRejectsMissingDependency(t *testing.T) {
	manager := NewManager()
	if err := manager.RegisterPlugin(testPlugin{name: "rbac", version: "0.1.0", dependsOn: []string{"user"}}); err != nil {
		t.Fatalf("register plugin: %v", err)
	}

	_, err := manager.Ordered()
	if err == nil {
		t.Fatal("expected missing dependency error")
	}
	if !strings.Contains(err.Error(), "depends on missing plugin user") {
		t.Fatalf("expected missing dependency error, got %v", err)
	}
}

// TestManagerOrderedRejectsDependencyCycle 验证循环依赖会被明确识别。
func TestManagerOrderedRejectsDependencyCycle(t *testing.T) {
	manager := NewManager()
	for _, current := range []Plugin{
		testPlugin{name: "user", version: "0.1.0", dependsOn: []string{"rbac"}},
		testPlugin{name: "rbac", version: "0.1.0", dependsOn: []string{"user"}},
	} {
		if err := manager.RegisterPlugin(current); err != nil {
			t.Fatalf("register plugin %s: %v", current.Name(), err)
		}
	}

	_, err := manager.Ordered()
	if err == nil {
		t.Fatal("expected dependency cycle error")
	}
	if !strings.Contains(err.Error(), "plugin dependency cycle detected") {
		t.Fatalf("expected dependency cycle error, got %v", err)
	}
}

// TestOrderModuleSpecsIsIndependentFromInputOrder 验证模块定义排序不依赖生成输入顺序。
func TestOrderModuleSpecsIsIndependentFromInputOrder(t *testing.T) {
	input := []ModuleSpec{
		{ID: "scheduler", ModuleVersion: "0.1.0", Builder: BuilderFunc(func(BuildContext) (Plugin, error) {
			return testPlugin{name: "scheduler", version: "0.1.0"}, nil
		})},
		{ID: "rbac", ModuleVersion: "0.1.0", Dependencies: []string{"user"}, Builder: BuilderFunc(func(BuildContext) (Plugin, error) {
			return testPlugin{name: "rbac", version: "0.1.0", dependsOn: []string{"user"}}, nil
		})},
		{ID: "audit", ModuleVersion: "0.1.0", Builder: BuilderFunc(func(BuildContext) (Plugin, error) {
			return testPlugin{name: "audit", version: "0.1.0"}, nil
		})},
		{ID: "user", ModuleVersion: "0.1.0", Builder: BuilderFunc(func(BuildContext) (Plugin, error) {
			return testPlugin{name: "user", version: "0.1.0"}, nil
		})},
	}

	ordered, err := OrderModuleSpecs(input)
	if err != nil {
		t.Fatalf("order module specs: %v", err)
	}

	got := make([]string, 0, len(ordered))
	for _, current := range ordered {
		got = append(got, current.Name())
	}

	expected := []string{"audit", "scheduler", "user", "rbac"}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}

// TestModuleSpecBuildWrapsCanonicalMetadata 验证模块定义构造出的运行时插件以
// 模块定义元数据为 canonical truth。
func TestModuleSpecBuildWrapsCanonicalMetadata(t *testing.T) {
	descriptor := ModuleSpec{
		ID:            "rbac",
		ModuleVersion: "0.2.0",
		Dependencies:  []string{"user"},
		Builder: BuilderFunc(func(BuildContext) (Plugin, error) {
			return testPlugin{name: "rbac", version: "0.2.0", dependsOn: []string{"user"}}, nil
		}),
	}

	built, err := descriptor.Build(BuildContext{})
	if err != nil {
		t.Fatalf("build descriptor: %v", err)
	}

	if built.Name() != "rbac" {
		t.Fatalf("expected descriptor name, got %q", built.Name())
	}
	if built.Version() != "0.2.0" {
		t.Fatalf("expected descriptor version, got %q", built.Version())
	}
	if !reflect.DeepEqual(built.DependsOn(), []string{"user"}) {
		t.Fatalf("expected descriptor dependencies, got %v", built.DependsOn())
	}
}

// TestModuleSpecBuildRejectsRuntimeMetadataDrift 验证运行时插件元数据一旦偏离
// 模块定义真相就会在构造阶段被阻断。
func TestModuleSpecBuildRejectsRuntimeMetadataDrift(t *testing.T) {
	descriptor := ModuleSpec{
		ID:            "rbac",
		ModuleVersion: "0.2.0",
		Dependencies:  []string{"user"},
		Builder: BuilderFunc(func(BuildContext) (Plugin, error) {
			return testPlugin{name: "rbac-v2", version: "0.2.0", dependsOn: []string{"user"}}, nil
		}),
	}

	_, err := descriptor.Build(BuildContext{})
	if err == nil {
		t.Fatal("expected module spec metadata drift error")
	}
	if !strings.Contains(err.Error(), "does not match module spec") {
		t.Fatalf("expected module spec mismatch error, got %v", err)
	}
}

func TestNewRuntimeMetadataPreservesOrderedDescriptorSnapshot(t *testing.T) {
	metadata := NewRuntimeMetadata([]ModuleSpec{
		{ID: "audit", ModuleVersion: "0.1.0"},
		{ID: "user", ModuleVersion: "0.2.0"},
		{ID: "rbac", ModuleVersion: "0.3.0", Dependencies: []string{"user"}},
	})

	got := metadata.OrderedModuleDescriptors()
	expected := []DescriptorSnapshot{
		{Name: "audit", Version: "0.1.0"},
		{Name: "user", Version: "0.2.0"},
		{Name: "rbac", Version: "0.3.0", DependsOn: []string{"user"}},
	}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("expected %v, got %v", expected, got)
	}

	got[0].Name = "mutated"
	got[2].DependsOn[0] = "mutated"

	unchanged := metadata.OrderedModuleDescriptors()
	if !reflect.DeepEqual(unchanged, expected) {
		t.Fatalf("expected runtime metadata to remain immutable, got %v", unchanged)
	}
}
