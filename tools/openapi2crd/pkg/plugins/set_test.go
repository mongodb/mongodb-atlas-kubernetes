package plugins

import (
	"errors"
	"fmt"
	"testing"

	configv1alpha1 "github.com/mongodb/atlas2crd/pkg/apis/config/v1alpha1"
	"github.com/stretchr/testify/require"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
)

func TestNewPluginSet(t *testing.T) {
	tests := map[string]struct {
		sets    []configv1alpha1.PluginSet
		want    []PluginSet
		wantErr error
	}{
		"fail to order sets": {
			sets: []configv1alpha1.PluginSet{
				{Name: "set1", InheritFrom: "set2", Plugins: []string{"crd"}},
				{Name: "set2", InheritFrom: "set1", Plugins: []string{"status"}},
			},
			want:    nil,
			wantErr: fmt.Errorf("circular dependency detected for plugin set: %s", "set1"),
		},
		"fail when having multiple default sets": {
			sets: []configv1alpha1.PluginSet{
				{Name: "set1", Default: true, Plugins: []string{"crd"}},
				{Name: "set2", Default: true, Plugins: []string{"status"}},
			},
			want:    nil,
			wantErr: errors.New("multiple default plugin sets defined"),
		},
		"fail when plugin is not on the catalog": {
			sets: []configv1alpha1.PluginSet{
				{Name: "set1", Default: true, Plugins: []string{"crd"}},
				{Name: "set2", Plugins: []string{"nonexistent"}},
			},
			want: nil,
			wantErr: fmt.Errorf(
				"failed to build plugin set %s: error getting plugin %s: %w", "set2", ""+
					"nonexistent",
				fmt.Errorf("plugin %s not found", "nonexistent"),
			),
		},
		"new plugin set": {
			sets: []configv1alpha1.PluginSet{
				{Name: "set1", Default: true, Plugins: []string{"crd"}},
				{Name: "set2", Plugins: []string{"status"}},
			},
			want: []PluginSet{
				{
					Name:    "set1",
					Default: true,
					Plugins: map[string]Plugin{
						"crd": NewCrdPlugin(),
					},
				},
				{
					Name: "set2",
					Plugins: map[string]Plugin{
						"status": NewStatusPlugin(&apiextensions.CustomResourceDefinition{}),
					},
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := NewPluginSet(tt.sets, NewPluginCatalog(nil))
			require.Equal(t, tt.wantErr, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestOrderPluginSet(t *testing.T) {
	tests := map[string]struct {
		sets    []configv1alpha1.PluginSet
		want    []configv1alpha1.PluginSet
		wantErr error
	}{
		"empty": {
			sets: []configv1alpha1.PluginSet{},
			want: []configv1alpha1.PluginSet{},
		},
		"single set": {
			sets: []configv1alpha1.PluginSet{
				{
					Name:    "set1",
					Plugins: []string{"plugin1", "plugin2"},
				},
			},
			want: []configv1alpha1.PluginSet{
				{
					Name:    "set1",
					Plugins: []string{"plugin1", "plugin2"},
				},
			},
		},
		"multiple sets with no dependency": {
			sets: []configv1alpha1.PluginSet{
				{
					Name:    "set1",
					Plugins: []string{"plugin1", "plugin2"},
				},
				{
					Name:    "set2",
					Plugins: []string{"plugin3", "plugin4"},
				},
			},
			want: []configv1alpha1.PluginSet{
				{
					Name:    "set1",
					Plugins: []string{"plugin1", "plugin2"},
				},
				{
					Name:    "set2",
					Plugins: []string{"plugin3", "plugin4"},
				},
			},
		},
		"multiple sets with dependency": {
			sets: []configv1alpha1.PluginSet{
				{
					Name:        "set1",
					Plugins:     []string{"plugin1", "plugin2"},
					InheritFrom: "set2",
				},
				{
					Name:    "set2",
					Plugins: []string{"plugin3", "plugin4"},
				},
			},
			want: []configv1alpha1.PluginSet{
				{
					Name:    "set2",
					Plugins: []string{"plugin3", "plugin4"},
				},
				{
					Name:        "set1",
					Plugins:     []string{"plugin1", "plugin2"},
					InheritFrom: "set2",
				},
			},
		},
		"multiple sets with depth dependency": {
			sets: []configv1alpha1.PluginSet{
				{
					Name:        "set1",
					Plugins:     []string{"plugin1", "plugin2"},
					InheritFrom: "set2",
				},
				{
					Name:        "set2",
					Plugins:     []string{"plugin3", "plugin4"},
					InheritFrom: "set3",
				},
				{
					Name:    "set3",
					Plugins: []string{"plugin5", "plugin6"},
				},
			},
			want: []configv1alpha1.PluginSet{
				{
					Name:    "set3",
					Plugins: []string{"plugin5", "plugin6"},
				},
				{
					Name:        "set2",
					Plugins:     []string{"plugin3", "plugin4"},
					InheritFrom: "set3",
				},
				{
					Name:        "set1",
					Plugins:     []string{"plugin1", "plugin2"},
					InheritFrom: "set2",
				},
			},
		},
		"multiple sets with depth mixed dependency": {
			sets: []configv1alpha1.PluginSet{
				{
					Name:        "set1",
					Plugins:     []string{"plugin1", "plugin2"},
					InheritFrom: "set3",
				},
				{
					Name:        "set2",
					Plugins:     []string{"plugin3", "plugin4"},
					InheritFrom: "set1",
				},
				{
					Name:    "set3",
					Plugins: []string{"plugin5", "plugin6"},
				},
			},
			want: []configv1alpha1.PluginSet{
				{
					Name:    "set3",
					Plugins: []string{"plugin5", "plugin6"},
				},
				{
					Name:        "set1",
					Plugins:     []string{"plugin1", "plugin2"},
					InheritFrom: "set3",
				},
				{
					Name:        "set2",
					Plugins:     []string{"plugin3", "plugin4"},
					InheritFrom: "set1",
				},
			},
		},
		"circular dependency": {
			sets: []configv1alpha1.PluginSet{
				{
					Name:        "set1",
					Plugins:     []string{"plugin1", "plugin2"},
					InheritFrom: "set2",
				},
				{
					Name:        "set2",
					Plugins:     []string{"plugin3", "plugin4"},
					InheritFrom: "set3",
				},
				{
					Name:        "set3",
					Plugins:     []string{"plugin5", "plugin6"},
					InheritFrom: "set1",
				},
			},
			want:    nil,
			wantErr: fmt.Errorf("circular dependency detected for plugin set: %s", "set1"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := orderPluginSet(tt.sets)
			require.Equal(t, tt.wantErr, err)
			require.Equal(t, tt.want, got)
		})
	}
}
